// Package middleware @Author: youngalone [2023/8/1]
package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CSRFTokenLength is the length of CSRF token
const CSRFTokenLength = 32

// CSRFStore stores CSRF tokens
type CSRFStore struct {
	tokens map[string]string
}

var csrfStore = &CSRFStore{
	tokens: make(map[string]string),
}

// CSRFMiddleware provides CSRF protection
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		// Skip CSRF check for safe methods and certain paths
		if method == "GET" || method == "HEAD" || method == "OPTIONS" ||
		   path == "/douyin/feed/" || path == "/douyin/user/register/" ||
		   path == "/douyin/user/login/" || path == "/douyin/user/" {
			c.Next()
			return
		}

		// Extract CSRF token from various headers
		csrfToken := extractCSRFToken(c)
		if csrfToken == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"status_code": 1,
				"status_msg":  "Missing CSRF token",
			})
			c.Abort()
			return
		}

		// Verify CSRF token
		sessionID := getSessionID(c)
		if !validateCSRFToken(sessionID, csrfToken) {
			zap.L().Error("CSRF token validation failed",
				zap.String("path", path),
				zap.String("method", method),
				zap.String("session_id", sessionID))
			c.JSON(http.StatusForbidden, gin.H{
				"status_code": 1,
				"status_msg":  "Invalid CSRF token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateCSRFToken generates a new CSRF token
func GenerateCSRFToken(sessionID string) string {
	token := generateRandomToken(CSRFTokenLength)
	csrfStore.tokens[sessionID] = token
	return token
}

// GetCSRFToken returns CSRF token for a session
func GetCSRFToken(sessionID string) string {
	return csrfStore.tokens[sessionID]
}

// ValidateCSRFToken validates a CSRF token
func ValidateCSRFToken(sessionID, token string) bool {
	storedToken, exists := csrfStore.tokens[sessionID]
	if !exists {
		return false
	}
	return storedToken == token
}

// validateCSRFToken validates CSRF token
func validateCSRFToken(sessionID, token string) bool {
	return ValidateCSRFToken(sessionID, token)
}

// extractCSRFToken extracts CSRF token from request
func extractCSRFToken(c *gin.Context) string {
	// Check Authorization header first (for mobile apps)
	auth := c.GetHeader("Authorization")
	if auth != "" {
		return auth // For simplicity, in real app you'd parse this properly
	}

	// Check X-CSRF-Token header
	csrfToken := c.GetHeader("X-CSRF-Token")
	if csrfToken != "" {
		return csrfToken
	}

	// Check query parameter
	csrfToken = c.Query("csrf_token")
	if csrfToken != "" {
		return csrfToken
	}

	// Check form value
	csrfToken = c.PostForm("csrf_token")
	if csrfToken != "" {
		return csrfToken
	}

	return ""
}

// generateRandomToken generates a random token
func generateRandomToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// getSessionID gets session ID from context
func getSessionID(c *gin.Context) string {
	// In a real application, this would extract from session cookie
	// For this demo, we'll use a simple approach
	userID := c.GetString("user_id")
	if userID == "" {
		return "anonymous"
	}
	return userID
}

// SetCSRFTokenHeader sets CSRF token in response headers
func SetCSRFTokenHeader(c *gin.Context, sessionID string) {
	token := GetCSRFToken(sessionID)
	if token == "" {
		token = GenerateCSRFToken(sessionID)
	}

	// Set CSRF token in response headers
	c.Header("X-CSRF-Token", token)

	// Set CSRF token in cookie (for web applications)
	c.SetCookie("csrf_token", token, 3600, "/", "", false, true)
}

// ValidateOrigin validates the origin header to prevent CSRF
func ValidateOrigin(c *gin.Context, allowedOrigins []string) bool {
	origin := c.GetHeader("Origin")
	if origin == "" {
		return true // No origin header for same-origin requests
	}

	parsedOrigin, err := url.Parse(origin)
	if err != nil {
		return false
	}

	for _, allowed := range allowedOrigins {
		if parsedOrigin.Host == allowed {
			return true
		}
	}

	return false
}

// CSRFWrapper wraps routes with CSRF protection
func CSRFWrapper(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !ValidateOrigin(c, allowedOrigins) {
			c.JSON(http.StatusForbidden, gin.H{
				"status_code": 1,
				"status_msg":  "Invalid origin",
			})
			c.Abort()
			return
		}

		// Generate and set CSRF token
		sessionID := getSessionID(c)
		SetCSRFTokenHeader(c, sessionID)

		c.Next()
	}
}