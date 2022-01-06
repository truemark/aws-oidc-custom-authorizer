package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/truemark/aws-oidc-custom-authorizer/config"
	"github.com/truemark/aws-oidc-custom-authorizer/logging"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lestrrat-go/jwx/jwk"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")

	// Path as post-fix for OpenID Configuration URL
	OpenidConfigUrlPostFix = ".well-known/openid-configuration"
)

func setupConfig() (*config.Config, error) {
	authorityEnv := os.Getenv("AUTHORITY")
	audienceEnv := os.Getenv("AUDIENCE")

	log.Debug().Str("Authority", authorityEnv)
	log.Debug().Str("Audience", audienceEnv)

	openIdConfigURL := authorityEnv + OpenidConfigUrlPostFix
	if strings.HasPrefix(openIdConfigURL, "http://") {
		return nil, errors.New("HTTP URL values for the AUTHORITY environment-variable is unsupported.")
	}
	openIdConfigURL = strings.Replace(openIdConfigURL, "https://", "", 1)
	if strings.Contains(openIdConfigURL, "//") {
		openIdConfigURL = strings.Replace(openIdConfigURL, "//", "/", -1)
	}
	openIdConfigURL = "https://" + openIdConfigURL
	log.Debug().
		Str("openIdConfigURL", openIdConfigURL).
		Msg("OpenID Configuration Data URL")
	openIdConfig := getOpenIDConfiguration(openIdConfigURL)

	config := config.Config{
		AuthorityEnv:    authorityEnv,
		AudienceEnv:     audienceEnv,
		OpenIDConfigURL: openIdConfigURL,
		OpenIDConfig:    openIdConfig,
	}
	return &config, nil
}

func getOpenIDConfiguration(url string) *config.OpenIDConfig {
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var openIDConfig config.OpenIDConfig
	json.Unmarshal(body, &openIDConfig)

	return &openIDConfig
}

func getPolicyDocument() {

}

func getToken(requestHeader string) (string, error) {
	return "", nil
}

func verifyToken(token string) (bool, error) {
	return true, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logging.LogRequest(request)

	// Setup our Config object
	config, err := setupConfig()
	if err != nil {
		fmt.Println("error")
	}
	logging.LogConfig(config)

	// Setup JWK and retreive our cached key...
	ctx := context.Background()
	ar := jwk.NewAutoRefresh(ctx)
	ar.Configure(config.OpenIDConfig.JWKS_URI)
	keyset, err := ar.Refresh(ctx, config.OpenIDConfig.JWKS_URI)
	if err != nil {
		fmt.Println("error")
		return events.APIGatewayProxyResponse{}, err
	}
	logging.LogKeySet(keyset)

	bearer := request.Headers["authorizationToken"]
	authToken, err := getToken(bearer)
	tokenVerified, err := verifyToken(authToken)
	fmt.Printf("Token Verified: %s\n", tokenVerified)

	return events.APIGatewayProxyResponse{
		Body:       "Hello JKW World",
		StatusCode: 200,
	}, nil
}

func main() {
	logging.Init()
	lambda.Start(handler)
}
