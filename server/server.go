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

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"

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
	/*
		metrics.APIRequests.With(prometheus.Labels{"url": request.RequestURI}).Inc()
		startTime := time.Now()
		nextHandler.ServeHTTP(writer, request)
		duration := time.Since(startTime)
		metrics.APIResponsesTime.With(prometheus.Labels{"url": request.RequestURI}).Observe(float64(duration.Microseconds()))
	*/
}

// LogRequest - middleware for loging requests
func (server *HTTPServer) LogRequest(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			logRequestHandler(writer, request, nextHandler)
		})
}

func (server *HTTPServer) mainEndpoint(writer http.ResponseWriter, _ *http.Request) {
	responses.SendResponse(writer, responses.BuildOkResponse())
}

// Initialize perform the server initialization
func (server *HTTPServer) Initialize(address string) http.Handler {
	log.Println("Initializing HTTP server at", address)

	router := mux.NewRouter().StrictSlash(true)
	router.Use(server.LogRequest)

	// common REST API endpoints
	router.HandleFunc("/", server.mainEndpoint).Methods(http.MethodGet)

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
