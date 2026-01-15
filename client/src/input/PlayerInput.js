import * as THREE from 'three';

export class PlayerInput {
  constructor(camera, renderer, onMove, onShoot, onBuyFromZone, onBulkBuyFromZone, onClaimTurret, onClaimBuyZone, onClaimBarracks, gameLoop) {
    this.camera = camera;
    this.renderer = renderer;
    this.onMove = onMove;
    this.onShoot = onShoot;
    this.onBuyFromZone = onBuyFromZone;
    this.onBulkBuyFromZone = onBulkBuyFromZone;
    this.onClaimTurret = onClaimTurret;
    this.onClaimBuyZone = onClaimBuyZone;
    this.onClaimBarracks = onClaimBarracks;
    this.gameLoop = gameLoop;

    // Movement keys state (using key codes for reliability across keyboard layouts)
    this.keys = {
      up: false,    // W or ArrowUp
      down: false,  // S or ArrowDown
      left: false,  // A or ArrowLeft
      right: false  // D or ArrowRight
    };

    // Map key codes to movement directions
    this.keyMap = {
      'KeyW': 'up',
      'KeyS': 'down',
      'KeyA': 'left',
      'KeyD': 'right',
      'ArrowUp': 'up',
      'ArrowDown': 'down',
      'ArrowLeft': 'left',
      'ArrowRight': 'right'
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
    this.handleBlur = this.handleBlur.bind(this);

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
    window.addEventListener('blur', this.handleBlur);
  }

  disable() {
    if (!this.enabled) return;
    this.enabled = false;

    window.removeEventListener('keydown', this.handleKeyDown);
    window.removeEventListener('keyup', this.handleKeyUp);
    window.removeEventListener('mousemove', this.handleMouseMove);
    window.removeEventListener('mousedown', this.handleMouseDown);
    window.removeEventListener('blur', this.handleBlur);

    // Reset keys
    this.resetKeys();
    this.sendMovement();
  }

  // Reset all keys (used when window loses focus)
  resetKeys() {
    this.keys = { up: false, down: false, left: false, right: false };
  }

  // Handle window losing focus
  handleBlur() {
    this.resetKeys();
    this.sendMovement();
  }

  handleKeyDown(event) {
    // Check for movement keys using key codes
    const direction = this.keyMap[event.code];
    if (direction && !this.keys[direction]) {
      event.preventDefault(); // Prevent page scrolling with arrow keys
      this.keys[direction] = true;
      this.sendMovement();
    }

    // X key for shooting
    if (event.code === 'KeyX') {
      this.updateTargetPosition();
      if (this.onShoot) {
        this.onShoot(this.targetPosition.x, this.targetPosition.z);
      }
    }

    // C key for buying from zones or claiming turrets
    if (event.code === 'KeyC') {
      this.handleInteraction();
    }

    // V key for bulk buying (10 units at 10% discount)
    if (event.code === 'KeyV') {
      this.handleBulkBuy();
    }
  }

  handleInteraction() {
    if (!this.gameLoop) return;

    // Priority: owned buy zones, then claimable buy zones, then turrets, then barracks
    const nearbyZone = this.gameLoop.getNearbyBuyZone();
    if (nearbyZone && this.onBuyFromZone) {
      this.onBuyFromZone(nearbyZone.id);
      return;
    }

    const nearbyClaimableZone = this.gameLoop.getNearbyClaimableBuyZone();
    if (nearbyClaimableZone && this.onClaimBuyZone) {
      this.onClaimBuyZone(nearbyClaimableZone.id);
      return;
    }

    const nearbyTurret = this.gameLoop.getNearbyTurret();
    if (nearbyTurret && this.onClaimTurret) {
      this.onClaimTurret(nearbyTurret.id);
      return;
    }

    const nearbyBarracks = this.gameLoop.getNearbyBarracks();
    if (nearbyBarracks && this.onClaimBarracks) {
      this.onClaimBarracks(nearbyBarracks.id);
    }
  }

  handleBulkBuy() {
    if (!this.gameLoop) return;

    // Only works for regular tank/helicopter zones (not super units)
    const nearbyZone = this.gameLoop.getNearbyBuyZone();
    if (nearbyZone && this.onBulkBuyFromZone) {
      // Only allow bulk buy for tank and airplane (helicopter)
      if (nearbyZone.unitType === 'tank' || nearbyZone.unitType === 'airplane') {
        this.onBulkBuyFromZone(nearbyZone.id);
      }
    }
  }

  handleKeyUp(event) {
    // Check for movement keys using key codes
    const direction = this.keyMap[event.code];
    if (direction && this.keys[direction]) {
      event.preventDefault();
      this.keys[direction] = false;
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
    // Calculate movement direction from keys (in camera-relative space)
    // forward/back is relative to camera view, left/right is perpendicular
    let forward = 0;
    let right = 0;

    if (this.keys.up) forward += 1;    // W = move forward (away from camera)
    if (this.keys.down) forward -= 1;  // S = move backward (toward camera)
    if (this.keys.left) right -= 1;    // A = move left
    if (this.keys.right) right += 1;   // D = move right

    // Get camera-relative directions
    const camForward = this.camera.getForwardDirection();
    const camRight = this.camera.getRightDirection();

    // Transform input to world space
    let dx = forward * camForward.x + right * camRight.x;
    let dz = forward * camForward.z + right * camRight.z;

    // Normalize if moving diagonally
    const len = Math.sqrt(dx * dx + dz * dz);
    if (len > 0) {
      dx /= len;
      dz /= len;
    }

    // Only send if direction changed (with small tolerance for floating point)
    const threshold = 0.001;
    if (Math.abs(dx - this.lastDirection.x) > threshold ||
        Math.abs(dz - this.lastDirection.z) > threshold) {
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

  // Check if any movement keys are currently pressed
  isMoving() {
    return this.keys.up || this.keys.down || this.keys.left || this.keys.right;
  }

  // Force recalculation of movement direction (called when camera rotates)
  updateMovementDirection() {
    if (this.enabled && this.isMoving()) {
      // Reset last direction to force a resend
      this.lastDirection.x = Infinity;
      this.lastDirection.z = Infinity;
      this.sendMovement();
    }
  }
}
