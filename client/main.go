// Copyright (C) 2017-Present Pivotal Software, Inc. All rights reserved.
//
// This program and the accompanying materials are made available under
// the terms of the under the Apache License, Version 2.0 (the "License”);
// you may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

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
