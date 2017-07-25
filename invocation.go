package binmock

import "strings"

// Invocation represents an invocation of the mock
type Invocation struct {
	args  []string
	env   map[string]string
	stdin []string
}

func newInvocation(args, env, stdin []string) Invocation {
	return Invocation{
		args:  args,
		env:   parseEnv(env),
		stdin: stdin,
	}
}

// Args represents the arguments passed to the mock when it was invoked
func (invocation Invocation) Args() []string {
	return invocation.args
}

// Env represents the environment at the time of invocation
func (invocation Invocation) Env() map[string]string {
	return invocation.env
}

// Stdin represents the standard input steam received by the mock as a slice of lines
func (invocation Invocation) Stdin() []string {
	return invocation.stdin
}

func parseEnv(envVars []string) map[string]string {
	parsedVars := map[string]string{}

	for _, v := range envVars {
		parts := strings.Split(v, "=")

		parsedVars[parts[0]] = parts[1]
	}
	return parsedVars
}
