import * as THREE from 'three';

export class BuyZonePopup {
  constructor() {
    this.popup = document.getElementById('buyzone-popup');
    this.popupText = document.getElementById('buyzone-popup-text');
    this.hideTimeout = null;
    this.camera = null;
    this.renderer = null;
    this.targetPosition = null;
  }

  setCamera(camera, renderer) {
    this.camera = camera;
    this.renderer = renderer;
  }

  show(message, worldPosition, duration = 2000) {
    if (!this.popup || !this.popupText) return;

    this.popupText.textContent = message;
    this.targetPosition = worldPosition ? new THREE.Vector3(
      worldPosition.x,
      worldPosition.y + 3, // Show above the buy zone
      worldPosition.z
    ) : null;

    // Update position immediately
    this.updatePosition();

    this.popup.classList.remove('hidden');
    this.popup.classList.add('visible');

    // Clear existing timeout
    if (this.hideTimeout) {
      clearTimeout(this.hideTimeout);
    }

    // Auto-hide after duration
    this.hideTimeout = setTimeout(() => {
      this.hide();
    }, duration);
  }

  hide() {
    if (!this.popup) return;
    this.popup.classList.remove('visible');
    this.popup.classList.add('hidden');
    this.targetPosition = null;
  }

  updatePosition() {
    if (!this.targetPosition || !this.camera || !this.renderer || !this.popup) return;

    const camera = this.camera.getCamera();
    const canvas = this.renderer.domElement;

    // Project 3D position to screen coordinates
    const projected = this.targetPosition.clone().project(camera);

    // Convert to screen coordinates
    const x = (projected.x * 0.5 + 0.5) * canvas.clientWidth;
    const y = (-projected.y * 0.5 + 0.5) * canvas.clientHeight;

    // Check if position is in front of camera
    if (projected.z > 1) {
      this.popup.classList.add('hidden');
      return;
    }

    this.popup.style.left = `${x}px`;
    this.popup.style.top = `${y}px`;
  }

  update() {
    // Update popup position each frame if visible
    if (this.targetPosition) {
      this.updatePosition();
    }
  }
}
