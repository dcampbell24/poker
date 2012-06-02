package bayes

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Table struct {
	Data       []float64
	Rows, Cols []string // row and column labels
}

func (this *Table) At(i, j int) float64 {
	return this.Data[int(i)*len(this.Cols)+j]
}

func (this *Table) Set(i, j int, v float64) {
	this.Data[int(i)*len(this.Cols)+j] = v
}

func sum(nums []float64) float64 {
	sum := 0.0
	for _, num := range nums {
		sum += num
	}
	return sum
}

// Sums returns the the row and column sums as slices and the grand total.
func (this *Table) Sums() ([]float64, []float64, float64) {
	rsums := make([]float64, len(this.Rows))
	csums := make([]float64, len(this.Cols))
	for i, val := range this.Data {
		rsums[int(i)/len(this.Cols)] += val
		csums[int(i)%len(this.Cols)] += val
	}

	total := 0.0
	if len(rsums) <= len(csums) {
		total = sum(rsums)
	} else {
		total = sum(csums)
	}
	return rsums, csums, total
}

// P returns a table showing the probabilities of the row field given the column
// field.
func (this *Table) P() *PTable {
	table := &Table{Data: make([]float64, len(this.Rows)*(len(this.Cols)+1)),
		Rows: this.Rows, Cols: append(this.Cols, "?")}
	rsums, csums, total := this.Sums()
	for i := range this.Rows {
		table.Set(i, len(table.Cols)-1, rsums[i] / total)
		for j := range this.Cols {
			table.Set(i, j, this.At(i, j) / csums[j])
		}
	}
	return &PTable{table}
}

// PCR returns a table showing the probabilities of the column field given the
// row field.
func (this *Table) PCR() *PTable {
	table := &Table{Data: make([]float64, (len(this.Rows)+1)*len(this.Cols)),
		Rows: append(this.Rows, "?"), Cols: this.Cols}
	rsums, csums, total := this.Sums()
	for j := range this.Cols {
		table.Set(len(table.Rows)-1, j, csums[j] / total)
		for i := range this.Rows {
			table.Set(i, j, this.At(i, j) / rsums[i])
		}
	}
	return &PTable{table}
}

func (this *Table) String() string {
	rsums, csums, total := this.Sums()
	width := len(fmt.Sprintf("%.0f", total))
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "%*s", 6, "")
	for _, col := range append(this.Cols, "totals") {
		fmt.Fprintf(b, "%-*s  ", width, col)
	}
	b.WriteString("\n")
	for line := range this.Rows {
		fmt.Fprintf(b, "%-4s  ", this.Rows[line])
		for i := range this.Cols {
			fmt.Fprintf(b, "%*.0f  ", width, this.At(line, i))
		}
		fmt.Fprintf(b, "%*.0f\n", width, rsums[line])
	}
	fmt.Fprintf(b, "%*s", 6, "totals")
	for _, col := range csums {
		fmt.Fprintf(b, "%*.0f  ", width, col)
	}
	fmt.Fprintf(b, "%*.0f\n", width, total)
	return b.String()
}

func ReadTable(file string) (*Table, error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(bs)), "\n")
	rlbs := strings.Fields(lines[0])
	table := &Table{Data: make([]float64, len(rlbs)*(len(lines)-1)),
		Rows: make([]string, len(lines)-1), Cols: rlbs}
	for i, line := range lines[1:] {
		fs := strings.Fields(line)
		table.Rows[i] = fs[0]
		for j, val := range fs[1:] {
			v, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
			table.Set(i, j, v)
		}
	}
	return table, nil
}

type PTable struct {
	*Table
}

func (this *PTable) String() string {
	_, _, total := this.Sums()
	width := len(fmt.Sprintf("%f", total))
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "%*s", 6, "")
	for _, col := range this.Cols {
		fmt.Fprintf(b, "%-*s  ", width, col)
	}
	b.WriteString("\n")
	for line := range this.Rows {
		fmt.Fprintf(b, "%-4s  ", this.Rows[line])
		for i := range this.Cols {
			fmt.Fprintf(b, "%*f  ", width, this.At(line, i))
		}
		b.WriteString("\n")
	}
	return b.String()
}
