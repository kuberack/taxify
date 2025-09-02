
# include /.$(PWD)/.env_unit_test

tidy:
	go mod tidy

codegen: | bin
	oapi-codegen -config openapi/cfg.yaml openapi/api.yaml; cp openapi/gen.go internal/api/

webapp: | bin
	go build -o bin/taxifyapp ./cmd/app

.PHONY: api
api: | bin
	go build -o bin/taxifyapi ./cmd/apiserver

bin:
	mkdir -p bin

integration_test:
	# setup the db; it is assumed mysql is installed on the machine and
	# has required username, password
	@echo "Running bash commands"
	for file in internal/models/*.sql; do \
	mysql -u shiv -p'shiv123' < $$file; \
	done
	# Load env, and run the integration test
	. $(PWD)/.env_integration_test; go test -run Integration ./internal/api/

unit_test:
	# it is assumed that prism is already installed
	# npm install -g @stoplight/prism-cli
	# run the prism mock of the twilio verify API
	prism mock https://raw.githubusercontent.com/twilio/twilio-oai/main/spec/json/twilio_verify_v2.json &
	# wait for the mock server to be up
	while ! lsof -i :4010 -sTCP:LISTEN >/dev/null 2>&1; do sleep 1; done; echo "Port 4010 is now listening"

	# load env variables, and run the unit test (Test functions with Unit in the name)
	. $(PWD)/.env_unit_test; go test -run Unit ./internal/api/ 2>&1

clean:
	rm -rf bin/*
	rm -rf openapi/gen.go
