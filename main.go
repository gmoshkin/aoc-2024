package main;

import (
    "os"
    "fmt"
    "reflect"
    "strings"
    "runtime"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "please provide a solution number")
        os.Exit(1)
    }

    solution_no := strings.TrimSpace(os.Args[1])

    // Here's some advanced magic, it's equivalent to something like this:
    //     switch "Solution" + solution_no {
    //         case "Solution1_1":
    //             dummy.Solution1_1();
    //         case "Solution1_2":
    //             dummy.Solution1_2();
    //         ...
    //     }
    // I just wrote it like this so I don't have to add another branch to that
    // switch statement every time I add another solution.
    //
    // You really don't need to know about this stuff to be solving Advent Of Code,
    // but if you're curious the `reflect` module gives you ability to inspect
    // stuff about the program at runtime. So here I lookup methods defined on
    // the struct by name and call them.
    type_ := reflect.TypeOf(dummy)
    method_name := "Solution" + solution_no
    method, ok := type_.MethodByName(method_name)
    if !ok {
        fmt.Fprintf(os.Stderr, "couldn't find solution named `%s`, which you requested\n", method_name)
        fmt.Fprintln(os.Stderr, "available solutions:")
        for i := range type_.NumMethod() {
            method := type_.Method(i)
            fmt.Fprintln(os.Stderr, "   ", method.Name)
        }
        os.Exit(1)
    }

    run_solution(method)
}

// Advanced topic:
// This is a special hack. Our solution functions don't actually logically
// accept any arguments, but because we want to do this `reflect` magic on them,
// they must be "methods" for us to be able to list them. So we add this extra
// struct which doesn't do anything expect that we implement methods for it. And
// all of it's methods will be our solutions to AOC tasks.
type ThisStructsMethodsAreSolutionFunctions struct {}
// As a result we will need to have an object of this type so that we can call
// the methods. We define a single global variable which will be used for that.
// It's a "dummy" because it doesn't do anything.
var dummy ThisStructsMethodsAreSolutionFunctions

func run_solution(method reflect.Method) {
    // Advanced topic:
    // Because golang is staically typed, every object must have a concrete type
    // during compilation. But `reflect` module does some magic stuff to make
    // dynamically typed objects (kinda like the python ones) with which we can
    // work without knowing their static type, the type will only be known at
    // run time. For example `method` is an object which can be called with some
    // arguments, but we don't know at compile time which arguments it expects.
    //
    // Well actually we do know in this case, because we create it, but in
    // general, when you see `method relfect.Method` theoretically it could be
    // any method with any list of arguments.
    //
    // So as a consequence we need to do extra actions to pass the arguments, we
    // first need to convert the statically typed object into a dynamically
    // typed one.
    type_erased_dummy := reflect.ValueOf(dummy)

    // And we also need to wrap it into an array, because in general methods
    // accept a list of arguments.
    type_erased_arguments := []reflect.Value{type_erased_dummy}

    // Finally we can call the method
    method.Func.Call(type_erased_arguments)
}

// Always a useful function to have for debugging. I think go has something like
// this builtin in `log` module, but it didn't work for me first try so I made
// my own.
func log(args ...any) {
    // This is the important part. It gives us the information about the "call
    // stack". runtime.Caller(0) would be this call itself, but we want info
    // about who called the `log` function, so we pass a `1`. The function
    // returns a bunch of stuff, but we care about the filename and the line
    // number.
    _, filename, line, _ := runtime.Caller(1)
    // Our log takes an aribtrary number of arguments because of `...` and to
    // pass the through to another function you also need to say `...`.
    message := fmt.Sprint(args...)
    // These `%s`, `%d` come from the C language. In C you had to specify the
    // type (%s is for string, %d is for decimal number) because the compiler
    // wouldn't tell the function what types the parameters are. But go doens't
    // have this problem, when you see `any` there's actually dynamic type
    // information inside it, so `fmt.Printf` knows what types of values you're
    // passing and your template string doesn't need to specify the types. I
    // think this template is equivalent to "[%v:%v] %v\n" or something, but I
    // wrote it like this by force of habbit
    fmt.Printf("[%s:%d] %s\n", filename, line, message)
}

func logf(format string, args ...any) {
    // Copy pasted because I'm lazy
    _, filename, line, _ := runtime.Caller(1)
    msg := fmt.Sprintf(format, args...)
    fmt.Printf("[%s:%d] %s\n", filename, line, msg)
}
