tidy:
	go mod tidy

codegen: | bin
#make -C openapi $@
	oapi-codegen -config openapi/cfg.yaml openapi/api.yaml; cp openapi/gen.go api/

webapp: | bin
	mkdir -p bin; go build -o bin/taxifyapp ./cmd/app

bin:
	mkdir -p bin

clean:
	rm -rf bin/*
	rm -rf openapi/gen.go