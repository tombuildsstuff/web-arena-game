// AuthService - handles GitHub/BlueSky login/logout and user state

export class AuthService {
  constructor() {
    this.user = null; // { userId, displayName, avatarUrl, isGuest }
    this.listeners = [];
  }

  // Check if user is logged in by fetching /api/me
  // This will also create a persistent guest session if not authenticated
  async checkAuth() {
    try {
      const response = await fetch('/api/me', {
        credentials: 'include' // Include cookies
      });

      if (response.ok) {
        this.user = await response.json();
        this.notifyListeners();
        return this.user;
      } else {
        this.user = null;
        this.notifyListeners();
        return null;
      }
    } catch (error) {
      console.error('Auth check failed:', error);
      this.user = null;
      this.notifyListeners();
      return null;
    }
  }

  // Check if user is a guest
  isGuest() {
    return this.user !== null && this.user.isGuest;
  }

  // Redirect to GitHub OAuth login
  loginWithGitHub() {
    window.location.href = '/auth/github';
  }

  // Login with BlueSky using handle and app password
  async loginWithBlueSky(handle, appPassword) {
    try {
      const response = await fetch('/auth/bluesky', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ handle, appPassword })
      });

      const data = await response.json();

      if (response.ok && data.success) {
        this.user = data.user;
        this.notifyListeners();
        return { success: true, user: data.user };
      } else {
        return { success: false, error: data.error || 'Authentication failed' };
      }
    } catch (error) {
      console.error('BlueSky login failed:', error);
      return { success: false, error: 'Connection failed. Please try again.' };
    }
  }

  // Legacy method for backwards compatibility
  login() {
    this.loginWithGitHub();
  }

  // Logout and clear session
  async logout() {
    try {
      await fetch('/auth/logout', {
        method: 'POST',
        credentials: 'include'
      });
    } catch (error) {
      console.error('Logout failed:', error);
    }

    this.user = null;
    this.notifyListeners();
  }

  // Get current user (may be null)
  getUser() {
    return this.user;
  }

  // Check if user is authenticated (not a guest)
  isAuthenticated() {
    return this.user !== null && !this.user.isGuest;
  }

  // Subscribe to auth state changes
  onAuthChange(callback) {
    this.listeners.push(callback);
    // Return unsubscribe function
    return () => {
      this.listeners = this.listeners.filter(cb => cb !== callback);
    };
  }

  // Notify all listeners of auth state change
  notifyListeners() {
    for (const listener of this.listeners) {
      listener(this.user);
    }
  }

  // Check URL for auth callback
  handleAuthCallback() {
    const params = new URLSearchParams(window.location.search);
    if (params.get('auth') === 'success') {
      // Clean up URL
      window.history.replaceState({}, '', window.location.pathname);
      // Re-check auth state
      this.checkAuth();
    }
  }
}
