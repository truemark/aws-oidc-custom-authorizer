build-AuthorizerFunction:
	GOOS=linux go build -o bootstrap
	cp ./bootstrap $(ARTIFACTS_DIR)/.

# ERIK: We need to figure out how to mix up the build for dev vs prod/deploy. 
# https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html
# GOARCH=arm64 GOOS=linux go build -o authorizer
# cp ./authorizer $(ARTIFACTS_DIR)/.