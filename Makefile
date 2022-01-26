.PHONY: clean build run

TOK := $(shell ./get_token.sh);

clean:
	rm -rf .aws-sam

build:
	sam build --debug

run:
	sam local invoke "AuthorizerFunction" -e ./event.json
	#$(./get_token.sh > sam local invoke "AuthorizerFunction" --debug --region "us-west-1" --event -)

