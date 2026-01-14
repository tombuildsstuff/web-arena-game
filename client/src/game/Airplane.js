import * as THREE from 'three';
import { HealthBar } from './HealthBar.js';

export class Airplane {
  constructor(scene, unit, color) {
    this.scene = scene;
    this.unit = unit;
    this.color = color;
    this.mesh = null;
    this.lastPosition = null;
    this.targetRotation = 0;
    this.currentBank = 0;
    this.healthBar = null;
    this.maxHealth = 30; // AirplaneHealth from server
    this.create();
  }

  create() {
    // Create parent group for positioning and yaw rotation
    this.mesh = new THREE.Group();

    // Create a yaw group (rotates around Y axis for heading)
    this.yawGroup = new THREE.Group();
    this.mesh.add(this.yawGroup);

    // Create a pitch/bank group for tilting
    this.bankGroup = new THREE.Group();
    this.yawGroup.add(this.bankGroup);

    // Airplane body (fuselage) - cone pointing along -Z axis
    const bodyGeometry = new THREE.ConeGeometry(0.8, 4, 8);
    const bodyMaterial = new THREE.MeshStandardMaterial({
      color: this.color,
      roughness: 0.4,
      metalness: 0.6
    });

    this.bodyMesh = new THREE.Mesh(bodyGeometry, bodyMaterial);
    // Rotate cone to point forward (-Z direction): rotate -90 degrees around X
    this.bodyMesh.rotation.x = -Math.PI / 2;
    this.bodyMesh.castShadow = true;
    this.bodyMesh.receiveShadow = true;
    this.bankGroup.add(this.bodyMesh);

    // Wings - attached to body mesh
    // After -90 degree X rotation: local Y points to +Z (forward), local Z points to +Y (up)
    // Wings should be wide in X, thin in world Y (local Z), with depth in world Z (local -Y)
    const wingGeometry = new THREE.BoxGeometry(8, 0.2, 2); // 8 wide (X), 0.2 thin (local Y->world Z), 2 depth (local Z->world Y)
    const wings = new THREE.Mesh(wingGeometry, bodyMaterial);
    wings.position.y = -0.5; // Local -Y = world +Z (slightly forward of center)
    this.bodyMesh.add(wings);

    // Tail wing - horizontal stabilizer
    const tailGeometry = new THREE.BoxGeometry(3, 0.2, 1);
    const tail = new THREE.Mesh(tailGeometry, bodyMaterial);
    tail.position.y = 1.8; // At the back of the fuselage (local +Y = world -Z)
    this.bodyMesh.add(tail);

    // Vertical stabilizer (tail fin)
    const finGeometry = new THREE.BoxGeometry(0.2, 1.5, 1);
    const fin = new THREE.Mesh(finGeometry, bodyMaterial);
    fin.position.y = 1.8;
    fin.position.z = 0.5; // Local +Z = world +Y (pointing up)
    this.bodyMesh.add(fin);

    // Create health bar above airplane (on parent group so it stays upright)
    this.healthBar = new HealthBar(this.maxHealth, 3, 0.35);
    this.healthBar.getGroup().position.y = 3;
    this.mesh.add(this.healthBar.getGroup());

    // Set initial position
    this.mesh.position.set(this.unit.position.x, this.unit.position.y, this.unit.position.z);
    this.lastPosition = { ...this.unit.position };

    this.scene.add(this.mesh);
  }

  updatePosition(position) {
    if (this.mesh && this.yawGroup && this.bankGroup) {
      // Calculate rotation based on movement direction
      if (this.lastPosition) {
        const dx = position.x - this.lastPosition.x;
        const dz = position.z - this.lastPosition.z;
        const distanceMoved = Math.sqrt(dx * dx + dz * dz);

        // Only update rotation if we've moved significantly
        if (distanceMoved > 0.01) {
          // Calculate target rotation to face movement direction
          // atan2(dx, dz) gives angle from +Z axis, we want to face the movement direction
          const newTargetRotation = Math.atan2(dx, dz);

          // Calculate turn rate for banking effect
          let turnRate = newTargetRotation - this.targetRotation;
          while (turnRate > Math.PI) turnRate -= Math.PI * 2;
          while (turnRate < -Math.PI) turnRate += Math.PI * 2;

          // Target bank angle based on turn rate (bank into turns)
          const targetBank = Math.max(-0.5, Math.min(0.5, -turnRate * 2));
          this.currentBank += (targetBank - this.currentBank) * 0.1;

          this.targetRotation = newTargetRotation;
        } else {
          // Reduce bank when not turning
          this.currentBank *= 0.95;
        }
      }

      // Smoothly interpolate yaw rotation (apply to yaw group)
      let rotationDiff = this.targetRotation - this.yawGroup.rotation.y;
      while (rotationDiff > Math.PI) rotationDiff -= Math.PI * 2;
      while (rotationDiff < -Math.PI) rotationDiff += Math.PI * 2;
      this.yawGroup.rotation.y += rotationDiff * 0.15;

      // Apply banking to bank group (roll around local Z axis)
      this.bankGroup.rotation.z = this.currentBank;

      // Store last position
      this.lastPosition = { x: position.x, y: position.y, z: position.z };
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
