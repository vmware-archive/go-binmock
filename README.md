# go-binmock

A library to mock interactions with collaborating executables. It works by creating a fake executable that can be orchestrated by the test. The path of the test executable can be injected into the system under test.


### Usage

Create a mock binary and getting its path.

```golang
mockMonit = binmock.NewBinMock(ginkgo.Fail)
monitPath := mockMonit.Path

thingToTest := NewThingToTest(monitPath)
```

Setup expected interactions with the binary

```golang
mockMonit.WhenCalledWith("start", "all").WillExitWith(0)
mockMonit.WhenCalledWith("summary").WillPrintToStdOut(output)
mockMonit.WhenCalledWith("summary").WillPrintToStdOut(output).WillPrintToStdErr("Noooo!").WillExitWith(1)
```

Assert on the interactions with the binary, after the fact

```golang
Expect(mockPGDump.Invocations()).To(HaveLen(1))
Expect(mockPGDump.Invocations()[0].Args()).To(Equal([]string{"dbname"}))
Expect(mockPGDump.Invocations()[0].Env()).To(HaveKeyWithValue("PGPASS", "we are going to build a wall"))
```
