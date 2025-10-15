
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

unit_test:
	# it is assumed that prism is already installed
	# npm install -g @stoplight/prism-cli
	./test_unit.sh

integration_test_baremetal:
	# setup the db; it is assumed mysql is installed on the machine and
	# has required username, password
	./test_integration_baremetal.sh

integration_test_docker:
	./test_integration_docker.sh

clean:
	rm -rf bin/*
	rm -rf openapi/gen.go
