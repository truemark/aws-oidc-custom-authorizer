package verify

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/rs/zerolog/log"
	"github.com/truemark/aws-oidc-custom-authorizer/base64"
)

func VerifyToken(tokenBuf string, keySet jwk.Set) (bool, string, error) {
	// protected, payload, signature, err := jws.SplitCompactString(tokenBuf)
	protected, _, _, err := jws.SplitCompactString(tokenBuf)
	if err != nil {
		return false, "", errors.New(`failed to split compact JWS message`)
	}

	decodedProtected, err := base64.Decode(protected)
	if err != nil {
		return false, "", errors.New(`failed to base64 decode protected headers`)
	}

	var protectedHeadersMap map[string]interface{}
	if err := json.Unmarshal(decodedProtected, &protectedHeadersMap); err != nil {
		return false, "", errors.New(`failed to decode protected headers`)
	}

	tokkid := protectedHeadersMap["kid"]
	strKid := fmt.Sprintf("%v", tokkid)

	i := 0
	for i < keySet.Len() {
		key, inRange := keySet.Get(i)
		log.Debug().Str("key.KeyID", key.KeyID())
		if inRange {
			if tokkid == key.KeyID() {
				log.Debug().Msgf("key matched on kid: %v\n", tokkid)
				return true, strKid, nil
			}
		}
		i = i + 1
	}

	return false, strKid, nil
}
