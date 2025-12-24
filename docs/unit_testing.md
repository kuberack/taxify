
# Unit Testing approach

 - Use mocks for the external services
 - For services such as twilio, the service's openapi specfication ([link](https://raw.githubusercontent.com/twilio/twilio-oai/main/spec/json/twilio_verify_v2.json)is used to run mocks http servers using prism.
 - For database, [sqlmock](https://github.com/DATA-DOG/go-sqlmock) is used.
 - Running the unit test

```
kuberack/taxify$ ./test_unit.sh
setting environment variables
starting prism mock server
prism mock server is listening on port 4010
running go test
ok  	kuberack.com/taxify/internal/api	0.305s
prism mock server terminated
```
