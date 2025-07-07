tidy:
	go mod tidy

codegen:
	oapi-codegen -config cfg.yaml ../../api.yml

