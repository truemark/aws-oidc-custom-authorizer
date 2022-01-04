package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/truemark/aws-oidc-custom-authorizer/logging"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

type Config interface {
	AuthorityEnv    string
	AudienceEnv     string
	OpenIDConfigURL string
	OpenIDConfig    OpenIDConfig
}

type OpenIDConfig struct {
	Issuer                      string `json:"issuer"`                        //: "https://auth.youngliving.com",
	JWKS_URI                    string `json:"jwks_uri"`                      //: "https://auth.youngliving.com/.well-known/openid-configuration/jwks",
	AuthorizationEndpoint       string `json:"authorization_endpoint"`        //: "https://auth.youngliving.com/connect/authorize",
	TokenEndpoint               string `json:"token_endpoint"`                //: "https://auth.youngliving.com/connect/token",
	UserInfoEndpoint            string `json:"userinfo_endpoint"`             //: "https://auth.youngliving.com/connect/userinfo",
	EndSessionEndpoint          string `json:"end_session_endpoint"`          //: "https://auth.youngliving.com/connect/endsession",
	CheckSessionIFrame          string `json:"check_session_iframe"`          //: "https://auth.youngliving.com/connect/checksession",
	RevocationEndpoint          string `json:"revocation_endpoint"`           //: "https://auth.youngliving.com/connect/revocation",
	IntrospectionEndpoint       string `json:"introspection_endpoint"`        //: "https://auth.youngliving.com/connect/introspect",
	DeviceAuthorizationEndpoint string `json:"device_authorization_endpoint"` //: "https://auth.youngliving.com/connect/deviceauthorization",

	FrontChannelLogoutSupported        bool `json:"frontchannel_logout_supported"`         //: true,
	FrontChannelLogoutSessionSupported bool `json:"frontchannel_logout_session_supported"` //: true,
	BackChannelLogoutSupported         bool `json:"backchannel_logout_supported"`          //: true,
	BackChannelLogoutSessionSupported  bool `json:"backchannel_logout_session_supported"`  //: true,

	ScopesSupported                   []string `json:"scopes_supported"`                      //: [],
	ClaimsSupported                   []string `json:"claims_supported"`                      //: [],
	GrantTypesSupported               []string `json:"grant_types_supported"`                 //: [],
	ResponseTypesSupported            []string `json:"response_types_supported"`              //: [],
	ResponseModesSupported            []string `json:"response_modes_supported"`              //: [],
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"` //: [],
	IdTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"` //: [],
	SubjectTypesSupported             []string `json:"subject_types_supported"`               //: [],
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`      //: [],
	RequestParameterSupported         bool     `json:"request_parameter_supported"`           //: true
}

func setupConfig() (Config, error) {

	authorityEnv := os.Getenv("AUTHORITY")
	audienceEnv := os.Getenv("AUDIENCE")

	openIdConfigURL := authorityEnv + OpenidConfigUrlPostFix
	if strings.HasPrefix(openIdConfigURL, "http://") {
		return nil, errors.New("HTTP URL values for the AUTHORITY environment-variable is unsupported.")
	}
	// if strings.Contains(openIdConfigURL, "//") {
	// 	openIdConfigURL = strings.Replace(openIdConfigURL, "//", "/", -1)
	// }
	openIdConfig := getOpenIDConfiguration(openIdConfigURL)

	config := Config{
		AuthorityEnv:    authorityEnv,
		AudienceEnv:     audienceEnv,
		OpenIDConfigURL: openIdConfigURL,
		OpenIDConfig:    openIdConfig,
	}
	return config, nil
}

func getOpenIDConfiguration(url string) OpenIDConfig {
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var openIDConfig OpenIDConfig

	json.Unmarshal(body, &openIDConfig)
	fmt.Printf("Results: %v\n", openIDConfig)

	return openIDConfig
}

func getPolicyDocument() {

}

func getToken() {
	// params {methodArn, authorizationToken}

}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logging.LogRequest(request)

	fmt.Printf("request: %v\n", request)

	config := setupConfig()
	fmt.Printf("config: %v\n", config)

	fmt.Printf("request.Headers: %v\n", request.Headers)
	fmt.Printf("request.Body: %v\n", request.Body)

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
