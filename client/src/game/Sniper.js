import * as THREE from 'three';
import { HealthBar } from './HealthBar.js';

export class Sniper {
  constructor(scene, unit, color) {
    this.scene = scene;
    this.unit = unit;
    this.color = color;
    this.mesh = null;
    this.lastPosition = null;
    this.targetRotation = 0;
    this.healthBar = null;
    this.maxHealth = 15; // SniperHealth from server
    this.create();
  }

  create() {
    // Sniper - smaller humanoid with long rifle
    const bodyGroup = new THREE.Group();
    const scale = 0.8; // Slightly smaller than player

    // Main body (torso) - wearing a cloak/ghillie suit
    const torsoGeometry = new THREE.BoxGeometry(1.2 * scale, 1.8 * scale, 0.8 * scale);
    const bodyMaterial = new THREE.MeshStandardMaterial({
      color: 0x4a5d4a, // Camouflage green
      roughness: 0.8,
      metalness: 0.1
    });
    const torso = new THREE.Mesh(torsoGeometry, bodyMaterial);
    torso.position.y = 1.3 * scale;
    torso.castShadow = true;
    torso.receiveShadow = true;
    bodyGroup.add(torso);

    // Hood/Head
    const headGeometry = new THREE.SphereGeometry(0.4 * scale, 16, 16);
    const hoodMaterial = new THREE.MeshStandardMaterial({
      color: 0x3d4d3d, // Darker green hood
      roughness: 0.9,
      metalness: 0.0
    });
    const head = new THREE.Mesh(headGeometry, hoodMaterial);
    head.position.y = 2.5 * scale;
    head.castShadow = true;
    bodyGroup.add(head);

    // Legs
    const legGeometry = new THREE.BoxGeometry(0.4 * scale, 1.0 * scale, 0.4 * scale);
    const leftLeg = new THREE.Mesh(legGeometry, bodyMaterial);
    leftLeg.position.set(-0.3 * scale, 0.4 * scale, 0);
    leftLeg.castShadow = true;
    bodyGroup.add(leftLeg);

    const rightLeg = new THREE.Mesh(legGeometry, bodyMaterial);
    rightLeg.position.set(0.3 * scale, 0.4 * scale, 0);
    rightLeg.castShadow = true;
    bodyGroup.add(rightLeg);

    // Arms
    const armGeometry = new THREE.BoxGeometry(0.3 * scale, 1.2 * scale, 0.3 * scale);
    const leftArm = new THREE.Mesh(armGeometry, bodyMaterial);
    leftArm.position.set(-0.85 * scale, 1.3 * scale, 0.2 * scale);
    leftArm.rotation.x = -0.3; // Holding rifle
    leftArm.castShadow = true;
    bodyGroup.add(leftArm);

    const rightArm = new THREE.Mesh(armGeometry, bodyMaterial);
    rightArm.position.set(0.85 * scale, 1.3 * scale, 0.2 * scale);
    rightArm.rotation.x = -0.3;
    rightArm.castShadow = true;
    bodyGroup.add(rightArm);

    // Long sniper rifle
    const rifleGeometry = new THREE.BoxGeometry(0.15 * scale, 0.15 * scale, 2.5 * scale);
    const rifleMaterial = new THREE.MeshStandardMaterial({
      color: 0x2a2a2a,
      roughness: 0.3,
      metalness: 0.7
    });
    const rifle = new THREE.Mesh(rifleGeometry, rifleMaterial);
    rifle.position.set(0.5 * scale, 1.1 * scale, 1.0 * scale);
    rifle.castShadow = true;
    bodyGroup.add(rifle);

    // Rifle scope
    const scopeGeometry = new THREE.CylinderGeometry(0.08 * scale, 0.08 * scale, 0.4 * scale, 8);
    const scopeMaterial = new THREE.MeshStandardMaterial({
      color: 0x1a1a1a,
      roughness: 0.2,
      metalness: 0.9
    });
    const scope = new THREE.Mesh(scopeGeometry, scopeMaterial);
    scope.rotation.x = Math.PI / 2;
    scope.position.set(0.5 * scale, 1.3 * scale, 0.5 * scale);
    bodyGroup.add(scope);

    // Team color indicator (armband)
    const armbandGeometry = new THREE.BoxGeometry(0.35 * scale, 0.2 * scale, 0.35 * scale);
    const armbandMaterial = new THREE.MeshStandardMaterial({
      color: this.color,
      roughness: 0.5,
      metalness: 0.3,
      emissive: this.color,
      emissiveIntensity: 0.2
    });
    const armband = new THREE.Mesh(armbandGeometry, armbandMaterial);
    armband.position.set(-0.85 * scale, 1.7 * scale, 0);
    bodyGroup.add(armband);

    // Create health bar
    this.healthBar = new HealthBar(this.maxHealth, 2 * scale, 0.25);
    this.healthBar.getGroup().position.y = 3.2 * scale;
    bodyGroup.add(this.healthBar.getGroup());

    this.mesh = bodyGroup;

    // Set initial position
    this.mesh.position.set(this.unit.position.x, this.unit.position.y, this.unit.position.z);
    this.lastPosition = { ...this.unit.position };

    this.scene.add(this.mesh);
  }

  updatePosition(position) {
    if (this.mesh) {
      // Calculate rotation based on movement direction
      if (this.lastPosition) {
        const dx = position.x - this.lastPosition.x;
        const dz = position.z - this.lastPosition.z;
        const distanceMoved = Math.sqrt(dx * dx + dz * dz);

        // Only update rotation if we've moved significantly
        if (distanceMoved > 0.01) {
          this.targetRotation = Math.atan2(dx, dz);
        }
      }

      // Smoothly interpolate rotation
      let rotationDiff = this.targetRotation - this.mesh.rotation.y;
      while (rotationDiff > Math.PI) rotationDiff -= Math.PI * 2;
      while (rotationDiff < -Math.PI) rotationDiff += Math.PI * 2;
      this.mesh.rotation.y += rotationDiff * 0.15;

      // Store last position
      this.lastPosition = { x: position.x, y: position.y, z: position.z };
    }
  }

  // Set target position for aiming
  setTarget(targetPosition) {
    if (this.mesh && targetPosition) {
      const dx = targetPosition.x - this.mesh.position.x;
      const dz = targetPosition.z - this.mesh.position.z;
      this.targetRotation = Math.atan2(dx, dz);
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
