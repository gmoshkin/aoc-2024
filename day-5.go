package main;

import (
    "os"
    "bufio"
    "strings"
    "strconv"
    "slices"
    "fmt"
    "time"
    "io"
)

func (ThisStructsMethodsAreSolutionFunctions) Solution5_1() {
    info := read_page_number_updates(os.Stdin)
    result := 0

    for _, updates_record := range info.updates_to_be {
        valid_update_record := true

    CheckUpdateRecord:
        for i, current := range updates_record {
            updates_so_far := updates_record[:i]
            must_go_after := info.must_go_after[current]
            for _, before := range updates_so_far {
                if slices.Contains(must_go_after, before) {
                    valid_update_record = false
                    break CheckUpdateRecord
                }
            }
        }

        if valid_update_record {
            n_updates := len(updates_record)
            result += updates_record[n_updates / 2]
        }
    }

    fmt.Printf("result: %d\n", result)
}

type PageNumberUpdates struct {
    // A mapping from int to array of int
    must_go_after map[int][]int
    // Array of arrays of int
    updates_to_be [][]int
}

func read_page_number_updates(reader io.Reader) PageNumberUpdates {
    scanner := bufio.NewScanner(os.Stdin)

    // Key goes before value
    must_go_after := make(map[int][]int)
    var updates_to_be [][]int

    const (
        ParsingRules = iota
        ParsingUpdates
    )
    state := ParsingRules

    t0 := time.Now()

    line_no := 0
    for scanner.Scan() {
        line_no += 1

        line := scanner.Text()
        if len(line) == 0 {
            assert(len(must_go_after) > 0, "%v", line_no)
            state = ParsingUpdates
            continue
        }

        // This is a common pattern when writing different kinds of parsers.
        // It's a state machine where inputs are handled differently based on
        // what state the machine is in, and when handling inputs the state may
        // change. This state machine is very simple, consisting of only 2
        // states. So in this case it would be equivalent to just having a
        // boolean variable `dependencies_parsed` or something, but I decided to
        // do a more elaborate thing just to showcase what a more robust
        // solution looks like.
        switch state {
            case ParsingRules:
                parts := strings.Split(line, "|")
                assert(len(parts) == 2, "%v", len(parts))

                left, err := strconv.Atoi(parts[0])
                assert(err == nil, "%v", err)

                right, err := strconv.Atoi(parts[1])
                assert(err == nil, "%v", err)

                was := must_go_after[left]
                must_go_after[left] = append(was, right)

            case ParsingUpdates:
                parts := strings.Split(line, ",")
                assert(len(parts) % 2 == 1, "%v", len(parts))

                var updates_record []int
                for _, part := range parts {
                    current, err := strconv.Atoi(part)
                    assert(err == nil, "%v: %v", part, err)

                    updates_record = append(updates_record, current)
                }

                updates_to_be = append(updates_to_be, updates_record)
        }
    }

    // When your program is doing something complicated and time consuming you
    // at some point want to start measuring the performance. For example to
    // understand what parts of your program are slow so you know where to
    // optimize, or to just understand how your program works.
    //
    // The simplest way to do this is to just check the time in 2 points of the
    // program. You want to make sure you're using the so-called "monotonic"
    // clock, which is different from the calendar or system clock. The system
    // clock in general may go backwards, for example if you change your system
    // time manually in the settings, or when there's day light savings, or the
    // system may even automatically adjust the time if it detects that your
    // clock is drifting. In any case such clocks are not applicable to
    // evaluating elapsed time. Fortunately computers provide also a
    // monotonically inrceasing clock which never goes back.
    //
    // Golang's `time.Now()` retunrs such monotonic time.
    logf("parsing input took %v", time.Now().Sub(t0))

    return PageNumberUpdates{
        must_go_after: must_go_after,
        updates_to_be: updates_to_be,
    }
}
