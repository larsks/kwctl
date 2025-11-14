package formatters

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type (
	TableFormatter struct {
		tw table.Writer
	}
)

func RowFromList(values []string) table.Row {
	row := make(table.Row, len(values))
	for i, s := range values {
		row[i] = s
	}

	return row
}

func NewTableFormatter(headers []string) *TableFormatter {
	formatter := TableFormatter{
		tw: table.NewWriter(),
	}
	formatter.tw.SetStyle(table.StyleLight)
	formatter.tw.AppendHeader(RowFromList(headers))

	return &formatter
}

func (f *TableFormatter) Update(rows [][]string) {
	for _, row := range rows {
		f.tw.AppendRow(RowFromList(row))
		f.tw.AppendSeparator()
	}
}

func (f *TableFormatter) Render(out *os.File) {
	if out == nil {
		out = os.Stdout
	}
	f.tw.SetOutputMirror(out)
	f.tw.Render()
}
