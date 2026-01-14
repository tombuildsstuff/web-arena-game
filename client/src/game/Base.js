import * as THREE from 'three';
import { BASE_SIZE } from '../utils/constants.js';

export class Base {
  constructor(scene, position, color) {
    this.scene = scene;
    this.position = position;
    this.color = color;
    this.meshes = [];
    this.create();
  }

  create() {
    // Shaded ground area - semi-transparent rectangle
    const groundGeometry = new THREE.PlaneGeometry(BASE_SIZE, BASE_SIZE);
    const groundMaterial = new THREE.MeshStandardMaterial({
      color: this.color,
      transparent: true,
      opacity: 0.3,
      roughness: 0.9,
      metalness: 0.1
    });

    const ground = new THREE.Mesh(groundGeometry, groundMaterial);
    ground.rotation.x = -Math.PI / 2;
    ground.position.set(this.position.x, 0.05, this.position.z);
    ground.receiveShadow = true;
    this.scene.add(ground);
    this.meshes.push(ground);

    // Border outline - slightly raised edges
    const borderMaterial = new THREE.MeshStandardMaterial({
      color: this.color,
      roughness: 0.5,
      metalness: 0.3,
      emissive: this.color,
      emissiveIntensity: 0.3
    });

    const borderThickness = 0.5;
    const borderHeight = 0.3;

    // Create 4 border edges
    const borders = [
      // North edge
      { pos: [this.position.x, borderHeight / 2, this.position.z - BASE_SIZE / 2], size: [BASE_SIZE, borderHeight, borderThickness] },
      // South edge
      { pos: [this.position.x, borderHeight / 2, this.position.z + BASE_SIZE / 2], size: [BASE_SIZE, borderHeight, borderThickness] },
      // West edge
      { pos: [this.position.x - BASE_SIZE / 2, borderHeight / 2, this.position.z], size: [borderThickness, borderHeight, BASE_SIZE] },
      // East edge
      { pos: [this.position.x + BASE_SIZE / 2, borderHeight / 2, this.position.z], size: [borderThickness, borderHeight, BASE_SIZE] }
    ];

    borders.forEach(border => {
      const geometry = new THREE.BoxGeometry(border.size[0], border.size[1], border.size[2]);
      const mesh = new THREE.Mesh(geometry, borderMaterial);
      mesh.position.set(border.pos[0], border.pos[1], border.pos[2]);
      mesh.receiveShadow = true;
      this.scene.add(mesh);
      this.meshes.push(mesh);
    });

    // Corner markers - small glowing posts
    const cornerPositions = [
      [this.position.x - BASE_SIZE / 2, this.position.z - BASE_SIZE / 2],
      [this.position.x + BASE_SIZE / 2, this.position.z - BASE_SIZE / 2],
      [this.position.x - BASE_SIZE / 2, this.position.z + BASE_SIZE / 2],
      [this.position.x + BASE_SIZE / 2, this.position.z + BASE_SIZE / 2]
    ];

    const cornerGeometry = new THREE.CylinderGeometry(0.4, 0.4, 2, 8);
    const cornerMaterial = new THREE.MeshStandardMaterial({
      color: this.color,
      emissive: this.color,
      emissiveIntensity: 0.5
    });

    cornerPositions.forEach(pos => {
      const corner = new THREE.Mesh(cornerGeometry, cornerMaterial);
      corner.position.set(pos[0], 1, pos[1]);
      corner.castShadow = true;
      this.scene.add(corner);
      this.meshes.push(corner);
    });

    // Add a subtle point light at the base center
    const light = new THREE.PointLight(this.color, 0.5, 20);
    light.position.set(this.position.x, 3, this.position.z);
    this.scene.add(light);
    this.meshes.push(light);
  }

  remove() {
    this.meshes.forEach(mesh => {
      this.scene.remove(mesh);
    });
    this.meshes = [];
  }
}
