import * as THREE from 'three';

export class Airplane {
  constructor(scene, unit, color) {
    this.scene = scene;
    this.unit = unit;
    this.color = color;
    this.mesh = null;
    this.lastPosition = null;
    this.targetRotation = 0;
    this.currentBank = 0;
    this.create();
  }

  create() {
    // Airplane body (fuselage)
    const bodyGeometry = new THREE.ConeGeometry(0.8, 4, 8);
    const bodyMaterial = new THREE.MeshStandardMaterial({
      color: this.color,
      roughness: 0.4,
      metalness: 0.6
    });

    this.mesh = new THREE.Mesh(bodyGeometry, bodyMaterial);
    this.mesh.rotation.x = Math.PI / 2; // Point forward
    this.mesh.castShadow = true;
    this.mesh.receiveShadow = true;

    // Wings
    const wingGeometry = new THREE.BoxGeometry(8, 0.2, 2);
    const wings = new THREE.Mesh(wingGeometry, bodyMaterial);
    wings.position.z = -0.5;
    this.mesh.add(wings);

    // Tail wing
    const tailGeometry = new THREE.BoxGeometry(3, 0.2, 1);
    const tail = new THREE.Mesh(tailGeometry, bodyMaterial);
    tail.position.z = -2;
    this.mesh.add(tail);

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
          const newTargetRotation = Math.atan2(dx, dz);

          // Calculate turn rate for banking effect
          let turnRate = newTargetRotation - this.targetRotation;
          while (turnRate > Math.PI) turnRate -= Math.PI * 2;
          while (turnRate < -Math.PI) turnRate += Math.PI * 2;

          // Target bank angle based on turn rate
          const targetBank = Math.max(-0.5, Math.min(0.5, turnRate * 2));
          this.currentBank += (targetBank - this.currentBank) * 0.1;

          this.targetRotation = newTargetRotation;
        } else {
          // Reduce bank when not turning
          this.currentBank *= 0.95;
        }
      }

      // Smoothly interpolate rotation
      let rotationDiff = this.targetRotation - this.mesh.rotation.y;
      while (rotationDiff > Math.PI) rotationDiff -= Math.PI * 2;
      while (rotationDiff < -Math.PI) rotationDiff += Math.PI * 2;
      this.mesh.rotation.y += rotationDiff * 0.15;

      // Apply banking
      this.mesh.rotation.z = this.currentBank;

      // Store last position
      this.lastPosition = { x: position.x, y: position.y, z: position.z };
    }
  }

  remove() {
    if (this.mesh) {
      this.scene.remove(this.mesh);
    }
  }
}
