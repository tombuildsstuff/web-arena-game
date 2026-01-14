import { Tank } from './Tank.js';
import { Airplane } from './Airplane.js';
import { Player } from './Player.js';
import { Obstacle } from './Obstacle.js';
import { Projectile } from './Projectile.js';
import { Explosion } from './Explosion.js';
import { BuyZone } from './BuyZone.js';
import { Turret } from './Turret.js';
import { PLAYER_COLORS } from '../utils/constants.js';

export class GameLoop {
  constructor(renderer, scene, camera, gameState) {
    this.renderer = renderer;
    this.scene = scene;
    this.camera = camera;
    this.gameState = gameState;
    this.isRunning = false;
    this.lastTime = 0;
    this.elapsedTime = 0; // Total elapsed time for animations
    this.unitMeshes = new Map(); // Map of unit ID to unit object (Tank/Airplane)
    this.obstacleMeshes = new Map(); // Map of obstacle ID to Obstacle object
    this.projectileMeshes = new Map(); // Map of projectile ID to Projectile object
    this.buyZoneMeshes = new Map(); // Map of buy zone ID to BuyZone object
    this.turretMeshes = new Map(); // Map of turret ID to Turret object
    this.explosions = []; // Active explosions
    this.obstaclesInitialized = false;
    this.buyZonesInitialized = false;
    this.turretsInitialized = false;
    this.previousUnitIDs = new Set(); // Track unit IDs for death detection
    this.nearbyBuyZone = null; // The buy zone the player is currently near
    this.nearbyTurret = null; // The turret the player is currently near
    this.buyZonePopup = null; // Reference to the popup for position updates
  }

  setBuyZonePopup(popup) {
    this.buyZonePopup = popup;
  }

  start() {
    this.isRunning = true;
    this.lastTime = performance.now();
    this.animate();
  }

  stop() {
    this.isRunning = false;
  }

  animate() {
    if (!this.isRunning) return;

    requestAnimationFrame(() => this.animate());

    const currentTime = performance.now();
    const deltaTime = (currentTime - this.lastTime) / 1000; // Convert to seconds
    this.lastTime = currentTime;

    this.update(deltaTime);
    this.render();
  }

  update(deltaTime) {
    // Track total elapsed time for animations
    this.elapsedTime += deltaTime;

    // Update camera controls
    this.camera.update();

    // Sync obstacles (once at game start)
    this.syncObstacles();

    // Sync buy zones (once at game start)
    this.syncBuyZones();

    // Sync turrets (once at game start, then update)
    this.syncTurrets();

    // Sync units from game state (with death detection)
    this.syncUnits();

    // Sync projectiles
    this.syncProjectiles();

    // Update explosions
    this.updateExplosions();

    // Update buy zone animations and check proximity
    this.updateBuyZones();

    // Update turret animations and check proximity
    this.updateTurrets();

    // Update unit positions with interpolation
    this.updateUnitPositions(deltaTime);
  }

  syncObstacles() {
    // Only initialize obstacles once
    if (this.obstaclesInitialized) return;

    const stateObstacles = this.gameState.obstacles || [];
    if (stateObstacles.length === 0) return;

    for (const obs of stateObstacles) {
      if (!this.obstacleMeshes.has(obs.id)) {
        const obstacle = new Obstacle(this.scene.getScene(), obs);
        this.obstacleMeshes.set(obs.id, obstacle);
      }
    }

    this.obstaclesInitialized = true;
  }

  syncBuyZones() {
    // Only initialize buy zones once
    if (this.buyZonesInitialized) return;

    const stateBuyZones = this.gameState.buyZones || [];
    if (stateBuyZones.length === 0) return;

    for (const zone of stateBuyZones) {
      if (!this.buyZoneMeshes.has(zone.id)) {
        const buyZone = new BuyZone(this.scene.getScene(), zone);
        this.buyZoneMeshes.set(zone.id, buyZone);
      }
    }

    this.buyZonesInitialized = true;
  }

  updateBuyZones() {
    const myPlayer = this.gameState.getMyPlayer();
    const playerMoney = myPlayer ? myPlayer.money : 0;

    // Update buy zone animations and affordability
    for (const buyZone of this.buyZoneMeshes.values()) {
      buyZone.update(this.elapsedTime);

      // Update affordability for player's own zones
      if (buyZone.zone.ownerId === this.gameState.playerId) {
        const canAfford = playerMoney >= buyZone.zone.cost;
        buyZone.setAffordable(canAfford);
      }
    }

    // Check if player is near any buy zone
    this.nearbyBuyZone = null;
    const playerUnit = this.gameState.getMyPlayerUnit();
    if (!playerUnit || playerUnit.isRespawning) return;

    const playerPos = playerUnit.position;
    for (const [zoneId, buyZone] of this.buyZoneMeshes.entries()) {
      // Only check zones owned by this player
      if (buyZone.zone.ownerId === this.gameState.playerId) {
        if (buyZone.isInRange(playerPos)) {
          this.nearbyBuyZone = buyZone.zone;
          break;
        }
      }
    }
  }

  getNearbyBuyZone() {
    return this.nearbyBuyZone;
  }

  syncTurrets() {
    const stateTurrets = this.gameState.turrets || [];
    if (stateTurrets.length === 0) return;

    // Initialize turrets once
    if (!this.turretsInitialized) {
      for (const turret of stateTurrets) {
        if (!this.turretMeshes.has(turret.id)) {
          const turretObj = new Turret(this.scene.getScene(), turret);
          this.turretMeshes.set(turret.id, turretObj);
        }
      }
      this.turretsInitialized = true;
    }
  }

  updateTurrets() {
    const stateTurrets = this.gameState.turrets || [];

    // Update turret states
    for (const turretData of stateTurrets) {
      const turretObj = this.turretMeshes.get(turretData.id);
      if (turretObj) {
        turretObj.update(turretData, this.elapsedTime);
      }
    }

    // Check if player is near any turret
    this.nearbyTurret = null;
    const playerUnit = this.gameState.getMyPlayerUnit();
    if (!playerUnit || playerUnit.isRespawning) return;

    const playerPos = playerUnit.position;
    for (const turretData of stateTurrets) {
      const turretObj = this.turretMeshes.get(turretData.id);
      if (turretObj && turretObj.isInRange(playerPos)) {
        // Check if can be claimed
        if (turretObj.canBeClaimed(this.gameState.playerId)) {
          this.nearbyTurret = turretData;
          break;
        }
      }
    }
  }

  getNearbyTurret() {
    return this.nearbyTurret;
  }

  syncUnits() {
    const stateUnits = this.gameState.units || [];
    const currentUnitIDs = new Set(stateUnits.map(u => u.id));

    // Detect destroyed units and create explosions (but not for players who are respawning)
    for (const [unitID, unitObj] of this.unitMeshes.entries()) {
      if (!currentUnitIDs.has(unitID)) {
        // Unit was destroyed - create explosion at its position
        if (unitObj.mesh) {
          const position = unitObj.mesh.position.clone();
          this.createExplosion(position);
        }
        unitObj.remove();
        this.unitMeshes.delete(unitID);
      }
    }

    // Add new units and update player respawn state
    for (const unit of stateUnits) {
      if (!this.unitMeshes.has(unit.id)) {
        this.createUnit(unit);
      } else if (unit.type === 'player') {
        // Update respawning state for existing player units
        const playerObj = this.unitMeshes.get(unit.id);
        if (playerObj && playerObj.setRespawning) {
          playerObj.setRespawning(unit.isRespawning || false);
        }
      }
    }

    // Update previous IDs for next frame
    this.previousUnitIDs = currentUnitIDs;
  }

  syncProjectiles() {
    const stateProjectiles = this.gameState.projectiles || [];
    const stateProjectileIDs = new Set(stateProjectiles.map(p => p.id));

    // Remove projectiles that no longer exist
    for (const [projID, projObj] of this.projectileMeshes.entries()) {
      if (!stateProjectileIDs.has(projID)) {
        projObj.remove();
        this.projectileMeshes.delete(projID);
      }
    }

    // Add new projectiles and update existing
    for (const proj of stateProjectiles) {
      if (!this.projectileMeshes.has(proj.id)) {
        // Find shooter to get color - could be unit or turret
        let color = 0xffffff;
        const shooter = (this.gameState.units || []).find(u => u.id === proj.shooterId);
        if (shooter) {
          color = PLAYER_COLORS[shooter.ownerId];
        } else {
          // Check if shooter is a turret
          const turret = (this.gameState.turrets || []).find(t => t.id === proj.shooterId);
          if (turret && turret.ownerId !== -1) {
            color = PLAYER_COLORS[turret.ownerId];
          }
        }
        const projectile = new Projectile(this.scene.getScene(), proj, color);
        this.projectileMeshes.set(proj.id, projectile);
      } else {
        // Update position
        this.projectileMeshes.get(proj.id).updatePosition(proj.position);
      }
    }
  }

  createExplosion(position) {
    const explosion = new Explosion(this.scene.getScene(), position, () => {
      // Cleanup handled in updateExplosions
    });
    this.explosions.push(explosion);
  }

  updateExplosions() {
    // Update all explosions and remove completed ones
    this.explosions = this.explosions.filter(explosion => explosion.update());
  }

  createUnit(unit) {
    const color = PLAYER_COLORS[unit.ownerId];
    let unitObj;

    if (unit.type === 'tank') {
      unitObj = new Tank(this.scene.getScene(), unit, color);
    } else if (unit.type === 'airplane') {
      unitObj = new Airplane(this.scene.getScene(), unit, color);
    } else if (unit.type === 'player') {
      const isLocalPlayer = unit.ownerId === this.gameState.playerId;
      unitObj = new Player(this.scene.getScene(), unit, color, isLocalPlayer);

      // Set camera to follow local player
      if (isLocalPlayer && unitObj.mesh) {
        this.camera.setTarget(unitObj.mesh);
      }
    }

    if (unitObj) {
      this.unitMeshes.set(unit.id, unitObj);
    }
  }

  updateUnitPositions(deltaTime) {
    const stateUnits = this.gameState.units || [];
    const stateProjectiles = this.gameState.projectiles || [];

    // Build a map of shooter ID to target position (from projectiles)
    const shooterTargets = new Map();
    for (const proj of stateProjectiles) {
      // Use the projectile's end position as the target
      shooterTargets.set(proj.shooterId, proj.endPos);
    }

    for (const unit of stateUnits) {
      const unitObj = this.unitMeshes.get(unit.id);
      if (unitObj && unitObj.mesh) {
        // Interpolate position for smooth movement
        const targetPos = unit.position;
        const currentPos = unitObj.mesh.position;

        // Lerp factor for smooth interpolation
        const lerpFactor = Math.min(1, deltaTime * 10);

        unitObj.mesh.position.x += (targetPos.x - currentPos.x) * lerpFactor;
        unitObj.mesh.position.y += (targetPos.y - currentPos.y) * lerpFactor;
        unitObj.mesh.position.z += (targetPos.z - currentPos.z) * lerpFactor;

        // Update unit's internal state for rotation
        unitObj.updatePosition(unit.position);

        // Update turret targeting for tanks
        if (unit.type === 'tank' && unitObj.setTarget) {
          // Check if this tank has a projectile (is currently shooting)
          const targetPosition = shooterTargets.get(unit.id);
          if (targetPosition) {
            unitObj.setTarget(targetPosition);
          } else {
            // If not shooting, aim at nearest enemy
            const nearestEnemy = this.findNearestEnemy(unit);
            if (nearestEnemy) {
              unitObj.setTarget(nearestEnemy.position);
            }
          }
        }
      }
    }
  }

  // Find nearest enemy unit to a given unit
  findNearestEnemy(unit) {
    const stateUnits = this.gameState.units || [];
    let nearest = null;
    let nearestDist = Infinity;

    for (const other of stateUnits) {
      if (other.ownerId !== unit.ownerId) {
        const dx = other.position.x - unit.position.x;
        const dz = other.position.z - unit.position.z;
        const dist = Math.sqrt(dx * dx + dz * dz);
        if (dist < nearestDist) {
          nearestDist = dist;
          nearest = other;
        }
      }
    }

    return nearest;
  }

  render() {
    this.renderer.render(this.scene.getScene(), this.camera.getCamera());

    // Update popup position (needs to happen after render for correct projection)
    if (this.buyZonePopup) {
      this.buyZonePopup.update();
    }
  }

  setUnitMeshes(meshes) {
    this.unitMeshes = meshes;
  }
}
