package errors_test

import (
	stderrors "errors"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

func ExampleNew() {
	err := errors.New("whoops")
	fmt.Println(err)

	// Output: whoops
}

func ExampleNew_printf() {
	err := errors.New("whoops")
	fmt.Printf("%+v", err)

	// Example output:
	// whoops
	// github.com/pkg/errors_test.ExampleNew_printf
	//         /home/dfc/src/github.com/pkg/errors/example_test.go:17
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/pkg/errors/_test/_testmain.go:106
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
}

func ExampleWithMessage() {
	cause := errors.New("whoops")
	err := errors.WithMessage(cause, "oh noes")
	fmt.Println(err)

	// Output: oh noes: whoops
}

func ExampleWrap() {
	cause := errors.New("New error")
	e1 := errors.Wrap(cause, "first wrap")
	e2 := errors.Wrap(e1, "second wrap")
	fmt.Printf("%v", e2)

	// Output:
	// second wrap: first wrap: New error
}

type myErr struct {
	msg string
}

func (e *myErr) Error() string { return "myErr" }

type myErr2 struct {
	msg string
}

func (e *myErr2) Error() string { return "myErr2" }

func ExampleWrap_errors_is_as() {
	eMine1 := &myErr{msg: "abc"}
	eMine2 := &myErr{msg: "bcd"}
	e1 := errors.Wrap(eMine1, "first wrap")
	e2 := errors.Wrap(e1, "second wrap")
	var unwrappedErr *myErr
	ok1 := stderrors.As(e2, &unwrappedErr)
	var unwrappedErr2 *myErr2
	ok2 := stderrors.As(e2, &unwrappedErr2)
	fmt.Printf("%v: %v: %v: %v: %v: %v: %v", e2,
		ok1, stderrors.Is(e2, eMine1), unwrappedErr == eMine1,
		ok2, stderrors.Is(e2, eMine2), unwrappedErr == eMine2,
	)

	// Output:
	// second wrap: first wrap: myErr: true: true: true: false: false: false
}

func ExampleWrap_printv() {
	cause := errors.New("New error")
	e1 := errors.Wrap(cause, "first wrap")
	e2 := errors.Wrap(e1, "second wrap")
	fmt.Printf("%+v", e2)

	// Example Output:
	// second wrap
	// github.com/pkg/errors_test.ExampleWrap_printv
	// 	 /home/fabstu/go/src/github.com/pkg/errors/example_test.go:104
	// first wrap
	// github.com/pkg/errors_test.ExampleWrap_printv
	// 	 /home/fabstu/go/src/github.com/pkg/errors/example_test.go:103
	// New error
	// github.com/pkg/errors_test.ExampleWrap_printv
	// 	 /home/fabstu/go/src/github.com/pkg/errors/example_test.go:102
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:62
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1118
	// main.main
	// 	_testmain.go:120
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:203
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1357
}

type myFormatterErr struct {
	msg string
}

func (e *myFormatterErr) Error() string { return e.msg }

func (e *myFormatterErr) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, "+v "+e.Error())
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func ExampleWrap_print_formatter() {
	cause := &myFormatterErr{msg: "myFormatterErr msg"}
	e1 := errors.Wrap(cause, "first wrap")
	e2 := errors.Wrap(e1, "second wrap")
	fmt.Printf("%v", e2)

	// Output:
	// second wrap: first wrap: myFormatterErr msg
}

func ExampleWrap_printv_formatter() {
	cause := &myFormatterErr{msg: "myFormatterErr msg"}
	e1 := errors.Wrap(cause, "first wrap")
	e2 := errors.Wrap(e1, "second wrap")
	fmt.Printf("%+v", e2)

	// Example Output:
	// second wrap
	// github.com/pkg/errors_test.ExampleWrap_printv_formatter
	// 	/home/fabstu/go/src/github.com/pkg/errors/example_test.go:189
	// first wrap
	// github.com/pkg/errors_test.ExampleWrap_printv_formatter
	// 	/home/fabstu/go/src/github.com/pkg/errors/example_test.go:188
	// +v myFormatterErr msg
}

func fn() error {
	e1 := errors.New("error")
	e2 := errors.Wrap(e1, "inner")
	e3 := errors.Wrap(e2, "middle")
	return errors.Wrap(e3, "outer")
}

func ExampleWrap_extended() {
	err := fn()
	fmt.Printf("%+v\n", err)

	// Example output:
	// error
	// github.com/pkg/errors_test.fn
	//         /home/dfc/src/github.com/pkg/errors/example_test.go:47
	// github.com/pkg/errors_test.ExampleCause_printf
	//         /home/dfc/src/github.com/pkg/errors/example_test.go:63
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/pkg/errors/_test/_testmain.go:104
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
	// github.com/pkg/errors_test.fn
	// 	  /home/dfc/src/github.com/pkg/errors/example_test.go:48: inner
	// github.com/pkg/errors_test.fn
	//        /home/dfc/src/github.com/pkg/errors/example_test.go:49: middle
	// github.com/pkg/errors_test.fn
	//      /home/dfc/src/github.com/pkg/errors/example_test.go:50: outer
}

func ExampleWrapf() {
	cause := errors.New("whoops")
	err := errors.Wrapf(cause, "oh noes #%d", 2)
	fmt.Println(err)

	// Output: oh noes #2: whoops
}

func ExampleErrorf_extended() {
	err := errors.Errorf("whoops: %s", "foo")
	fmt.Printf("%+v", err)

	// Example output:
	// whoops: foo
	// github.com/pkg/errors_test.ExampleErrorf
	//         /home/dfc/src/github.com/pkg/errors/example_test.go:101
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/pkg/errors/_test/_testmain.go:102
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
}

func ExampleCause_printf() {
	err := errors.Wrap(func() error {
		return func() error {
			return errors.Errorf("hello %s", fmt.Sprintf("world"))
		}()
	}(), "failed")

	fmt.Printf("%v", err)

	// Output: failed: hello world
}
