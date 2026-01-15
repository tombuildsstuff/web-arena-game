import * as THREE from 'three';

export class HealthPack {
  constructor(scene, healthPack) {
    this.scene = scene;
    this.healthPack = healthPack;
    this.mesh = null;
    this.bobOffset = Math.random() * Math.PI * 2; // Random starting phase for bobbing
    this.create();
  }

  create() {
    // Create a group to hold all parts
    this.mesh = new THREE.Group();

    // White box base
    const boxGeometry = new THREE.BoxGeometry(2, 1.5, 2);
    const whiteMaterial = new THREE.MeshStandardMaterial({
      color: 0xffffff,
      roughness: 0.3,
      metalness: 0.1
    });
    const box = new THREE.Mesh(boxGeometry, whiteMaterial);
    box.castShadow = true;
    box.receiveShadow = true;
    this.mesh.add(box);

    // Red cross material
    const redMaterial = new THREE.MeshStandardMaterial({
      color: 0xff0000,
      roughness: 0.4,
      metalness: 0.2,
      emissive: 0xff0000,
      emissiveIntensity: 0.2
    });

    // Horizontal part of cross (on top)
    const crossHorizontal = new THREE.BoxGeometry(1.6, 0.3, 0.5);
    const crossH = new THREE.Mesh(crossHorizontal, redMaterial);
    crossH.position.y = 0.8;
    this.mesh.add(crossH);

    // Vertical part of cross (on top)
    const crossVertical = new THREE.BoxGeometry(0.5, 0.3, 1.6);
    const crossV = new THREE.Mesh(crossVertical, redMaterial);
    crossV.position.y = 0.8;
    this.mesh.add(crossV);

    // Add crosses on the sides too
    // Front side cross
    const sideCrossH1 = new THREE.BoxGeometry(1.2, 0.4, 0.1);
    const sideCrossV1 = new THREE.BoxGeometry(0.4, 1.0, 0.1);
    const frontH = new THREE.Mesh(sideCrossH1, redMaterial);
    const frontV = new THREE.Mesh(sideCrossV1, redMaterial);
    frontH.position.set(0, 0, 1.01);
    frontV.position.set(0, 0, 1.01);
    this.mesh.add(frontH);
    this.mesh.add(frontV);

    // Back side cross
    const backH = new THREE.Mesh(sideCrossH1, redMaterial);
    const backV = new THREE.Mesh(sideCrossV1, redMaterial);
    backH.position.set(0, 0, -1.01);
    backV.position.set(0, 0, -1.01);
    this.mesh.add(backH);
    this.mesh.add(backV);

    // Left side cross
    const sideCrossH2 = new THREE.BoxGeometry(0.1, 0.4, 1.2);
    const sideCrossV2 = new THREE.BoxGeometry(0.1, 1.0, 0.4);
    const leftH = new THREE.Mesh(sideCrossH2, redMaterial);
    const leftV = new THREE.Mesh(sideCrossV2, redMaterial);
    leftH.position.set(-1.01, 0, 0);
    leftV.position.set(-1.01, 0, 0);
    this.mesh.add(leftH);
    this.mesh.add(leftV);

    // Right side cross
    const rightH = new THREE.Mesh(sideCrossH2, redMaterial);
    const rightV = new THREE.Mesh(sideCrossV2, redMaterial);
    rightH.position.set(1.01, 0, 0);
    rightV.position.set(1.01, 0, 0);
    this.mesh.add(rightH);
    this.mesh.add(rightV);

    // Set initial position
    this.mesh.position.set(
      this.healthPack.position.x,
      this.healthPack.position.y,
      this.healthPack.position.z
    );

    this.scene.add(this.mesh);
  }

  // Update animation (bobbing and rotation)
  update(time) {
    if (this.mesh) {
      // Gentle bobbing motion
      const bobHeight = Math.sin(time * 2 + this.bobOffset) * 0.3;
      this.mesh.position.y = this.healthPack.position.y + bobHeight;

      // Slow rotation
      this.mesh.rotation.y = time * 0.5;
    }
  }

  remove() {
    if (this.mesh) {
      this.scene.remove(this.mesh);
    }
  }
}
