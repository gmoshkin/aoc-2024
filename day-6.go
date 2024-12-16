package main

import (
	"os"
	"strings"

	// "strings"
	// "strconv"
	"errors"
	"fmt"
	"slices"

	// "time"
	"io"
)

type WhichPart int
const (
    Part1 = WhichPart(iota)
    Part2
)

func (ThisStructsMethodsAreSolutionFunctions) Solution6_1() {
    day6_generic(Part1)
}

func (ThisStructsMethodsAreSolutionFunctions) Solution6_2() {
    day6_generic(Part2)
}

func day6_generic(which_part WhichPart) {
    var m PatrolMap
    // Bad design! We silently read a file only if it exists without telling
    // anything to the user. I do this because it's simpler that way to work in
    // the framework that I've constructed for myself, but real programs should
    // not do this.
    if f, err := os.Open("patrol_map.txt"); err == nil {
        defer f.Close()
        m = read_guard_patrol_map(f)
        m.do_step(which_part)
        m.dump(os.Stdout)
        m.dump_to_file("patrol_map.txt")
    } else {
        m = read_guard_patrol_map(os.Stdin)
        guard_i_was := m.guard_i
        for m.guard_i != -1 {
            m.do_step(which_part)
            assert(guard_i_was != m.guard_i, "")
            guard_i_was = m.guard_i
        }
        m.dump(os.Stdout)
    }

    result := 0
    switch which_part {
        case Part1:
            for _, cell := range m.cells {
                if cell == 'X' { result += 1 }
            }
            fmt.Printf("total cells covered: %d\n", result)
        case Part2:
            fmt.Println("potential loop generation obstacles:", m.n_potential_loops)
    }

}

type DirectionTrace byte
const (
    TraceUp    DirectionTrace = 0b0001
    TraceRight DirectionTrace = 0b0010
    TraceLeft  DirectionTrace = 0b0100
    TraceDown  DirectionTrace = 0b1000
)

type PatrolMap struct {
    cells []byte
    line_length int
    // height is deduced from len(cells) / width
    guard_i int
    guard_direction byte

    trace []DirectionTrace
    n_potential_loops int
}

func read_guard_patrol_map(reader io.Reader) PatrolMap {
    data, err := io.ReadAll(reader)
    assert(err == nil, "%v", err)
    assert(len(data) > 0, "")

    var m PatrolMap

    // +1 because this '\n' is also part of the line
    m.line_length = 1 + slices.Index(data, '\n')
    if data[len(data) - 1] != '\n' {
        data = append(data, '\n')
    }

    //
    // Figure a patch for the guard's current position
    //
    last_line_start := len(data) - 1
    // Last character is '\n', we made sure of it
    last_line_start -= 1
    for {
        if data[last_line_start - 1] == '\n' { break }
        last_line_start -= 1
    }
    last_line := string(data[last_line_start:])
    trace_under_guard := byte('.')
    if strings.HasPrefix(last_line, "trace under guard: ") {
        data = data[:last_line_start]
        n := len("trace under guard: ")
        trace_under_guard = last_line[n]
    }

    assert(len(data) % m.line_length == 0, "%v, %v", len(data), m.line_length)

    m.cells = data
    m.guard_i = slices.IndexFunc(m.cells, func(b byte) bool {
        return b == '^' || b == '>' || b == 'v' || b == '<'
    })
    if m.guard_i != -1 {
        m.guard_direction = m.cells[m.guard_i]
        m.cells[m.guard_i] = trace_under_guard
    }

    m.trace = make([]DirectionTrace, len(m.cells))

    return m
}

func (m PatrolMap) dump_to_file(path string) {
    f, err := os.Create(path)
    assert(err == nil, "%v", err)
    defer f.Close()

    m.dump(f)
}

func (m PatrolMap) patch(trace_under_guard byte) {
    m.cells[m.guard_i] = trace_under_guard
}

func (m PatrolMap) dump(w io.Writer) {
    var trace_under_guard byte = '?'
    if m.guard_i != -1 {
        trace_under_guard = m.cells[m.guard_i]
        m.cells[m.guard_i] = m.guard_direction
        defer m.patch(trace_under_guard)
    }

    err := write_all_data(w, m.cells)
    assert(err == nil, "%v", err)

    // By the way, this interface is retarded. Fprintf returns the number of
    // bytes written, and it's possible that the call doesn't actually write the
    // whole buffer to the file, but there's no way we can figure this out over
    // here. Fprintf constructs the final buffer and the buffer size is known
    // only inside that function, so there's no point in them returning the
    // number of bytes, we can't do anything with it. It should realy just
    // either write everything or return error.
    _, err = fmt.Fprintf(w, "trace under guard: %c\n", trace_under_guard)
    assert(err == nil, "%v", err)
}

func (m *PatrolMap) do_step(which_part WhichPart) {
    switch which_part {
        case Part1: m.do_step_v1()
        case Part2: m.do_step_v2()
    }
}

func (m *PatrolMap) do_step_v1() {
    if m.guard_i == -1 { return }

    new_guard_i := -1
    for {
        guard_x := m.guard_i % m.line_length
        switch m.guard_direction {
            case '^':
                new_guard_i = m.guard_i - m.line_length
                if new_guard_i < 0 {
                    new_guard_i = -1
                }

            case 'v':
                new_guard_i = m.guard_i + m.line_length
                if new_guard_i > len(m.cells) {
                    new_guard_i = -1
                }

            case '>':
                if guard_x == m.line_length - 2 {
                    new_guard_i = -1
                } else {
                    new_guard_i = m.guard_i + 1
                }

            case '<':
                if guard_x == 0 {
                    new_guard_i = -1
                } else {
                    new_guard_i = m.guard_i - 1
                }
        }

        if new_guard_i == -1 {
            // New position is outside map bounds
            break
        }

        if m.cells[new_guard_i] == '#' {
            // Obstacle in front of the guard
            switch m.guard_direction {
                case '^': m.guard_direction = '>'
                case '>': m.guard_direction = 'v'
                case 'v': m.guard_direction = '<'
                case '<': m.guard_direction = '^'
            }

            // Can't move yet, must make sure the new direction is also clear,
            // do another iteration of the loop
            continue
        }

        // Path is clear, move the guard and update the map representation
        break
    }

    // Mark the seen places on the map
    m.cells[m.guard_i] = 'X'

    // Finally move the guard
    m.guard_i = new_guard_i
}

func (m *PatrolMap) do_step_v2() {
    if m.guard_i == -1 { return }

    new_guard_i := -1
    did_a_turn := false
    for {
        guard_x := m.guard_i % m.line_length
        switch m.guard_direction {
            case '^':
                new_guard_i = m.guard_i - m.line_length
                if new_guard_i < 0 {
                    new_guard_i = -1
                }

            case 'v':
                new_guard_i = m.guard_i + m.line_length
                if new_guard_i > len(m.cells) {
                    new_guard_i = -1
                }

            case '>':
                if guard_x == m.line_length - 2 {
                    new_guard_i = -1
                } else {
                    new_guard_i = m.guard_i + 1
                }

            case '<':
                if guard_x == 0 {
                    new_guard_i = -1
                } else {
                    new_guard_i = m.guard_i - 1
                }
        }

        if new_guard_i == -1 {
            // New position is outside map bounds
            break
        }

        if m.cells[new_guard_i] == '#' {
            // Obstacle in front of the guard
            switch m.guard_direction {
                case '^': m.guard_direction = '>'
                case '>': m.guard_direction = 'v'
                case 'v': m.guard_direction = '<'
                case '<': m.guard_direction = '^'
            }

            // Can't move yet, must make sure the new direction is also clear,
            // do another iteration of the loop
            did_a_turn = true
            continue
        }

        // Path is clear, move the guard and update the map representation
        break
    }

    // If we go over a cell in which we have already been, but previously we
    // walked right (relative to our current direction), then if there was an
    // obstacle in front of us, we would go into an infinite loop
    trace_direction := char_to_trace_direction(m.guard_direction)
    direction_right := trace_direction.turn_right()
    if (m.trace[m.guard_i] & direction_right) != 0 {
        m.n_potential_loops += 1
    }
    // Mark in which direction we're passing the current cell, so that we can
    // use this info in future steps
    m.trace[m.guard_i] |= trace_direction

    // Pick a replacement character for the map
    var new_floor byte
    switch m.guard_direction {
        case '^': new_floor = '|'
        case 'v': new_floor = '|'
        case '>': new_floor = '-'
        case '<': new_floor = '-'
    }
    old_floor := m.cells[m.guard_i]
    switch old_floor {
        case '-': if new_floor == '|' { new_floor = '2' }
        case '|': if new_floor == '-' { new_floor = '2' }
        case '+': new_floor = '2'
    }
    if is_digit(old_floor) {
        new_floor = old_floor + 1
    } else if did_a_turn {
        new_floor = '+'
    }
    m.cells[m.guard_i] = new_floor

    // Finally move the guard
    m.guard_i = new_guard_i
}

func char_to_trace_direction(c byte) DirectionTrace {
    switch c {
        case '^': return TraceUp
        case 'v': return TraceDown
        case '>': return TraceRight
        case '<': return TraceLeft
    }
    panic("unreachable")
}

func (d DirectionTrace) turn_right() DirectionTrace {
    switch d {
        case TraceUp:    return TraceRight
        case TraceRight: return TraceDown
        case TraceDown:  return TraceLeft
        case TraceLeft:  return TraceUp
    }
    panic("unreachable")
}

func (m PatrolMap) index_to_position(i int) (x, y int) {
    y = i / m.line_length
    x = i - y * m.line_length
    assert(x != m.line_length - 1, "'\n' is not part of the map")

    return x, y
}

func (m PatrolMap) get(x, y int) byte {
    // +1 for an extra '\n' at the end
    return m.cells[y * m.line_length + x]
}

func file_exists(path string) bool {
    _, err := os.Stat(path)
    if err == nil { return true }
    // Not the best way to write this function. In reality it should return an
    // error instead of asserting.
    assert(errors.Is(err, os.ErrNotExist), "%v", err)
    return false
}

func file_reader(path string) (io.Reader, error) {
    f, err := os.Open(path)
    if err != nil { return nil, err }

    return f, nil
}

func write_all_data(w io.Writer, data []byte) error {
    n := 0
    for n < len(data) {
        var err error
        n, err = w.Write(data[n:])
        if err != nil { return err }
    }
    return nil
}

func min(a, b int) int {
    if a < b { return a }
    return b
}
