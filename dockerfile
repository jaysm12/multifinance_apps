# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
EXPOSE 32001
RUN go build -o server ./cmd/multifinance-apps
CMD ["./server"]