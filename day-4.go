package main;

import (
    "os"
    "bufio"
    // "strings"
    "fmt"
)

func (ThisStructsMethodsAreSolutionFunctions) Solution4_1() {
    scanner := bufio.NewScanner(os.Stdin)

    var lines []string
    line_length := 0
    for scanner.Scan() {
        line := scanner.Text()
        if line_length == 0 {
            line_length = len(line)
        }
        if line_length != len(line) {
            panic("all lines must be the same length")
        }

        lines = append(lines, line)
    }

    row_size = line_length
    mask_size := len(lines) * line_length
    mask = make([]bool, mask_size, mask_size)
    for i := range mask {
        mask[i] = false
    }

    counts := make(map[rune]int)

    // horizontal forward & backward
    for line_i := range lines {
        counts['→'] += count_generic_direction(lines, line_i, 0, 0, +1, "XMAS")
        counts['←'] += count_generic_direction(lines, line_i, 0, 0, +1, "SAMX")
        // The above is equivalent to this. I only use my function so that
        // `mark_match` is called for horizontal matches.
        // counts['→'] += strings.Count(lines[line_i], "XMAS")
        // counts['←'] += strings.Count(lines[line_i], "SAMX")
    }

    // vertical forward & backward
    for column_i := range line_length {
        counts['↓'] += count_vertically(lines, column_i, "XMAS")
        counts['↑'] += count_vertically(lines, column_i, "SAMX")
    }

    // south-eastern diagonal
    //
    // This is the order:
    //    3 2 1 0
    //    ↘ ↘ ↘ ↘
    //  4 ↘ ↘ ↘ ↘
    //  5 ↘ ↘ ↘ ↘
    //  6 ↘ ↘ ↘ ↘
    //  7 ↘ ↘ ↘ ↘
    // first line
    column_i := line_length - 1
    line_i := 0
    for column_i > 0 {
        counts['↘'] += count_south_east_diagonal(lines, line_i, column_i, "XMAS")
        counts['↖'] += count_south_east_diagonal(lines, line_i, column_i, "SAMX")
        column_i -= 1
    }
    column_i = 0
    for line_i < len(lines) {
        counts['↘'] += count_south_east_diagonal(lines, line_i, column_i, "XMAS")
        counts['↖'] += count_south_east_diagonal(lines, line_i, column_i, "SAMX")
        line_i += 1
    }

    // north-eastern diagonal
    //
    // This is the order:
    //  0 ↗ ↗ ↗ ↗
    //  1 ↗ ↗ ↗ ↗
    //  2 ↗ ↗ ↗ ↗
    //  3 ↗ ↗ ↗ ↗
    //  4 ↗ ↗ ↗ ↗
    //      5 6 7
    line_i = 0
    column_i = 0
    for line_i < len(lines) {
        counts['↗'] += count_north_east_diagonal(lines, line_i, column_i, "XMAS")
        counts['↙'] += count_north_east_diagonal(lines, line_i, column_i, "SAMX")
        line_i += 1
    }
    line_i = len(lines) - 1
    column_i = 1
    for column_i < line_length {
        counts['↗'] += count_north_east_diagonal(lines, line_i, column_i, "XMAS")
        counts['↙'] += count_north_east_diagonal(lines, line_i, column_i, "SAMX")
        column_i += 1
    }

    // debug. I was going to delete this before committing, but decided to keep it...
    for line_i, line := range lines {
        for column_i, c := range line {
            if is_marked(line_i, column_i) {
                fmt.Printf("\x1b[31m%c\x1b[0m", c)
            } else {
                fmt.Printf("%c", c)
            }
        }
        fmt.Println()
    }

    total := 0
    for dir, count := range counts {
        logf("'%c': %v", dir, count)
        total += count
    }
    fmt.Println("result:", total)
}

func count_vertically(lines []string, column_i int, substring string) int {
    return count_generic_direction(lines, 0, column_i, +1, 0, substring)
}

// Goes from top to bottom, left to right starting from cell in line l0, column c0
// in the south-east direction
func count_south_east_diagonal(lines []string, l0, c0 int, substring string) int {
    return count_generic_direction(lines, l0, c0, +1, +1, substring)
}

// Goes from bottom to top, left to right starting from cell in line l0, column c0
// in the north-east direction
func count_north_east_diagonal(lines []string, l0, c0 int, substring string) int {
    return count_generic_direction(lines, l0, c0, -1, +1, substring)
}

func count_generic_direction(lines []string, l_start, c_start, l_step, c_step int, substring string) int {
    found_count := 0
    n_matched := 0
    li := l_start
    ci := c_start
    l_match_start, c_match_start := 0, 0 // debug
    // Some of these checks are redundant
    for li >= 0 && li < len(lines) && ci >= 0 && ci < len(lines[li]) {
        if lines[li][ci] == substring[n_matched] {
            if n_matched == 0 { l_match_start, c_match_start = li, ci } // debug
            n_matched += 1
            if n_matched == len(substring) {
                found_count += 1
                n_matched = 0
                mark_match(l_match_start, c_match_start, li, ci) // debug
            }

            li += l_step
            ci += c_step
        } else if n_matched == 0 {
            li += l_step
            ci += c_step
        } else {
            // Possible match failed. Must not move the cursor, because there
            // could be another match starting from this position, e.g.
            // XMXMAS  | n_matched == 2
            //   ^     | ci = 2
            n_matched = 0
        }
    }

    return found_count
}

// this is just for debugging
var mask []bool
var row_size int

// XXX: there's a bug in this function, it doesn't work for north-eastern
// diagonals for some reason. I thought there was a bug in
// `count_generic_direction`, but it turns out the bug is here. I found another
// bug using this function though, fixed it and passed the task. But this
// function still doesn't work for north-eastern diagonal...
func mark_match(l_start, c_start, l_end, c_end int) {
    l_step := l_end - l_start
    c_step := c_end - c_start
    if l_step != 0 { l_step /= abs(l_step) }
    if c_step != 0 { c_step /= abs(c_step) }

    li, ci := l_start, c_start
    for li >= 0 && li <= l_end && ci >= 0 && ci <= c_end {
        mark_cell(li, ci)
        li += l_step
        ci += c_step
    }
}

func mark_cell(li, ci int) {
    mask[li * row_size + ci] = true
}

func is_marked(li, ci int) bool {
    return mask[li * row_size + ci]
}

func (ThisStructsMethodsAreSolutionFunctions) Solution4_2() {
    scanner := bufio.NewScanner(os.Stdin)

    var lines []string
    line_length := 0
    for scanner.Scan() {
        line := scanner.Text()
        if line_length == 0 {
            line_length = len(line)
        }
        if line_length != len(line) {
            panic("all lines must be the same length")
        }

        lines = append(lines, line)
    }

    row_size = line_length
    mask_size := len(lines) * line_length
    mask = make([]bool, mask_size, mask_size)
    for i := range mask {
        mask[i] = false
    }

    total := 0
    for line_i := range lines[:len(lines)-2] {
        for column_i := range line_length-2 {
            if matches_x_mas(lines, line_i, column_i) {
                total += 1
            }
        }
    }

    // debug. I was going to delete this before committing, but decided to keep it...
    for line_i, line := range lines {
        for column_i, c := range line {
            if is_marked(line_i, column_i) {
                fmt.Printf("\x1b[31m%c\x1b[0m", c)
            } else {
                fmt.Printf("%c", c)
            }
        }
        fmt.Println()
    }

    fmt.Println("result:", total)
}

func matches_x_mas(lines []string, l, c int) bool {
    if lines[l + 1][c + 1] != 'A' { return false }
    if !((lines[l][c] == 'M' && lines[l + 2][c + 2] == 'S') ||
         (lines[l][c] == 'S' && lines[l + 2][c + 2] == 'M')) { return false }
    if !((lines[l][c + 2] == 'M' && lines[l + 2][c] == 'S') ||
         (lines[l][c + 2] == 'S' && lines[l + 2][c] == 'M')) { return false }

    mark_cell(l, c)
    mark_cell(l+1, c+1)
    mark_cell(l+2, c+2)
    mark_cell(l, c+2)
    mark_cell(l+2, c)

    return true
}
