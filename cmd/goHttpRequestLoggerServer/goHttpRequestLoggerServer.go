package main

import (
	"encoding/json"
	"fmt"
	"github.com/lao-tseu-is-alive/go-http-request-logger/pkg/version"
	"github.com/rs/xid"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	initCallMsg            = "INITIAL CALL TO %s()\n"
	charsetUTF8            = "charset=UTF-8"
	MIMEAppJSON            = "application/json"
	MIMEAppJSONCharsetUTF8 = MIMEAppJSON + "; " + charsetUTF8
	HeaderContentType      = "Content-Type"
	defaultReadTimeout     = 10 * time.Second // max time to read request from the client
	defaultWriteTimeout    = 10 * time.Second // max time to write response to the client
	defaultIdleTimeout     = 2 * time.Minute  // max time for connections using TCP Keep-Alive
	defaultPort            = 8888
	defaultListenIP        = "localhost"
	defaultLogName         = "stderr"
)

type RequestLog struct {
	Id            string              `json:"Id"`
	IsoDateTime   time.Time           `json:"isoDateTime"`
	Protocol      string              `json:"protocol"`
	Method        string              `json:"method"`
	Url           string              `json:"url"`
	ContentLength int64               `json:"contentLength"`
	IpClient      string              `json:"ipClient"`
	Headers       map[string][]string `json:"headers"`
	Body          string              `json:"body,omitempty"`
}

// getPortFromEnvOrPanic returns a valid TCP/IP port from the environment or a default.
func getPortFromEnvOrPanic(defaultPort int) int {
	srvPort := defaultPort
	if val, exist := os.LookupEnv("PORT"); exist {
		if p, err := strconv.Atoi(val); err == nil {
			srvPort = p
		} else {
			panic(fmt.Errorf("💥💥 ERROR: CONFIG ENV PORT should contain a valid integer. %v", err))
		}
	}
	if srvPort < 1 || srvPort > 65535 {
		panic(fmt.Errorf("💥💥 ERROR: PORT should contain an integer between 1 and 65535"))
	}
	return srvPort
}

// GetLogWriterFromEnvOrPanic returns the name of the filename to use for LOG from the content of the env variable :
// LOG_FILE : string containing the filename to use for LOG, use DISCARD for no log, default is STDERR
func GetLogWriterFromEnvOrPanic(defaultLogName string) io.Writer {
	logFileName := defaultLogName
	val, exist := os.LookupEnv("LOG_FILE")
	if exist {
		logFileName = val
	}
	if utf8.RuneCountInString(logFileName) < 5 {
		panic(fmt.Sprintf("💥💥 error env LOG_FILE filename should contain at least %d characters (got %d).",
			5, utf8.RuneCountInString(val)))
	}
	switch logFileName {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	case "DISCARD":
		return io.Discard
	default:
		// Open the file with append, create, and write permissions.
		// The 0644 permission allows the owner to read/write and others to read.
		file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			// Return an error if the file cannot be opened (e.g., due to permissions).
			panic(fmt.Sprintf("💥💥 ERROR: LOG_FILE %q could not be open : %v", logFileName, err))
		}
		return file
	}
}

// requestHandler captures and logs all details of an incoming HTTP request.
func requestHandler(l *log.Logger) http.HandlerFunc {
	handlerName := "requestHandler"
	l.Printf(initCallMsg, handlerName)
	return func(w http.ResponseWriter, r *http.Request) {
		guid := strings.ToUpper(xid.New().String())
		msgRequest := fmt.Sprintf("'%s %s', %d bytes, from %s", r.Method, r.URL.String(), r.ContentLength, r.RemoteAddr)
		l.Printf("## ----- New Request %s ----- ##", guid)
		l.Printf("%s\tRequest : \t%s", guid, msgRequest)

		// Log request headers
		l.Printf("Headers:")
		for name, values := range r.Header {
			for _, value := range values {
				l.Printf("\t%s: %s", name, value)
			}
		}
		myRequest := RequestLog{
			Id:            guid,
			IsoDateTime:   time.Now().UTC(),
			Protocol:      r.Proto,
			Method:        r.Method,
			Url:           r.URL.String(),
			ContentLength: r.ContentLength,
			IpClient:      r.RemoteAddr,
			Headers:       r.Header,
			Body:          "",
		}
		var responseBody []byte
		var err error
		// Read and log the request body if present
		if r.Body != nil {
			responseBody, err = io.ReadAll(r.Body)
			if err != nil {
				l.Printf("Error reading request body: %v", err)
			} else if len(responseBody) > 0 {
				myRequest.Body = string(responseBody)
				l.Printf("%s\tBody: %s", guid, myRequest.Body)
			}

		}
		jsonRequest, err := json.Marshal(myRequest)
		if err != nil {
			l.Printf("💥💥 json.marshal failed. Error: %v", err)
		}
		fmt.Println(string(jsonRequest))
		// Respond to the client to confirm receipt (optional, but good for testing)
		w.Header().Set(HeaderContentType, MIMEAppJSONCharsetUTF8)
		w.Write([]byte("OK\n"))
		l.Printf("## ----- End Request %s ----- ##\n", guid)
	}
}

// Helper function to pretty print bytes (not used directly in this example, but useful for binary data)
func dumpBytes(data []byte) string {
	var sb strings.Builder
	for i, b := range data {
		if i > 0 && i%16 == 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("%02x ", b))
	}
	return sb.String()
}

func main() {

	l := log.New(GetLogWriterFromEnvOrPanic(defaultLogName), fmt.Sprintf("%s, ", version.APP), log.Ldate|log.Ltime|log.Lshortfile)
	l.Printf("🚀🚀 Starting App: %s, version: %s, build: %s", version.APP, version.VERSION, version.BuildStamp)

	myServerMux := http.NewServeMux()
	listenAddress := fmt.Sprintf("%s:%d", defaultListenIP, getPortFromEnvOrPanic(defaultPort))

	myServerMux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./favicon.ico")
	})

	// Register the handler function for all other paths
	myServerMux.Handle("/", requestHandler(l))

	server := http.Server{
		Addr:         listenAddress,
		Handler:      myServerMux,
		ErrorLog:     l,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	}

	l.Printf("Server starting on %s", listenAddress)
	if err := server.ListenAndServe(); err != nil {
		l.Fatalf("💥💥 Server failed to start: %v", err)
	}
}
