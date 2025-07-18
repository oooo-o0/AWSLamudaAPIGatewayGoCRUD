package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey []byte
	ttl       time.Duration
}

type customClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// NewJWTService はJWT認証サービスを初期化
func NewJWTService(secret string, ttlMinutes int) AuthService {
	return &jwtService{
		secretKey: []byte(secret),
		ttl:       time.Duration(ttlMinutes) * time.Minute,
	}
}

// Login はユーザー認証成功時にJWTトークンを生成
func (j *jwtService) Login(ctx context.Context, username, password string) (string, error) {
	// ここは本来DBでユーザーパスワード検証ロジックを呼ぶが省略
	// 例: successなら
	userID := username        // ユーザーIDはusernameで代用
	roles := []string{"user"} // 例としてuserロールを付与

	claims := customClaims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "kaigo-insurance-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// VerifyToken はJWTトークンの検証を行いユーザーIDを返す
func (j *jwtService) VerifyToken(ctx context.Context, tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*customClaims)
	if !ok || !token.Valid {
		return "", ErrUnauthorized
	}
	return claims.UserID, nil
}

// RefreshToken はトークンのリフレッシュ（実装省略、必要に応じ追加）
func (j *jwtService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	return "", errors.New("not implemented")
}

// GetRoles はユーザーのロールを返す（JWTに含めているので本実装では空でOK）
func (j *jwtService) GetRoles(ctx context.Context, userID string) ([]string, error) {
	return []string{"user"}, nil
}
