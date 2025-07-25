tidy:
	go mod tidy

codegen: | bin
#make -C openapi $@
	oapi-codegen -config openapi/cfg.yaml openapi/api.yaml; cp openapi/gen.go api/

webapp: | bin
	go build -o bin/taxifyapp ./cmd/app

api: | bin
	go build -o bin/taxifyapi ./cmd/apiserver

bin:
	mkdir -p bin

clean:
	rm -rf bin/*
	rm -rf openapi/gen.go