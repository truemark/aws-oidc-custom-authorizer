.PHONY: setupevt clean build run

TOK := $(shell ./get_token.sh);
EVENT_FILE := "./event2.json"

setupevt:
	./get_token.sh > $(EVENT_FILE)


clean:
	rm -rf .aws-sam
	rm -f $(EVENT_FILE)
	
build:
	sam build --debug

run:
	sam local invoke "AuthorizerFunction" -e $(EVENT_FILE)