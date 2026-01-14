import * as THREE from 'three';

export class Camera {
  constructor(renderer) {
    this.renderer = renderer;
    this.target = null; // The object to follow (player mesh)
    this.targetPosition = new THREE.Vector3();
    this.currentPosition = new THREE.Vector3();
    this.currentLookAt = new THREE.Vector3();

    // Camera orbit parameters (spherical coordinates)
    this.distance = 45; // Distance from player
    this.minDistance = 15;
    this.maxDistance = 100;
    this.azimuth = 0; // Horizontal rotation (radians)
    this.polar = Math.PI / 4; // Vertical angle (radians) - 45 degrees from horizontal
    this.minPolar = 0.2; // Minimum vertical angle (prevent going below ground)
    this.maxPolar = Math.PI / 2.2; // Maximum vertical angle (prevent going too high)

    this.lookAtOffset = new THREE.Vector3(0, 2, 0); // Look slightly above player

    // Smoothing factor for camera movement
    this.smoothing = 0.08;

    // Mouse control state
    this.isDragging = false;
    this.lastMouseX = 0;
    this.lastMouseY = 0;
    this.rotationSpeed = 0.005;
    this.zoomSpeed = 0.1;

    this.create();
    this.setupControls();
  }

  create() {
    // Create perspective camera
    this.camera = new THREE.PerspectiveCamera(
      60,
      window.innerWidth / window.innerHeight,
      0.1,
      1000
    );

    // Initial position (will be updated when following player)
    this.camera.position.set(0, 80, 80);
    this.camera.lookAt(0, 0, 0);

    this.currentPosition.copy(this.camera.position);
    this.currentLookAt.set(0, 0, 0);
  }

  setupControls() {
    const canvas = this.renderer.domElement;

    // Right-click or middle-click drag to rotate
    canvas.addEventListener('mousedown', (e) => {
      if (e.button === 2 || e.button === 1) { // Right or middle click
        this.isDragging = true;
        this.lastMouseX = e.clientX;
        this.lastMouseY = e.clientY;
        e.preventDefault();
      }
    });

    window.addEventListener('mouseup', (e) => {
      if (e.button === 2 || e.button === 1) {
        this.isDragging = false;
      }
    });

    window.addEventListener('mousemove', (e) => {
      if (this.isDragging) {
        const deltaX = e.clientX - this.lastMouseX;
        const deltaY = e.clientY - this.lastMouseY;

        // Update azimuth (horizontal rotation)
        this.azimuth -= deltaX * this.rotationSpeed;

        // Update polar (vertical rotation)
        this.polar += deltaY * this.rotationSpeed;
        this.polar = Math.max(this.minPolar, Math.min(this.maxPolar, this.polar));

        this.lastMouseX = e.clientX;
        this.lastMouseY = e.clientY;
      }
    });

    // Mouse wheel to zoom
    canvas.addEventListener('wheel', (e) => {
      e.preventDefault();
      const zoomDelta = e.deltaY * this.zoomSpeed;
      this.distance += zoomDelta;
      this.distance = Math.max(this.minDistance, Math.min(this.maxDistance, this.distance));
    }, { passive: false });

    // Prevent context menu on right-click
    canvas.addEventListener('contextmenu', (e) => {
      e.preventDefault();
    });
  }

  // Calculate camera offset from spherical coordinates
  calculateOffset() {
    const x = this.distance * Math.sin(this.polar) * Math.sin(this.azimuth);
    const y = this.distance * Math.cos(this.polar);
    const z = this.distance * Math.sin(this.polar) * Math.cos(this.azimuth);
    return new THREE.Vector3(x, y, z);
  }

  // Set the target to follow (player mesh)
  setTarget(targetMesh) {
    this.target = targetMesh;
    if (targetMesh) {
      // Initialize camera position behind player
      const offset = this.calculateOffset();
      this.targetPosition.copy(targetMesh.position).add(offset);
      this.currentPosition.copy(this.targetPosition);
      this.camera.position.copy(this.currentPosition);
    }
  }

  update() {
    if (this.target && this.target.position) {
      // Calculate offset based on current spherical coordinates
      const offset = this.calculateOffset();

      // Calculate desired camera position
      this.targetPosition.copy(this.target.position).add(offset);

      // Smoothly interpolate camera position
      this.currentPosition.lerp(this.targetPosition, this.smoothing);
      this.camera.position.copy(this.currentPosition);

      // Calculate look-at target (player position + offset)
      const lookAtTarget = this.target.position.clone().add(this.lookAtOffset);

      // Smoothly interpolate look-at
      this.currentLookAt.lerp(lookAtTarget, this.smoothing);
      this.camera.lookAt(this.currentLookAt);
    }
  }

  handleResize() {
    this.camera.aspect = window.innerWidth / window.innerHeight;
    this.camera.updateProjectionMatrix();
  }

  getCamera() {
    return this.camera;
  }

  // Get the current look-at position (for raycasting ground plane)
  getLookAtPosition() {
    return this.currentLookAt.clone();
  }
}
