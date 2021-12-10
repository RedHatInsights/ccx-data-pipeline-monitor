/*
Copyright © 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/ccx-data-pipeline-monitor/packages/server/server.html

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/RedHatInsights/insights-operator-utils/responses"

	"github.com/RedHatInsights/ccx-data-pipeline-monitor/config"
)

// HTTPServer in an implementation of Server interface
type HTTPServer struct {
	Config config.ServerConfig
	Serv   *http.Server
}

// New constructs new implementation of Server interface
func New(config config.ServerConfig) *HTTPServer {
	return &HTTPServer{
		Config: config,
	}
}

func logRequestHandler(writer http.ResponseWriter, request *http.Request, nextHandler http.Handler) {
	log.Println("Request URI: " + request.RequestURI)
	log.Println("Request method: " + request.Method)
	nextHandler.ServeHTTP(writer, request)
}

func staticPage(filename string) func(writer http.ResponseWriter, request *http.Request) {
	log.Println("Serving static file", filename)
	return func(writer http.ResponseWriter, request *http.Request) {
		sendStaticPage(writer, filename)
	}
}

func sendStaticPage(writer http.ResponseWriter, filename string) {
	// disable "G304 (CWE-22): Potential file inclusion via variable"
	// #nosec G304
	body, err := ioutil.ReadFile(filename)
	if err == nil {
		writer.Header().Set("Server", "A Go Web Server")
		writer.Header().Set("Content-Type", getContentType(filename))
		_, err = fmt.Fprint(writer, string(body))
		if err != nil {
			log.Println("Error sending response body", err)
		}
	} else {
		writer.WriteHeader(http.StatusNotFound)
		notFoundResponse(writer)
	}
}

func getContentType(filename string) string {
	// TODO: to map
	if strings.HasSuffix(filename, ".html") {
		return "text/html"
	} else if strings.HasSuffix(filename, ".js") {
		return "application/javascript"
	} else if strings.HasSuffix(filename, ".css") {
		return "text/css"
	}
	return "text/html"
}

func writeResponse(writer http.ResponseWriter, message string) {
	_, err := fmt.Fprint(writer, message)
	if err != nil {
		log.Println("Error sending response", err)
	}
}

func notFoundResponse(writer http.ResponseWriter) {
	writeResponse(writer, "Not found!")
}

// LogRequest - middleware for loging requests
func (server *HTTPServer) LogRequest(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			logRequestHandler(writer, request, nextHandler)
		})
}

func (server *HTTPServer) mainEndpoint(writer http.ResponseWriter, _ *http.Request) {
	err := responses.SendOK(writer, responses.BuildOkResponse())
	if err != nil {
		log.Println("Error sending response in main endpoint", err)
	}
}

// Initialize perform the server initialization
func (server *HTTPServer) Initialize(address string) http.Handler {
	log.Println("Initializing HTTP server at", address)

	router := mux.NewRouter().StrictSlash(false)
	router.Use(server.LogRequest)

	// HTML etc.
	router.HandleFunc("/", staticPage("html/index.html")).Methods(http.MethodGet)
	router.HandleFunc("/bootstrap.min.css", staticPage("html/bootstrap.min.css"))
	router.HandleFunc("/bootstrap.min.js", staticPage("html/bootstrap.min.js"))
	router.HandleFunc("/ccx.css", staticPage("html/ccx.css"))

	// common REST API endpoints
	return router
}

// Start starts server
func (server *HTTPServer) Start() error {
	address := server.Config.Address
	log.Println("Starting HTTP server at", address)
	router := server.Initialize(address)
	server.Serv = &http.Server{Addr: address, Handler: router}

	err := server.Serv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Printf("Unable to start HTTP server %v", err)
		return err
	}

	return nil
}

// Stop stops server's execution
func (server *HTTPServer) Stop(ctx context.Context) error {
	return server.Serv.Shutdown(ctx)
}
