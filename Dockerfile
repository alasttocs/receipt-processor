FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o fetchAPI

EXPOSE 8080

# Set a default CMD instruction
CMD ["sh", "-c", "if [ \"$NOAUTH\" = \"true\" ]; then ./fetchAPI -noauth; else ./fetchAPI; fi"]
