package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func makeJWT(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return s
}

func TestUserMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("options preflight", func(t *testing.T) {
		r := gin.New()
		r.Use(UserMiddleware())
		r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

		req := httptest.NewRequest(http.MethodOptions, "/x", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", w.Code)
		}
	})

	t.Run("internal secret bypass", func(t *testing.T) {
		t.Setenv("SECRET_X_INTERNAL_SECRET", "internal-123")
		r := gin.New()
		r.Use(UserMiddleware())
		r.GET("/x", func(c *gin.Context) {
			v, _ := c.Get("userID")
			c.JSON(http.StatusOK, gin.H{"userID": v})
		})

		req := httptest.NewRequest(http.MethodGet, "/x", nil)
		req.Header.Set("X-Internal-Secret", "internal-123")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("missing token", func(t *testing.T) {
		r := gin.New()
		r.Use(UserMiddleware())
		r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("valid id claim", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "secret")
		r := gin.New()
		r.Use(UserMiddleware())
		r.GET("/x", func(c *gin.Context) {
			v, _ := c.Get("user_id")
			c.JSON(http.StatusOK, gin.H{"user_id": v})
		})

		token := makeJWT(t, "secret", jwt.MapClaims{"id": "u1"})
		req := httptest.NewRequest(http.MethodGet, "/x", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var body map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &body)
		if body["user_id"] != "u1" {
			t.Fatalf("expected user_id u1, got %v", body)
		}
	})

	t.Run("missing id and sub", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "secret")
		r := gin.New()
		r.Use(UserMiddleware())
		r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })
		token := makeJWT(t, "secret", jwt.MapClaims{"role": "x"})
		req := httptest.NewRequest(http.MethodGet, "/x", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})
}
