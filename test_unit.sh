#!/bin/bash

# set environment variables
echo "setting environment variables"
set -a
source .env_unit_test
source .env_secrets
set +a

# Launch prism mock and wait for it to finish
echo "starting prism mock server"
prism mock https://raw.githubusercontent.com/twilio/twilio-oai/main/spec/json/twilio_verify_v2.json > /dev/null &

# wait for the mock server to be up
while ! lsof -i :4010 -sTCP:LISTEN > /dev/null; do 
	sleep 1; 
done

echo "prism mock server is listening on port 4010"
prism_pid=$!

# run the test
echo "running go test"
go test -run Unit ./internal/api/ 2>&1

# clean up prism mock
kill -9 $prism_pid
echo "prism mock server terminated"

