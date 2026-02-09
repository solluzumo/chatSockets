package middleware

import (
	"context"
	"crypto/rsa"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

type JWTMiddleware struct {
	publicKey *rsa.PublicKey
	issuer    string
	mLogger   *zap.Logger
}

func NewJWTMiddleware(publicKey *rsa.PublicKey, issuer string, appLogger *zap.Logger) *JWTMiddleware {
	return &JWTMiddleware{
		publicKey: publicKey,
		issuer:    issuer,
		mLogger:   appLogger.Named("auth_middleware"),
	}
}

func (m *JWTMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			m.mLogger.Warn("нет authorization заголовка")
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			m.mLogger.Warn("не правильный authorization заголовок")

			return
		}

		tokenStr := parts[1]

		claims := jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(
			tokenStr,
			&claims,
			func(token *jwt.Token) (any, error) {
				return m.publicKey, nil
			},
			jwt.WithIssuer(m.issuer),
		)
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			m.mLogger.Warn("токен не верный")

			return
		}

		userID := claims.Subject
		if userID == "" {
			http.Error(w, "missing subject", http.StatusUnauthorized)
			m.mLogger.Warn("userID отсутствует в токене")

			return
		}

		uid, err := strconv.Atoi(claims.Subject)
		if err != nil {
			http.Error(w, "invalid subject", http.StatusUnauthorized)
			m.mLogger.Warn("userID не int значение")

			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, uid)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (m *JWTMiddleware) HandleWS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.URL.Query().Get("token")
		if tokenStr == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		claims := jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(
			tokenStr,
			&claims,
			func(token *jwt.Token) (any, error) {
				return m.publicKey, nil
			},
			jwt.WithIssuer(m.issuer),
		)
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		uid, _ := strconv.Atoi(claims.Subject)
		ctx := context.WithValue(r.Context(), userIDKey, uid)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}
