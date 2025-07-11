tidy:
	go mod tidy

codegen: | bin
	make -C openapi $@

webapp: | bin
	mkdir -p bin; go build -o bin/taxifyapp ./cmd/app

bin:
	mkdir -p bin

clean:
	rm -rf bin/*
	rm -rf openapi/gen.go