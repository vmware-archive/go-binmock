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

	mappings    []*InvocationStub
	invocations []Invocation
}

type FailHandler func(message string, callerSkip ...int)

// Creates a new binMock
func NewBinMock(failHandler FailHandler) *Mock {
	server := getCurrentServer()

	identifier := strconv.FormatInt(time.Now().UnixNano(), 10)
	binaryPath, err := buildBinary(identifier, server.listener.Addr().String())
	if err != nil {
		failHandler(fmt.Sprintf("cant build binary %v", err))
	}

	mock := &Mock{identifier: identifier, Path: binaryPath, failHandler: failHandler}

	server.monitor(mock)
	return mock
}

func (mock *Mock) invoke(args, env, stdin []string) (int, string, string) {
	if mock.currentMappingIndex >= len(mock.mappings) {
		mock.failHandler(fmt.Sprintf("Too many calls to the mock! Last call with %v", args))
		return 1, "", ""
	}
	currentMapping := mock.mappings[mock.currentMappingIndex]
	mock.currentMappingIndex = mock.currentMappingIndex + 1
	if currentMapping.expectedArgs != nil && !reflect.DeepEqual(currentMapping.expectedArgs, args) {
		mock.failHandler(fmt.Sprintf("Expected %v to equal %v", args, currentMapping.expectedArgs))
		return 1, "", ""
	}
	mock.invocations = append(mock.invocations, newInvocation(args, env, stdin))
	return currentMapping.exitCode, currentMapping.stdout, currentMapping.stderr
}

// Sets up a stub for a possible invocation of the mock, accepting any arguments
func (mock *Mock) WhenCalled() *InvocationStub {
	return mock.createMapping(&InvocationStub{})
}

// Sets up a stub for a possible invocation of the mock, with specific arguments
// If args don't match the actual arguments to the mock then it fails
func (mock *Mock) WhenCalledWith(args ...string) *InvocationStub {
	invocation := &InvocationStub{}
	invocation.expectedArgs = args
	return mock.createMapping(invocation)
}

func (mock *Mock) createMapping(mapping *InvocationStub) *InvocationStub {
	mock.mappings = append(mock.mappings, mapping)
	return mapping
}

// Invocations returns the list of invocations of the mock till now
func (mock *Mock) Invocations() []Invocation {
	return mock.invocations
}
