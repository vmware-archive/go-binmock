package binmock

import "strings"

type Invocation struct {
	args  []string
	env   map[string]string
	stdin []string
}

func NewInvocation(args, env, stdin []string) Invocation {
	return Invocation{
		args:  args,
		env:   parseEnv(env),
		stdin: stdin,
	}
}

func (invocation Invocation) Args() []string {
	return invocation.args
}

func (invocation Invocation) Env() map[string]string {
	return invocation.env
}

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
