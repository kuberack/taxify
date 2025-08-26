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
	# setup the db
	@echo "Running bash commands"
	for file in internal/models/*.sql; do \
	mysql -u shiv -p'shiv123' < $$file; \
	done
	# test
	go test -run Integration ./internal/api/

unit_test:
	go test -run Unit ./internal/api/

clean:
	rm -rf bin/*
	rm -rf openapi/gen.go
