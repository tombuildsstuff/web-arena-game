package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// Config holds OAuth configuration
type Config struct {
	GitHubClientID     string
	GitHubClientSecret string
	BaseURL            string
}

// UserInfo represents authenticated user information
type UserInfo struct {
	UserID      string `json:"userId"`      // Persistent UUID for the user
	DisplayName string `json:"displayName"` // GitHub username or Guest_XXXX
	AvatarURL   string `json:"avatarUrl,omitempty"`
	IsGuest     bool   `json:"isGuest"`
}

// Handler handles authentication routes
type Handler struct {
	config      *Config
	oauthConfig *oauth2.Config

	// Server-side session storage: authToken -> UserInfo
	sessions   map[string]*UserInfo
	sessionsMu sync.RWMutex

	// GitHub ID to UserID mapping for consistent user IDs across logins
	githubToUserID   map[int64]string
	githubToUserIDMu sync.RWMutex

	// BlueSky DID to UserID mapping for consistent user IDs across logins
	blueskyToUserID   map[string]string
	blueskyToUserIDMu sync.RWMutex
}

// LoadConfig loads OAuth configuration from environment variables
func LoadConfig() *Config {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}

	return &Config{
		GitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		BaseURL:            baseURL,
	}
}

// NewHandler creates a new auth handler
func NewHandler(cfg *Config) *Handler {
	var oauthConfig *oauth2.Config

	// Only configure OAuth if credentials are provided
	if cfg.GitHubClientID != "" && cfg.GitHubClientSecret != "" {
		oauthConfig = &oauth2.Config{
			ClientID:     cfg.GitHubClientID,
			ClientSecret: cfg.GitHubClientSecret,
			Scopes:       []string{"read:user"},
			Endpoint:     github.Endpoint,
			RedirectURL:  cfg.BaseURL + "/auth/github/callback",
		}
	}

	return &Handler{
		config:          cfg,
		oauthConfig:     oauthConfig,
		sessions:        make(map[string]*UserInfo),
		githubToUserID:  make(map[int64]string),
		blueskyToUserID: make(map[string]string),
	}
}

// HandleLogin redirects to GitHub OAuth
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if h.oauthConfig == nil {
		http.Error(w, "GitHub OAuth not configured", http.StatusServiceUnavailable)
		return
	}

	// Generate state for CSRF protection
	state := generateRandomString(16)

	// Set state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		MaxAge:   300,   // 5 minutes
		SameSite: http.SameSiteLaxMode,
	})

	url := h.oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleCallback handles the OAuth callback from GitHub
func (h *Handler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	if h.oauthConfig == nil {
		http.Error(w, "GitHub OAuth not configured", http.StatusServiceUnavailable)
		return
	}

	// Verify state
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Exchange code for token
	code := r.URL.Query().Get("code")
	token, err := h.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Get user info from GitHub
	userInfo, githubID, err := h.getGitHubUser(token)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Get or create persistent UserID for this GitHub user
	userInfo.UserID = h.getOrCreateUserID(githubID)

	// Generate auth token and store session
	authToken := uuid.New().String()
	h.sessionsMu.Lock()
	h.sessions[authToken] = userInfo
	h.sessionsMu.Unlock()

	// Set auth token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    authToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		MaxAge:   604800, // 7 days
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect back to app
	http.Redirect(w, r, "/?auth=success", http.StatusTemporaryRedirect)
}

// HandleLogout clears the auth cookie and removes session
func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Remove session from server
	if cookie, err := r.Cookie("auth_token"); err == nil {
		h.sessionsMu.Lock()
		delete(h.sessions, cookie.Value)
		h.sessionsMu.Unlock()
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "auth_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

// HandleMe returns the current user info
func (h *Handler) HandleMe(w http.ResponseWriter, r *http.Request) {
	userInfo := h.GetUserFromRequest(r)
	if userInfo == nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "not authenticated"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

// GetUserFromRequest extracts user info from the request (auth token cookie)
func (h *Handler) GetUserFromRequest(r *http.Request) *UserInfo {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return nil
	}

	return h.GetUserFromToken(cookie.Value)
}

// GetUserFromToken looks up user info from auth token
func (h *Handler) GetUserFromToken(authToken string) *UserInfo {
	h.sessionsMu.RLock()
	defer h.sessionsMu.RUnlock()

	userInfo, exists := h.sessions[authToken]
	if !exists {
		return nil
	}

	return userInfo
}

// GenerateGuestName creates a random guest name
func GenerateGuestName() string {
	return "Guest_" + generateRandomString(4)
}

// GenerateGuestUser creates a new guest user with unique IDs
func GenerateGuestUser() *UserInfo {
	return &UserInfo{
		UserID:      uuid.New().String(),
		DisplayName: GenerateGuestName(),
		IsGuest:     true,
	}
}

// BlueSkyLoginRequest represents the request body for BlueSky login
type BlueSkyLoginRequest struct {
	Handle      string `json:"handle"`
	AppPassword string `json:"appPassword"`
}

// HandleBlueSkyLogin authenticates a user with BlueSky using their handle and app password
func (h *Handler) HandleBlueSkyLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req BlueSkyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Handle == "" || req.AppPassword == "" {
		http.Error(w, "Handle and app password are required", http.StatusBadRequest)
		return
	}

	// Clean up handle (remove @ if present)
	handle := strings.TrimPrefix(req.Handle, "@")

	// Authenticate with BlueSky
	userInfo, did, err := h.authenticateBluesky(handle, req.AppPassword)
	if err != nil {
		log.Printf("BlueSky auth failed for %s: %v", handle, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Get or create persistent UserID for this BlueSky user
	userInfo.UserID = h.getOrCreateBlueSkyUserID(did)

	// Generate auth token and store session
	authToken := uuid.New().String()
	h.sessionsMu.Lock()
	h.sessions[authToken] = userInfo
	h.sessionsMu.Unlock()

	// Set auth token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    authToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		MaxAge:   604800, // 7 days
		SameSite: http.SameSiteLaxMode,
	})

	log.Printf("BlueSky login successful for %s (DID: %s)", handle, did)

	// Return success with user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    userInfo,
	})
}

// authenticateBluesky authenticates with BlueSky using the AT Protocol
func (h *Handler) authenticateBluesky(handle, appPassword string) (*UserInfo, string, error) {
	// Resolve the handle to find the PDS
	pdsURL, err := h.resolveBlueSkyPDS(handle)
	if err != nil {
		return nil, "", fmt.Errorf("failed to resolve handle: %v", err)
	}

	// Create session with AT Protocol
	sessionURL := pdsURL + "/xrpc/com.atproto.server.createSession"

	reqBody, _ := json.Marshal(map[string]string{
		"identifier": handle,
		"password":   appPassword,
	})

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(sessionURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, "", fmt.Errorf("failed to connect to BlueSky: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		json.Unmarshal(body, &errResp)
		if errResp.Message != "" {
			return nil, "", fmt.Errorf("%s", errResp.Message)
		}
		return nil, "", fmt.Errorf("authentication failed")
	}

	var session struct {
		DID    string `json:"did"`
		Handle string `json:"handle"`
	}
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Get the user's profile for avatar
	avatarURL := h.getBlueSkyAvatar(pdsURL, session.DID)

	return &UserInfo{
		DisplayName: "@" + session.Handle,
		AvatarURL:   avatarURL,
		IsGuest:     false,
	}, session.DID, nil
}

// resolveBlueSkyPDS resolves a BlueSky handle to its PDS URL
func (h *Handler) resolveBlueSkyPDS(handle string) (string, error) {
	// For most users on bsky.social, we can use the main PDS
	// For custom domains, we'd need to do DNS/HTTP resolution

	// Try to resolve via the public API first
	client := &http.Client{Timeout: 10 * time.Second}

	// Check if handle ends with .bsky.social or similar known hosts
	if strings.HasSuffix(handle, ".bsky.social") ||
		!strings.Contains(handle, ".") {
		// Use the main bsky.social PDS
		return "https://bsky.social", nil
	}

	// For custom domains, try to resolve via DID document
	// First, try to get the DID from the public API
	resolveURL := fmt.Sprintf("https://bsky.social/xrpc/com.atproto.identity.resolveHandle?handle=%s", handle)
	resp, err := client.Get(resolveURL)
	if err != nil {
		// Fall back to main PDS
		return "https://bsky.social", nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result struct {
			DID string `json:"did"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err == nil && result.DID != "" {
			// For now, use main bsky.social as the PDS
			// In a full implementation, we'd resolve the DID document to find the actual PDS
			return "https://bsky.social", nil
		}
	}

	// Default to main PDS
	return "https://bsky.social", nil
}

// getBlueSkyAvatar fetches the user's avatar from BlueSky
func (h *Handler) getBlueSkyAvatar(pdsURL, did string) string {
	client := &http.Client{Timeout: 10 * time.Second}

	profileURL := fmt.Sprintf("%s/xrpc/app.bsky.actor.getProfile?actor=%s", pdsURL, did)
	resp, err := client.Get(profileURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var profile struct {
		Avatar string `json:"avatar"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return ""
	}

	return profile.Avatar
}

// getOrCreateBlueSkyUserID returns a consistent UserID for a BlueSky user
func (h *Handler) getOrCreateBlueSkyUserID(did string) string {
	h.blueskyToUserIDMu.Lock()
	defer h.blueskyToUserIDMu.Unlock()

	if userID, exists := h.blueskyToUserID[did]; exists {
		return userID
	}

	userID := uuid.New().String()
	h.blueskyToUserID[did] = userID
	return userID
}

// getGitHubUser fetches user info from GitHub API
func (h *Handler) getGitHubUser(token *oauth2.Token) (*UserInfo, int64, error) {
	client := h.oauthConfig.Client(context.Background(), token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var ghUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return nil, 0, err
	}

	return &UserInfo{
		DisplayName: ghUser.Login,
		AvatarURL:   ghUser.AvatarURL,
		IsGuest:     false,
	}, ghUser.ID, nil
}

// getOrCreateUserID returns a consistent UserID for a GitHub user
func (h *Handler) getOrCreateUserID(githubID int64) string {
	h.githubToUserIDMu.Lock()
	defer h.githubToUserIDMu.Unlock()

	if userID, exists := h.githubToUserID[githubID]; exists {
		return userID
	}

	userID := uuid.New().String()
	h.githubToUserID[githubID] = userID
	return userID
}

// generateRandomString creates a random hex string
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}
