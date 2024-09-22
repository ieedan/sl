package table

import (
	"fmt"
	"log"
	"strings"

	"github.com/ieedan/sl/util"
)

var DEFAULT_OPTIONS = TableOptions{
	VerticalSeparator:     '│',
	HorizontalSeparator:   '─',
	Junction4Way:          '┼',
	JunctionLeftThreeWay:  '┤',
	JunctionRightThreeWay: '├',
	Cap:                   false,
	LeftPadding:           2,
}

type Table struct {
	header  *Row
	rows    []Row
	options TableOptions
}

func New(options TableOptions) Table {
	return Table{rows: []Row{}, options: options}
}

func (table *Table) AddHeader(columns ...string) {
	table.header = &Row{Columns: columns}
}

func (table *Table) AddRow(columns ...string) {
	if table.header != nil {
		// we always want the columns to have the same length to ensure the table prints correctly
		if len(table.header.Columns) != len(columns) {
			log.Fatalf("You provided a different number of columns than from the table header! Header: %v Provided: %v\n", len(table.header.Columns), len(columns))
		}
	}

	table.rows = append(table.rows, Row{
		Columns: columns,
		transform: func(str string) string {
			return str
		},
	})
}

func (table *Table) AddRowTransform(transform func(string) string, columns ...string) {
	if table.header != nil {
		// we always want the columns to have the same length to ensure the table prints correctly
		if len(table.header.Columns) != len(columns) {
			log.Fatalf("You provided a different number of columns than from the table header! Header: %v Provided: %v\n", len(table.header.Columns), len(columns))
		}
	}

	table.rows = append(table.rows, Row{Columns: columns, transform: transform})
}

func (table Table) String() string {
	columnMins := make(map[int]int)

	options := table.options

	if table.header != nil {
		for i, col := range table.header.Columns {
			colMin, ok := columnMins[i]

			if !ok {
				columnMins[i] = len(col)
				continue
			}

			if len(col) > colMin {
				columnMins[i] = len(col)
			}
		}
	}

	for _, row := range table.rows {
		for i, col := range row.Columns {
			colMin, ok := columnMins[i]

			if !ok {
				columnMins[i] = len(col)
				continue
			}

			if len(col) > colMin {
				columnMins[i] = len(col)
			}
		}
	}

	str := "\n"

	if table.header != nil {
		heading := string(options.VerticalSeparator)
		bottomBorder := string(options.JunctionRightThreeWay)

		for i, col := range table.header.Columns {
			colMin := columnMins[i]

			var rightSeparatorBorder rune

			if i+1 < len(table.header.Columns) {
				rightSeparatorBorder = options.Junction4Way
			} else {
				rightSeparatorBorder = options.JunctionLeftThreeWay
			}

			heading += fmt.Sprintf(" %v %v", util.PadRightMin(col, colMin), string(options.VerticalSeparator))
			bottomBorder += fmt.Sprintf("%v%v", strings.Repeat(string(options.HorizontalSeparator), colMin+2), string(rightSeparatorBorder))
		}

		str += util.LPad(heading, options.LeftPadding) + "\n" + util.LPad(bottomBorder, options.LeftPadding) + "\n"
	}

	for _, row := range table.rows {
		rowStr := string(options.VerticalSeparator)

		for i, col := range row.Columns {
			colMin := columnMins[i]

			rowStr += fmt.Sprintf(" %v %v", util.PadRightMin(col, colMin), string(options.VerticalSeparator))
		}

		str += util.LPad(row.transform(rowStr), options.LeftPadding) + "\n"
	}

	return str
}

type Row struct {
	Columns   []string
	transform func(string) string
}

type TableOptions struct {
	VerticalSeparator     rune
	HorizontalSeparator   rune
	Junction4Way          rune
	JunctionLeftThreeWay  rune
	JunctionRightThreeWay rune
	Cap                   bool
	LeftPadding           int
}
