import * as THREE from 'three';
import { PLAYER_COLORS } from '../utils/constants.js';

export class Turret {
  constructor(scene, turret) {
    this.scene = scene;
    this.turret = turret;
    this.meshGroup = new THREE.Group();
    this.baseMesh = null;
    this.gunMesh = null;
    this.rangeMesh = null;
    this.create();
  }

  create() {
    // Determine color based on owner
    const color = this.getColor();

    // Base platform
    const baseGeometry = new THREE.CylinderGeometry(2, 2.5, 1, 16);
    const baseMaterial = new THREE.MeshStandardMaterial({
      color: 0x444444,
      roughness: 0.6,
      metalness: 0.4
    });
    this.baseMesh = new THREE.Mesh(baseGeometry, baseMaterial);
    this.baseMesh.position.y = 0.5;
    this.baseMesh.castShadow = true;
    this.baseMesh.receiveShadow = true;
    this.meshGroup.add(this.baseMesh);

    // Turret head (sphere)
    const headGeometry = new THREE.SphereGeometry(1.2, 16, 12);
    const headMaterial = new THREE.MeshStandardMaterial({
      color: color,
      roughness: 0.3,
      metalness: 0.7,
      emissive: color,
      emissiveIntensity: 0.2
    });
    this.headMesh = new THREE.Mesh(headGeometry, headMaterial);
    this.headMesh.position.y = 1.5;
    this.headMesh.castShadow = true;
    this.meshGroup.add(this.headMesh);

    // Gun barrel
    const barrelGeometry = new THREE.CylinderGeometry(0.25, 0.25, 2.5, 8);
    const barrelMaterial = new THREE.MeshStandardMaterial({
      color: 0x333333,
      roughness: 0.4,
      metalness: 0.8
    });
    this.gunMesh = new THREE.Mesh(barrelGeometry, barrelMaterial);
    this.gunMesh.rotation.x = Math.PI / 2;
    this.gunMesh.position.set(0, 1.5, 1.5);
    this.gunMesh.castShadow = true;
    this.meshGroup.add(this.gunMesh);

    // Range indicator (semi-transparent ring)
    const rangeGeometry = new THREE.RingGeometry(
      this.turret.claimRadius - 0.3,
      this.turret.claimRadius,
      32
    );
    const rangeMaterial = new THREE.MeshBasicMaterial({
      color: color,
      side: THREE.DoubleSide,
      transparent: true,
      opacity: 0.3
    });
    this.rangeMesh = new THREE.Mesh(rangeGeometry, rangeMaterial);
    this.rangeMesh.rotation.x = -Math.PI / 2;
    this.rangeMesh.position.y = 0.05;
    this.meshGroup.add(this.rangeMesh);

    // Position the group
    this.meshGroup.position.set(
      this.turret.position.x,
      this.turret.position.y,
      this.turret.position.z
    );

    this.scene.add(this.meshGroup);
  }

  getColor() {
    if (this.turret.ownerId === -1) {
      return 0x888888; // Gray for unclaimed
    }
    // Parse hex color string to number
    const colorStr = PLAYER_COLORS[this.turret.ownerId];
    return parseInt(colorStr.replace('#', ''), 16);
  }

  update(turretData, time) {
    this.turret = turretData;

    // Update color if owner changed
    const color = this.getColor();
    if (this.headMesh) {
      this.headMesh.material.color.setHex(color);
      this.headMesh.material.emissive.setHex(color);
    }
    if (this.rangeMesh) {
      this.rangeMesh.material.color.setHex(color);
    }

    // Handle destroyed state
    if (this.turret.isDestroyed) {
      this.meshGroup.visible = false;
    } else {
      this.meshGroup.visible = true;

      // Animate the turret (slow rotation when active)
      if (this.turret.ownerId !== -1) {
        // Rotate when owned
        this.gunMesh.rotation.y = Math.sin(time * 0.5) * 0.3;
      }

      // Pulse effect for unclaimed turrets
      if (this.turret.ownerId === -1 && this.rangeMesh) {
        const pulse = 0.2 + Math.sin(time * 2) * 0.1;
        this.rangeMesh.material.opacity = pulse;
      } else if (this.rangeMesh) {
        this.rangeMesh.material.opacity = 0.3;
      }
    }
  }

  // Check if player is in claiming range
  isInRange(position) {
    const dx = position.x - this.turret.position.x;
    const dz = position.z - this.turret.position.z;
    const distSq = dx * dx + dz * dz;
    return distSq <= this.turret.claimRadius * this.turret.claimRadius;
  }

  // Check if turret can be claimed by player
  canBeClaimed(playerId) {
    if (this.turret.isDestroyed) return false;
    if (this.turret.ownerId === playerId) return false;
    return true;
  }

  remove() {
    this.scene.remove(this.meshGroup);
  }
}
