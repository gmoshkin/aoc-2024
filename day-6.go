package main;

import (
    "os"
    // "strings"
    // "strconv"
    "slices"
    "errors"
    "fmt"
    // "time"
    "io"
)

func (ThisStructsMethodsAreSolutionFunctions) Solution6_1() {
    var patrol_map PatrolMap
    // Bad design! We silently read a file only if it exists without telling
    // anything to the user. I do this because it's simpler that way to work in
    // the framework that I've constructed for myself, but real programs should
    // not do this.
    if f, err := os.Open("patrol_map.txt"); err == nil {
        patrol_map = read_guard_patrol_map(f)
        patrol_map.do_step()
        patrol_map.dump(os.Stdout)
        patrol_map.dump_to_file("patrol_map.txt")
    } else {
        patrol_map = read_guard_patrol_map(os.Stdin)
        guard_i_was := patrol_map.guard_i
        for patrol_map.guard_i != -1 {
            patrol_map.do_step()
            assert(guard_i_was != patrol_map.guard_i, "")
        }
        patrol_map.dump(os.Stdout)
    }

    result := 0
    for _, cell := range patrol_map.cells {
        if cell == 'X' {
            result += 1
        }
    }

    fmt.Printf("result: %d\n", result)
}

type PatrolMap struct {
    cells []byte
    line_length int
    // height is deduced from len(cells) / width
    guard_i int
}

func read_guard_patrol_map(reader io.Reader) PatrolMap {
    data, err := io.ReadAll(reader)
    assert(err == nil, "%v", err)

    // +1 because this '\n' is also part of the line
    line_length := 1 + slices.Index(data, '\n')
    if data[len(data) - 1] != '\n' {
        data = append(data, '\n')
    }
    assert(len(data) % line_length == 0, "%v, %v", len(data), line_length)

    guard_i := slices.IndexFunc(data, func(b byte) bool {
        return b == '^' || b == '>' || b == 'v' || b == '<'
    })

    return PatrolMap{
        cells: data,
        line_length: line_length,
        guard_i: guard_i,
    }
}

func (m PatrolMap) dump_to_file(path string) {
    f, err := os.Create(path)
    assert(err == nil, "%v", err)

    m.dump(f)
}

func (m PatrolMap) dump(w io.Writer) {
    err := write_all_data(w, m.cells)
    assert(err == nil, "%v", err)
}

func (m *PatrolMap) do_step() {
    if m.guard_i == -1 { return }
    guard_direction := m.cells[m.guard_i]
    m.cells[m.guard_i] = 'X'

    for {
        new_guard_i := -1
        guard_x := m.guard_i % m.line_length
        switch guard_direction {
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
                if guard_x == m.line_length - 1 {
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
            m.guard_i = -1
            // New position is outside map bounds
            return
        }

        if m.cells[new_guard_i] == '#' {
            // Obstacle in front of the guard
            switch guard_direction {
                case '^': guard_direction = '>'
                case '>': guard_direction = 'v'
                case 'v': guard_direction = '<'
                case '<': guard_direction = '^'
            }

            // Can't move yet, must make sure the new direction is also clear,
            // do another iteration of the loop
            continue
        }

        // Path is clear, move the guard and update the map representation
        m.guard_i = new_guard_i
        m.cells[m.guard_i] = guard_direction
        break
    }
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
