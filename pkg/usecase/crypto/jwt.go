package crypto

import (
	"context"

	"fin-asg/pkg/domain/claim"
	"github.com/kataras/jwt"
)

var (
	sharedKey = []byte("sercrethatmaycontainch@r$32chars")
)

func CreateJWT(ctx context.Context, claim any) (string, error) {

	token, err := jwt.Sign(jwt.HS256, sharedKey, claim)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func VerifyJWT(ctx context.Context, token string) (claims claim.Access, err error) {

	verifiedToken, err := jwt.Verify(jwt.HS256, sharedKey, []byte(token))
	if err != nil {
		return claims, err
	}

	err = verifiedToken.Claims(&claims)
	if err != nil {
		return claims, err
	}
	return claims, err
}
