
# Integration testing
 - Integration testing uses the actual external services
 - There are two scenarios
   - baremetal
   - docker

# Baremetal
 - The taxify service, and the database are run on the baremetal
 - For other external service, the actual service is used.

# Docker
 - In a docker integration scenario, docker compose is used
 - The taxify service is run in a docker container
 - For external services such as twilio, the actual service is used
 - For database, a mysql docker container is used
