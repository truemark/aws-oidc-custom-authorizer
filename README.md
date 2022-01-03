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

## License

The contents of this repository are under the BSD 3-Clause license. See the
license [here](https://github.com/truemark/aws-cli-docker/blob/main/LICENSE.txt).
