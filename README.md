# aws-oidc-custom-authorizer

An AWS API Gateway [Lambda Authorizer](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-use-lambda-authorizer.html) that authorizes API requests against an OIDC provider.
This authorizer is written in golang and compiled to arm64 to maximize performance and minimize execution costs.

## Requirements
    - This will accept a AUTHORITY and AUDIENCE environment variable and be able to take OUATH2 JWT tokens and verify them.

## Libraries
    - https://github.com/lestrrat-go/jwx
    - https://github.com/hkra/go-jwks

## Maintainers

- [erikrj](https://github.com/erikrj)
- [briancabbott](https://github.com/briancabbott)

## License

The contents of this repository are under the BSD 3-Clause license. See the
license [here](https://github.com/truemark/aws-cli-docker/blob/main/LICENSE.txt).


## Dev Notes (temp section)

### x86_64 / ARM64 targets for dev vs deploy
We need to setup the build to produce targets depending on dev vs deploy - the local builds need 
the make target setup to build a x86_64 target for a local docker container deployment. See:
   - https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html

   
## TODOs

### Multi-targeted builds
We still need to figure out how to nicely sense the HW were compiling on and pic 
the appropriate metal for it... 

### Use Letstrat's verification API instead of manaul verification
Use the jws.VerifySet() from "github.com/lestrrat-go/jwx/jws" instead of manually 
iterating through each key. This keeps us congruent with the standard/free updates 
from Lestrat, etc. (just better practice)

### Remove/update printfs
Either take out the extraneous prints or move them to use Error-Logging if they are useful.

### Single Package/Module
Move everything into the authorizer package and out of individual sub-modules.

### Use Viper instead of manual for envs
Perform "out of the box" type setup activities with Viper instead of manual env-varing.
