package binmock

import (
	"fmt"
	"strconv"
	"time"

	"reflect"
)

//go:generate go-bindata -pkg binmock -o packaged_client.go client/
type Mock struct {
	Path                string
	identifier          string
	currentMappingIndex int
	failHandler         FailHandler

	mappings    []*MockMapping
	invocations []Invocation
}

type FailHandler func(message string, callerSkip ...int)

func NewBinMock(name string, failHandler FailHandler) *Mock {
	server := CurrentServer()

	identifier := strconv.FormatInt(time.Now().UnixNano(), 10)
	binaryPath, err := buildBinary(identifier)
	if err != nil {
		failHandler(fmt.Sprintf("cant build binary %v", err))
	}

	mock := &Mock{identifier: identifier, Path: binaryPath, failHandler: failHandler}

	server.Monitor(mock)
	return mock
}

func (mock *Mock) invoke(args, env, stdin []string) (int, string, string) {
	if mock.currentMappingIndex >= len(mock.mappings) {
		mock.failHandler(fmt.Sprintf("Too many calls to the mock! Last call with %v", args))
	}
	currentMapping := mock.mappings[mock.currentMappingIndex]
	mock.currentMappingIndex = mock.currentMappingIndex + 1
	if currentMapping.expectedArgs != nil && !reflect.DeepEqual(currentMapping.expectedArgs, args) {
		mock.failHandler(fmt.Sprintf("Expected %v to equal %v", args, currentMapping.expectedArgs))
	}
	mock.invocations = append(mock.invocations, NewInvocation(args, env, stdin))
	return currentMapping.exitCode, currentMapping.stdout, currentMapping.stderr
}

func (mock *Mock) WhenCalled() *MockMapping {
	return mock.createMapping(&MockMapping{})
}

func (mock *Mock) WhenCalledWith(args ...string) *MockMapping {
	invocation := &MockMapping{}
	invocation.expectedArgs = args
	return mock.createMapping(invocation)
}

func (mock *Mock) createMapping(mapping *MockMapping) *MockMapping {
	mock.mappings = append(mock.mappings, mapping)
	return mapping
}

func (mock *Mock) Invocations() []Invocation {
	return mock.invocations
}
