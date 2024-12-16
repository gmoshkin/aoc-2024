package main

import (
	"os"
	// "strings"

	// "strings"
	// "strconv"
	// "errors"
	"fmt"
	"slices"

	// "time"
	"io"
)

func (ThisStructsMethodsAreSolutionFunctions) Solution15_1() {
    var m RobotAndBoxesMap
    // Bad design! We silently read a file only if it exists without telling
    // anything to the user. I do this because it's simpler that way to work in
    // the framework that I've constructed for myself, but real programs should
    // not do this.
    if f, err := os.Open("robot_and_boxes_map.txt"); err == nil {
        defer f.Close()
        m.parse(f)
        m.do_step()
        m.dump(os.Stdout)
        m.dump_to_file("robot_and_boxes_map.txt")
    } else {
        m.parse(os.Stdin)
        previous_move := m.current_move
        for m.current_move < len(m.moves) {
            m.do_step()
            assert(previous_move != m.current_move, "")
            previous_move = m.current_move
        }
        m.dump(os.Stdout)
    }

    result := 0
    result = m.count_up_gps()
    fmt.Println("result:", result)
}

type RobotAndBoxesMap struct {
    cells []byte
    line_length int
    robot_i int
    current_move int
    moves []byte
}

func (m *RobotAndBoxesMap) do_step() {
    move_char := byte('|')
    for {
        if m.current_move >= len(m.moves) { return }
        move_char = m.moves[m.current_move]
        if is_white_space(move_char) {
            m.current_move += 1
            continue
        }
        break
    }
    assert(move_char != '|', "")
    m.current_move += 1

    next_i := -1
    offset := 0
    switch move_char {
        case '<':
            offset = -1
        case '^':
            offset = -m.line_length
        case '>':
            offset = +1
        case 'v':
            offset = +m.line_length
    }
    next_i = m.robot_i + offset
    m.maybe_push_boxes(next_i, offset)
    if m.cells[next_i] == '.' {
        m.cells[m.robot_i] = '.'
        m.robot_i = next_i
        m.cells[m.robot_i] = '@'
    }
}

func (m *RobotAndBoxesMap) maybe_push_boxes(box_i int, offset int) {
    // box_x := box_i % m.line_length
    // box_y := box_i / m.line_length
    // logf("box_i: %d, box_x: %d, box_y: %d, cell: '%c'", box_i, box_x, box_y, m.cells[box_i])
    if m.cells[box_i] != 'O' { return }

    i := box_i
    for m.cells[i] == 'O' { i += offset }

    // logf("m.cells[%v] = '%c'", i, m.cells[i])
    if m.cells[i] == '#' { return }

    m.cells[box_i] = '.'
    m.cells[i] = 'O'
}

func (m *RobotAndBoxesMap) parse(reader io.Reader) {
    data, err := io.ReadAll(reader)
    assert(err == nil, "%v", err)
    assert(len(data) > 0, "")

    // +1 because this '\n' is also part of the line
    line_end := slices.Index(data, '\n')
    m.line_length = 1 + line_end
    cursor := m.line_length
    for {
        line_end = slices.Index(data[cursor:], '\n') + cursor
        assert(line_end != -1, "")
        if cursor == line_end {
            // Found empty line
            break
        }
        cursor = line_end + 1
    }

    m.cells = data[:cursor]
    assert(m.cells[len(m.cells) - 1] == '\n', "")

    m.robot_i = slices.Index(m.cells, '@')
    assert(m.robot_i > 0, "")

    for is_white_space(data[cursor]) { cursor += 1 }

    data = data[cursor:]

    m.current_move = 0
    i := slices.Index(data, '|')
    if i >= 0 {
        m.moves = append(m.moves, data[:i]...)
        m.moves = append(m.moves, data[i+1:]...)
        m.current_move = i
    } else {
        m.moves = data
    }
}

func (m *RobotAndBoxesMap) dump_to_file(path string) {
    f, err := os.Create(path)
    assert(err == nil, "%v", err)
    defer f.Close()

    m.dump(f)
}

func (m *RobotAndBoxesMap) dump(w io.Writer) {
    err := write_all_data(w, m.cells)
    assert(err == nil, "%v", err)

    err = write_all_data(w, []byte("\n"))
    assert(err == nil, "%v", err)

    err = write_all_data(w, m.moves[:m.current_move])
    assert(err == nil, "%v", err)

    err = write_all_data(w, []byte("|"))
    assert(err == nil, "%v", err)

    if m.current_move < len(m.moves) {
        err = write_all_data(w, m.moves[m.current_move:])
        assert(err == nil, "%v", err)
    }
}

func (m *RobotAndBoxesMap) count_up_gps() int {
    score := 0
    for i, cell := range m.cells {
        if cell != 'O' { continue }

        x := i % m.line_length
        y := i / m.line_length
        score += x + 100 * y
    }

    return score
}

func is_white_space(c byte) bool {
    return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
