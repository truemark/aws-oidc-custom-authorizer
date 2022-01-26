package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	// Path as post-fix for OpenID Configuration URL
	OpenidConfigUrlPostFix = ".well-known/openid-configuration"
)

type OpenIDConfig struct {
	Issuer                             string   `json:"issuer"`
	JWKS_URI                           string   `json:"jwks_uri"`
	AuthorizationEndpoint              string   `json:"authorization_endpoint"`
	TokenEndpoint                      string   `json:"token_endpoint"`
	UserInfoEndpoint                   string   `json:"userinfo_endpoint"`
	EndSessionEndpoint                 string   `json:"end_session_endpoint"`
	CheckSessionIFrame                 string   `json:"check_session_iframe"`
	RevocationEndpoint                 string   `json:"revocation_endpoint"`
	IntrospectionEndpoint              string   `json:"introspection_endpoint"`
	DeviceAuthorizationEndpoint        string   `json:"device_authorization_endpoint"`
	FrontChannelLogoutSupported        bool     `json:"frontchannel_logout_supported"`
	FrontChannelLogoutSessionSupported bool     `json:"frontchannel_logout_session_supported"`
	BackChannelLogoutSupported         bool     `json:"backchannel_logout_supported"`
	BackChannelLogoutSessionSupported  bool     `json:"backchannel_logout_session_supported"`
	ScopesSupported                    []string `json:"scopes_supported"`
	ClaimsSupported                    []string `json:"claims_supported"`
	GrantTypesSupported                []string `json:"grant_types_supported"`
	ResponseTypesSupported             []string `json:"response_types_supported"`
	ResponseModesSupported             []string `json:"response_modes_supported"`
	TokenEndpointAuthMethodsSupported  []string `json:"token_endpoint_auth_methods_supported"`
	IdTokenSigningAlgValuesSupported   []string `json:"id_token_signing_alg_values_supported"`
	SubjectTypesSupported              []string `json:"subject_types_supported"`
	CodeChallengeMethodsSupported      []string `json:"code_challenge_methods_supported"`
	RequestParameterSupported          bool     `json:"request_parameter_supported"`
}

type Config struct {
	AuthorityEnv    string
	AudienceEnv     string
	OpenIDConfigURL string // TODO: not sure we need this...
	OpenIDConfig    *OpenIDConfig
}

func SetupConfig() (*Config, error) {
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

	config := Config{
		AuthorityEnv:    authorityEnv,
		AudienceEnv:     audienceEnv,
		OpenIDConfigURL: openIdConfigURL,
		OpenIDConfig:    openIdConfig,
	}
	return &config, nil
}

func getOpenIDConfiguration(url string) *OpenIDConfig {
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	log.Debug().Msg("HTTP GET Request Completed")

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var openIDConfig OpenIDConfig
	json.Unmarshal(body, &openIDConfig)
	log.Debug().
		Str("body", string(body)).
		Msg("OpenID (body) Configuration Struct successfully unmarshalled")

	return &openIDConfig
}
