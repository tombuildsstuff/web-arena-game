export class SoundManager {
  constructor() {
    this.audioContext = null;
    this.masterGain = null;
    this.buffers = new Map();
    this.isMuted = false;
    this.isInitialized = false;
    this.listenerPosition = { x: 0, y: 0, z: 0 };

    // Sound configuration (supports .mp3, .wav, .ogg, .aac)
    this.sounds = {
      player_death: { url: '/audio/player_death.wav', volume: 0.7 },
      shot_fired: { url: '/audio/shot_fired.wav', volume: 0.4 },
      player_respawn: { url: '/audio/player_respawn.ogg', volume: 0.6 },
      turret_acquired: { url: '/audio/turret_acquired.wav', volume: 0.6 },
      turret_destroyed: { url: '/audio/turret_destroyed.ogg', volume: 0.6 },
      match_start: { url: '/audio/match_start.wav', volume: 0.8 },
      match_end: { url: '/audio/match_end.wav', volume: 0.8 }
    };

    // Spatial audio configuration
    this.maxDistance = 150; // Maximum hearing distance
    this.refDistance = 10;  // Distance at which volume is at full
    this.rolloffFactor = 1; // How quickly sound fades with distance
  }

  async init() {
    if (this.isInitialized) return;

    try {
      // Create AudioContext (needs user interaction on some browsers)
      this.audioContext = new (window.AudioContext || window.webkitAudioContext)();

      // Create master gain node
      this.masterGain = this.audioContext.createGain();
      this.masterGain.connect(this.audioContext.destination);

      // Load all sound buffers
      await this.loadAllSounds();

      this.isInitialized = true;
      console.log('SoundManager initialized');
    } catch (error) {
      console.error('Failed to initialize SoundManager:', error);
    }
  }

  async loadAllSounds() {
    const loadPromises = Object.entries(this.sounds).map(async ([name, config]) => {
      try {
        const response = await fetch(config.url);
        if (!response.ok) {
          console.warn(`Sound file not found: ${config.url}`);
          return;
        }
        const arrayBuffer = await response.arrayBuffer();
        const audioBuffer = await this.audioContext.decodeAudioData(arrayBuffer);
        this.buffers.set(name, audioBuffer);
        console.log(`Loaded sound: ${name}`);
      } catch (error) {
        console.warn(`Failed to load sound ${name}:`, error);
      }
    });

    await Promise.all(loadPromises);
  }

  // Resume AudioContext (required after user interaction on some browsers)
  async resume() {
    if (this.audioContext && this.audioContext.state === 'suspended') {
      await this.audioContext.resume();
    }
  }

  // Update listener position (call every frame with player/camera position)
  updateListener(position, forward = { x: 0, y: 0, z: -1 }, up = { x: 0, y: 1, z: 0 }) {
    if (!this.audioContext) return;

    this.listenerPosition = { ...position };

    const listener = this.audioContext.listener;

    // Set listener position
    if (listener.positionX) {
      // Modern API
      listener.positionX.setValueAtTime(position.x, this.audioContext.currentTime);
      listener.positionY.setValueAtTime(position.y, this.audioContext.currentTime);
      listener.positionZ.setValueAtTime(position.z, this.audioContext.currentTime);

      listener.forwardX.setValueAtTime(forward.x, this.audioContext.currentTime);
      listener.forwardY.setValueAtTime(forward.y, this.audioContext.currentTime);
      listener.forwardZ.setValueAtTime(forward.z, this.audioContext.currentTime);

      listener.upX.setValueAtTime(up.x, this.audioContext.currentTime);
      listener.upY.setValueAtTime(up.y, this.audioContext.currentTime);
      listener.upZ.setValueAtTime(up.z, this.audioContext.currentTime);
    } else {
      // Legacy API
      listener.setPosition(position.x, position.y, position.z);
      listener.setOrientation(forward.x, forward.y, forward.z, up.x, up.y, up.z);
    }
  }

  // Play a sound at a specific position in 3D space
  playPositional(soundName, position) {
    if (!this.isInitialized || this.isMuted) return;

    const buffer = this.buffers.get(soundName);
    if (!buffer) {
      console.warn(`Sound not loaded: ${soundName}`);
      return;
    }

    const config = this.sounds[soundName];

    // Create source
    const source = this.audioContext.createBufferSource();
    source.buffer = buffer;

    // Create panner for 3D positioning
    const panner = this.audioContext.createPanner();
    panner.panningModel = 'HRTF';
    panner.distanceModel = 'inverse';
    panner.refDistance = this.refDistance;
    panner.maxDistance = this.maxDistance;
    panner.rolloffFactor = this.rolloffFactor;
    panner.coneInnerAngle = 360;
    panner.coneOuterAngle = 360;
    panner.coneOuterGain = 1;

    // Set position
    if (panner.positionX) {
      panner.positionX.setValueAtTime(position.x, this.audioContext.currentTime);
      panner.positionY.setValueAtTime(position.y, this.audioContext.currentTime);
      panner.positionZ.setValueAtTime(position.z, this.audioContext.currentTime);
    } else {
      panner.setPosition(position.x, position.y, position.z);
    }

    // Create gain for volume control
    const gainNode = this.audioContext.createGain();
    gainNode.gain.value = config.volume;

    // Connect: source -> panner -> gain -> master
    source.connect(panner);
    panner.connect(gainNode);
    gainNode.connect(this.masterGain);

    source.start();

    // Cleanup when done
    source.onended = () => {
      source.disconnect();
      panner.disconnect();
      gainNode.disconnect();
    };
  }

  // Play a non-positional sound (for UI/global events)
  playGlobal(soundName) {
    if (!this.isInitialized || this.isMuted) return;

    const buffer = this.buffers.get(soundName);
    if (!buffer) {
      console.warn(`Sound not loaded: ${soundName}`);
      return;
    }

    const config = this.sounds[soundName];

    // Create source
    const source = this.audioContext.createBufferSource();
    source.buffer = buffer;

    // Create gain for volume
    const gainNode = this.audioContext.createGain();
    gainNode.gain.value = config.volume;

    // Connect directly to master (no panner)
    source.connect(gainNode);
    gainNode.connect(this.masterGain);

    source.start();

    // Cleanup when done
    source.onended = () => {
      source.disconnect();
      gainNode.disconnect();
    };
  }

  // Mute/unmute
  setMuted(muted) {
    this.isMuted = muted;
    if (this.masterGain) {
      this.masterGain.gain.value = muted ? 0 : 1;
    }
  }

  toggleMute() {
    this.setMuted(!this.isMuted);
    return this.isMuted;
  }

  getMuted() {
    return this.isMuted;
  }
}
