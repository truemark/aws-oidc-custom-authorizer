package logging

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/truemark/aws-oidc-custom-authorizer/config"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes logging reading the LOG_LEVEL from the environment
func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel := os.Getenv("LOG_LEVEL")
	// TODO Log level should be set up in template.yaml as a variable
	if logLevel == "Trace" {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else if logLevel == "Debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else if logLevel == "Info" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if logLevel == "Warn" {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if logLevel == "Error" {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if logLevel == "Fatal" {
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: true, // CloudWatch doesn't like colors
	})
	log.Trace().Msgf("logging initialized successfully")
}

// LogRequest will log the entire lambda request to stdout if trace is enabled
func LogRequest(request events.APIGatewayProxyRequest) {
	log.Debug().
		Str("requestId", request.RequestContext.RequestID).
		Msg("received request")
	if log.Trace().Enabled() {
		j, _ := json.MarshalIndent(request, "", "  ")
		fmt.Println(string(j))
	}
}

func LogConfig(config *config.Config) {
	fmt.Printf("config: %v\n", config)
	fmt.Printf("config.JWKS_URI: %v\n", config.OpenIDConfig.JWKS_URI)
}

func LogKeySet(keyset jwk.Set) {
	fmt.Printf("KeySet as JSON:\n")
	jsonKeyset, _ := json.MarshalIndent(keyset, "", "  ")
	fmt.Println(string(jsonKeyset))
}
