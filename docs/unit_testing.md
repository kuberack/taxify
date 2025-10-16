
# Unit Testing approach

 - Use mocks for the external services
 - For services such as twilio, the service's openapi specfiication is used to run mocks http servers using prism.
 - For database, [sqlmock](https://github.com/DATA-DOG/go-sqlmock) is used.
