// Package validation @Author: youngalone [2023/8/1]
package middleware

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ValidationMiddleware validates user input to prevent injection attacks
func ValidationMiddleware(c *gin.Context) {
	path := c.Request.URL.Path

	// Skip validation for certain endpoints
	switch path {
	case "/douyin/user/register/", "/douyin/user/login/", "/douyin/feed/":
		c.Next()
		return
	}

	// Validate all query parameters
	for key, values := range c.Request.URL.Query() {
		for _, value := range values {
			if !isValidInput(value) {
				zap.L().Error("Invalid input detected",
					zap.String("path", path),
					zap.String("key", key),
					zap.String("value", value))
				c.JSON(400, gin.H{
					"status_code": 1,
					"status_msg":  "Invalid input",
				})
				c.Abort()
				return
			}
		}
	}

	// Validate form data
	if c.Request.Form != nil {
		for key, values := range c.Request.Form {
			for _, value := range values {
				if !isValidInput(value) {
					zap.L().Error("Invalid input detected",
						zap.String("path", path),
						zap.String("key", key),
						zap.String("value", value))
					c.JSON(400, gin.H{
						"status_code": 1,
						"status_msg":  "Invalid input",
					})
					c.Abort()
					return
				}
			}
		}
	}

	c.Next()
}

// isValidInput checks if input contains potentially dangerous characters
func isValidInput(input string) bool {
	// Check for SQL injection patterns
	sqlPatterns := []string{
		"(?i)(SELECT|INSERT|UPDATE|DELETE|DROP|CREATE|ALTER|TRUNCATE|EXEC|UNION|JOIN)",
		"(?i)(OR|AND|NOT|XOR)",
		"(?i)(1=1|1=0)",
		"(?i)(--)",
		"(?i)(#)",
		"(?i)(/*)",
	}

	for _, pattern := range sqlPatterns {
		matched, _ := regexp.MatchString(pattern, input)
		if matched {
			return false
		}
	}

	// Check for XSS patterns
	xssPatterns := []string{
		"<script.*?>.*?</script>",
		"javascript:",
		"onload=",
		"onclick=",
		"onerror=",
		"<iframe",
		"<object",
		"<embed",
		"<link",
	}

	for _, pattern := range xssPatterns {
		matched, _ := regexp.MatchString(pattern, input)
		if matched {
			return false
		}
	}

	// Check for path traversal
	pathTraversalPatterns := []string{
		"../",
	"..\\",
	}

	for _, pattern := range pathTraversalPatterns {
		if strings.Contains(input, pattern) {
			return false
		}
	}

	// Check for command injection
	cmdPatterns := []string{
		";",
		"|",
		"&",
		"$",
		"`",
		"(",
		")",
	}

	for _, pattern := range cmdPatterns {
		if len(input) > 100 && strings.Contains(input, pattern) {
			return false
		}
	}

	return true
}

// SanitizeInput sanitizes input by removing potentially harmful characters
func SanitizeInput(input string) string {
	// Replace SQL injection patterns
	input = regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE|DROP|CREATE|ALTER|TRUNCATE|EXEC|UNION|JOIN|OR|AND|NOT|XOR|--|#|/\*|\*/)`).ReplaceAllString(input, "")

	// Replace XSS patterns
	input = regexp.MustCompile(`<script.*?>.*?</script>|javascript:|onload=|onclick=|onerror=|<iframe|<object|<embed|<link`).ReplaceAllString(input, "")

	// Escape HTML entities
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#39;")

	return input
}

// ValidateUserID validates user ID format
func ValidateUserID(userID string) bool {
	_, err := strconv.ParseInt(userID, 10, 64)
	return err == nil && len(userID) <= 10
}

// ValidateUsername validates username format
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	// Only allow letters, numbers, and underscore
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	return matched
}

// ValidatePassword validates password strength
func ValidatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	// Check for at least one letter and one number
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasLetter && hasNumber
}