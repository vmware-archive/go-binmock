package binmock

type MockMapping struct {
	expectedArgs []string

	exitCode int
	stdout   string
	stderr   string
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
