import * as THREE from 'three';
import { PLAYER_COLORS } from '../utils/constants.js';

export class Turret {
  constructor(scene, turret) {
    this.scene = scene;
    this.turret = turret;
    this.meshGroup = new THREE.Group();
    this.headGroup = null; // Rotating part (head + gun)
    this.baseMesh = null;
    this.gunMesh = null;
    this.rangeMesh = null;
    this.targetRotation = 0;
    this.currentRotation = 0;
    this.create();
  }

  create() {
    // Determine color based on owner
    const color = this.getColor();

    // Base platform (stationary)
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

    // Create rotating head group
    this.headGroup = new THREE.Group();
    this.headGroup.position.y = 1.5;
    this.meshGroup.add(this.headGroup);

    // Turret head (sphere) - part of rotating group
    const headGeometry = new THREE.SphereGeometry(1.2, 16, 12);
    const headMaterial = new THREE.MeshStandardMaterial({
      color: color,
      roughness: 0.3,
      metalness: 0.7,
      emissive: color,
      emissiveIntensity: 0.2
    });
    this.headMesh = new THREE.Mesh(headGeometry, headMaterial);
    this.headMesh.castShadow = true;
    this.headGroup.add(this.headMesh);

    // Gun barrel - part of rotating group
    const barrelGeometry = new THREE.CylinderGeometry(0.25, 0.25, 2.5, 8);
    const barrelMaterial = new THREE.MeshStandardMaterial({
      color: 0x333333,
      roughness: 0.4,
      metalness: 0.8
    });
    this.gunMesh = new THREE.Mesh(barrelGeometry, barrelMaterial);
    this.gunMesh.rotation.x = Math.PI / 2;
    this.gunMesh.position.set(0, 0, 1.5);
    this.gunMesh.castShadow = true;
    this.headGroup.add(this.gunMesh);

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

  update(turretData, time, units = []) {
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

      // Find nearest enemy and rotate toward them
      if (this.turret.ownerId !== -1 && this.headGroup) {
        const nearestEnemy = this.findNearestEnemy(units);
        if (nearestEnemy) {
          // Calculate angle to enemy
          const dx = nearestEnemy.position.x - this.turret.position.x;
          const dz = nearestEnemy.position.z - this.turret.position.z;
          this.targetRotation = Math.atan2(dx, dz);
        }

        // Smoothly interpolate rotation
        let rotationDiff = this.targetRotation - this.currentRotation;
        // Normalize to [-PI, PI]
        while (rotationDiff > Math.PI) rotationDiff -= Math.PI * 2;
        while (rotationDiff < -Math.PI) rotationDiff += Math.PI * 2;
        this.currentRotation += rotationDiff * 0.1;

        this.headGroup.rotation.y = this.currentRotation;

        // Visual feedback for tracking state
        if (this.headMesh && this.turret.isTracking) {
          // Pulse emissive intensity while tracking, brighter as it locks on
          const baseIntensity = 0.2;
          const trackingPulse = Math.sin(time * 8) * 0.15; // Fast pulse while tracking
          const lockOnBoost = this.turret.trackingProgress * 0.4; // Gets brighter as it locks
          this.headMesh.material.emissiveIntensity = baseIntensity + trackingPulse + lockOnBoost;
        } else if (this.headMesh) {
          // Normal emissive when not tracking
          this.headMesh.material.emissiveIntensity = 0.2;
        }
      } else if (this.turret.ownerId === -1 && this.headGroup) {
        // Slow idle rotation for unclaimed turrets
        this.headGroup.rotation.y = time * 0.2;
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

  // Find nearest enemy unit within attack range
  findNearestEnemy(units) {
    if (!units || units.length === 0) return null;

    const attackRange = 20; // Match server TurretAttackRange (~4 squares)
    let nearest = null;
    let nearestDist = attackRange + 1;

    for (const unit of units) {
      // Only target enemy units
      if (unit.ownerId === this.turret.ownerId) continue;
      // Skip dead/respawning units
      if (unit.isRespawning) continue;

      const dx = unit.position.x - this.turret.position.x;
      const dz = unit.position.z - this.turret.position.z;
      const dist = Math.sqrt(dx * dx + dz * dz);

      if (dist < nearestDist) {
        nearestDist = dist;
        nearest = unit;
      }
    }

    return nearest;
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
