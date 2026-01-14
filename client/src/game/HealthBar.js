import * as THREE from 'three';

export class HealthBar {
  constructor(maxHealth, width = 2, height = 0.3) {
    this.maxHealth = maxHealth;
    this.currentHealth = maxHealth;
    this.width = width;
    this.height = height;
    this.group = new THREE.Group();
    this.backgroundMesh = null;
    this.foregroundMesh = null;
    this.visible = false;
    this.create();
  }

  create() {
    // Background (dark red/gray)
    const bgGeometry = new THREE.PlaneGeometry(this.width, this.height);
    const bgMaterial = new THREE.MeshBasicMaterial({
      color: 0x333333,
      side: THREE.DoubleSide,
      transparent: true,
      opacity: 0.8,
      depthTest: false
    });
    this.backgroundMesh = new THREE.Mesh(bgGeometry, bgMaterial);
    this.group.add(this.backgroundMesh);

    // Foreground (health - green to red gradient based on health)
    const fgGeometry = new THREE.PlaneGeometry(this.width, this.height);
    const fgMaterial = new THREE.MeshBasicMaterial({
      color: 0x44ff44,
      side: THREE.DoubleSide,
      transparent: true,
      opacity: 0.9,
      depthTest: false
    });
    this.foregroundMesh = new THREE.Mesh(fgGeometry, fgMaterial);
    this.foregroundMesh.position.z = 0.01; // Slightly in front of background
    this.group.add(this.foregroundMesh);

    // Border
    const borderGeometry = new THREE.EdgesGeometry(bgGeometry);
    const borderMaterial = new THREE.LineBasicMaterial({
      color: 0x000000,
      transparent: true,
      opacity: 0.5
    });
    const border = new THREE.LineSegments(borderGeometry, borderMaterial);
    border.position.z = 0.02;
    this.group.add(border);

    // Initially hidden
    this.group.visible = false;
  }

  update(currentHealth, camera) {
    this.currentHealth = currentHealth;
    const healthPercent = Math.max(0, currentHealth / this.maxHealth);

    // Show health bar only when damaged
    const isDamaged = currentHealth < this.maxHealth;
    this.group.visible = isDamaged;
    this.visible = isDamaged;

    if (!isDamaged) return;

    // Update foreground scale and position
    this.foregroundMesh.scale.x = healthPercent;
    // Offset to align left edge
    this.foregroundMesh.position.x = -(this.width * (1 - healthPercent)) / 2;

    // Update color based on health (green -> yellow -> red)
    const color = this.getHealthColor(healthPercent);
    this.foregroundMesh.material.color.setHex(color);

    // Make health bar face the camera (billboard effect)
    if (camera) {
      this.group.quaternion.copy(camera.quaternion);
    }
  }

  getHealthColor(percent) {
    if (percent > 0.6) {
      return 0x44ff44; // Green
    } else if (percent > 0.3) {
      return 0xffff44; // Yellow
    } else {
      return 0xff4444; // Red
    }
  }

  getGroup() {
    return this.group;
  }

  dispose() {
    if (this.backgroundMesh) {
      this.backgroundMesh.geometry.dispose();
      this.backgroundMesh.material.dispose();
    }
    if (this.foregroundMesh) {
      this.foregroundMesh.geometry.dispose();
      this.foregroundMesh.material.dispose();
    }
  }
}
