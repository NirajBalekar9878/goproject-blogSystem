package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cms-blog-api/middleware"
	"cms-blog-api/utils"

	"github.com/gin-gonic/gin"
)

func TestJWTGenerationAndValidation(t *testing.T) {
	userID := uint(42)
	username := "testuser"
	role := "admin"

	// 1. Generate Token
	tokenString, err := utils.GenerateToken(userID, username, role)
	if err != nil {
		t.Fatalf("Expected no error generating token, got: %v", err)
	}
	if tokenString == "" {
		t.Fatal("Expected token string to be non-empty")
	}

	// 2. Validate Token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("Expected token to be valid, got error: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}
	if claims.Username != username {
		t.Errorf("Expected Username %s, got %s", username, claims.Username)
	}
	if claims.Role != role {
		t.Errorf("Expected Role %s, got %s", role, claims.Role)
	}

	// 3. Validate Tampered/Invalid Token
	invalidToken := tokenString + "tampered"
	_, err = utils.ValidateToken(invalidToken)
	if err == nil {
		t.Error("Expected error for tampered token, got nil")
	}
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setupRouter := func() *gin.Engine {
		r := gin.New()
		protected := r.Group("/protected")
		protected.Use(middleware.AuthMiddleware())
		protected.GET("/test", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			username, _ := c.Get("username")
			role, _ := c.Get("role")
			c.JSON(http.StatusOK, gin.H{
				"user_id":  userID,
				"username": username,
				"role":     role,
			})
		})
		return r
	}

	r := setupRouter()

	t.Run("Missing Authorization Header returns 401", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected/test", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized, got %d", resp.Code)
		}
	})

	t.Run("Malformed Authorization Header returns 401", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected/test", nil)
		req.Header.Set("Authorization", "Basic sometoken")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized, got %d", resp.Code)
		}
	})

	t.Run("Invalid JWT Token returns 401", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.jwt.token")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized, got %d", resp.Code)
		}
	})

	t.Run("Valid JWT Token returns 200 and injects claims", func(t *testing.T) {
		validToken, err := utils.GenerateToken(100, "niraj", "admin")
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		req, _ := http.NewRequest(http.MethodGet, "/protected/test", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected 200 OK, got %d: %s", resp.Code, resp.Body.String())
		}
	})
}
