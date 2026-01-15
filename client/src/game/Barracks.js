import * as THREE from 'three';
import { PLAYER_COLORS } from '../utils/constants.js';

export class Barracks {
  constructor(scene, barracks) {
    this.scene = scene;
    this.barracks = barracks;
    this.meshGroup = new THREE.Group();
    this.buildingMesh = null;
    this.roofMesh = null;
    this.flagMesh = null;
    this.rangeMesh = null;
    this.healthBarGroup = null;
    this.healingIndicator = null;
    this.create();
  }

  create() {
    // Determine color based on owner
    const color = this.getColor();

    // Base platform
    const platformGeometry = new THREE.BoxGeometry(6, 0.3, 6);
    const platformMaterial = new THREE.MeshStandardMaterial({
      color: 0x555555,
      roughness: 0.8,
      metalness: 0.2
    });
    const platform = new THREE.Mesh(platformGeometry, platformMaterial);
    platform.position.y = 0.15;
    platform.castShadow = true;
    platform.receiveShadow = true;
    this.meshGroup.add(platform);

    // Main building structure (military tent/bunker style)
    const buildingGeometry = new THREE.BoxGeometry(4, 2.5, 4);
    const buildingMaterial = new THREE.MeshStandardMaterial({
      color: 0x6b6b4c, // Military olive
      roughness: 0.7,
      metalness: 0.1
    });
    this.buildingMesh = new THREE.Mesh(buildingGeometry, buildingMaterial);
    this.buildingMesh.position.y = 1.55;
    this.buildingMesh.castShadow = true;
    this.buildingMesh.receiveShadow = true;
    this.meshGroup.add(this.buildingMesh);

    // Roof (pyramid style)
    const roofGeometry = new THREE.ConeGeometry(3.5, 2, 4);
    const roofMaterial = new THREE.MeshStandardMaterial({
      color: 0x4a5a4a,
      roughness: 0.8,
      metalness: 0.1
    });
    this.roofMesh = new THREE.Mesh(roofGeometry, roofMaterial);
    this.roofMesh.position.y = 3.8;
    this.roofMesh.rotation.y = Math.PI / 4; // Rotate to align with building
    this.roofMesh.castShadow = true;
    this.meshGroup.add(this.roofMesh);

    // Entrance (doorway)
    const doorGeometry = new THREE.BoxGeometry(1.2, 2, 0.2);
    const doorMaterial = new THREE.MeshStandardMaterial({
      color: 0x333333,
      roughness: 0.5,
      metalness: 0.3
    });
    const door = new THREE.Mesh(doorGeometry, doorMaterial);
    door.position.set(0, 1.3, 2.1);
    this.meshGroup.add(door);

    // Flag pole
    const poleGeometry = new THREE.CylinderGeometry(0.08, 0.08, 4, 8);
    const poleMaterial = new THREE.MeshStandardMaterial({
      color: 0x555555,
      roughness: 0.5,
      metalness: 0.6
    });
    const pole = new THREE.Mesh(poleGeometry, poleMaterial);
    pole.position.set(2.5, 2, 2.5);
    this.meshGroup.add(pole);

    // Flag (shows ownership)
    const flagGeometry = new THREE.PlaneGeometry(1.5, 1);
    const flagMaterial = new THREE.MeshStandardMaterial({
      color: color,
      roughness: 0.8,
      metalness: 0.1,
      side: THREE.DoubleSide,
      emissive: color,
      emissiveIntensity: 0.2
    });
    this.flagMesh = new THREE.Mesh(flagGeometry, flagMaterial);
    this.flagMesh.position.set(3.25, 3.5, 2.5);
    this.meshGroup.add(this.flagMesh);

    // Range indicator (shows claim radius)
    const rangeGeometry = new THREE.RingGeometry(7.5, 8, 32);
    const rangeMaterial = new THREE.MeshBasicMaterial({
      color: color,
      side: THREE.DoubleSide,
      transparent: true,
      opacity: 0.25
    });
    this.rangeMesh = new THREE.Mesh(rangeGeometry, rangeMaterial);
    this.rangeMesh.rotation.x = -Math.PI / 2;
    this.rangeMesh.position.y = 0.05;
    this.meshGroup.add(this.rangeMesh);

    // Health bar background
    this.healthBarGroup = new THREE.Group();
    this.healthBarGroup.position.y = 5.5;

    const healthBarBg = new THREE.Mesh(
      new THREE.BoxGeometry(4, 0.3, 0.1),
      new THREE.MeshBasicMaterial({ color: 0x333333 })
    );
    this.healthBarGroup.add(healthBarBg);

    // Health bar fill
    this.healthBarFill = new THREE.Mesh(
      new THREE.BoxGeometry(4, 0.3, 0.12),
      new THREE.MeshBasicMaterial({ color: 0x00ff00 })
    );
    this.healthBarFill.position.z = 0.01;
    this.healthBarGroup.add(this.healthBarFill);

    this.meshGroup.add(this.healthBarGroup);

    // Healing indicator (glowing cross icon) - hidden by default
    this.healingIndicator = new THREE.Group();
    this.healingIndicator.position.y = 6;
    this.healingIndicator.visible = false;

    // Cross shape using two boxes
    const crossMaterial = new THREE.MeshBasicMaterial({
      color: 0x00ff00,
      transparent: true,
      opacity: 0.8
    });
    const crossVertical = new THREE.Mesh(
      new THREE.BoxGeometry(0.3, 1.2, 0.1),
      crossMaterial
    );
    const crossHorizontal = new THREE.Mesh(
      new THREE.BoxGeometry(1.2, 0.3, 0.1),
      crossMaterial
    );
    this.healingIndicator.add(crossVertical);
    this.healingIndicator.add(crossHorizontal);

    // Outer glow ring
    const glowGeometry = new THREE.RingGeometry(0.8, 1, 16);
    const glowMaterial = new THREE.MeshBasicMaterial({
      color: 0x00ff00,
      transparent: true,
      opacity: 0.4,
      side: THREE.DoubleSide
    });
    this.healingGlow = new THREE.Mesh(glowGeometry, glowMaterial);
    this.healingIndicator.add(this.healingGlow);

    this.meshGroup.add(this.healingIndicator);

    // Position the group
    this.meshGroup.position.set(
      this.barracks.position.x,
      this.barracks.position.y,
      this.barracks.position.z
    );

    this.scene.add(this.meshGroup);
  }

  getColor() {
    if (this.barracks.ownerId === -1) {
      return 0x888888; // Gray for neutral
    }
    // Parse hex color string to number
    const colorStr = PLAYER_COLORS[this.barracks.ownerId];
    return parseInt(colorStr.replace('#', ''), 16);
  }

  update(barracksData, time) {
    this.barracks = barracksData;

    // Update color if owner changed
    const color = this.getColor();
    if (this.flagMesh) {
      this.flagMesh.material.color.setHex(color);
      this.flagMesh.material.emissive.setHex(color);
      // Wave the flag
      this.flagMesh.rotation.y = Math.sin(time * 3) * 0.2;
    }
    if (this.rangeMesh) {
      this.rangeMesh.material.color.setHex(color);
    }

    // Handle destroyed state
    if (this.barracks.isDestroyed) {
      // Make building look damaged
      if (this.buildingMesh) {
        this.buildingMesh.material.color.setHex(0x333333);
        this.buildingMesh.rotation.z = 0.1; // Tilted
      }
      if (this.roofMesh) {
        this.roofMesh.visible = false; // Roof destroyed
      }
      if (this.flagMesh) {
        this.flagMesh.visible = false;
      }
      // Show respawn timer through opacity
      const alpha = 0.3 + 0.2 * Math.sin(time * 4);
      this.meshGroup.traverse((child) => {
        if (child.isMesh && child.material) {
          child.material.transparent = true;
          child.material.opacity = alpha;
        }
      });
    } else {
      // Reset building appearance
      if (this.buildingMesh) {
        this.buildingMesh.material.color.setHex(0x6b6b4c);
        this.buildingMesh.rotation.z = 0;
      }
      if (this.roofMesh) {
        this.roofMesh.visible = true;
      }
      if (this.flagMesh) {
        this.flagMesh.visible = true;
      }
      this.meshGroup.traverse((child) => {
        if (child.isMesh && child.material) {
          child.material.transparent = false;
          child.material.opacity = 1.0;
        }
      });

      // Pulse range indicator for neutral barracks
      if (this.barracks.ownerId === -1 && this.rangeMesh) {
        const pulse = 0.15 + Math.sin(time * 2) * 0.1;
        this.rangeMesh.material.opacity = pulse;
      } else if (this.rangeMesh) {
        this.rangeMesh.material.opacity = 0.25;
      }
    }

    // Update health bar
    this.updateHealthBar();

    // Update healing indicator based on occupant count
    this.updateHealingIndicator(time);
  }

  updateHealingIndicator(time) {
    if (!this.healingIndicator) return;

    const hasOccupants = this.barracks.occupantCount > 0;
    this.healingIndicator.visible = hasOccupants && !this.barracks.isDestroyed;

    if (hasOccupants && !this.barracks.isDestroyed) {
      // Pulse the indicator
      const pulse = 0.6 + Math.sin(time * 4) * 0.4;
      this.healingIndicator.children.forEach(child => {
        if (child.material) {
          child.material.opacity = pulse;
        }
      });

      // Rotate the glow
      if (this.healingGlow) {
        this.healingGlow.rotation.z = time * 2;
      }

      // Scale pulse
      const scale = 1 + Math.sin(time * 3) * 0.1;
      this.healingIndicator.scale.set(scale, scale, scale);
    }
  }

  updateHealthBar() {
    if (!this.healthBarFill) return;

    const healthPercent = this.barracks.health / this.barracks.maxHealth;
    this.healthBarFill.scale.x = Math.max(0.01, healthPercent);
    this.healthBarFill.position.x = (healthPercent - 1) * 2; // Shift left as it decreases

    // Change color based on health
    if (healthPercent > 0.6) {
      this.healthBarFill.material.color.setHex(0x00ff00); // Green
    } else if (healthPercent > 0.3) {
      this.healthBarFill.material.color.setHex(0xffff00); // Yellow
    } else {
      this.healthBarFill.material.color.setHex(0xff0000); // Red
    }

    // Hide health bar when at full health
    this.healthBarGroup.visible = healthPercent < 1;
  }

  // Check if unit is in claiming range
  isInRange(position) {
    const dx = position.x - this.barracks.position.x;
    const dz = position.z - this.barracks.position.z;
    const distSq = dx * dx + dz * dz;
    return distSq <= this.barracks.claimRadius * this.barracks.claimRadius;
  }

  // Check if barracks can be claimed by player
  canBeClaimed(playerId) {
    if (this.barracks.isDestroyed) return false;
    if (this.barracks.ownerId === playerId) return false;
    return true;
  }

  remove() {
    this.scene.remove(this.meshGroup);
  }
}
