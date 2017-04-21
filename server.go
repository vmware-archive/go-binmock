package binmock

import (
	"encoding/json"
	"net/http"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Server struct {
	mocks map[string]*Mock
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
	go http.ListenAndServe("0.0.0.0:5555", http.HandlerFunc(server.Serve))
}

type InvocationRequest struct {
	Id   string
	Args []string
	Env  []string
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
	defer ginkgo.GinkgoRecover()
	jsonInvocationRequest := InvocationRequest{}
	decodingError := json.NewDecoder(req.Body).Decode(&jsonInvocationRequest)
	Expect(decodingError).NotTo(HaveOccurred())

	currentMock := server.mocks[jsonInvocationRequest.Id]
	jsonInvocationResponse := NewInvocationResponse(currentMock.invoke(jsonInvocationRequest.Args, jsonInvocationRequest.Env))
	json.NewEncoder(resp).Encode(jsonInvocationResponse)
}

func (server *Server) Monitor(mock *Mock) {
	server.mocks[mock.identifier] = mock
}
