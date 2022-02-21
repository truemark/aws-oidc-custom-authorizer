package main

import (
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

func VerifyToken(tokenBuf string, keySet jwk.Set) (bool, string, error) {
	token, err := jwt.ParseString(
		tokenBuf,
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
		jwt.UseDefaultKey(true),
		jwt.InferAlgorithmFromKey(true))
	LogToken(token)
	fmt.Println(token)
	if err != nil {
		fmt.Printf("token-verification: failed to parse payload: %s\n", err)
		return false, "", err
	}

	return true, "{KID-VALUE-HERE}", nil
}
