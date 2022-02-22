package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
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
		j, err := json.MarshalIndent(request, "", "  ")
		if err != nil {
			LogError(err)
		}
		log.Trace().Msg(fmt.Sprintf("request: %s\n", string(j)))
	}
}

func LogConfig(config *Config) {
	jsonConfig, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		LogError(err)
	}
	log.Debug().Msg("Configuration: " + string(jsonConfig))
}

func LogKeySet(keyset jwk.Set) {
	if log.Debug().Enabled() {
		msg := "KeySet as JSON:\n"
		jsonKeyset, err := json.MarshalIndent(keyset, "", "  ")
		if err != nil {
			LogError(err)
		}
		msg += string(jsonKeyset)
		log.Debug().Msg(msg)
	}
}

func LogKey(key jwk.Key) {
	if log.Debug().Enabled() {
		msg := "Key as JSON:\n"
		jsonKey, err := json.MarshalIndent(key, "", "  ")
		if err != nil {
			LogError(err)
		}
		msg += string(jsonKey)
		log.Debug().Msg(msg)
	}
}

func LogToken(tok jwt.Token) {
	if log.Debug().Enabled() {
		msg := "Token as JSON:\n"
		jsonTok, err := json.MarshalIndent(tok, "", "  ")
		if err != nil {
			LogError(err)
		}
		msg += string(jsonTok)
		log.Debug().Msg(msg)
	}
}

func LogError(err error) {
	log.Info().
		Str("err", err.Error()).
		Msg("Error Message Logged")
}
