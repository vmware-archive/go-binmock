package binmock

type InvocationStub struct {
	expectedArgs []string

	exitCode int
	stdout   string
	stderr   string
}

// WillPrintToStdOut sets up what the mock will print to standard out on invocation
func (stub *InvocationStub) WillPrintToStdOut(out string) *InvocationStub {
	stub.stdout = out
	return stub
}

// WillPrintToStdErr sets up what the mock will print to standard error on invocation
func (stub *InvocationStub) WillPrintToStdErr(err string) *InvocationStub {
	stub.stderr = err
	return stub
}

// WillExitWith sets up the exit code of the mock invocation
func (stub *InvocationStub) WillExitWith(exitCode int) *InvocationStub {
	stub.exitCode = exitCode
	return stub
}
