#!/bin/bash

# Setup the environment vars
echo "setting environment variables"
set -a
source .env_integration_test_baremetal
source .env_secrets
set +a

# Setup the sql db
echo "setting sql db"
for file in internal/models/*.sql
do
	mysql -u shiv -p'shiv123' --port=3306 < $file
done

# Run the go test
echo "running go test"
go test -v -run Integration ./internal/api/
