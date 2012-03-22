package equity

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Table struct {
	Data       []float64
	Stride     int
	Rows, Cols []string // row and column labels
}

func (this *Table) NumRows() int {
	return len(this.Data) / this.Stride
}

func (this *Table) At(i, j int) float64 {
	return this.Data[int(i)*this.Stride+j]
}

func (this *Table) Set(i, j int, v float64) {
	this.Data[int(i)*this.Stride+j] = v
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
	rsums := make([]float64, this.NumRows())
	csums := make([]float64, this.Stride)
	for i, val := range this.Data {
		rsums[int(i)/this.Stride] += val
		csums[int(i)%this.Stride] += val
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
	table := &Table{Data: make([]float64, this.NumRows()*(this.Stride+1)),
		Stride: this.Stride+1, Rows: this.Rows, Cols: append(this.Cols, "?")}
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
	table := &Table{Data: make([]float64, (this.NumRows()+1)*this.Stride),
		Stride: this.Stride, Rows: append(this.Rows, "?"), Cols: this.Cols}
	rsums, csums, total := this.Sums()
	for j := range this.Cols {
		table.Set(len(table.Rows)-1, j, csums[j] / total)
		for i := range this.Rows {
			table.Set(i, j, this.At(i, j) / rsums[i])
		}
	}
	return &PTable{table}
}

// Ppoker

func (this *Table) String() string {
	rsums, csums, total := this.Sums()
	width := len(fmt.Sprintf("%.0f", total))
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "%*s", 6, "")
	for _, col := range append(this.Cols, "totals") {
		fmt.Fprintf(b, "%-*s  ", width, col)
	}
	b.WriteString("\n")
	rows := this.NumRows()
	for line := 0; line < rows; line++ {
		fmt.Fprintf(b, "%-4s  ", this.Rows[line])
		for i := 0; i < this.Stride; i++ {
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

func ReadTable(file string) *Table {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(strings.TrimSpace(string(bs)), "\n")
	rlbs := strings.Fields(lines[0])
	table := &Table{Data: make([]float64, len(rlbs)*(len(lines)-1)),
		Stride: len(rlbs), Rows: make([]string, len(lines)-1), Cols: rlbs}
	for i, line := range lines[1:] {
		fs := strings.Fields(line)
		table.Rows[i] = fs[0]
		for j, val := range fs[1:] {
			//v, err := strconv.Atoi(val)
			v, err := strconv.ParseFloat(val, 64)
			if err != nil {
				panic(err)
			}
			table.Set(i, j, v)
		}
	}
	return table
}

func mul(a, b *Table) *Table {
	rows := a.NumRows()
	cols := b.Stride
	c := &Table{Data: make([]float64, rows*cols), Stride: cols}
	for m := 0; m < rows; m++ {
		for n := 0; n < cols; n++ {
			for i := 0; i < a.Stride; i++ {
				c.Set(m, n, c.At(m, n)+a.At(m, i)*b.At(i, n))
			}
		}
	}
	return c
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
	rows := this.NumRows()
	for line := 0; line < rows; line++ {
		fmt.Fprintf(b, "%-4s  ", this.Rows[line])
		for i := 0; i < this.Stride; i++ {
			fmt.Fprintf(b, "%*f  ", width, this.At(line, i))
		}
		b.WriteString("\n")
	}
	return b.String()
}
