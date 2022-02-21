package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes logging reading the LOG_LEVEL from the environment
func Init() {
	fmt.Println("Logging INIT")
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
	log.Trace().Msg("logging initialized successfully")
}

// LogRequest will log the entire lambda request to stdout if trace is enabled
func LogRequest(request events.APIGatewayCustomAuthorizerRequest) {
	if log.Trace().Enabled() {
		j, _ := json.MarshalIndent(request, "", "  ")
		log.Trace().Msg(fmt.Sprintf("request: %s\n", string(j)))
	}
}

func LogConfig(config *Config) {
	log.Debug().
		Str("config", fmt.Sprintln("%v", config)).
		Msg(fmt.Sprintf("config.JWKS_URI: %s\n", config.OpenIDConfig.JWKS_URI))
}

func LogKeySet(keyset jwk.Set) {
	fmt.Printf("KeySet as JSON:\n")
	jsonKeyset, _ := json.MarshalIndent(keyset, "", "  ")
	fmt.Println(string(jsonKeyset))
}

func LogKey(key jwk.Key) {
	fmt.Printf("Key as JSON:\n")
	jsonKey, _ := json.MarshalIndent(key, "", "  ")
	fmt.Println(string(jsonKey))
}

func LogToken(tok jwt.Token) {
	fmt.Printf("Token as JSON:\n")
	jsonTok, err := json.MarshalIndent(tok, "", "  ")
	if err != nil {
		LogError(err)
	}
	fmt.Println(string(jsonTok))
}

func LogError(err error) {
	log.Debug().
		Str("err", err.Error()).
		Msg("Error Message Logged")
}
