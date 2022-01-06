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

// func getJWKStr(url string) string {
// 	res, err := http.Get(url)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return string(body)
// }

func getPolicyDocument() {

}

func getToken() {
	// params {methodArn, authorizationToken}

}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logging.LogRequest(request)

	fmt.Printf("request: %v\n", request)
	config, err := setupConfig()
	if err != nil {
		fmt.Println("error")
	}
	logging.LogConfig(config)

	ctx := context.Background()
	ar := jwk.NewAutoRefresh(ctx)
	ar.Configure(config.OpenIDConfig.JWKS_URI)
	keyset, err := ar.Refresh(ctx, config.OpenIDConfig.JWKS_URI)
	if err != nil {
		fmt.Println("error")
		return events.APIGatewayProxyResponse{}, err
	}
	logging.LogKeySet(keyset)

	// authToken := getToken()
	// resp, err := http.Get(DefaultHTTPGetAddress)
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{}, err
	// }
	// if resp.StatusCode != 200 {
	// 	return events.APIGatewayProxyResponse{}, ErrNon200Response
	// }
	// ip, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{}, err
	// }
	// if len(ip) == 0 {
	// 	return events.APIGatewayProxyResponse{}, ErrNoIP
	// }
	ip := "myIP"
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, %v", string(ip)),
		StatusCode: 200,
	}, nil
}

func main() {
	logging.Init()
	lambda.Start(handler)
}
