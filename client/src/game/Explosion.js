import * as THREE from 'three';

export class Explosion {
  constructor(scene, position, onComplete) {
    this.scene = scene;
    this.position = position.clone();
    this.onComplete = onComplete;
    this.particles = [];
    this.startTime = performance.now();
    this.duration = 600; // 600ms explosion
    this.isComplete = false;
    this.create();
  }

  create() {
    // Create particle system for explosion
    const particleCount = 25;
    const colors = [0xff4400, 0xff6600, 0xff8800, 0xffaa00, 0xffcc00];

    for (let i = 0; i < particleCount; i++) {
      const size = 0.3 + Math.random() * 0.5;
      const geometry = new THREE.SphereGeometry(size, 6, 6);
      const material = new THREE.MeshBasicMaterial({
        color: colors[Math.floor(Math.random() * colors.length)],
        transparent: true,
        opacity: 1
      });

      const particle = new THREE.Mesh(geometry, material);
      particle.position.copy(this.position);

      // Random velocity direction (explosion outward)
      const theta = Math.random() * Math.PI * 2;
      const phi = Math.random() * Math.PI;
      const speed = 10 + Math.random() * 15;

      particle.userData.velocity = new THREE.Vector3(
        Math.sin(phi) * Math.cos(theta) * speed,
        Math.abs(Math.cos(phi)) * speed * 0.8 + Math.random() * 5, // Bias upward
        Math.sin(phi) * Math.sin(theta) * speed
      );

      particle.userData.rotationSpeed = new THREE.Vector3(
        (Math.random() - 0.5) * 10,
        (Math.random() - 0.5) * 10,
        (Math.random() - 0.5) * 10
      );

      this.particles.push(particle);
      this.scene.add(particle);
    }

    // Add flash light
    this.light = new THREE.PointLight(0xff6600, 3, 30);
    this.light.position.copy(this.position);
    this.scene.add(this.light);

    // Add central bright flash
    const flashGeometry = new THREE.SphereGeometry(2, 8, 8);
    const flashMaterial = new THREE.MeshBasicMaterial({
      color: 0xffffff,
      transparent: true,
      opacity: 1
    });
    this.flash = new THREE.Mesh(flashGeometry, flashMaterial);
    this.flash.position.copy(this.position);
    this.scene.add(this.flash);
  }

  update() {
    if (this.isComplete) return false;

    const elapsed = performance.now() - this.startTime;
    const progress = Math.min(elapsed / this.duration, 1);

    if (progress >= 1) {
      this.remove();
      this.isComplete = true;
      if (this.onComplete) {
        this.onComplete();
      }
      return false;
    }

    const deltaTime = 0.016; // Approximate 60fps

    // Update particles
    for (const particle of this.particles) {
      // Apply velocity
      particle.position.add(
        particle.userData.velocity.clone().multiplyScalar(deltaTime)
      );

      // Apply gravity
      particle.userData.velocity.y -= 25 * deltaTime;

      // Apply rotation
      particle.rotation.x += particle.userData.rotationSpeed.x * deltaTime;
      particle.rotation.y += particle.userData.rotationSpeed.y * deltaTime;
      particle.rotation.z += particle.userData.rotationSpeed.z * deltaTime;

      // Fade out
      particle.material.opacity = 1 - progress;

      // Shrink
      const scale = 1 - progress * 0.5;
      particle.scale.setScalar(scale);
    }

    // Update flash
    if (this.flash) {
      this.flash.material.opacity = Math.max(0, 1 - progress * 3);
      this.flash.scale.setScalar(1 + progress * 3);
    }

    // Fade light
    if (this.light) {
      this.light.intensity = 3 * (1 - progress);
    }

    return true;
  }

  remove() {
    // Remove particles
    for (const particle of this.particles) {
      this.scene.remove(particle);
      particle.geometry.dispose();
      particle.material.dispose();
    }
    this.particles = [];

    // Remove light
    if (this.light) {
      this.scene.remove(this.light);
      this.light = null;
    }

    // Remove flash
    if (this.flash) {
      this.scene.remove(this.flash);
      this.flash.geometry.dispose();
      this.flash.material.dispose();
      this.flash = null;
    }
  }
}
