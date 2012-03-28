package equity

import (
	"testing"
	"fmt"
)

const distFile = "actionDistributions/AD1.txt"

func TestReadTable (test *testing.T) {
	table1, err := ReadTable(distFile)
	if err != nil {
		test.Fatal(err)
	}
	fmt.Println(table1)
	fmt.Println(table1.P())
	fmt.Println(table1.PCR())

	table2, err := ReadTable("snowday.txt")
	if err != nil {
		test.Fatal(err)
	}
	fmt.Println(table2)
	fmt.Println(table2.P())
	fmt.Println(table2.PCR())
}
