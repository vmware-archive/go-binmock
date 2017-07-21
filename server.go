package binmock

import (
	"encoding/json"
	"net"
	"net/http"
)

type Server struct {
	mocks map[string]*Mock
	*http.Server
	listener net.Listener
}

var currentServer *Server

func CurrentServer() *Server {
	if currentServer == nil {
		currentServer = &Server{
			mocks: map[string]*Mock{},
		}
		currentServer.Start()
	}
	return currentServer
}

func (server *Server) Start() {
	server.Server = &http.Server{Addr: ":0", Handler: http.HandlerFunc(server.Serve)}
	server.listener, _ = net.Listen("tcp", "127.0.0.1:0")
	go server.Server.Serve(server.listener)
}

type InvocationRequest struct {
	Id    string
	Args  []string
	Env   []string
	Stdin []string
}

type InvocationResponse struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func NewInvocationResponse(exitCode int, stdout, stderr string) InvocationResponse {
	return InvocationResponse{
		ExitCode: exitCode,
		Stdout:   stdout,
		Stderr:   stderr,
	}
}

func (server *Server) Serve(resp http.ResponseWriter, req *http.Request) {
	invocationRequest := InvocationRequest{}
	json.NewDecoder(req.Body).Decode(&invocationRequest)
	currentMock := server.mocks[invocationRequest.Id]
	invocationResponse := NewInvocationResponse(currentMock.invoke(invocationRequest.Args, invocationRequest.Env, invocationRequest.Stdin))
	json.NewEncoder(resp).Encode(invocationResponse)
}

func (server *Server) Monitor(mock *Mock) {
	server.mocks[mock.identifier] = mock
}
