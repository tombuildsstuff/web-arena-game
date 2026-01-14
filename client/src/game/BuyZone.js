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

    // Add an icon/symbol on top (simple box for tank, cone for airplane)
    if (this.zone.unitType === 'tank') {
      const iconGeometry = new THREE.BoxGeometry(1.5, 0.8, 2);
      const iconMaterial = new THREE.MeshStandardMaterial({
        color: 0x333333,
        roughness: 0.5,
        metalness: 0.5,
        transparent: true,
        opacity: 1.0
      });
      const icon = new THREE.Mesh(iconGeometry, iconMaterial);
      icon.position.set(
        this.zone.position.x,
        baseY + 0.7,
        this.zone.position.z
      );
      icon.castShadow = true;
      this.scene.add(icon);
      this.iconMesh = icon;
    } else if (this.zone.unitType === 'airplane') {
      const iconGeometry = new THREE.ConeGeometry(0.6, 2, 8);
      const iconMaterial = new THREE.MeshStandardMaterial({
        color: 0x333333,
        roughness: 0.5,
        metalness: 0.5,
        transparent: true,
        opacity: 1.0
      });
      const icon = new THREE.Mesh(iconGeometry, iconMaterial);
      icon.rotation.x = Math.PI / 2;
      icon.position.set(
        this.zone.position.x,
        baseY + 0.7,
        this.zone.position.z
      );
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
