package main

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/truemark/aws-oidc-custom-authorizer/config"
	"github.com/truemark/aws-oidc-custom-authorizer/logging"
	"github.com/truemark/aws-oidc-custom-authorizer/verify"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lestrrat-go/jwx/jwk"
	// "github.com/pkg/errors"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func getPolicyDocument() {

}

func getToken(requestHeader string) (*string, error) {
	//r, _ := regexp.Compile(`^Bearer (.*)$`)
	//match := r.MatchString(requestHeader)

	match, _ := regexp.Match(`Bearer (.*)`, []byte(requestHeader))
	log.Debug().
		Bool("match", match).
		Str("authorizationToken", requestHeader).
		Msg("Header Matched for AuthorizationToken on Bearer")

	if !match {
		return nil, errors.New(fmt.Sprintf("Invalid Authorization token - %s does not match \"Bearer .*\"\n", requestHeader))
	}

	r, _ := regexp.Compile(`Bearer (.*)`)
	matchedToken := r.FindString(requestHeader)
	keyToken := strings.Replace(matchedToken, "Bearer ", "", 1)
	//keyToken = strings.Replace(keyToken, "'", "", -1)

	log.Debug().
		Str("matchedToken", keyToken).
		Msg("DELETE-ME: matchedToken Found is:")

	return &keyToken, nil
}

func handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("ENTERING HANDLER...")
	logging.LogRequest(event)

	// Setup our Config object
	config, err := config.SetupConfig()
	if err != nil {
		fmt.Println("error")
	}
	logging.LogConfig(config)

	// Setup JWK and retreive our cached key...
	ar := jwk.NewAutoRefresh(ctx)
	ar.Configure(config.OpenIDConfig.JWKS_URI)
	keyset, err := ar.Refresh(ctx, config.OpenIDConfig.JWKS_URI)
	if err != nil {
		logging.LogError(err)
		return events.APIGatewayProxyResponse{
			Body:       "ERROR on Setting up JWK AUTO-REFRESH TODO::MAKE ME BETTER",
			StatusCode: 500,
		}, err
	}
	logging.LogKeySet(keyset)

	bearer := event.AuthorizationToken
	fmt.Println("BEARER IS:")
	fmt.Println(bearer)
	authToken, err := getToken(bearer)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "ERROR on GetToken TODO::MAKE ME BETTER",
			StatusCode: 500,
		}, err
	}
	tokenVerified, kidVerified, err := verify.VerifyToken(*authToken, keyset)
	if err != nil {
		return events.APIGatewayProxyResponse{
			// TODO: What do we want to supply in the body for error handling, etc...
			Body:       "ERROR on VerifyToken TODO::MAKE ME BETTER",
			StatusCode: 500,
		}, err
	}
	fmt.Printf("Token Verified: %s\n", tokenVerified)

	msg := fmt.Sprintf("{\"kid\": \"%v\", \"verified\": %v, \"verificationMethod\": \"KID_VERIFICATION\"}", kidVerified, tokenVerified)
	return events.APIGatewayProxyResponse{
		Body:       msg,
		StatusCode: 200,
	}, nil
}

// TODO: Since this is a secure/sensitive app, we need to determine what can and cant be logged into AWS, etc

func main() {
	logging.Init()
	lambda.Start(handler)
}
