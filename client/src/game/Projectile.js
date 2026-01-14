import * as THREE from 'three';

export class Projectile {
  constructor(scene, projectileData, color) {
    this.scene = scene;
    this.data = projectileData;
    this.mesh = null;
    this.trail = null;
    this.trailPositions = [];
    this.maxTrailLength = 10;
    this.create(color);
  }

  create(color) {
    // Create projectile as a glowing sphere
    const geometry = new THREE.SphereGeometry(0.4, 8, 8);
    const material = new THREE.MeshBasicMaterial({
      color: color,
      transparent: true,
      opacity: 0.9
    });

    this.mesh = new THREE.Mesh(geometry, material);
    this.updatePosition(this.data.position);

    // Add a point light for glow effect
    this.light = new THREE.PointLight(color, 0.5, 5);
    this.mesh.add(this.light);

    // Create trail
    this.createTrail(color);

    this.scene.add(this.mesh);
  }

  createTrail(color) {
    // Create a trail using a line
    const trailGeometry = new THREE.BufferGeometry();
    const positions = new Float32Array(this.maxTrailLength * 3);
    trailGeometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));

    const trailMaterial = new THREE.LineBasicMaterial({
      color: color,
      transparent: true,
      opacity: 0.4
    });

    this.trail = new THREE.Line(trailGeometry, trailMaterial);
    this.scene.add(this.trail);
  }

  updatePosition(position) {
    if (this.mesh) {
      this.mesh.position.set(position.x, position.y, position.z);

      // Update trail
      this.updateTrail(position);
    }
  }

  updateTrail(position) {
    if (!this.trail) return;

    // Add current position to trail history
    this.trailPositions.push({ x: position.x, y: position.y, z: position.z });

    // Limit trail length
    while (this.trailPositions.length > this.maxTrailLength) {
      this.trailPositions.shift();
    }

    // Update trail geometry
    const positions = this.trail.geometry.attributes.position.array;
    for (let i = 0; i < this.maxTrailLength; i++) {
      if (i < this.trailPositions.length) {
        const pos = this.trailPositions[i];
        positions[i * 3] = pos.x;
        positions[i * 3 + 1] = pos.y;
        positions[i * 3 + 2] = pos.z;
      } else {
        // Fill remaining with last position
        const lastPos = this.trailPositions[this.trailPositions.length - 1] || position;
        positions[i * 3] = lastPos.x;
        positions[i * 3 + 1] = lastPos.y;
        positions[i * 3 + 2] = lastPos.z;
      }
    }
    this.trail.geometry.attributes.position.needsUpdate = true;
    this.trail.geometry.setDrawRange(0, this.trailPositions.length);
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

    if (this.trail) {
      this.scene.remove(this.trail);
      if (this.trail.geometry) {
        this.trail.geometry.dispose();
      }
      if (this.trail.material) {
        this.trail.material.dispose();
      }
      this.trail = null;
    }
  }
}
