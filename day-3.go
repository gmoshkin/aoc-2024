package main;

import (
    "os"
    "bufio"
    "strings"
    "fmt"
    "strconv"
)

func (ThisStructsMethodsAreSolutionFunctions) Solution3_1() {
    scanner := bufio.NewScanner(os.Stdin)
    line_no := 0

    total := 0
    for scanner.Scan() {
        line_no += 1
        line := scanner.Text()
        line = strings.TrimSpace(line)
        if len(line) == 0 {
            continue
        }

        tail := line
        for len(tail) > 0 {
            i := strings.Index(tail, "mul(")
            if i == -1 { break }
            tail = tail[i:]

            res, ok := parse_and_eval_mul_expr(&tail)
            if ok {
                total += res
            }
        }
    }

    fmt.Println(total)
}

func (ThisStructsMethodsAreSolutionFunctions) Solution3_2() {
    scanner := bufio.NewScanner(os.Stdin)
    line_no := 0

    total := 0
    multiply_enabled := true
    for scanner.Scan() {
        line_no += 1
        line := scanner.Text()
        line = strings.TrimSpace(line)
        if len(line) == 0 {
            continue
        }

        tail := line
        for len(tail) > 0 {
            switch tail[0] {
                case 'd':
                    do, ok := parse_do_or_dont(&tail)
                    if ok {
                        multiply_enabled = do
                    }
                case 'm':
                    res, ok := parse_and_eval_mul_expr(&tail)
                    if ok && multiply_enabled {
                        total += res
                    }
                default:
                    // Advance the string by one character so that parsing continues with
                    // the next character
                    chop(&tail, 1)
            }
        }
    }

    fmt.Println(total)
}

// This parser is pretty shit, it does the job for the task at hand, but in
// real programs you want something more robust.
func parse_and_eval_mul_expr(text *string) (int, bool) {
    if !strings.HasPrefix(*text, "mul(") {
        // Advance the string by one character so that parsing continues with
        // the next character
        chop(text, 1)
        return 0, false
    }
    chop(text, 4)

    i := 0
    for is_digit((*text)[i]) { i += 1 }
    if i == 0 {
        return 0, false
    }

    left_text := chop(text, i)
    left, err := strconv.Atoi(left_text)
    if err != nil {
        return 0, false
    }

    comma := chop(text, 1)
    if comma != "," {
        return 0, false
    }

    i = 0
    for is_digit((*text)[i]) { i += 1 }
    if i == 0 {
        return 0, false
    }

    right_text := chop(text, i)
    right, err := strconv.Atoi(right_text)
    if err != nil {
        return 0, false
    }

    closing_paren := chop(text, 1)
    if closing_paren != ")" {
        return 0, false
    }

    return left * right, true
}

func parse_do_or_dont(text *string) (bool, bool) {
    if strings.HasPrefix(*text, "do()") {
        chop(text, len("do()"))
        return true, true
    }

    if strings.HasPrefix(*text, "don't()") {
        chop(text, len("don't()"))
        return false, true
    }

    // Advance the string by one character so that parsing continues with
    // the next character
    chop(text, 1)

    return false, false
}

/// Splits the string at given number of "runes", returns the head and advances
/// `text` to the tail.
func chop(text *string, count int) string {
    count = min(len(*text), count)

    head := (*text)[:count]
    *text = (*text)[count:]
    return head
}

func max(a, b int) int {
    if a > b { return a }
    return b
}

func is_digit(c byte) bool {
    return c >= '0' && c <= '9'
}
