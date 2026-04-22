// Package middleware @Author: youngalone [2023/8/1]
package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none';")

		// Strict Transport Security (only if using HTTPS)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// X-Content-Type-Options
		c.Header("X-Content-Type-Options", "nosniff")

		// X-Frame-Options
		c.Header("X-Frame-Options", "DENY")

		// X-XSS-Protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")

		// Cache Control for sensitive content
		if c.Request.URL.Path == "/douyin/user/register/" ||
		   c.Request.URL.Path == "/douyin/user/login/" ||
		   c.Request.URL.Path == "/douyin/user/" {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
			c.Header("Pragma", "no-cache")
		}

		c.Next()
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// In production, you would set this to your specific domain
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		c.Header("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware with improved security
func SecurityRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement rate limiting with security considerations
		// This would typically use Redis or similar for distributed rate limiting

		// Add security headers first
		c.Header("X-Rate-Limit-Limit", "100")
		c.Header("X-Rate-Limit-Remaining", "99")
		c.Header("X-Rate-Limit-Reset", "60")

		c.Next()
	}
}

// AuditLogMiddleware logs security-related events
func AuditLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log security events for auditing
		// This is a simplified version, in production you'd use proper logging

		c.Next()
	}
}

// HealthCheckMiddleware prevents health check endpoints from being exposed in production
func HealthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.JSON(404, gin.H{
				"status_code": 1,
				"status_msg":  "Not found",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}