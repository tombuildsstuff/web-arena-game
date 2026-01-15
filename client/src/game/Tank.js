import * as THREE from 'three';
import { HealthBar } from './HealthBar.js';

export class Tank {
  constructor(scene, unit, color, isSuper = false) {
    this.scene = scene;
    this.unit = unit;
    this.color = color;
    this.isSuper = isSuper;
    this.mesh = null;
    this.turretGroup = null; // Separate group for turret rotation
    this.lastPosition = null;
    this.targetRotation = 0;
    this.turretTargetRotation = 0;
    this.currentTurretRotation = 0;
    this.healthBar = null;
    this.maxHealth = isSuper ? 90 : 30; // SuperTankHealth vs TankHealth
    this.create();
  }

  create() {
    // Scale factor for super tanks (1.3x larger)
    const scale = this.isSuper ? 1.3 : 1.0;

    // Tank body
    const bodyGeometry = new THREE.BoxGeometry(3 * scale, 2 * scale, 4 * scale);
    const bodyMaterial = new THREE.MeshStandardMaterial({
      color: this.color,
      roughness: 0.6,
      metalness: 0.4
    });

    this.mesh = new THREE.Mesh(bodyGeometry, bodyMaterial);
    this.mesh.castShadow = true;
    this.mesh.receiveShadow = true;

    // Add gold stripe for super tanks
    if (this.isSuper) {
      const stripeGeometry = new THREE.BoxGeometry(3.1 * scale, 0.3 * scale, 4.1 * scale);
      const stripeMaterial = new THREE.MeshStandardMaterial({
        color: 0xffd700, // Gold
        roughness: 0.3,
        metalness: 0.8,
        emissive: 0xffd700,
        emissiveIntensity: 0.2
      });
      const stripe = new THREE.Mesh(stripeGeometry, stripeMaterial);
      stripe.position.y = 0.5 * scale;
      this.mesh.add(stripe);
    }

    // Create turret group for independent rotation
    this.turretGroup = new THREE.Group();
    this.turretGroup.position.y = 1.75 * scale;
    this.mesh.add(this.turretGroup);

    // Turret base
    const turretGeometry = new THREE.BoxGeometry(2 * scale, 1.5 * scale, 2 * scale);
    const turret = new THREE.Mesh(turretGeometry, bodyMaterial);
    this.turretGroup.add(turret);

    // Cannon - attached to turret group so it rotates with it
    const cannonGeometry = new THREE.CylinderGeometry(0.3 * scale, 0.3 * scale, 3 * scale, 8);
    const cannon = new THREE.Mesh(cannonGeometry, bodyMaterial);
    cannon.rotation.z = Math.PI / 2;
    cannon.position.set(1.5 * scale, 0, 0); // Relative to turret group
    this.turretGroup.add(cannon);

    // Create health bar above tank
    this.healthBar = new HealthBar(this.maxHealth, 3 * scale, 0.35);
    this.healthBar.getGroup().position.y = 3.5 * scale;
    this.mesh.add(this.healthBar.getGroup());

    // Set initial position
    this.mesh.position.set(this.unit.position.x, this.unit.position.y, this.unit.position.z);
    this.lastPosition = { ...this.unit.position };

    this.scene.add(this.mesh);
  }

  updatePosition(position) {
    if (this.mesh) {
      // Calculate body rotation based on movement direction
      if (this.lastPosition) {
        const dx = position.x - this.lastPosition.x;
        const dz = position.z - this.lastPosition.z;
        const distanceMoved = Math.sqrt(dx * dx + dz * dz);

        // Only update body rotation if we've moved significantly
        if (distanceMoved > 0.01) {
          this.targetRotation = Math.atan2(dx, dz);
        }
      }

      // Smoothly interpolate body rotation
      const currentRotation = this.mesh.rotation.y;
      let rotationDiff = this.targetRotation - currentRotation;

      // Handle rotation wrapping
      while (rotationDiff > Math.PI) rotationDiff -= Math.PI * 2;
      while (rotationDiff < -Math.PI) rotationDiff += Math.PI * 2;

      this.mesh.rotation.y += rotationDiff * 0.1;

      // Smoothly interpolate turret rotation (relative to body)
      if (this.turretGroup) {
        // Calculate turret rotation relative to body
        const relativeTargetRotation = this.turretTargetRotation - this.mesh.rotation.y;
        let turretDiff = relativeTargetRotation - this.turretGroup.rotation.y;

        // Handle rotation wrapping
        while (turretDiff > Math.PI) turretDiff -= Math.PI * 2;
        while (turretDiff < -Math.PI) turretDiff += Math.PI * 2;

        this.turretGroup.rotation.y += turretDiff * 0.15;
      }

      // Store last position for next frame
      this.lastPosition = { x: position.x, y: position.y, z: position.z };
    }
  }

  // Set target position for turret to aim at
  setTarget(targetPosition) {
    if (this.mesh && targetPosition) {
      const dx = targetPosition.x - this.mesh.position.x;
      const dz = targetPosition.z - this.mesh.position.z;
      this.turretTargetRotation = Math.atan2(dx, dz);
    }
  }

  // Aim turret at a world position
  aimAt(worldX, worldZ) {
    if (this.mesh) {
      const dx = worldX - this.mesh.position.x;
      const dz = worldZ - this.mesh.position.z;
      this.turretTargetRotation = Math.atan2(dx, dz);
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
