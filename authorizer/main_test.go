package main

import (
	"testing"
)

func TestHandler(t *testing.T) {

	//t.Run("Example", func(t *testing.T) {
	//	ctx, cancel := context.WithCancel(context.Background())
	//	const jwksUrl = "<<URL>>"
	//	ar := jwk.NewAutoRefresh(ctx)
	//	ar.Configure(jwksUrl, jwk.WithMinRefreshInterval(15*time.Minute))
	//	keySet, err := ar.Refresh(ctx, jwksUrl)
	//	if err != nil {
	//		t.Fatalf("failed to refresh JWKS: %s\n", err)
	//		return
	//	}
	//	fmt.Println("fetch JWKS")
	//
	//	const key = "<<JWT>>"
	//	token, err := jwt.Parse(
	//		[]byte(key), // payload
	//		jwt.WithKeySet(keySet))
	//	if err != nil {
	//		t.Fatalf("failed to parse payload: %s\n", err)
	//		return
	//	}
	//	_ = token
	//
	//	fmt.Println("parsed jwt")
	//	cancel()
	//})

}
