# Specifies a parent image
FROM golang:1.24.5-bookworm
 
RUN apt-get update && apt-get install -y make

# Creates an app directory to hold your app’s source code
WORKDIR /app
 
# Copies everything from your root directory into /app
COPY . .
 
# Installs Go dependencies
RUN make tidy
 
# Builds your app with optional configuration
RUN make api
 
# Tells Docker which network port your container listens on
EXPOSE 8080
 
# Specifies the executable command that runs when the container starts
CMD [ “./bin/taxifyapi” ]
