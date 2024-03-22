FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o fetchAPI

EXPOSE 8080

# Set a default CMD instruction
CMD ["./fetchAPI"]

ARG NOAUTH
# Conditionally set the CMD instruction based on the value of the NOAUTH build argument
CMD ["./fetchAPI", "-noauth"]
