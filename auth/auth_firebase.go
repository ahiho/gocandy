package auth

import (
	"context"
	"fmt"

	firebaseApp "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	firAuth *firebaseAuth.Client
)

type FirebaseClaims struct {
	Token *firebaseAuth.Token
	Sub   string
	Email string
	Phone string
}

func InitFirebase(ctx context.Context, opt option.ClientOption) error {
	firApp, err := firebaseApp.NewApp(ctx, nil, opt)
	if err != nil {
		return err
	}
	firAuth, err = firApp.Auth(ctx)
	if err != nil {
		return err
	}
	return nil
}

func ValidateFirebaseToken(token string) (*FirebaseClaims, error) {
	ftoken, err := firAuth.VerifyIDToken(context.Background(), token)
	if err != nil {
		return nil, err
	}
	fc := &FirebaseClaims{
		Token: ftoken,
		Sub:   fmt.Sprintf("firebase|%v", ftoken.UID),
	}

	claims := ftoken.Claims

	if email, ok := claims["email"]; ok {
		fc.Email = email.(string)
	}

	if phone, ok := claims["phone_number"]; ok {
		fc.Phone = phone.(string)
	}

	return fc, nil
}
