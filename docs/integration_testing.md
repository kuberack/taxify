
# Integration testing
 - Integration testing uses the actual external services
 - There are two scenarios
   - baremetal
   - docker

# Baremetal
 - The taxify service, and the database are run on the baremetal
 - For other external service, the actual service is used.
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

```
```
