package main

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog/log"
)

var (
	ErrStatusCode     = 500
	SuccessStatusCode = 200
)

func getToken(requestHeader string) (string, error) {
	match, _ := regexp.Match(`Bearer (.*)`, []byte(requestHeader))
	log.Debug().
		Bool("match", match).
		Str("authorizationToken", requestHeader).
		Msg("Header Matched for AuthorizationToken on Bearer")

	if !match {
		errMsg := fmt.Sprintf("Invalid Authorization token - %s does not match \"Bearer .*\"\n", requestHeader)
		return "", errors.New(errMsg)
	}
	r, _ := regexp.Compile(`Bearer (.*)`)
	matchedToken := r.FindString(requestHeader)
	keyToken := strings.Replace(matchedToken, "Bearer ", "", 1)
	log.Debug().
		Str("matchedToken", keyToken).
		Msg("DELETE-ME: matchedToken Found is:")

	return keyToken, nil
}

func handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayProxyResponse, error) {
	LogRequest(event)

	// Setup our Config object
	config, err := SetupConfig()
	if err != nil {
		LogError(err)
		return events.APIGatewayProxyResponse{
			Body:       "ERROR on Setting up Configuration data-structure for aws-oidc-authorizer-lambda",
			StatusCode: ErrStatusCode,
		}, err
	}
	LogConfig(config)

	// Setup JWK and retreive our cached key...
	ar := jwk.NewAutoRefresh(ctx)
	ar.Configure(config.OpenIDConfig.JWKS_URI, jwk.WithMinRefreshInterval(15*time.Minute))
	keySet, err := ar.Refresh(ctx, config.OpenIDConfig.JWKS_URI)
	if err != nil {
		LogError(err)
		return events.APIGatewayProxyResponse{
			// TODO: Update/Enhance messaging for err handling - do we need a JSON struct here in the body? etc...

			Body:       "ERROR on Setting up JWK AUTO-REFRESH TODO::MAKE ME BETTER",
			StatusCode: ErrStatusCode,
		}, err
	}
	LogKeySet(keySet)

	bearer := event.AuthorizationToken
	authToken, err := getToken(bearer)
	if err != nil {
		LogError(err)
		return events.APIGatewayProxyResponse{
			// TODO: Update/Enhance messaging for err handling - do we need a JSON struct here in the body? etc...
			Body:       "ERROR on GetToken TODO::MAKE ME BETTER",
			StatusCode: ErrStatusCode,
		}, err
	}

	// Perform main verification
	token, err := jwt.Parse(
		[]byte(authToken),
		jwt.WithValidate(true),
		jwt.WithKeySet(keySet),
		jwt.WithAudience(config.AudienceEnv),
	)
	if err != nil {
		LogError(err)
		return events.APIGatewayProxyResponse{
			// TODO: What do we want to supply in the body for error handling, etc...
			Body:       "ERROR on VerifyToken TODO::MAKE ME BETTER",
			StatusCode: ErrStatusCode,
		}, err
	}
	LogToken(token)
	_ = token

	return events.APIGatewayProxyResponse{
		Body:       "Verification Success",
		StatusCode: SuccessStatusCode,
	}, nil
}

// TODO: Since this is a secure/sensitive app, we need to determine what can and cant be logged into AWS, etc

func main() {
	Init()
	lambda.Start(handler)
}
