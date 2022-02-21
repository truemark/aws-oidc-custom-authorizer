package main

import (
	"context"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"testing"
	"time"
)

// TODO: Make me a test runner, data driven....
func TestHandler(t *testing.T) {

	t.Run("Example", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		const jwksUrl = "url"
		ar := jwk.NewAutoRefresh(ctx)
		ar.Configure(jwksUrl, jwk.WithMinRefreshInterval(15*time.Minute))
		keySet, err := ar.Refresh(ctx, jwksUrl)
		if err != nil {
			t.Fatalf("failed to refresh JWKS: %s\n", err)
			return
		}
		fmt.Println("fetch JWKS")

		audience := ""
		key := "key"
		token, err := jwt.Parse(
			[]byte(key), // payload
			jwt.WithValidate(true),
			jwt.WithKeySet(keySet),
			jwt.WithAudience(audience))
		if err != nil {
			t.Fatalf("failed to parse payload: %s\n", err)
			return
		}

		// TODO: Validate Audience:
		fmt.Printf("aud: %s\n", token.Audience())
		_ = token

		fmt.Println("parsed jwt")
		cancel()
	})

}
