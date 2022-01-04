.PHONY: clean build run

clean:
	rm -rf .aws-sam

build:
	sam build --debug

run:
	sam local invoke "AuthorizerFunction" -e ./event.json 