# Go HTTP Request Logger
[![cve-trivy-scan](https://github.com/lao-tseu-is-alive/go-http-request-logger/actions/workflows/trivy-scan.yml/badge.svg)](https://github.com/lao-tseu-is-alive/go-http-request-logger/actions/workflows/trivy-scan.yml)
[![Go-Test](https://github.com/lao-tseu-is-alive/go-http-request-logger/actions/workflows/go.yml/badge.svg)](https://github.com/lao-tseu-is-alive/go-http-request-logger/actions/workflows/go.yml)


A simple and configurable Go HTTP server designed to capture and log every incoming request for debugging and development. It logs details such as the request method, URL, headers, and body, and outputs this information in a structured JSON format.

-----

### 🚀 Features

* **Comprehensive Logging**: Captures and logs the request method, URL, client IP, headers, and body.
* **Structured Output**: Outputs logged requests in a structured **JSON** format, making them easy to parse and analyze.
* **Configurable Port**: The listening port can be set via an environment variable.
* **Flexible Log Output**: Logs can be directed to `stdout`, `stderr`, a file, or completely discarded.
* **Unique Request IDs**: Assigns a unique ID to each request for easy tracking.

-----

### 🛠️ Getting Started


#### Running the Server from the source code

###### Prerequisites

*  Install **Go 1.24.5** or higher from [official Go.dev web](https://go.dev/dl/). 

###### Instructions

1.  **Clone this repo**
2.  **Run the application**: Open your terminal and execute the following command:
    ```bash
    git clone https://github.com/lao-tseu-is-alive/go-http-request-logger.git .
    cd go-http-request-logger
    go run cmd/goHttpRequestLoggerServer/goHttpRequestLoggerServer.go 
    
    goHttpRequestLogger, 2025/08/14 18:02:22 goHttpRequestLoggerServer.go:159: 🚀🚀 Starting App: goHttpRequestLogger, version: 0.0.1, build: unknown
    goHttpRequestLogger, 2025/08/14 18:02:22 goHttpRequestLoggerServer.go:94: INITIAL CALL TO requestHandler()
    goHttpRequestLogger, 2025/08/14 18:02:22 goHttpRequestLoggerServer.go:180: Server starting on localhost:8888


    ```
    By default, the server will start on `http://localhost:8888` and log output to `stderr`.

-----

### ⚙️ Configuration

You can customize the server's behavior using the following environment variables:

| Environment Variable | Description                                                                                                                                                             | Default Value | Example                            |
| :------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------ | :--------------------------------- |
| `PORT`               | The TCP/IP port the server should listen on.                                                                                                                            | `8888`              | `export PORT=9000`                 |
| `LOG_FILE`           | The destination for log output. It can be `stdout`, `stderr`, `DISCARD` (for no logging), or a specific file path.                                                        | `stderr`            | `export LOG_FILE=requests.log`     |

-----

### 🧪 How to Test

Once the server is running, you can send HTTP requests to it using tools like `curl`, Postman, or your web browser.

#### Example [curl](https://curl.se/) Commands

* **GET request**:

  ```bash
  curl http://localhost:8888/
  ```
will produce this on the server side console :   
```bash
  go run cmd/goHttpRequestLoggerServer/goHttpRequestLoggerServer.go 
goHttpRequestLogger, 2025/08/14 18:02:22 goHttpRequestLoggerServer.go:159: 🚀🚀 Starting App: goHttpRequestLogger, version: 0.0.1, build: unknown
goHttpRequestLogger, 2025/08/14 18:02:22 goHttpRequestLoggerServer.go:94: INITIAL CALL TO requestHandler()
goHttpRequestLogger, 2025/08/14 18:02:22 goHttpRequestLoggerServer.go:180: Server starting on localhost:8888

goHttpRequestLogger, 2025/08/14 18:03:32 goHttpRequestLoggerServer.go:98: ## ----- New Request D2F0HL7882250NMU7TUG ----- ##
goHttpRequestLogger, 2025/08/14 18:03:32 goHttpRequestLoggerServer.go:99: D2F0HL7882250NMU7TUG	Request : 	'GET /', 0 bytes, from 127.0.0.1:57932
goHttpRequestLogger, 2025/08/14 18:03:32 goHttpRequestLoggerServer.go:102: Headers:
goHttpRequestLogger, 2025/08/14 18:03:32 goHttpRequestLoggerServer.go:105: 	User-Agent: curl/8.5.0
goHttpRequestLogger, 2025/08/14 18:03:32 goHttpRequestLoggerServer.go:105: 	Accept: */*
{"Id":"D2F0HL7882250NMU7TUG","isoDateTime":"2025-08-14T16:03:32.742398378Z","protocol":"HTTP/1.1","method":"GET","url":"/","contentLength":0,"ipClient":"127.0.0.1:57932","headers":{"Accept":["*/*"],"User-Agent":["curl/8.5.0"]}}
goHttpRequestLogger, 2025/08/14 18:03:32 goHttpRequestLoggerServer.go:140: ## ----- End Request D2F0HL7882250NMU7TUG ----- ##

  ```
 and if you did discard the log you can directly parse the received requests logged as json with jq: 
```bash
LOG_FILE=DISCARD go run cmd/goHttpRequestLoggerServer/goHttpRequestLoggerServer.go |jq
{
  "Id": "D2F0J778822587HUODT0",
  "isoDateTime": "2025-08-14T16:06:52.172147895Z",
  "protocol": "HTTP/1.1",
  "method": "GET",
  "url": "/",
  "contentLength": 0,
  "ipClient": "127.0.0.1:40630",
  "headers": {
    "Accept": [
      "*/*"
    ],
    "User-Agent": [
      "curl/8.5.0"
    ]
  }
}

```
* **POST request with JSON body**:

  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"username":"devuser", "action":"login"}' http://localhost:8888/api/login
  ```

Each request you send will be logged in the terminal (or the file specified by `LOG_FILE`), showing the structured JSON output with all the captured details.