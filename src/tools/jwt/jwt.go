package jwt

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

var (
	tokenEncodeString   = []byte(os.Getenv("TOKEN_HASH"))
	accessTokenDuration = time.Minute * 5
)

// Claims  claims
type Claims struct {
	UserID        int    `json:"user_id"`
	Login         string `json:"login"`
	RefreshToken  string `json:"refresh_token"`
	AccessExpired int64  `json:"access_expired"`
	jwt.StandardClaims
}

// NewToken создание нового токена
func NewToken(userID int, login, refreshToken string) (string, error) {

	// Create the Claims
	claims := Claims{
		userID,
		login,
		refreshToken,
		time.Now().Add(accessTokenDuration).Unix(),
		jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(tokenEncodeString)
}

// NewExpiredToken создание нового токена просроченного (для теста)
func NewExpiredToken(userID int, login, refreshToken string) (string, error) {

	// Create the Claims
	claims := Claims{
		userID,
		login,
		refreshToken,
		time.Now().Unix(),
		jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(tokenEncodeString)
}

// ParseToken parse token
func ParseToken(unparsedToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(unparsedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return tokenEncodeString, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid

}

// TokenIsExpired время допуска без проверки истек
func (c Claims) TokenIsExpired() bool {
	return time.Unix(c.AccessExpired, 0).Before(time.Now())
}

// GenerateRefreshToken генерация refreshtoken в uuid
func GenerateRefreshToken() string {
	return uuid.NewV4().String()
}
