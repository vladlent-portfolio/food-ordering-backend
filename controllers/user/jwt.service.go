package user

import (
	"errors"
	"github.com/dgrijalva/jwt-go/v4"
)

const JWTSecret = "secret string"
const JWTUserIDKey = "uid"

type AuthClaims struct {
	UserID uint `json:"uid"`
}

type JWTService struct{}

func ProvideJWTService() *JWTService {
	return &JWTService{}
}

func (c AuthClaims) Valid(v *jwt.ValidationHelper) error {
	if c.UserID == 0 {
		return errors.New("invalid user id")
	}

	return nil
}

// Generate uses provided user id to create a token with AuthClaims and signs it with JWTSecret.
// Returns signed token.
func (s *JWTService) Generate(uid uint) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims{UserID: uid})
	ss, _ := token.SignedString([]byte(JWTSecret))
	return ss
}

// Parse reads provided token string and parses it with AuthClaims and JWTSecret.
func (s *JWTService) Parse(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
}

// AuthClaimsFromToken extracts AuthClaims from provided token.
func (s *JWTService) AuthClaimsFromToken(token string) (*AuthClaims, error) {
	t, err := s.Parse(token)

	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(*AuthClaims)

	if !ok {
		return nil, &jwt.InvalidClaimsError{}
	}

	return claims, claims.Valid(jwt.DefaultValidationHelper)
}
