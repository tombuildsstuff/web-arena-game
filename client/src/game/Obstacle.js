import * as THREE from 'three';

export class Obstacle {
  constructor(scene, obstacleData) {
    this.scene = scene;
    this.data = obstacleData;
    this.mesh = null;
    this.create();
  }

  create() {
    const { type, position, size, rotation } = this.data;

    let geometry, material;

    switch (type) {
      case 'wall':
        geometry = new THREE.BoxGeometry(size.x, size.y, size.z);
        material = new THREE.MeshStandardMaterial({
          color: 0x555555,
          roughness: 0.8,
          metalness: 0.2
        });
        break;

      case 'pillar':
        // Use cylinder for pillars
        geometry = new THREE.CylinderGeometry(size.x / 2, size.x / 2, size.y, 12);
        material = new THREE.MeshStandardMaterial({
          color: 0x666666,
          roughness: 0.7,
          metalness: 0.3
        });
        break;

      case 'cover':
        geometry = new THREE.BoxGeometry(size.x, size.y, size.z);
        material = new THREE.MeshStandardMaterial({
          color: 0x4a4a4a,
          roughness: 0.9,
          metalness: 0.1
        });
        break;

      case 'ramp':
        geometry = this.createRampGeometry(size, this.data.elevationStart, this.data.elevationEnd);
        material = new THREE.MeshStandardMaterial({
          color: 0x3d3d3d,
          roughness: 0.6,
          metalness: 0.2
        });
        break;

      default:
        // Default to box
        geometry = new THREE.BoxGeometry(size.x, size.y, size.z);
        material = new THREE.MeshStandardMaterial({
          color: 0x444444,
          roughness: 0.8,
          metalness: 0.2
        });
    }

    this.mesh = new THREE.Mesh(geometry, material);

    // Position - center at position, with Y offset for height
    if (type === 'ramp') {
      // Ramps are positioned differently
      this.mesh.position.set(position.x, position.y, position.z);
    } else {
      this.mesh.position.set(position.x, position.y + size.y / 2, position.z);
    }

    // Apply rotation around Y axis
    if (rotation) {
      this.mesh.rotation.y = rotation;
    }

    this.mesh.castShadow = true;
    this.mesh.receiveShadow = true;

    this.scene.add(this.mesh);
  }

  createRampGeometry(size, elevStart, elevEnd) {
    // Create a ramp using BufferGeometry
    // The ramp goes from elevStart at -Z to elevEnd at +Z
    const halfX = size.x / 2;
    const halfZ = size.z / 2;

    // Define vertices for the ramp shape
    const vertices = new Float32Array([
      // Bottom face (y = 0)
      -halfX, 0, -halfZ,
       halfX, 0, -halfZ,
       halfX, 0,  halfZ,
      -halfX, 0,  halfZ,

      // Top face (sloped)
      -halfX, elevStart, -halfZ,
       halfX, elevStart, -halfZ,
       halfX, elevEnd,    halfZ,
      -halfX, elevEnd,    halfZ,
    ]);

    // Define faces using indices
    const indices = [
      // Bottom
      0, 1, 2,  0, 2, 3,
      // Top (sloped surface)
      4, 6, 5,  4, 7, 6,
      // Front (low end)
      0, 4, 5,  0, 5, 1,
      // Back (high end)
      2, 6, 7,  2, 7, 3,
      // Left side
      0, 3, 7,  0, 7, 4,
      // Right side
      1, 5, 6,  1, 6, 2,
    ];

    const geometry = new THREE.BufferGeometry();
    geometry.setAttribute('position', new THREE.BufferAttribute(vertices, 3));
    geometry.setIndex(indices);
    geometry.computeVertexNormals();

    return geometry;
  }

  remove() {
    if (this.mesh) {
      this.scene.remove(this.mesh);
      if (this.mesh.geometry) {
        this.mesh.geometry.dispose();
      }
      if (this.mesh.material) {
        this.mesh.material.dispose();
      }
      this.mesh = null;
    }
  }
}
