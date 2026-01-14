import * as THREE from 'three';

export class Player {
  constructor(scene, unit, color, isLocalPlayer = false) {
    this.scene = scene;
    this.unit = unit;
    this.color = color;
    this.isLocalPlayer = isLocalPlayer;
    this.mesh = null;
    this.lastPosition = null;
    this.targetRotation = 0;
    this.respawnOverlay = null;
    this.create();
  }

  create() {
    // Player body - humanoid shape
    const bodyGroup = new THREE.Group();

    // Main body (torso)
    const torsoGeometry = new THREE.BoxGeometry(1.5, 2, 1);
    const bodyMaterial = new THREE.MeshStandardMaterial({
      color: this.color,
      roughness: 0.5,
      metalness: 0.3
    });
    const torso = new THREE.Mesh(torsoGeometry, bodyMaterial);
    torso.position.y = 1.5;
    torso.castShadow = true;
    torso.receiveShadow = true;
    bodyGroup.add(torso);

    // Head
    const headGeometry = new THREE.SphereGeometry(0.5, 16, 16);
    const head = new THREE.Mesh(headGeometry, bodyMaterial);
    head.position.y = 3;
    head.castShadow = true;
    bodyGroup.add(head);

    // Legs
    const legGeometry = new THREE.BoxGeometry(0.5, 1.2, 0.5);
    const leftLeg = new THREE.Mesh(legGeometry, bodyMaterial);
    leftLeg.position.set(-0.4, 0.4, 0);
    leftLeg.castShadow = true;
    bodyGroup.add(leftLeg);

    const rightLeg = new THREE.Mesh(legGeometry, bodyMaterial);
    rightLeg.position.set(0.4, 0.4, 0);
    rightLeg.castShadow = true;
    bodyGroup.add(rightLeg);

    // Arms
    const armGeometry = new THREE.BoxGeometry(0.4, 1.5, 0.4);
    const leftArm = new THREE.Mesh(armGeometry, bodyMaterial);
    leftArm.position.set(-1.1, 1.5, 0);
    leftArm.castShadow = true;
    bodyGroup.add(leftArm);

    const rightArm = new THREE.Mesh(armGeometry, bodyMaterial);
    rightArm.position.set(1.1, 1.5, 0);
    rightArm.castShadow = true;
    bodyGroup.add(rightArm);

    // Gun (held by right arm)
    const gunGeometry = new THREE.BoxGeometry(0.2, 0.2, 1.5);
    const gunMaterial = new THREE.MeshStandardMaterial({
      color: 0x333333,
      roughness: 0.3,
      metalness: 0.8
    });
    const gun = new THREE.Mesh(gunGeometry, gunMaterial);
    gun.position.set(1.1, 1.2, 0.7);
    gun.castShadow = true;
    bodyGroup.add(gun);

    // Indicator for local player
    if (this.isLocalPlayer) {
      const indicatorGeometry = new THREE.RingGeometry(1.5, 2, 32);
      const indicatorMaterial = new THREE.MeshBasicMaterial({
        color: 0x00ff00,
        side: THREE.DoubleSide,
        transparent: true,
        opacity: 0.5
      });
      const indicator = new THREE.Mesh(indicatorGeometry, indicatorMaterial);
      indicator.rotation.x = -Math.PI / 2;
      indicator.position.y = 0.1;
      bodyGroup.add(indicator);
    }

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

  // Set respawning state (dims the player)
  setRespawning(isRespawning) {
    if (!this.mesh) return;

    this.mesh.traverse((child) => {
      if (child.isMesh && child.material) {
        if (isRespawning) {
          child.material.transparent = true;
          child.material.opacity = 0.3;
        } else {
          child.material.transparent = false;
          child.material.opacity = 1.0;
        }
      }
    });
  }

  // Aim at a target position
  aimAt(targetX, targetZ) {
    if (this.mesh) {
      const dx = targetX - this.mesh.position.x;
      const dz = targetZ - this.mesh.position.z;
      this.targetRotation = Math.atan2(dx, dz);
    }
  }

  remove() {
    if (this.mesh) {
      this.scene.remove(this.mesh);
    }
  }
}
