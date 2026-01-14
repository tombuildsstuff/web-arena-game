import * as THREE from 'three';
import { Arena } from './Arena.js';
import { Base } from './Base.js';
import { PLAYER_COLORS, BASE_POSITIONS } from '../utils/constants.js';

export class Scene {
  constructor() {
    this.scene = new THREE.Scene();
    this.bases = [];
    this.setupLighting();
    this.createArena();
  }

  setupLighting() {
    // Ambient light
    const ambientLight = new THREE.AmbientLight(0x404040, 1.5);
    this.scene.add(ambientLight);

    // Directional light (sun)
    const directionalLight = new THREE.DirectionalLight(0xffffff, 1);
    directionalLight.position.set(50, 100, 50);
    directionalLight.castShadow = true;

    // Shadow camera settings
    directionalLight.shadow.camera.left = -150;
    directionalLight.shadow.camera.right = 150;
    directionalLight.shadow.camera.top = 150;
    directionalLight.shadow.camera.bottom = -150;
    directionalLight.shadow.camera.near = 0.5;
    directionalLight.shadow.camera.far = 500;
    directionalLight.shadow.mapSize.width = 2048;
    directionalLight.shadow.mapSize.height = 2048;

    this.scene.add(directionalLight);

    // Hemisphere light for better ambient
    const hemiLight = new THREE.HemisphereLight(0x8888ff, 0x444444, 0.5);
    this.scene.add(hemiLight);
  }

  createArena() {
    this.arena = new Arena(this.scene);
  }

  createBases(players) {
    // Clear existing bases
    this.bases.forEach(base => base.remove());
    this.bases = [];

    // Create bases for each player
    players.forEach((player, index) => {
      if (player) {
        const position = BASE_POSITIONS[index];
        const color = PLAYER_COLORS[index];
        const base = new Base(this.scene, position, color);
        this.bases.push(base);
      }
    });
  }

  getScene() {
    return this.scene;
  }
}
