package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bufio"
	"bytes"
	"os"
)

var identifier string
var serverUrl string

func main() {
	jsonInvocationRequest := InvocationRequest{}
	jsonInvocationRequest.Id = identifier
	jsonInvocationRequest.Args = os.Args[1:]
	jsonInvocationRequest.Env = os.Environ()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		jsonInvocationRequest.Stdin = append(jsonInvocationRequest.Stdin, scanner.Text())
	}

	buffer := bytes.NewBufferString("")
	if err := json.NewEncoder(buffer).Encode(jsonInvocationRequest); err != nil {
		panic(err)
	}

	response, err := http.Post("http://"+serverUrl, "", buffer)
	if err != nil {
		panic(err)
	}

	jsonInvocationResponse := InvocationResponse{}
	if err := json.NewDecoder(response.Body).Decode(&jsonInvocationResponse); err != nil {
		panic(err)
	}

	fmt.Fprint(os.Stdout, jsonInvocationResponse.Stdout)
	fmt.Fprint(os.Stderr, jsonInvocationResponse.Stderr)
	os.Exit(jsonInvocationResponse.ExitCode)
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
