package output

import (
	"fmt"
	"strings"
)

type Table struct {
	columns []string
	rows    [][]string
}

func NewTable() *Table {
	return &Table{}
}

func (t *Table) AddColumn(column string) {
	t.columns = append(t.columns, column)
}

func (t *Table) AddColumns(columns ...string) {
	t.columns = append(t.columns, columns...)
}

func (t *Table) AddRow(row ...string) error {
	if len(row) != len(t.columns) {
		return fmt.Errorf("row length (%d) does not match column length (%d)", len(row), len(t.columns))
	}
	t.rows = append(t.rows, row)

	return nil
}

func (t *Table) Print() {
	widths := make([]int, len(t.columns))

	for _, row := range t.rows {
		for j, column := range row {
			if len(column) > widths[j] {
				widths[j] = len(column)
			}
		}
	}

	for i, column := range t.columns {
		if len(column) > widths[i] {
			widths[i] = len(column)
		}
	}

	for i, column := range t.columns {
		fmt.Printf("%-*s ", widths[i], column)
	}
	fmt.Println()

	for i := range t.columns {

		fmt.Printf("%-*s ", widths[i], strings.Repeat("-", widths[i]))
	}
	fmt.Println()

	for _, row := range t.rows {
		for i, column := range row {
			fmt.Printf("%-*s ", widths[i], column)
		}
		fmt.Println()
	}

	fmt.Println()
}
