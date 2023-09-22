# use official Golang image
FROM golang:1.21

# set working directory
WORKDIR /app

# Copy the source code
COPY . . 

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o pay ./cmd/main.go 

#EXPOSE the port
EXPOSE 8000

# Run the executable
CMD ["./pay"]