import * as THREE from 'three';
import { PLAYER_COLORS } from '../utils/constants.js';

// Neutral/claimable zone color (grey, matching unclaimed turrets)
const NEUTRAL_COLOR = 0x888888;

export class BuyZone {
  constructor(scene, zone) {
    this.scene = scene;
    this.zone = zone;
    this.mesh = null;
    this.label = null;
    this.glowMesh = null;
    this.canAfford = true;
    this.baseColor = this.getZoneColor(zone.ownerId);
    this.create();
  }

  getZoneColor(ownerId) {
    if (ownerId === -1) {
      return NEUTRAL_COLOR;
    }
    return PLAYER_COLORS[ownerId];
  }

  create() {
    const color = this.getZoneColor(this.zone.ownerId);

    // Create a circular platform for the buy zone
    const platformGeometry = new THREE.CylinderGeometry(
      this.zone.radius,
      this.zone.radius,
      0.3,
      32
    );
    const platformMaterial = new THREE.MeshStandardMaterial({
      color: color,
      roughness: 0.3,
      metalness: 0.7,
      transparent: true,
      opacity: 0.8
    });

    this.mesh = new THREE.Mesh(platformGeometry, platformMaterial);
    const baseY = this.zone.position.y || 0;
    this.mesh.position.set(
      this.zone.position.x,
      baseY + 0.15,
      this.zone.position.z
    );
    this.mesh.receiveShadow = true;
    this.scene.add(this.mesh);

    // Add a glowing ring around the platform
    const ringGeometry = new THREE.RingGeometry(
      this.zone.radius - 0.2,
      this.zone.radius + 0.2,
      32
    );
    const ringMaterial = new THREE.MeshBasicMaterial({
      color: color,
      side: THREE.DoubleSide,
      transparent: true,
      opacity: 0.6
    });
    this.glowMesh = new THREE.Mesh(ringGeometry, ringMaterial);
    this.glowMesh.rotation.x = -Math.PI / 2;
    this.glowMesh.position.set(
      this.zone.position.x,
      baseY + 0.35,
      this.zone.position.z
    );
    this.scene.add(this.glowMesh);

    // Add an icon/symbol on top based on unit type
    this.createUnitIcon(baseY);
  }

  createUnitIcon(baseY) {
    const unitType = this.zone.unitType;
    if (!unitType) return; // No icon for base zones

    // Determine if this is a super unit (larger, with gold accent)
    const isSuper = unitType.startsWith('super_');
    const scale = isSuper ? 1.3 : 1.0;

    // Base icon color (dark grey, or gold tinted for super units)
    const iconColor = isSuper ? 0x555544 : 0x333333;

    const iconMaterial = new THREE.MeshStandardMaterial({
      color: iconColor,
      roughness: 0.5,
      metalness: 0.5,
      transparent: true,
      opacity: 1.0
    });

    let icon;

    if (unitType === 'tank' || unitType === 'super_tank') {
      // Tank icon - box shape
      const iconGeometry = new THREE.BoxGeometry(1.5 * scale, 0.8 * scale, 2 * scale);
      icon = new THREE.Mesh(iconGeometry, iconMaterial);
      icon.position.set(
        this.zone.position.x,
        baseY + 0.7 * scale,
        this.zone.position.z
      );

      // Add gold stripe for super tanks
      if (isSuper) {
        const stripeGeometry = new THREE.BoxGeometry(1.6 * scale, 0.2 * scale, 2.1 * scale);
        const stripeMaterial = new THREE.MeshStandardMaterial({
          color: 0xffd700,
          roughness: 0.3,
          metalness: 0.8,
          emissive: 0xffd700,
          emissiveIntensity: 0.3
        });
        const stripe = new THREE.Mesh(stripeGeometry, stripeMaterial);
        stripe.position.y = 0.3 * scale;
        icon.add(stripe);
      }
    } else if (unitType === 'airplane' || unitType === 'super_helicopter') {
      // Airplane/helicopter icon - cone shape
      const iconGeometry = new THREE.ConeGeometry(0.6 * scale, 2 * scale, 8);
      icon = new THREE.Mesh(iconGeometry, iconMaterial);
      icon.rotation.x = Math.PI / 2;
      icon.position.set(
        this.zone.position.x,
        baseY + 0.7 * scale,
        this.zone.position.z
      );

      // Add gold ring for super helicopters
      if (isSuper) {
        const ringGeometry = new THREE.TorusGeometry(0.5 * scale, 0.1 * scale, 8, 16);
        const ringMaterial = new THREE.MeshStandardMaterial({
          color: 0xffd700,
          roughness: 0.3,
          metalness: 0.8,
          emissive: 0xffd700,
          emissiveIntensity: 0.3
        });
        const ring = new THREE.Mesh(ringGeometry, ringMaterial);
        ring.rotation.x = Math.PI / 2;
        icon.add(ring);
      }
    }

    if (icon) {
      icon.castShadow = true;
      this.scene.add(icon);
      this.iconMesh = icon;
    }
  }

  // Update for animation (pulsing glow)
  update(time) {
    if (this.glowMesh) {
      if (this.canAfford) {
        const pulse = 0.4 + Math.sin(time * 3) * 0.2;
        this.glowMesh.material.opacity = pulse;
      } else {
        // Dimmed when can't afford
        this.glowMesh.material.opacity = 0.1;
      }
    }
  }

  // Set whether the player can afford this zone's unit
  setAffordable(canAfford) {
    if (this.canAfford === canAfford) return;

    this.canAfford = canAfford;

    if (canAfford) {
      // Restore normal appearance
      if (this.mesh) {
        this.mesh.material.opacity = 0.8;
        this.mesh.material.color.setHex(this.baseColor);
      }
      if (this.iconMesh) {
        this.iconMesh.material.opacity = 1.0;
      }
    } else {
      // Dim the zone
      if (this.mesh) {
        this.mesh.material.opacity = 0.3;
        this.mesh.material.color.setHex(0x444444);
      }
      if (this.iconMesh) {
        this.iconMesh.material.opacity = 0.3;
      }
    }
  }

  // Check if a position is within this buy zone
  isInRange(position) {
    const dx = position.x - this.zone.position.x;
    const dz = position.z - this.zone.position.z;
    const distSq = dx * dx + dz * dz;
    return distSq <= this.zone.radius * this.zone.radius;
  }

  // Update the owner of this zone (when claimed)
  updateOwner(newOwnerId) {
    if (this.zone.ownerId === newOwnerId) return;

    this.zone.ownerId = newOwnerId;
    this.baseColor = this.getZoneColor(newOwnerId);

    // Update mesh colors
    if (this.mesh) {
      this.mesh.material.color.setHex(this.baseColor);
    }
    if (this.glowMesh) {
      this.glowMesh.material.color.setHex(this.baseColor);
    }
  }

  remove() {
    if (this.mesh) {
      this.scene.remove(this.mesh);
    }
    if (this.glowMesh) {
      this.scene.remove(this.glowMesh);
    }
    if (this.iconMesh) {
      this.scene.remove(this.iconMesh);
    }
  }
}
