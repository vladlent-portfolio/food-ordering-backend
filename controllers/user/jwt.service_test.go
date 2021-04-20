package user

import (
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestAuthClaims_Valid(t *testing.T) {
	t.Run("should return error if UserID == 0", func(t *testing.T) {
		it := assert.New(t)
		claims := AuthClaims{UserID: 0}

		it.Error(claims.Valid(jwt.DefaultValidationHelper))
	})
}

func TestJWTService_Generate(t *testing.T) {
	service := &JWTService{}
	t.Run("should generate a valid token string with encoded AuthClaims", func(t *testing.T) {
		it := assert.New(t)

		for _, uid := range generateUserIDs(100) {
			if uid == 0 {
				uid++
			}
			ss := service.Generate(uid)
			it.NotZero(ss)

			token, err := jwt.ParseWithClaims(ss, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(JWTSecret), nil
			})

			if it.NoError(err) {
				validateToken(t, token, uid)
			}
		}

	})
}

func TestJWTService_Parse(t *testing.T) {
	service := &JWTService{}
	t.Run("should parse the token", func(t *testing.T) {
		it := assert.New(t)

		for _, id := range generateUserIDs(100) {
			ss := service.Generate(id)
			token, err := service.Parse(ss)

			if it.NoError(err) {
				validateToken(t, token, id)
			}
		}
	})
}

func TestJWTService_AuthClaimsFromToken(t *testing.T) {
	service := &JWTService{}

	t.Run("should extract AuthClaims from provided token", func(t *testing.T) {
		it := assert.New(t)

		for _, id := range generateUserIDs(100) {
			ss := service.Generate(id)
			claims, err := service.AuthClaimsFromToken(ss)

			if it.NoError(err) {
				it.Equal(id, claims.UserID)
			}
		}
	})
}

func validateToken(t *testing.T, token *jwt.Token, uid uint) {
	it := assert.New(t)
	claims, ok := token.Claims.(*AuthClaims)

	if it.True(ok) {
		it.NoError(claims.Valid(jwt.DefaultValidationHelper))
		it.Equal(uid, claims.UserID)
	}
}

func generateUserIDs(amount int) []uint {
	ids := make([]uint, amount)
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range ids {
		ids[i] = uint(generator.Uint64() + 1)
	}

	return ids
}
