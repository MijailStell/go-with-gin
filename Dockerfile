# Start from the latest goland base image
FROM golang:latest

# Add maintainer info
LABEL maintainer="devops@company.com"

# Set the current working directory inside the container
WORKDIR /app

# Copy Go Modules dependency requirements file
COPY go.mod .

# Copy Go Domules expected hashes file
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy all the app source (recursively)
COPY . .

# Set http port
ENV PORT=5000

# Build the app
RUN go build

# Remove source file
RUN find . -name "*.go" -type f -delete
RUN find . -name "*.mod" -type f -delete
RUN find . -name "*.sum" -type f -delete

# Make port 5000 available to the world outside this countainer
EXPOSE $PORT

# Run the app
CMD [ "./microservices" ]
