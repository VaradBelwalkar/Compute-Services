# Use a lightweight Alpine image
FROM alpine:latest

RUN apt-get update && apt-get install -y ca-certificates
# Set the working directory
WORKDIR /app

# Copy the Go application source code to the container
COPY bin/golang_server .

# Set the default command to run the Go application
CMD ["./golang_server"]
