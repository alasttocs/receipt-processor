# Receipt Processor Challenge

This project is a web service that processes receipts according to the provided API(s). Since the challenge documentation states that the api need not be ready for production, it was designed to simplify aspects such as authentication, logging, port routing, etc.

# Command Line Arguements

- `-noauth`: Runs the application without authentication.
- `-debug`: Enables debug mode for additional logging to assist with troubleshooting.
- `-log`: Enables logging to a file.
- `-logfile`: Overrides the name of the default log file.

_Note: for challenge simplicity logToFile/logFileName options are not fully supported when running in a docker container. I wanted to avoid the need for the reviewer to mount disks, copy additional files, etc._

# Installation and Usage

The application will be accessible at http://localhost:8080

To build the Docker Image: `docker build -t fetchapi .`

To run the Docker container **without** authentication: `docker run -p 8080:8080 fetchapi -noauth`

To run the Docker container **with** authentication simply remove the flag: `docker run -p 8080:8080 fetchapi`

Running Locally can be acheived with standard go commands: `go build -o fetchAPI` & `./fetchAPI`

# Running Tests

Tests are configured to run locally. Before running tests, ensure that the API server is running (either within Docker or locally) without the `-noauth` flag provided (there is a test for authentication). Then simply execute the standard go command `go test`

_Note: Logging to tests are written to `logs/testlogfile.log`_

# File Descriptions

- **apiAuth.go:** Handles authentication for the API.
- **main.go:** Entry point of the application. Sets up routes and handles HTTP requests.
- **utils.go:** Provides utility functions for processing receipts and calculating points.

- **api_test.go:** Contains test cases for the API endpoints (including the provided example requests).
- **utils_unit_test.go:** Test cases for the utility functions that help to caclulate receipt points.
