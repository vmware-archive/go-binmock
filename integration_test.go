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

package binmock_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"os/exec"

	"bytes"

	"github.com/onsi/gomega/gbytes"
	"github.com/pivotal-cf/go-binmock"
)

var _ = Describe("go-binmock", func() {
	var binMock *binmock.Mock
	var currentMockFailure *mockFailure

	BeforeEach(func() {
		currentMockFailure = &mockFailure{}
		binMock = binmock.NewBinMock(currentMockFailure.Fail)
	})

	Describe("when no stubs are defined", func() {
		It("fails", func() {
			session := RunCommand(binMock.Path)

			Expect(session).To(gexec.Exit(1))
			Expect(currentMockFailure.called).To(BeTrue())
			Expect(currentMockFailure.lastMessage).To(ContainSubstring("Too many calls to the mock"))
		})
	})

	Describe("when a stub with no expected args is defined", func() {
		It("returns the defined exit status", func() {
			binMock.WhenCalled().WillExitWith(42)

			session := RunCommand(binMock.Path)

			Expect(session).To(gexec.Exit(42))
		})

		It("returns the defined stdout", func() {
			binMock.WhenCalled().WillPrintToStdOut("stdout")

			session := RunCommand(binMock.Path)

			Expect(session.Out).To(gbytes.Say("stdout"))
		})

		It("returns the defined stderr", func() {
			binMock.WhenCalled().WillPrintToStdErr("stderr")

			session := RunCommand(binMock.Path)

			Expect(session.Err).To(gbytes.Say("stderr"))
		})

		It("ignores any additional args", func() {
			binMock.WhenCalled().WillExitWith(1)

			session := RunCommand(binMock.Path, "foo", "bar")

			Expect(session).To(gexec.Exit(1))
		})

		It("supports combination of return values", func() {
			binMock.WhenCalled().WillExitWith(42).WillPrintToStdErr("err").WillPrintToStdOut("out")

			session := RunCommand(binMock.Path)

			Expect(session).To(gexec.Exit(42))
			Expect(session.Out).To(gbytes.Say("out"))
			Expect(session.Err).To(gbytes.Say("err"))
		})
	})

	Describe("when a stub with expected args is defined", func() {
		Describe("and it's called with the right args", func() {
			It("matches the correct args", func() {
				binMock.WhenCalledWith("foo", "bar").WillExitWith(42).WillPrintToStdErr("err").WillPrintToStdOut("out")

				session := RunCommand(binMock.Path, "foo", "bar")

				Expect(session).To(gexec.Exit(42))
				Expect(session.Out).To(gbytes.Say("out"))
				Expect(session.Err).To(gbytes.Say("err"))
			})
		})

		Describe("and it's called with the wrong args", func() {
			It("fails", func() {
				binMock.WhenCalledWith("foo", "bar").WillExitWith(42).WillPrintToStdErr("err").WillPrintToStdOut("out")

				RunCommand(binMock.Path, "foo", "baz")

				Expect(currentMockFailure.called).To(BeTrue())
				Expect(currentMockFailure.lastMessage).To(ContainSubstring("Expected [foo baz] to equal [foo bar]"))
			})
		})
	})

	Describe("when invoked", func() {
		It("captures the args", func() {
			binMock.WhenCalled()
			RunCommand(binMock.Path, "foo", "bar")

			Expect(binMock.Invocations()[0].Args()).To(Equal([]string{"foo", "bar"}))
		})

		It("captures the environment", func() {
			binMock.WhenCalled()

			command := MakeCommand(binMock.Path, "foo", "bar")
			command.Env = append(command.Env, "foo=bar")
			StartCommand(command)

			Expect(binMock.Invocations()[0].Env()).To(HaveKeyWithValue("foo", "bar"))
		})

		It("captures stdin", func() {
			binMock.WhenCalled()

			command := MakeCommand(binMock.Path, "foo", "bar")
			command.Stdin = bytes.NewBufferString("stdin\nnextStdin")
			StartCommand(command)

			Expect(binMock.Invocations()[0].Stdin()).To(ConsistOf("stdin", "nextStdin"))
		})
	})

	Describe("when multiple stubs are defined", func() {
		Describe("and it gets invoked correctly", func() {
			It("returns the correct mocked values", func() {
				binMock.WhenCalledWith("one").WillPrintToStdOut("first")
				binMock.WhenCalledWith("two").WillPrintToStdOut("second")

				session := RunCommand(binMock.Path, "one")
				Expect(session.Out).To(gbytes.Say("first"))

				session = RunCommand(binMock.Path, "two")
				Expect(session.Out).To(gbytes.Say("second"))
			})

			It("captures all invocations", func() {
				binMock.WhenCalledWith("one")
				binMock.WhenCalledWith("two")

				RunCommand(binMock.Path, "one")
				RunCommand(binMock.Path, "two")

				Expect(binMock.Invocations()[0].Args()).To(ConsistOf("one"))
				Expect(binMock.Invocations()[1].Args()).To(ConsistOf("two"))
			})
		})

		Describe("and it gets invoked in a different order", func() {
			It("fails", func() {
				binMock.WhenCalledWith("one")
				binMock.WhenCalledWith("two")

				RunCommand(binMock.Path, "two")

				Expect(currentMockFailure.called).To(BeTrue())
				Expect(currentMockFailure.lastMessage).To(ContainSubstring("Expected [two] to equal [one]"))
			})
		})

		Describe("and more invocations occur", func() {
			It("fails", func() {
				binMock.WhenCalledWith("one")
				binMock.WhenCalledWith("two")

				RunCommand(binMock.Path, "one")
				RunCommand(binMock.Path, "two")
				RunCommand(binMock.Path, "three")

				Expect(currentMockFailure.called).To(BeTrue())
				Expect(currentMockFailure.lastMessage).To(ContainSubstring("Too many calls to the mock"))
			})
		})
	})

	Describe("when multiple mock binaries are created", func() {
		It("returns the response from the correct mock", func() {
			firstMock := binmock.NewBinMock(currentMockFailure.Fail)
			secondMock := binmock.NewBinMock(currentMockFailure.Fail)

			firstMock.WhenCalled().WillPrintToStdOut("first")
			firstMock.WhenCalled().WillPrintToStdOut("first again")
			secondMock.WhenCalled().WillPrintToStdOut("second")
			secondMock.WhenCalled().WillPrintToStdOut("second again")

			Expect(RunCommand(firstMock.Path).Out).To(gbytes.Say("first"))
			Expect(RunCommand(secondMock.Path).Out).To(gbytes.Say("second"))
			Expect(RunCommand(firstMock.Path).Out).To(gbytes.Say("first again"))
			Expect(RunCommand(secondMock.Path).Out).To(gbytes.Say("second again"))
		})
	})
})

type mockFailure struct {
	lastMessage string
	called      bool
}

func (m *mockFailure) Fail(message string, callerSkip ...int) {
	m.called = true
	m.lastMessage = message
}

func RunCommand(binPath string, args ...string) *gexec.Session {
	cmd := MakeCommand(binPath, args...)
	return StartCommand(cmd)
}

func StartCommand(cmd *exec.Cmd) *gexec.Session {
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gexec.Exit())
	return session
}

func MakeCommand(binPath string, args ...string) *exec.Cmd {
	return exec.Command(binPath, args...)
}
