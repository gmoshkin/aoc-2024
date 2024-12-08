package main;

import (
    "os"
    "bufio"
    "strings"
    "fmt"
    "strconv"
    "sort"
)

func (ThisStructsMethodsAreSolutionFunctions) Solution1_1() {
    // The Scanner struct implements methods for "parsing" text. We're only
    // going to be using it to read the text line by line. But you could also
    // use `string.Split` or something, it doesn't really matter in this case.
    //
    // And `os.Stdin` is the standard input of the program. So you can type or
    // copy paste the input data right into the terminal once the program has
    // started. Or you can pipe the input from the text file.
    scanner := bufio.NewScanner(os.Stdin)

    // A couple of arrays. One for the left column and one for the right one.
    var left, right [] int

    // Line number is only used for debugging.
    line_no := 0
    for scanner.Scan() {
        line_no += 1
        // Get next line of text. Bad name, because the `Scanner` provides a
        // generic interface to scanning text, and uses the scanning line by
        // line by default.
        line := scanner.Text()

        line = strings.TrimSpace(line)

        // Actual AOC inputs don't have any empty lines, but when doing manual
        // testing it's more convenient to ignore them rather than panicking.
        if len(line) == 0 {
            continue
        }

        // Split the string into pieces separated by whitespace.
        parts := strings.Fields(line)
        if len(parts) != 2 {
            // This error handling style is pretty simple and is actually pretty
            // good when starting the project. If you have a big complicated
            // program you will probably want to do somethin more complicated,
            // like returning this error up the call chain so that it can be
            // handled more gracefully, but for a simple program like that it's
            // perfeclty fine to just panic as soon as a problem happens.
            //
            // It's always a good idea to add as much useful information in the
            // error message as possible. For example here we provide the line
            // number which makes it easy to debug the problem, because the user
            // knows where to go look. We also additionall print out the whole
            // line and the separated parts. This is probably an overkill in
            // this case, but I didn't pay too much though to it and just added
            // them all by force of habbit. You can always remove stuff later,
            // but just having them at the begginig may come in handy in a pinch.
            panic(fmt.Sprintf("line #%d: expected 2 numbers per line, got %d in `%s`", line_no, len(parts), line))
        }

        // Atoi is another bad name. The meaning is "A to I", 'A' standing for "alpha" or
        // "letter", and 'I' standing for "integer". It really should be named `string_to_int`.
        // The name `atoi` comes from the C language, where it was added a long long
        // time ago (probably some time in the 1970s), because in those days it was very
        // important to have short names because of low disk capacity, low CPU clock rates, etc.
        //
        // Today it really doesn't make a lot of sense to have such short names
        // especially for library functions which will be used by a lot of people.
        l, err := strconv.Atoi(parts[0])
        // Not the best error handling over here, in real life probably want to
        // add more info to the error message
        if err != nil { panic(err) }
        // This shit is straight up retarded. `append` doesn't modify the array
        // in place, but rather returns another one which we must assign to the
        // variable explicitly. There's been literally 0 times in my life where
        // I wanted this behavior.
        // They must have been smoking something weird to come up with this bullshit.
        left = append(left, l)

        r, err := strconv.Atoi(parts[1])
        if err != nil { panic(err) }
        right = append(right, r)
    }

    //
    // The actual solution, nothing fancy
    //
    sort.Ints(left)
    sort.Ints(right)

    total_distance := 0
    // Golang's "for range" situation is kind of a mess, you have to remember
    // that the first variable is always the index and the second one is always
    // the value. And there's almost no flexibility, no reverse iteration, no
    // starting from the second element, etc.
    for i, l := range left {
        r := right[i]
        total_distance += difference(l, r)
    }

    fmt.Println(total_distance)
}

func (ThisStructsMethodsAreSolutionFunctions) Solution1_2() {
    scanner := bufio.NewScanner(os.Stdin)
    // This code is mostly copy-pasted from above and modified to suit the other
    // task. I didn't care about reusing any code (extracting common functions)
    // because I know that this code will not live for long, it's a one-time thing.
    //
    // But even in real life I don't always jump to refactoring out common
    // pieces. It's often fine to have the same piece of code copy-pasted 2 or 3
    // times before you want to start refactoring. Firstly if the code is simple
    // it's often better to have it spelled out, rather than seeing a function
    // call. Secondly it's sometimes hard to factor out code without making a
    // huge mess. And thirdly you never know how long your code we live for.
    // Maybe you type it 3 times and never return to it, it just does the job,
    // or you even don't need it anymore, so refactoring it would just be a
    // waste of time.
    //
    // Anyway you'll need some experience writing and refactoring code before
    // you get that intuition about when it's worth it and when it's not to
    // refactor. Just keep in mind that sometimes it ok to copy-paste.
    var left [] int
    // In this case I don't even collect the right column into an array, because
    // the task only needs the counts. I just use the hash table (also sometimes
    // called "hash map", and python calls them "dictionaries") to count up the values.
    right := make(map[int]int)

    line_no := 0
    for scanner.Scan() {
        line_no += 1
        line := scanner.Text()
        line = strings.TrimSpace(line)
        if len(line) == 0 {
            continue
        }

        parts := strings.Fields(line)
        if len(parts) != 2 {
            panic(fmt.Sprintf("line #%d: expected 2 numbers per line, got %d in `%s`", line_no, len(parts), line))
        }

        l, err := strconv.Atoi(parts[0])
        if err != nil {
            panic(err)
        }
        left = append(left, l)

        r, err := strconv.Atoi(parts[1])
        if err != nil {
            panic(err)
        }

        // I'm not sure, but `right[r]++` may also work. But even if it does,
        // it's not the best idea to always use the one-liners, often it's
        // better to spell out the code. For example if in the future we need to
        // debug the code it's really easy to add `log.Println(count)` after
        // this line, but if it was a one-liner, we would need to refactor it.
        count := right[r]
        right[r] = count + 1
    }

    //
    // The actual solution, nothing fancy
    //
    total_similarity := 0
    // Back to the messiness of `for range`, here go forces us to ignore the
    // index variable by naming it `_`
    for _, l := range left {
        total_similarity += l * right[l]
    }

    fmt.Println(total_similarity)
}

func difference(a, b int) int {
    if a < b {
        return b - a
    } else {
        return a - b
    }
}
