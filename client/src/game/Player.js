import * as THREE from 'three';
import { HealthBar } from './HealthBar.js';

export class Player {
  constructor(scene, unit, color, isLocalPlayer = false, displayName = null) {
    this.scene = scene;
    this.unit = unit;
    this.color = color;
    this.isLocalPlayer = isLocalPlayer;
    this.displayName = displayName || `Player ${unit.ownerId + 1}`;
    this.mesh = null;
    this.lastPosition = null;
    this.targetRotation = 0;
    this.respawnOverlay = null;
    this.nametagSprite = null;
    this.healthBar = null;
    this.maxHealth = 10; // PlayerUnitHealth from server
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

    // Create floating nametag
    this.nametagSprite = this.createNametag();
    this.nametagSprite.position.y = 4.5; // Above the head
    bodyGroup.add(this.nametagSprite);

    // Create health bar (between head and nametag)
    this.healthBar = new HealthBar(this.maxHealth, 2.5, 0.3);
    this.healthBar.getGroup().position.y = 3.8;
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

  // Update health bar
  updateHealth(health, camera) {
    if (this.healthBar) {
      this.healthBar.update(health, camera);
    }
  }

  // Create a floating nametag sprite
  createNametag() {
    const name = this.displayName;

    // Create canvas for the nametag
    const canvas = document.createElement('canvas');
    const context = canvas.getContext('2d');

    // Set canvas size
    canvas.width = 256;
    canvas.height = 64;

    // Get background color from player color
    const bgColor = typeof this.color === 'string' ? this.color : `#${this.color.toString(16).padStart(6, '0')}`;

    // Draw rounded rectangle background
    const padding = 10;
    const radius = 12;
    context.fillStyle = bgColor;
    context.beginPath();
    context.roundRect(padding, padding, canvas.width - padding * 2, canvas.height - padding * 2, radius);
    context.fill();

    // Add slight border
    context.strokeStyle = 'rgba(255, 255, 255, 0.5)';
    context.lineWidth = 2;
    context.stroke();

    // Draw text
    context.fillStyle = '#ffffff';
    context.font = 'bold 28px Arial, sans-serif';
    context.textAlign = 'center';
    context.textBaseline = 'middle';
    context.fillText(name, canvas.width / 2, canvas.height / 2);

    // Create texture from canvas
    const texture = new THREE.CanvasTexture(canvas);
    texture.needsUpdate = true;

    // Create sprite material
    const material = new THREE.SpriteMaterial({
      map: texture,
      transparent: true,
      depthTest: false, // Always render on top
      depthWrite: false
    });

    // Create sprite
    const sprite = new THREE.Sprite(material);
    sprite.scale.set(4, 1, 1); // Scale to match aspect ratio

    return sprite;
  }

  remove() {
    if (this.nametagSprite) {
      if (this.nametagSprite.material.map) {
        this.nametagSprite.material.map.dispose();
      }
      this.nametagSprite.material.dispose();
    }
    if (this.healthBar) {
      this.healthBar.dispose();
    }
    if (this.mesh) {
      this.scene.remove(this.mesh);
    }
  }
}
