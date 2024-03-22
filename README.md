# Receipt Processor Challenge

This project is a web service that processes receipts according to the provided API(s). Since the challenge documentation states that the api need not be ready for production, it was designed to simplify aspects such as authentication, logging, port routing, etc.

# Command Line Arguements

- `-noAuthMode`: Runs the application in test mode without authentication.
- debugMode: Enables debug mode for logging to assist with troubleshooting.
- logToFile: Enables logging to a file.
- logFileName: Overrides the name of the log file.

Note: for challenge simplicity logToFile/logFileName options are not fully supported when running in a docker container. I wanted to avoid the need for the reviewer to mount disks, copy additional files, etc.

# Installation and Usage

The application will be accessible at http://localhost:8080

To build the Docker Image: `docker build -t fetchapi .`

To run the Docker container without authorization: `docker run -p 8080:8080 fetchapi -noauth`

To run the Docker container with authorization simply remove the flag: `docker run -p 8080:8080 fetchapi`

Running Locally can be acheived with standard go commands: `go build -o fetchAPI` & `./fetchAPI`

# Running Tests

Before running tests, ensure that the API server is running. Then simply execute the standard go command `go test`

# File Descriptions

- apiAuth.go: Handles authentication for the API.
- main.go: Entry point of the application. Sets up routes and handles HTTP requests.
- utils.go: Provides utility functions for processing receipts and calculating points.

- api_test.go: Contains test cases for the API endpoints (including the provided example requests).
- utils_unit_test.go: Test cases for the utility functions that help to caclulate receipt points.
