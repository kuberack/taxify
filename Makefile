tidy:
	go mod tidy

codegen:
	make -C openapi $@

clean:
	rm -rf openapi/gen.go