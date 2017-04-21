package binmock

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gexec"

	"strings"

	"io/ioutil"

	. "github.com/onsi/gomega"
)

//go:generate go-bindata -pkg binmock -o packaged_client.go client/
type Mock struct {
	Path                string
	identifier          string
	currentMappingIndex int

	mappings    []*MockMapping
	invocations []Invocation
}

func (mock *Mock) invoke(args, env []string) (int, string, string) {
	if mock.currentMappingIndex >= len(mock.mappings) {
		ginkgo.Fail(fmt.Sprintf("Too many calls to the mock! Last call with %v", args))
	}
	currentMapping := mock.mappings[mock.currentMappingIndex]
	mock.currentMappingIndex = mock.currentMappingIndex + 1
	if currentMapping.expectedArgs != nil {
		Expect(currentMapping.expectedArgs).To(Equal(args))
	}
	mock.invocations = append(mock.invocations, NewInvocation(args, env))
	return currentMapping.exitCode, currentMapping.stdout, currentMapping.stderr
}

func getSourceFile() string {
	data, err := Asset("client/main.go")
	Expect(err).NotTo(HaveOccurred())

	tempFile, err := ioutil.TempFile("", "go-bindata-client")
	Expect(err).NotTo(HaveOccurred())

	_, err = tempFile.Write(data)
	Expect(err).NotTo(HaveOccurred())
	Expect(tempFile.Close()).To(Succeed())

	sourceFilePath := tempFile.Name() + ".go"
	Expect(os.Rename(tempFile.Name(), sourceFilePath)).To(Succeed())

	return sourceFilePath
}

func NewBinMock(name string) *Mock {
	server := CurrentServer()

	identifier := strconv.FormatInt(time.Now().UnixNano(), 10)
	clientPath := getSourceFile()
	binaryPath, err := gexec.Build(clientPath, "-ldflags", "-X main.serverUrl=0.0.0.0:5555 -X main.identifier="+identifier)
	Expect(err).ToNot(HaveOccurred())
	Expect(os.Remove(clientPath)).To(Succeed())

	mock := &Mock{identifier: identifier, Path: binaryPath}

	server.Monitor(mock)
	return mock
}

type MockMapping struct {
	expectedArgs []string

	exitCode int
	stdout   string
	stderr   string
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

func (mapping *MockMapping) WillPrintToStdOut(out string) *MockMapping {
	mapping.stdout = out
	return mapping
}

func (mapping *MockMapping) WillPrintToStdErr(err string) *MockMapping {
	mapping.stderr = err
	return mapping
}

func (mapping *MockMapping) WillExitWith(exitCode int) *MockMapping {
	mapping.exitCode = exitCode
	return mapping
}

type Invocation struct {
	args []string
	env  map[string]string
}

func NewInvocation(args, env []string) Invocation {
	return Invocation{
		args: args,
		env:  parseEnv(env),
	}
}

func parseEnv(envVars []string) map[string]string {
	parsedVars := map[string]string{}

	for _, v := range envVars {
		parts := strings.Split(v, "=")

		parsedVars[parts[0]] = parts[1]
	}
	return parsedVars
}

func (invocation Invocation) Args() []string {
	return invocation.args
}

func (invocation Invocation) Env() map[string]string {
	return invocation.env
}

func (mock *Mock) Invocations() []Invocation {
	return mock.invocations
}
