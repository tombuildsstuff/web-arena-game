import * as THREE from 'three';
import { HealthBar } from './HealthBar.js';

export class Airplane {
  constructor(scene, unit, color, isSuper = false) {
    this.scene = scene;
    this.unit = unit;
    this.color = color;
    this.isSuper = isSuper;
    this.mesh = null;
    this.lastPosition = null;
    this.targetRotation = 0;
    this.currentBank = 0;
    this.healthBar = null;
    this.maxHealth = isSuper ? 90 : 30; // SuperHelicopterHealth vs AirplaneHealth
    this.create();
  }

  create() {
    // Scale factor for super helicopters (1.3x larger)
    const scale = this.isSuper ? 1.3 : 1.0;

    // Create parent group for positioning and yaw rotation
    this.mesh = new THREE.Group();

    // Create a yaw group (rotates around Y axis for heading)
    this.yawGroup = new THREE.Group();
    this.mesh.add(this.yawGroup);

    // Create a pitch/bank group for tilting
    this.bankGroup = new THREE.Group();
    this.yawGroup.add(this.bankGroup);

    // Airplane body (fuselage) - cone pointing along -Z axis
    const bodyGeometry = new THREE.ConeGeometry(0.8 * scale, 4 * scale, 8);
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

    // Add gold stripe for super helicopters
    if (this.isSuper) {
      const stripeGeometry = new THREE.ConeGeometry(0.85 * scale, 0.8 * scale, 8);
      const stripeMaterial = new THREE.MeshStandardMaterial({
        color: 0xffd700, // Gold
        roughness: 0.3,
        metalness: 0.8,
        emissive: 0xffd700,
        emissiveIntensity: 0.2
      });
      const stripe = new THREE.Mesh(stripeGeometry, stripeMaterial);
      stripe.position.y = 0; // Center of fuselage in local space
      this.bodyMesh.add(stripe);
    }

    // Wings - attached to body mesh
    // After -90 degree X rotation: local Y points to +Z (forward), local Z points to +Y (up)
    // Wings should be wide in X, thin in world Y (local Z), with depth in world Z (local -Y)
    const wingGeometry = new THREE.BoxGeometry(8 * scale, 0.2 * scale, 2 * scale);
    const wings = new THREE.Mesh(wingGeometry, bodyMaterial);
    wings.position.y = -0.5 * scale; // Local -Y = world +Z (slightly forward of center)
    this.bodyMesh.add(wings);

    // Tail wing - horizontal stabilizer
    const tailGeometry = new THREE.BoxGeometry(3 * scale, 0.2 * scale, 1 * scale);
    const tail = new THREE.Mesh(tailGeometry, bodyMaterial);
    tail.position.y = 1.8 * scale; // At the back of the fuselage (local +Y = world -Z)
    this.bodyMesh.add(tail);

    // Vertical stabilizer (tail fin)
    const finGeometry = new THREE.BoxGeometry(0.2 * scale, 1.5 * scale, 1 * scale);
    const fin = new THREE.Mesh(finGeometry, bodyMaterial);
    fin.position.y = 1.8 * scale;
    fin.position.z = 0.5 * scale; // Local +Z = world +Y (pointing up)
    this.bodyMesh.add(fin);

    // Create health bar above airplane (on parent group so it stays upright)
    this.healthBar = new HealthBar(this.maxHealth, 3 * scale, 0.35);
    this.healthBar.getGroup().position.y = 3 * scale;
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
