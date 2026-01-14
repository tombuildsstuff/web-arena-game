import * as THREE from 'three';

export class PlayerInput {
  constructor(camera, renderer, onMove, onShoot, onBuyFromZone, onClaimTurret, gameLoop) {
    this.camera = camera;
    this.renderer = renderer;
    this.onMove = onMove;
    this.onShoot = onShoot;
    this.onBuyFromZone = onBuyFromZone;
    this.onClaimTurret = onClaimTurret;
    this.gameLoop = gameLoop;

    // Movement keys state
    this.keys = {
      w: false,
      a: false,
      s: false,
      d: false
    };

    // Raycaster for mouse position
    this.raycaster = new THREE.Raycaster();
    this.mouse = new THREE.Vector2();
    // Ground plane at player height (Y=1.5) for accurate aiming
    this.groundPlane = new THREE.Plane(new THREE.Vector3(0, 1, 0), -1.5);

    // Bind event handlers
    this.handleKeyDown = this.handleKeyDown.bind(this);
    this.handleKeyUp = this.handleKeyUp.bind(this);
    this.handleMouseMove = this.handleMouseMove.bind(this);
    this.handleMouseDown = this.handleMouseDown.bind(this);

    // Last sent direction (to avoid spamming)
    this.lastDirection = { x: 0, z: 0 };

    // Target position for shooting
    this.targetPosition = { x: 0, z: 0 };

    this.enabled = false;
  }

  enable() {
    if (this.enabled) return;
    this.enabled = true;

    window.addEventListener('keydown', this.handleKeyDown);
    window.addEventListener('keyup', this.handleKeyUp);
    window.addEventListener('mousemove', this.handleMouseMove);
    window.addEventListener('mousedown', this.handleMouseDown);
  }

  disable() {
    if (!this.enabled) return;
    this.enabled = false;

    window.removeEventListener('keydown', this.handleKeyDown);
    window.removeEventListener('keyup', this.handleKeyUp);
    window.removeEventListener('mousemove', this.handleMouseMove);
    window.removeEventListener('mousedown', this.handleMouseDown);

    // Reset keys
    this.keys = { w: false, a: false, s: false, d: false };
    this.sendMovement();
  }

  handleKeyDown(event) {
    const key = event.key.toLowerCase();
    if (this.keys.hasOwnProperty(key) && !this.keys[key]) {
      this.keys[key] = true;
      this.sendMovement();
    }

    // E key for buying from zones or claiming turrets
    if (key === 'e') {
      this.handleInteraction();
    }
  }

  handleInteraction() {
    if (!this.gameLoop) return;

    // Priority: buy zones first, then turrets
    const nearbyZone = this.gameLoop.getNearbyBuyZone();
    if (nearbyZone && this.onBuyFromZone) {
      this.onBuyFromZone(nearbyZone.id);
      return;
    }

    const nearbyTurret = this.gameLoop.getNearbyTurret();
    if (nearbyTurret && this.onClaimTurret) {
      this.onClaimTurret(nearbyTurret.id);
    }
  }

  handleKeyUp(event) {
    const key = event.key.toLowerCase();
    if (this.keys.hasOwnProperty(key) && this.keys[key]) {
      this.keys[key] = false;
      this.sendMovement();
    }
  }

  handleMouseMove(event) {
    // Update mouse position
    const rect = this.renderer.domElement.getBoundingClientRect();
    this.mouse.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
    this.mouse.y = -((event.clientY - rect.top) / rect.height) * 2 + 1;

    // Update target position
    this.updateTargetPosition();
  }

  handleMouseDown(event) {
    // Left click to shoot
    if (event.button === 0 && this.enabled) {
      this.updateTargetPosition();
      if (this.onShoot) {
        this.onShoot(this.targetPosition.x, this.targetPosition.z);
      }
    }
  }

  updateTargetPosition() {
    const camera = this.camera.getCamera();
    this.raycaster.setFromCamera(this.mouse, camera);

    // Find intersection with ground plane
    const intersection = new THREE.Vector3();
    if (this.raycaster.ray.intersectPlane(this.groundPlane, intersection)) {
      this.targetPosition.x = intersection.x;
      this.targetPosition.z = intersection.z;
    }
  }

  sendMovement() {
    // Calculate movement direction from WASD keys
    let dx = 0;
    let dz = 0;

    if (this.keys.w) dz -= 1;
    if (this.keys.s) dz += 1;
    if (this.keys.a) dx -= 1;
    if (this.keys.d) dx += 1;

    // Normalize if moving diagonally
    if (dx !== 0 && dz !== 0) {
      const len = Math.sqrt(dx * dx + dz * dz);
      dx /= len;
      dz /= len;
    }

    // Only send if direction changed
    if (dx !== this.lastDirection.x || dz !== this.lastDirection.z) {
      this.lastDirection.x = dx;
      this.lastDirection.z = dz;

      if (this.onMove) {
        this.onMove({ x: dx, y: 0, z: dz });
      }
    }
  }

  getTargetPosition() {
    return this.targetPosition;
  }
}
