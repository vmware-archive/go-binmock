package binmock

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gexec"

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

func (mock *Mock) invoke(args []string) (int, string, string) {
	if mock.currentMappingIndex >= len(mock.mappings) {
		ginkgo.Fail(fmt.Sprintf("Too many calls to the mock! Last call with %v", args))
	}
	currentMapping := mock.mappings[mock.currentMappingIndex]
	mock.currentMappingIndex = mock.currentMappingIndex + 1
	Expect(currentMapping.expectedArgs).To(Equal(args))
	mock.invocations = append(mock.invocations, Invocation{args: args})
	return currentMapping.exitCode, currentMapping.stdout, currentMapping.stderr
}

func getSourceFile(sourcePath string) string {
	data, err := Asset("client/main.go")
	Expect(err).NotTo(HaveOccurred())

	pathInProject := path.Join(sourcePath, "bin_mock_client.go")

	file, err := os.Create(pathInProject)
	Expect(err).NotTo(HaveOccurred())

	_, err = file.Write(data)
	Expect(err).NotTo(HaveOccurred())
	Expect(file.Close()).To(Succeed())

	return pathInProject
}

func NewBinMock(name, path string) *Mock {
	server := CurrentServer()

	identifier := strconv.FormatInt(time.Now().UnixNano(), 10)
	clientPath := getSourceFile(path)
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

func (mock *Mock) WhenCalledWith(args ...string) *MockMapping {
	invocation := &MockMapping{}
	invocation.expectedArgs = args
	mock.mappings = append(mock.mappings, invocation)
	return invocation
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
}

func (invocation Invocation) Args() []string {
	return invocation.args
}

func (mock *Mock) Invocations() []Invocation {
	return mock.invocations
}
