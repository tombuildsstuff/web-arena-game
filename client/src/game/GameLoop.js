import { Tank } from './Tank.js';
import { Airplane } from './Airplane.js';
import { Player } from './Player.js';
import { Obstacle } from './Obstacle.js';
import { Projectile } from './Projectile.js';
import { Explosion } from './Explosion.js';
import { BuyZone } from './BuyZone.js';
import { Turret } from './Turret.js';
import { HealthPack } from './HealthPack.js';
import { PLAYER_COLORS } from '../utils/constants.js';

export class GameLoop {
  constructor(renderer, scene, camera, gameState, soundManager = null) {
    this.renderer = renderer;
    this.scene = scene;
    this.camera = camera;
    this.gameState = gameState;
    this.soundManager = soundManager;
    this.isRunning = false;
    this.lastTime = 0;
    this.elapsedTime = 0; // Total elapsed time for animations
    this.unitMeshes = new Map(); // Map of unit ID to unit object (Tank/Airplane)
    this.obstacleMeshes = new Map(); // Map of obstacle ID to Obstacle object
    this.projectileMeshes = new Map(); // Map of projectile ID to Projectile object
    this.buyZoneMeshes = new Map(); // Map of buy zone ID to BuyZone object
    this.turretMeshes = new Map(); // Map of turret ID to Turret object
    this.healthPackMeshes = new Map(); // Map of health pack ID to HealthPack object
    this.explosions = []; // Active explosions
    this.obstaclesInitialized = false;
    this.buyZonesInitialized = false;
    this.turretsInitialized = false;
    this.previousUnitIDs = new Set(); // Track unit IDs for death detection
    this.previousProjectileIDs = new Set(); // Track projectile IDs for shot detection
    this.previousTurretStates = new Map(); // Track turret states for sound triggers
    this.previousRespawnStates = new Map(); // Track respawn states for sound triggers
    this.nearbyBuyZone = null; // The buy zone the player is currently near
    this.nearbyClaimableBuyZone = null; // A claimable buy zone the player is near
    this.nearbyTurret = null; // The turret the player is currently near
    this.buyZonePopup = null; // Reference to the popup for position updates
    this.isSpectating = false; // Spectator mode flag
  }

  setBuyZonePopup(popup) {
    this.buyZonePopup = popup;
  }

  setSoundManager(soundManager) {
    this.soundManager = soundManager;
  }

  updateSoundListener() {
    if (!this.soundManager) return;

    // Get the player's position for the listener
    const playerUnit = this.gameState.getMyPlayerUnit();
    if (!playerUnit) return;

    const position = playerUnit.position;

    // Get camera direction for listener orientation
    const camera = this.camera.getCamera();
    const forward = { x: 0, y: 0, z: -1 };
    const up = { x: 0, y: 1, z: 0 };

    // Extract forward direction from camera
    if (camera && camera.getWorldDirection) {
      try {
        // getWorldDirection modifies the passed vector and returns it
        const dir = camera.getWorldDirection(camera.position.clone());
        forward.x = dir.x;
        forward.y = dir.y;
        forward.z = dir.z;
      } catch (e) {
        // Fallback to default forward direction
      }
    }

    this.soundManager.updateListener(position, forward, up);
  }

  setSpectatorMode(isSpectating) {
    this.isSpectating = isSpectating;
    // Enable free camera mode for spectators
    this.camera.setSpectatorMode(isSpectating);
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
    this.camera.update(deltaTime);

    // Update sound listener position
    this.updateSoundListener();

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

    // Sync and update health packs
    this.syncHealthPacks();

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
    const stateBuyZones = this.gameState.buyZones || [];
    if (stateBuyZones.length === 0) return;

    // Initialize buy zones once
    if (!this.buyZonesInitialized) {
      for (const zone of stateBuyZones) {
        if (!this.buyZoneMeshes.has(zone.id)) {
          const buyZone = new BuyZone(this.scene.getScene(), zone);
          this.buyZoneMeshes.set(zone.id, buyZone);
        }
      }
      this.buyZonesInitialized = true;
    }

    // Update buy zone ownership (in case zones were claimed)
    for (const zoneData of stateBuyZones) {
      const buyZone = this.buyZoneMeshes.get(zoneData.id);
      if (buyZone && buyZone.zone.ownerId !== zoneData.ownerId) {
        buyZone.updateOwner(zoneData.ownerId);
      }
    }
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
    this.nearbyClaimableBuyZone = null;
    const playerUnit = this.gameState.getMyPlayerUnit();
    if (!playerUnit || playerUnit.isRespawning) return;

    const playerPos = playerUnit.position;
    for (const [zoneId, buyZone] of this.buyZoneMeshes.entries()) {
      if (!buyZone.isInRange(playerPos)) continue;

      // Check if this is a zone owned by this player (can buy from it)
      // Skip zones with no unitType - those are claimable base zones, not purchase zones
      if (buyZone.zone.ownerId === this.gameState.playerId && buyZone.zone.unitType) {
        this.nearbyBuyZone = buyZone.zone;
        break;
      }

      // Check if this is a claimable neutral zone
      if (buyZone.zone.isClaimable && buyZone.zone.ownerId === -1) {
        this.nearbyClaimableBuyZone = buyZone.zone;
        // Continue checking - owned zones take priority
      }
    }
  }

  getNearbyBuyZone() {
    return this.nearbyBuyZone;
  }

  getNearbyClaimableBuyZone() {
    return this.nearbyClaimableBuyZone;
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
        // Initialize previous state tracking
        this.previousTurretStates.set(turret.id, {
          ownerId: turret.ownerId,
          isDestroyed: turret.isDestroyed || false
        });
      }
      this.turretsInitialized = true;
    }

    // Check for turret state changes to trigger sounds
    if (this.soundManager) {
      for (const turret of stateTurrets) {
        const prevState = this.previousTurretStates.get(turret.id);
        if (!prevState) continue;

        // Check for turret acquisition (ownership changed from -1 or different owner)
        if (turret.ownerId !== -1 && turret.ownerId !== prevState.ownerId) {
          this.soundManager.playPositional('turret_acquired', turret.position);
        }

        // Check for turret destruction
        if (turret.isDestroyed && !prevState.isDestroyed) {
          this.soundManager.playPositional('turret_destroyed', turret.position);
        }

        // Update previous state
        this.previousTurretStates.set(turret.id, {
          ownerId: turret.ownerId,
          isDestroyed: turret.isDestroyed || false
        });
      }
    }
  }

  updateTurrets() {
    const stateTurrets = this.gameState.turrets || [];
    const stateUnits = this.gameState.units || [];

    // Update turret states and pass units for target tracking
    for (const turretData of stateTurrets) {
      const turretObj = this.turretMeshes.get(turretData.id);
      if (turretObj) {
        turretObj.update(turretData, this.elapsedTime, stateUnits);
      }
    }

    // Check if player is near any turret
    this.nearbyTurret = null;
    const playerUnit = this.gameState.getMyPlayerUnit();
    if (!playerUnit || playerUnit.isRespawning) return;

    const playerPos = playerUnit.position;
    const playerId = this.gameState.playerId;

    for (const turretData of stateTurrets) {
      const turretObj = this.turretMeshes.get(turretData.id);
      if (!turretObj) continue;

      // Check if player is in range of turret
      if (!turretObj.isInRange(playerPos)) continue;

      // Use turretData directly for claim check to ensure we have latest state
      // Can't claim destroyed turrets
      if (turretData.isDestroyed) continue;
      // Can't claim your own turret
      if (turretData.ownerId === playerId) continue;
      // Can only claim neutral turrets - must destroy enemy turrets first
      if (turretData.ownerId !== -1) continue;

      // Can claim neutral turrets
      this.nearbyTurret = turretData;
      break;
    }
  }

  getNearbyTurret() {
    return this.nearbyTurret;
  }

  syncHealthPacks() {
    const stateHealthPacks = this.gameState.healthPacks || [];
    const stateHealthPackIDs = new Set(stateHealthPacks.map(hp => hp.id));

    // Remove health packs that no longer exist (collected)
    for (const [packID, packObj] of this.healthPackMeshes.entries()) {
      if (!stateHealthPackIDs.has(packID)) {
        packObj.remove();
        this.healthPackMeshes.delete(packID);
      }
    }

    // Add new health packs and update existing ones
    for (const pack of stateHealthPacks) {
      if (!this.healthPackMeshes.has(pack.id)) {
        const healthPack = new HealthPack(this.scene.getScene(), pack);
        this.healthPackMeshes.set(pack.id, healthPack);
      } else {
        // Update animation
        this.healthPackMeshes.get(pack.id).update(this.elapsedTime);
      }
    }
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

          // Play death sound at the unit's position
          if (this.soundManager) {
            this.soundManager.playPositional('player_death', position);
          }
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

        // Check for respawn sound trigger (was respawning, now not respawning)
        const wasRespawning = this.previousRespawnStates.get(unit.id);
        const isRespawning = unit.isRespawning || false;
        if (wasRespawning && !isRespawning && this.soundManager) {
          this.soundManager.playPositional('player_respawn', unit.position);
        }
        this.previousRespawnStates.set(unit.id, isRespawning);
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

        // Play shot sound at projectile's starting position (if new)
        if (!this.previousProjectileIDs.has(proj.id) && this.soundManager) {
          this.soundManager.playPositional('shot_fired', proj.position);
        }
      } else {
        // Update position
        this.projectileMeshes.get(proj.id).updatePosition(proj.position);
      }
    }

    // Update previous projectile IDs for next frame
    this.previousProjectileIDs = stateProjectileIDs;
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
      unitObj = new Tank(this.scene.getScene(), unit, color, false);
    } else if (unit.type === 'super_tank') {
      unitObj = new Tank(this.scene.getScene(), unit, color, true);
    } else if (unit.type === 'airplane') {
      unitObj = new Airplane(this.scene.getScene(), unit, color, false);
    } else if (unit.type === 'super_helicopter') {
      unitObj = new Airplane(this.scene.getScene(), unit, color, true);
    } else if (unit.type === 'player') {
      const isLocalPlayer = !this.isSpectating && unit.ownerId === this.gameState.playerId;
      // Get display name from player data
      const playerData = this.gameState.players[unit.ownerId];
      const displayName = playerData ? playerData.displayName : null;
      unitObj = new Player(this.scene.getScene(), unit, color, isLocalPlayer, displayName);

      // Set camera to follow local player (not in spectator mode - spectators use free camera)
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
    const camera = this.camera.getCamera();

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

        // Update health bar (show when damaged)
        if (unitObj.updateHealth && unit.health !== undefined) {
          unitObj.updateHealth(unit.health, camera);
        }

        // Update turret targeting for tanks (including super tanks)
        if ((unit.type === 'tank' || unit.type === 'super_tank') && unitObj.setTarget) {
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

  // Reset the game loop for a new game
  reset() {
    // Remove all unit meshes
    for (const unitObj of this.unitMeshes.values()) {
      unitObj.remove();
    }
    this.unitMeshes.clear();

    // Remove all obstacle meshes
    for (const obstacle of this.obstacleMeshes.values()) {
      obstacle.remove();
    }
    this.obstacleMeshes.clear();

    // Remove all projectile meshes
    for (const projectile of this.projectileMeshes.values()) {
      projectile.remove();
    }
    this.projectileMeshes.clear();

    // Remove all buy zone meshes
    for (const buyZone of this.buyZoneMeshes.values()) {
      buyZone.remove();
    }
    this.buyZoneMeshes.clear();

    // Remove all turret meshes
    for (const turret of this.turretMeshes.values()) {
      turret.remove();
    }
    this.turretMeshes.clear();

    // Remove all health pack meshes
    for (const healthPack of this.healthPackMeshes.values()) {
      healthPack.remove();
    }
    this.healthPackMeshes.clear();

    // Clear explosions
    this.explosions = [];

    // Reset initialization flags
    this.obstaclesInitialized = false;
    this.buyZonesInitialized = false;
    this.turretsInitialized = false;

    // Reset tracking
    this.previousUnitIDs.clear();
    this.previousProjectileIDs.clear();
    this.previousTurretStates.clear();
    this.previousRespawnStates.clear();
    this.nearbyBuyZone = null;
    this.nearbyClaimableBuyZone = null;
    this.nearbyTurret = null;
    this.elapsedTime = 0;

    // Reset spectator state
    this.isSpectating = false;
  }
}
