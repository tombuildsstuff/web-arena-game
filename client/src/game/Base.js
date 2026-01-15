import * as THREE from 'three';
import { BASE_SIZE, ARENA_SIZE } from '../utils/constants.js';

export class Base {
  constructor(scene, position, color) {
    this.scene = scene;
    this.position = position;
    this.color = color;
    this.meshes = [];
    this.create();
  }

  create() {
    // Calculate clipped base boundaries to stay within arena
    const arenaHalf = ARENA_SIZE / 2;
    const halfBaseWidth = BASE_SIZE / 2;
    const halfBaseDepth = (BASE_SIZE * 1.15) / 2; // 15% taller in Z direction

    // Calculate the actual bounds, clamped to arena edges
    const westEdge = Math.max(this.position.x - halfBaseWidth, -arenaHalf);
    const eastEdge = Math.min(this.position.x + halfBaseWidth, arenaHalf);
    const northEdge = Math.max(this.position.z - halfBaseDepth, -arenaHalf);
    const southEdge = Math.min(this.position.z + halfBaseDepth, arenaHalf);

    // Calculate actual width and depth after clamping
    const actualWidth = eastEdge - westEdge;
    const actualDepth = southEdge - northEdge;
    const centerX = (westEdge + eastEdge) / 2;
    const centerZ = (northEdge + southEdge) / 2;

    // Shaded ground area - semi-transparent rectangle (clipped to arena)
    const groundGeometry = new THREE.PlaneGeometry(actualWidth, actualDepth);
    const groundMaterial = new THREE.MeshStandardMaterial({
      color: this.color,
      transparent: true,
      opacity: 0.3,
      roughness: 0.9,
      metalness: 0.1
    });

    const ground = new THREE.Mesh(groundGeometry, groundMaterial);
    ground.rotation.x = -Math.PI / 2;
    ground.position.set(centerX, 0.05, centerZ);
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

    // Create 4 border edges (using clipped dimensions)
    const borders = [
      // North edge
      { pos: [centerX, borderHeight / 2, northEdge], size: [actualWidth, borderHeight, borderThickness] },
      // South edge
      { pos: [centerX, borderHeight / 2, southEdge], size: [actualWidth, borderHeight, borderThickness] },
      // West edge
      { pos: [westEdge, borderHeight / 2, centerZ], size: [borderThickness, borderHeight, actualDepth] },
      // East edge
      { pos: [eastEdge, borderHeight / 2, centerZ], size: [borderThickness, borderHeight, actualDepth] }
    ];

    borders.forEach(border => {
      const geometry = new THREE.BoxGeometry(border.size[0], border.size[1], border.size[2]);
      const mesh = new THREE.Mesh(geometry, borderMaterial);
      mesh.position.set(border.pos[0], border.pos[1], border.pos[2]);
      mesh.receiveShadow = true;
      this.scene.add(mesh);
      this.meshes.push(mesh);
    });

    // Corner markers - small glowing posts (using clipped corners)
    const cornerPositions = [
      [westEdge, northEdge],
      [eastEdge, northEdge],
      [westEdge, southEdge],
      [eastEdge, southEdge]
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
