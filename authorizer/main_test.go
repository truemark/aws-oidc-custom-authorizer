package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"io/ioutil"
	"testing"
	"time"
)

type TestObj struct {
	key      string
	jwksUrl  string
	audience string
}

func LoadTestConfig() ([]TestObj, error) {
	file, _ := ioutil.ReadFile("../main_test.json")
	data := []TestObj{}
	err := json.Unmarshal([]byte(file), &data)
	return data, err
}

// TODO: Make me a test runner, data driven....
func TestHandler(t *testing.T) {
	fmt.Println("   ---- DOING (PRE) TEST ----   ")
	testConfig, err := LoadTestConfig()
	if err != nil {
		fmt.Println("Error occurred reading test-def file: ../main_test.json")
		LogError(err)
	}
	for i := 0; i < len(testConfig); i++ {
		testObj := testConfig[i]
		fmt.Println("   ---- DOING TEST ----   ")
		fmt.Println("\t Key: " + testObj.key)
		fmt.Println("\t JWKS-Url: " + testObj.jwksUrl)
		fmt.Println("\t Audience: " + testObj.audience)

		t.Run("Example", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			//const jwksUrl = "url"
			ar := jwk.NewAutoRefresh(ctx)
			ar.Configure(testObj.jwksUrl, jwk.WithMinRefreshInterval(15*time.Minute))
			keySet, err := ar.Refresh(ctx, testObj.jwksUrl)
			if err != nil {
				t.Fatalf("failed to refresh JWKS: %s\n", err)
				return
			}
			fmt.Println("fetch JWKS")

			token, err := jwt.Parse(
				[]byte(testObj.key), // payload
				jwt.WithValidate(true),
				jwt.WithKeySet(keySet),
				jwt.WithAudience(testObj.audience))
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

}
