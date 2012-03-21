package equity

import (
	"testing"
	"fmt"
)

const distFile = "actionDistributions/AD1.txt"

func TestReadTable (_ *testing.T) {
	table1 := ReadTable(distFile)
	fmt.Println(table1)
	fmt.Println(table1.P())
	fmt.Println(table1.PCR())


	table2 := ReadTable("snowday.txt")
	fmt.Println(table2)
	fmt.Println(table2.P())
	fmt.Println(table2.PCR())

}
