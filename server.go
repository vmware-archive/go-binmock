package binmock

import (
	"encoding/json"
	"net"
	"net/http"
)

type server struct {
	mocks map[string]*Mock
	*http.Server
	listener net.Listener
}

var currentServer *server

func getCurrentServer() *server {
	if currentServer == nil {
		currentServer = &server{
			mocks: map[string]*Mock{},
		}
		currentServer.Start()
	}
	return currentServer
}

func (server *server) Start() {
	server.Server = &http.Server{Addr: ":0", Handler: http.HandlerFunc(server.Serve)}
	server.listener, _ = net.Listen("tcp", "127.0.0.1:0")
	go server.Server.Serve(server.listener)
}

type invocationRequest struct {
	Id    string
	Args  []string
	Env   []string
	Stdin []string
}

type invocationResponse struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func newInvocationResponse(exitCode int, stdout, stderr string) invocationResponse {
	return invocationResponse{
		ExitCode: exitCode,
		Stdout:   stdout,
		Stderr:   stderr,
	}
}

func (server *server) Serve(resp http.ResponseWriter, req *http.Request) {
	invocationRequest := invocationRequest{}
	json.NewDecoder(req.Body).Decode(&invocationRequest)
	currentMock := server.mocks[invocationRequest.Id]
	invocationResponse := newInvocationResponse(currentMock.invoke(invocationRequest.Args, invocationRequest.Env, invocationRequest.Stdin))
	json.NewEncoder(resp).Encode(invocationResponse)
}

func (server *server) Monitor(mock *Mock) {
	server.mocks[mock.identifier] = mock
}
