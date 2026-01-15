import * as THREE from 'three';
import { HealthBar } from './HealthBar.js';

export class RocketLauncher {
  constructor(scene, unit, color) {
    this.scene = scene;
    this.unit = unit;
    this.color = color;
    this.mesh = null;
    this.lastPosition = null;
    this.targetRotation = 0;
    this.healthBar = null;
    this.maxHealth = 20; // RocketLauncherHealth from server
    this.create();
  }

  create() {
    // Rocket Launcher soldier - bulkier than sniper with heavy weapon
    const bodyGroup = new THREE.Group();
    const scale = 0.85;

    // Main body (torso) - heavier armor
    const torsoGeometry = new THREE.BoxGeometry(1.4 * scale, 2 * scale, 1 * scale);
    const bodyMaterial = new THREE.MeshStandardMaterial({
      color: 0x5c5c5c, // Gray armor
      roughness: 0.6,
      metalness: 0.4
    });
    const torso = new THREE.Mesh(torsoGeometry, bodyMaterial);
    torso.position.y = 1.4 * scale;
    torso.castShadow = true;
    torso.receiveShadow = true;
    bodyGroup.add(torso);

    // Armored helmet
    const helmetGeometry = new THREE.SphereGeometry(0.45 * scale, 16, 16);
    const helmetMaterial = new THREE.MeshStandardMaterial({
      color: 0x4a4a4a,
      roughness: 0.5,
      metalness: 0.6
    });
    const helmet = new THREE.Mesh(helmetGeometry, helmetMaterial);
    helmet.position.y = 2.6 * scale;
    helmet.scale.y = 0.9; // Slightly flattened
    helmet.castShadow = true;
    bodyGroup.add(helmet);

    // Visor
    const visorGeometry = new THREE.BoxGeometry(0.5 * scale, 0.15 * scale, 0.2 * scale);
    const visorMaterial = new THREE.MeshStandardMaterial({
      color: 0x222222,
      roughness: 0.2,
      metalness: 0.9
    });
    const visor = new THREE.Mesh(visorGeometry, visorMaterial);
    visor.position.set(0, 2.55 * scale, 0.35 * scale);
    bodyGroup.add(visor);

    // Legs (armored)
    const legGeometry = new THREE.BoxGeometry(0.45 * scale, 1.1 * scale, 0.45 * scale);
    const leftLeg = new THREE.Mesh(legGeometry, bodyMaterial);
    leftLeg.position.set(-0.35 * scale, 0.4 * scale, 0);
    leftLeg.castShadow = true;
    bodyGroup.add(leftLeg);

    const rightLeg = new THREE.Mesh(legGeometry, bodyMaterial);
    rightLeg.position.set(0.35 * scale, 0.4 * scale, 0);
    rightLeg.castShadow = true;
    bodyGroup.add(rightLeg);

    // Arms (armored)
    const armGeometry = new THREE.BoxGeometry(0.4 * scale, 1.3 * scale, 0.4 * scale);
    const leftArm = new THREE.Mesh(armGeometry, bodyMaterial);
    leftArm.position.set(-0.95 * scale, 1.3 * scale, 0.3 * scale);
    leftArm.rotation.x = -0.4; // Holding launcher
    leftArm.castShadow = true;
    bodyGroup.add(leftArm);

    const rightArm = new THREE.Mesh(armGeometry, bodyMaterial);
    rightArm.position.set(0.95 * scale, 1.2 * scale, 0);
    rightArm.castShadow = true;
    bodyGroup.add(rightArm);

    // Rocket launcher tube
    const tubeGeometry = new THREE.CylinderGeometry(0.2 * scale, 0.25 * scale, 2 * scale, 12);
    const tubeMaterial = new THREE.MeshStandardMaterial({
      color: 0x3d6b3d, // Military green
      roughness: 0.5,
      metalness: 0.5
    });
    const tube = new THREE.Mesh(tubeGeometry, tubeMaterial);
    tube.rotation.x = Math.PI / 2;
    tube.position.set(0.7 * scale, 1.4 * scale, 0.8 * scale);
    tube.castShadow = true;
    bodyGroup.add(tube);

    // Launcher front opening
    const frontGeometry = new THREE.CylinderGeometry(0.2 * scale, 0.2 * scale, 0.1 * scale, 12);
    const frontMaterial = new THREE.MeshStandardMaterial({
      color: 0x1a1a1a,
      roughness: 0.2,
      metalness: 0.8
    });
    const front = new THREE.Mesh(frontGeometry, frontMaterial);
    front.rotation.x = Math.PI / 2;
    front.position.set(0.7 * scale, 1.4 * scale, 1.85 * scale);
    bodyGroup.add(front);

    // Launcher grip
    const gripGeometry = new THREE.BoxGeometry(0.15 * scale, 0.4 * scale, 0.15 * scale);
    const grip = new THREE.Mesh(gripGeometry, tubeMaterial);
    grip.position.set(0.7 * scale, 1.1 * scale, 0.5 * scale);
    bodyGroup.add(grip);

    // Backpack (ammo storage)
    const backpackGeometry = new THREE.BoxGeometry(0.8 * scale, 1 * scale, 0.5 * scale);
    const backpackMaterial = new THREE.MeshStandardMaterial({
      color: 0x4a5a4a,
      roughness: 0.7,
      metalness: 0.2
    });
    const backpack = new THREE.Mesh(backpackGeometry, backpackMaterial);
    backpack.position.set(0, 1.5 * scale, -0.6 * scale);
    backpack.castShadow = true;
    bodyGroup.add(backpack);

    // Team color indicator (shoulder pad)
    const padGeometry = new THREE.BoxGeometry(0.5 * scale, 0.3 * scale, 0.5 * scale);
    const padMaterial = new THREE.MeshStandardMaterial({
      color: this.color,
      roughness: 0.5,
      metalness: 0.3,
      emissive: this.color,
      emissiveIntensity: 0.2
    });
    const shoulderPad = new THREE.Mesh(padGeometry, padMaterial);
    shoulderPad.position.set(-0.95 * scale, 2 * scale, 0);
    bodyGroup.add(shoulderPad);

    // Create health bar
    this.healthBar = new HealthBar(this.maxHealth, 2.2 * scale, 0.25);
    this.healthBar.getGroup().position.y = 3.3 * scale;
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
