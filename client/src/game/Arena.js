import * as THREE from 'three';
import { ARENA_SIZE } from '../utils/constants.js';

export class Arena {
  constructor(scene) {
    this.scene = scene;
    this.create();
  }

  create() {
    // Ground plane
    const groundGeometry = new THREE.PlaneGeometry(ARENA_SIZE, ARENA_SIZE);
    const groundMaterial = new THREE.MeshStandardMaterial({
      color: 0x1a1a2e,
      roughness: 0.8,
      metalness: 0.2
    });
    this.ground = new THREE.Mesh(groundGeometry, groundMaterial);
    this.ground.rotation.x = -Math.PI / 2;
    this.ground.receiveShadow = true;
    this.scene.add(this.ground);

    // Grid helper
    const gridHelper = new THREE.GridHelper(ARENA_SIZE, 20, 0x444444, 0x222222);
    this.scene.add(gridHelper);

    // Arena boundaries (walls)
    this.createWalls();
  }

  createWalls() {
    const wallHeight = 5;
    const wallThickness = 1;
    const wallMaterial = new THREE.MeshStandardMaterial({
      color: 0x333333,
      roughness: 0.7,
      metalness: 0.3
    });

    const halfSize = ARENA_SIZE / 2;

    // North wall
    const northWall = new THREE.Mesh(
      new THREE.BoxGeometry(ARENA_SIZE, wallHeight, wallThickness),
      wallMaterial
    );
    northWall.position.set(0, wallHeight / 2, -halfSize);
    northWall.castShadow = true;
    this.scene.add(northWall);

    // South wall
    const southWall = new THREE.Mesh(
      new THREE.BoxGeometry(ARENA_SIZE, wallHeight, wallThickness),
      wallMaterial
    );
    southWall.position.set(0, wallHeight / 2, halfSize);
    southWall.castShadow = true;
    this.scene.add(southWall);

    // East wall
    const eastWall = new THREE.Mesh(
      new THREE.BoxGeometry(wallThickness, wallHeight, ARENA_SIZE),
      wallMaterial
    );
    eastWall.position.set(halfSize, wallHeight / 2, 0);
    eastWall.castShadow = true;
    this.scene.add(eastWall);

    // West wall
    const westWall = new THREE.Mesh(
      new THREE.BoxGeometry(wallThickness, wallHeight, ARENA_SIZE),
      wallMaterial
    );
    westWall.position.set(-halfSize, wallHeight / 2, 0);
    westWall.castShadow = true;
    this.scene.add(westWall);
  }
}
