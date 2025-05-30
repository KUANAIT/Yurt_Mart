package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"api-gateway/internal/config/order"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrBlacklistedToken = errors.New("token has been blacklisted")
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTAuth struct {
	config     *order.AuthConfig
	authClient AuthServiceClient
}

func NewJWTAuth(config *order.AuthConfig, authClient AuthServiceClient) *JWTAuth {
	return &JWTAuth{
		config:     config,
		authClient: authClient,
	}
}

func (j *JWTAuth) GenerateToken(userID, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.config.TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.JWTSecret))
}

func (j *JWTAuth) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.config.JWTSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Check if token is blacklisted
		isBlacklisted, err := j.authClient.IsTokenBlacklisted(tokenString)
		if err != nil {
			return nil, err
		}
		if isBlacklisted {
			return nil, ErrBlacklistedToken
		}
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (j *JWTAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		claims, err := j.ValidateToken(parts[1])
		if err != nil {
			switch {
			case errors.Is(err, ErrExpiredToken):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has expired"})
			case errors.Is(err, ErrBlacklistedToken):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has been invalidated"})
			default:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			}
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func (j *JWTAuth) AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}

		c.Next()
	}
}
