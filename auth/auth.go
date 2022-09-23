package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Option struct {
	JwtSecret       string
	DefaultDuration *time.Duration
}

type claimKey struct {
	Key string
}

type TokenClaims struct {
	UID    int64
	Sub    string
	Iat    int64
	Exp    int64
	Claims map[string]interface{}
}

type CreateTokenOption struct {
	UID      int64
	Sub      string
	Claims   map[string]string
	ExpireIn *time.Duration
}

type JWTToken struct {
	AccessToken string
	Iat         int64
	Exp         int64
}

const AnonymousUserID = "Anonymous"

type TokenValidator = func(context.Context, *TokenClaims) error

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrTokenInvalid            = errors.New("invalid jwt token")

	User = claimKey{
		Key: "auth.user",
	}
	Anonymous = TokenClaims{Sub: AnonymousUserID}
)

var (
	jwtSecret       []byte
	defaultDuration time.Duration
	tokenValidators []TokenValidator
)

func Init(op Option) error {
	if len(op.JwtSecret) == 0 {
		return errors.New("jwt secret is required")
	}
	jwtSecret = []byte(op.JwtSecret)
	if op.DefaultDuration != nil && *op.DefaultDuration > 0 {
		defaultDuration = *op.DefaultDuration
	} else {
		defaultDuration = time.Hour * 24 * 365 // 1 year
	}
	return nil
}

func CreateJWT(op *CreateTokenOption) (token *JWTToken, err error) {
	t := jwt.New(jwt.SigningMethodHS256)
	clms := t.Claims.(jwt.MapClaims)
	now := time.Now()
	iat := now.Unix()
	var exp int64
	if op.ExpireIn != nil && *op.ExpireIn > 0 {
		exp = now.Add(*op.ExpireIn).Unix()
	} else {
		exp = now.Add(defaultDuration).Unix()
	}
	for k, v := range op.Claims {
		clms[k] = v
	}
	clms["uid"] = fmt.Sprint(op.UID)
	clms["sub"] = op.Sub
	clms["iat"] = iat
	clms["exp"] = exp

	jwtToken, err := t.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}
	return &JWTToken{
		AccessToken: jwtToken,
		Iat:         iat,
		Exp:         exp,
	}, nil
}

func ParseJWT(jwtToken string) (u *TokenClaims, err error) {
	token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrUnexpectedSigningMethod
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrTokenInvalid
	}
	claims := token.Claims.(jwt.MapClaims)
	u = &TokenClaims{}
	u.UID, _ = strconv.ParseInt(claims["uid"].(string), 10, 64)
	u.Claims = claims
	u.Exp, _ = claims["exp"].(int64)
	u.Iat, _ = claims["iat"].(int64)

	return u, nil
}

func VerifyJWTToken(ctx context.Context) (context.Context, error) {
	claims, err := verifyJwtToken(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, User, claims), nil
}

func VerifyJWTTokenOptional(ctx context.Context) (context.Context, error) {
	claims, err := verifyJwtToken(ctx)
	if err != nil {
		return context.WithValue(ctx, User, Anonymous), nil
	}
	return context.WithValue(ctx, User, claims), nil
}

func verifyJwtToken(ctx context.Context) (*TokenClaims, error) {
	jwtToken, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token required: %v", err)
	}
	claims, err := ParseJWT(jwtToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	if len(tokenValidators) > 0 {
		for _, validator := range tokenValidators {
			err = validator(ctx, claims)
			if err != nil {
				return nil, err
			}
		}
	}

	return claims, nil
}
