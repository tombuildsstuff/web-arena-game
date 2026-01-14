import * as THREE from 'three';
import { HealthBar } from './HealthBar.js';

export class Airplane {
  constructor(scene, unit, color) {
    this.scene = scene;
    this.unit = unit;
    this.color = color;
    this.mesh = null;
    this.lastPosition = null;
    this.targetRotation = 0;
    this.currentBank = 0;
    this.healthBar = null;
    this.maxHealth = 30; // AirplaneHealth from server
    this.create();
  }

  create() {
    // Create parent group for positioning
    this.mesh = new THREE.Group();

    // Airplane body (fuselage)
    const bodyGeometry = new THREE.ConeGeometry(0.8, 4, 8);
    const bodyMaterial = new THREE.MeshStandardMaterial({
      color: this.color,
      roughness: 0.4,
      metalness: 0.6
    });

    this.bodyMesh = new THREE.Mesh(bodyGeometry, bodyMaterial);
    this.bodyMesh.rotation.x = Math.PI / 2; // Point forward
    this.bodyMesh.castShadow = true;
    this.bodyMesh.receiveShadow = true;
    this.mesh.add(this.bodyMesh);

    // Wings
    const wingGeometry = new THREE.BoxGeometry(8, 0.2, 2);
    const wings = new THREE.Mesh(wingGeometry, bodyMaterial);
    wings.position.z = -0.5;
    this.bodyMesh.add(wings);

    // Tail wing
    const tailGeometry = new THREE.BoxGeometry(3, 0.2, 1);
    const tail = new THREE.Mesh(tailGeometry, bodyMaterial);
    tail.position.z = -2;
    this.bodyMesh.add(tail);

    // Create health bar above airplane (on parent group so it stays upright)
    this.healthBar = new HealthBar(this.maxHealth, 3, 0.35);
    this.healthBar.getGroup().position.y = 3;
    this.mesh.add(this.healthBar.getGroup());

    // Set initial position
    this.mesh.position.set(this.unit.position.x, this.unit.position.y, this.unit.position.z);
    this.lastPosition = { ...this.unit.position };

    this.scene.add(this.mesh);
  }

  updatePosition(position) {
    if (this.mesh && this.bodyMesh) {
      // Calculate rotation based on movement direction
      if (this.lastPosition) {
        const dx = position.x - this.lastPosition.x;
        const dz = position.z - this.lastPosition.z;
        const distanceMoved = Math.sqrt(dx * dx + dz * dz);

        // Only update rotation if we've moved significantly
        if (distanceMoved > 0.01) {
          const newTargetRotation = Math.atan2(dx, dz);

          // Calculate turn rate for banking effect
          let turnRate = newTargetRotation - this.targetRotation;
          while (turnRate > Math.PI) turnRate -= Math.PI * 2;
          while (turnRate < -Math.PI) turnRate += Math.PI * 2;

          // Target bank angle based on turn rate
          const targetBank = Math.max(-0.5, Math.min(0.5, turnRate * 2));
          this.currentBank += (targetBank - this.currentBank) * 0.1;

          this.targetRotation = newTargetRotation;
        } else {
          // Reduce bank when not turning
          this.currentBank *= 0.95;
        }
      }

      // Smoothly interpolate rotation (apply to body mesh, not parent group)
      let rotationDiff = this.targetRotation - this.bodyMesh.rotation.y;
      while (rotationDiff > Math.PI) rotationDiff -= Math.PI * 2;
      while (rotationDiff < -Math.PI) rotationDiff += Math.PI * 2;
      this.bodyMesh.rotation.y += rotationDiff * 0.15;

      // Apply banking to body mesh
      this.bodyMesh.rotation.z = this.currentBank;

      // Store last position
      this.lastPosition = { x: position.x, y: position.y, z: position.z };
    }
  }

  // Update health bar
  updateHealth(health, camera) {
    if (this.healthBar) {
      this.healthBar.update(health, camera);
    }
  }

  remove() {
    if (this.healthBar) {
      this.healthBar.dispose();
    }
    if (this.mesh) {
      this.scene.remove(this.mesh);
    }
  }
}
