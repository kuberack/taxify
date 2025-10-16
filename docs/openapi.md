
Refs: 
 - https://github.com/oapi-codegen/oapi-codegen
 - https://github.com/oapi-codegen/nethttp-middleware
 - https://stoplight.io/open-source/prism

# Usage
 - The taxify API is captured as an openapi yaml file ([taxify](../openapi/api.yaml)).
 - Server-side code generation for standard net/http server is done using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) 
 - Additionally, request validation is done using [nethttp-middleware](https://github.com/oapi-codegen/nethttp-middleware)
 - For external services, such as twilio, the service's openapi specification is used by [prism](https://stoplight.io/open-source/prism) to generate a mock http server.
