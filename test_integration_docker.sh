#!/bin/bash

# Delete all current containers
docker stop $(docker ps -aq)
docker rm $(docker ps -aq)

# Remove the mysql data directory. Only then the docker compose
# will execute the sql files in initdb directory
sudo rm -rf deployment/docker/mysql_data/

# setup the environment
echo "setting environment variables"
set -a
source .env_integration_test_docker
source .env_secrets
set +a

# Bring up the containers
echo "docker compose"
docker compose -f deployment/docker/docker-compose.yaml up -d

# wait until the health of the api service is fine
while true; do
  curl --silent --fail http://localhost:8080/healthz
  if [ $? -eq 0 ]; then
    echo "api service docker container is listening on port 8080, exiting loop"
    break
  else
    echo "api service docker container not yet ready"
  fi
  sleep 1
done

# run the tests now
echo "launching curl"
curl -w "%{http_code}" -X POST --header "Content-Type: application/json" http://localhost:8080/signup/phone?type=driver --data '{"phone":9886240527}' -o /dev/null

# Check status
echo "checking response status"
if [ "$http_code" != "200" ]; then
  echo "Request failed with status code $http_code"
  exit 1
fi

# Bring down the containers
# docker compose -f deployment/docker/docker-compose.yaml up -d

