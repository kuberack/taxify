
# Integration testing
 - Integration testing uses the actual external services
 - There are two scenarios
   - baremetal
   - docker

# Baremetal
 - The taxify service, and the database are run on the baremetal
 - For other external service, the actual service is used.
 - Start the mysql service in the VM
 - Running the test

```
$ ./test_integration_baremetal.sh
setting environment variables
setting sql db
mysql: [Warning] Using a password on the command line interface can be insecure.
mysql: [Warning] Using a password on the command line interface can be insecure.
mysql: [Warning] Using a password on the command line interface can be insecure.
mysql: [Warning] Using a password on the command line interface can be insecure.
running go test
=== RUN   TestPostSignupPhoneIntegration
Caller: /home/hima/repos/kuberack/taxify/internal/api/server_middleware.go
Environment: baremetal
2025/12/23 20:36:12 INFO request received user.ip="" request.method=POST request.url="/signup/phone?type=driver" request.proto=HTTP/1.1
getDb: INTEGRATION_TEST_BAREMETAL
--- PASS: TestPostSignupPhoneIntegration (6.55s)
PASS
ok  	kuberack.com/taxify/internal/api	6.576s
```

# Docker
 - In a docker integration scenario, docker compose is used
 - The taxify service is run in a docker container
 - For external services such as twilio, the actual service is used
 - For database, a mysql docker container is used
 - Stop any mysql service running in the VM host since mysql in docker will attempt to open the same port.

```
kuberack/taxify$ ./test_integration_docker.sh
"docker stop" requires at least 1 argument.
See 'docker stop --help'.

Usage:  docker stop [OPTIONS] CONTAINER [CONTAINER...]

Stop one or more running containers
"docker rm" requires at least 1 argument.
See 'docker rm --help'.

Usage:  docker rm [OPTIONS] CONTAINER [CONTAINER...]

Remove one or more containers
setting environment variables
docker compose
[+] Building 0.0s (0/0)
[+] Running 3/3
 ✔ Network docker_app-network  Created                                     0.7s
 ✔ Container docker-db-1       Healthy                                    52.6s
 ✔ Container docker-api-1      Started                                    54.4s
api service docker container not yet ready
api service docker container not yet ready
api service docker container not yet ready
api service docker container not yet ready
api service docker container not yet ready
api service docker container is listening on port 8080, exiting loop
sleep for a minute
launching curl
checking response status
Test 1 success
```
