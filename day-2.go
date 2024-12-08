package main;

import (
    "os"
    "bufio"
    "strings"
    "fmt"
    "strconv"
)

func (ThisStructsMethodsAreSolutionFunctions) Solution2_1() {
    scanner := bufio.NewScanner(os.Stdin)
    line_no := 0

    safe_count := 0
Next_Line:
    for scanner.Scan() {
        line_no += 1
        line := scanner.Text()
        line = strings.TrimSpace(line)
        if len(line) == 0 {
            continue
        }

        parts := strings.Fields(line)

        if len(parts) < 2 {
            safe_count += 1
            continue
        }

        first, err := strconv.Atoi(parts[0])
        if err != nil { panic(err) }

        second, err := strconv.Atoi(parts[1])
        if err != nil { panic(err) }

        diff := second - first
        if diff == 0     { continue }
        if abs(diff) > 3 { continue }

        prev := second
        prev_diff := diff
        for i := 2; i < len(parts); i++ {
            next, err := strconv.Atoi(parts[i])
            if err != nil { panic(err) }

            diff := next - prev
            if diff == 0            { continue Next_Line }
            if abs(diff) > 3        { continue Next_Line }
            if diff * prev_diff < 0 { continue Next_Line }

            prev = next
        }

        safe_count += 1
    }

    fmt.Println(safe_count)
}

func (ThisStructsMethodsAreSolutionFunctions) Solution2_2() {
    scanner := bufio.NewScanner(os.Stdin)
    line_no := 0

    safe_count := 0

    for scanner.Scan() {
        line_no += 1
        line := scanner.Text()
        line = strings.TrimSpace(line)
        if len(line) == 0 {
            continue
        }

        parts := strings.Fields(line)

        if len(parts) < 2 {
            safe_count += 1
            continue
        }

        var readings [] int
        for _, part := range parts {
            reading, err := strconv.Atoi(part)
            if err != nil { panic(err) }
            readings = append(readings, reading)
        }

        safe := readings_are_safeish(readings)

        var verdict string
        if safe {
            safe_count += 1

            verdict = "safe"
        } else {
            verdict = "unsafe"
        }

        if !safe {
            fmt.Fprintf(os.Stderr, "`%s` %s\n", line, verdict)
        }
    }

    fmt.Println(safe_count)
}

func readings_are_safeish(readings []int) bool {
    prev := readings[0]
    prev_diff := 0
    had_unsafe := false
    for i := 1; i < len(readings); i++ {
        curr := readings[i]
        diff := curr - prev
        if !readings_pair_is_safe(diff, prev_diff) {
            had_unsafe = true
            break
        }
        prev_diff = diff
        prev = curr
    }

    if !had_unsafe {
        return true
    }

    for cutout := range readings {
        ia, ib := 0, 1
        if cutout == 0 { ia, ib = 1, 2 }
        prev := readings[ia]
        prev_diff := 0
        had_unsafe := false
        for ib < len(readings) {
            curr := readings[ib]
            diff := curr - prev
            if !readings_pair_is_safe(diff, prev_diff) {
                had_unsafe = true
                break;
            }
            prev_diff = diff
            prev = curr
            ib += 1
            if ib == cutout { ib += 1 }
        }

        if !had_unsafe {
            return true
        }
    }

    return false
}

func readings_pair_is_safe(diff int, prev_diff int) bool {
    if diff == 0            { return false }
    if abs(diff) > 3        { return false }
    if diff * prev_diff < 0 { return false }

    return true
}

func abs(a int) int {
    if a < 0 { return -a } else { return a }
}
