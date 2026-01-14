import nipplejs from 'nipplejs';

/**
 * TouchControls - Completely isolated mobile input handler
 * Does not modify any existing input systems - uses callbacks directly
 */
export class TouchControls {
  constructor(options) {
    this.camera = options.camera;
    this.gameState = options.gameState;
    this.onMove = options.onMove;
    this.onShoot = options.onShoot;
    this.onInteract = options.onInteract;
    this.canvas = options.canvas;

    this.moveJoystick = null;
    this.aimJoystick = null;

    // Firing state
    this.isAiming = false;
    this.aimDirection = { x: 0, z: 0 };
    this.fireInterval = null;
    this.fireRate = 150;

    // Camera touch state (for two-finger gestures)
    this.cameraTouch = {
      active: false,
      lastTouches: null,
      initialDistance: null
    };

    this.enabled = false;
    this.boundHandleInteract = this.handleInteract.bind(this);
    this.boundHandleCameraTouchStart = this.handleCameraTouchStart.bind(this);
    this.boundHandleCameraTouchMove = this.handleCameraTouchMove.bind(this);
    this.boundHandleCameraTouchEnd = this.handleCameraTouchEnd.bind(this);
  }

  enable() {
    if (this.enabled) return;
    this.enabled = true;

    document.body.classList.add('mobile');
    this.createJoysticks();
    this.setupInteractButton();
    this.setupCameraTouch();
  }

  disable() {
    if (!this.enabled) return;
    this.enabled = false;

    document.body.classList.remove('mobile');
    this.destroyJoysticks();
    this.removeInteractButton();
    this.removeCameraTouch();
    this.stopFiring();

    // Send stop movement
    if (this.onMove) {
      this.onMove({ x: 0, y: 0, z: 0 });
    }
  }

  createJoysticks() {
    const leftZone = document.getElementById('joystick-left');
    const rightZone = document.getElementById('joystick-right');

    if (leftZone) {
      this.moveJoystick = nipplejs.create({
        zone: leftZone,
        mode: 'static',
        position: { left: '50%', top: '50%' },
        color: 'rgba(255, 255, 255, 0.5)',
        size: 100
      });

      this.moveJoystick.on('move', (evt, data) => {
        if (!data.direction || !this.enabled) return;
        const { dx, dz } = this.joystickToWorldDirection(data.angle.radian, data.force);
        if (this.onMove) {
          this.onMove({ x: dx, y: 0, z: dz });
        }
      });

      this.moveJoystick.on('end', () => {
        if (this.onMove) {
          this.onMove({ x: 0, y: 0, z: 0 });
        }
      });
    }

    if (rightZone) {
      this.aimJoystick = nipplejs.create({
        zone: rightZone,
        mode: 'static',
        position: { left: '50%', top: '50%' },
        color: 'rgba(239, 68, 68, 0.5)',
        size: 100
      });

      this.aimJoystick.on('move', (evt, data) => {
        if (!data.direction || !this.enabled) return;
        const { dx, dz } = this.joystickToWorldDirection(data.angle.radian, data.force);
        this.aimDirection = { x: dx, z: dz };
        if (!this.isAiming) {
          this.startFiring();
        }
      });

      this.aimJoystick.on('end', () => {
        this.stopFiring();
      });
    }
  }

  destroyJoysticks() {
    if (this.moveJoystick) {
      this.moveJoystick.destroy();
      this.moveJoystick = null;
    }
    if (this.aimJoystick) {
      this.aimJoystick.destroy();
      this.aimJoystick = null;
    }
  }

  joystickToWorldDirection(angle, force) {
    const clampedForce = Math.min(force, 1);
    const camForward = this.camera.getForwardDirection();
    const camRight = this.camera.getRightDirection();

    // nipplejs: angle 0 = right, PI/2 = up
    const joyX = Math.cos(angle) * clampedForce;
    const joyY = Math.sin(angle) * clampedForce;

    // Transform: joyY = forward/back, joyX = right/left
    const dx = joyY * camForward.x + joyX * camRight.x;
    const dz = joyY * camForward.z + joyX * camRight.z;

    return { dx, dz };
  }

  startFiring() {
    this.isAiming = true;
    this.fire();
    this.fireInterval = setInterval(() => this.fire(), this.fireRate);
  }

  stopFiring() {
    this.isAiming = false;
    this.aimDirection = { x: 0, z: 0 };
    if (this.fireInterval) {
      clearInterval(this.fireInterval);
      this.fireInterval = null;
    }
  }

  fire() {
    if (!this.enabled || !this.onShoot) return;

    const player = this.gameState?.getMyPlayer();
    if (!player) return;

    const range = 50;
    const targetX = player.position.x + this.aimDirection.x * range;
    const targetZ = player.position.z + this.aimDirection.z * range;

    this.onShoot(targetX, targetZ);
  }

  setupInteractButton() {
    const btn = document.getElementById('interact-btn');
    if (btn) {
      btn.addEventListener('touchstart', this.boundHandleInteract, { passive: false });
    }
  }

  removeInteractButton() {
    const btn = document.getElementById('interact-btn');
    if (btn) {
      btn.removeEventListener('touchstart', this.boundHandleInteract);
    }
  }

  handleInteract(e) {
    e.preventDefault();
    if (this.onInteract) {
      this.onInteract();
    }
  }

  // Camera touch controls (two-finger rotation and pinch zoom)
  setupCameraTouch() {
    if (!this.canvas) return;
    this.canvas.addEventListener('touchstart', this.boundHandleCameraTouchStart, { passive: false });
    this.canvas.addEventListener('touchmove', this.boundHandleCameraTouchMove, { passive: false });
    this.canvas.addEventListener('touchend', this.boundHandleCameraTouchEnd);
    this.canvas.addEventListener('touchcancel', this.boundHandleCameraTouchEnd);
  }

  removeCameraTouch() {
    if (!this.canvas) return;
    this.canvas.removeEventListener('touchstart', this.boundHandleCameraTouchStart);
    this.canvas.removeEventListener('touchmove', this.boundHandleCameraTouchMove);
    this.canvas.removeEventListener('touchend', this.boundHandleCameraTouchEnd);
    this.canvas.removeEventListener('touchcancel', this.boundHandleCameraTouchEnd);
  }

  handleCameraTouchStart(e) {
    if (e.touches.length === 2) {
      e.preventDefault();
      this.cameraTouch.active = true;
      this.cameraTouch.lastTouches = this.getTouchData(e.touches);
      this.cameraTouch.initialDistance = this.cameraTouch.lastTouches.distance;
    }
  }

  handleCameraTouchMove(e) {
    if (!this.cameraTouch.active || e.touches.length !== 2) return;
    e.preventDefault();

    const current = this.getTouchData(e.touches);
    const last = this.cameraTouch.lastTouches;

    // Rotation from center point movement
    const deltaX = current.centerX - last.centerX;
    const deltaY = current.centerY - last.centerY;

    this.camera.azimuth -= deltaX * 0.005;
    this.camera.polar += deltaY * 0.005;
    this.camera.polar = Math.max(this.camera.minPolar, Math.min(this.camera.maxPolar, this.camera.polar));

    // Zoom from pinch
    const pinchDelta = this.cameraTouch.initialDistance - current.distance;
    this.camera.distance += pinchDelta * 0.1;
    this.camera.distance = Math.max(this.camera.minDistance, Math.min(this.camera.maxDistance, this.camera.distance));

    this.cameraTouch.lastTouches = current;
    this.cameraTouch.initialDistance = current.distance;
  }

  handleCameraTouchEnd(e) {
    if (e.touches.length < 2) {
      this.cameraTouch.active = false;
      this.cameraTouch.lastTouches = null;
      this.cameraTouch.initialDistance = null;
    }
  }

  getTouchData(touches) {
    const t0 = touches[0];
    const t1 = touches[1];
    const dx = t1.clientX - t0.clientX;
    const dy = t1.clientY - t0.clientY;
    return {
      centerX: (t0.clientX + t1.clientX) / 2,
      centerY: (t0.clientY + t1.clientY) / 2,
      distance: Math.sqrt(dx * dx + dy * dy)
    };
  }

  destroy() {
    this.disable();
  }
}
