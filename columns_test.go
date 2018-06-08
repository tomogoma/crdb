package crdb_test

import (
	"github.com/tomogoma/crdb"
	"fmt"
)

func ExampleColDesc() {
	colDesc := crdb.ColDesc("name", "unit_price", "quantity")
	fmt.Print(colDesc)
	// output: name, unit_price, quantity
}

func ExampleColDescTbl() {
	tblName := "items"
	colDesc := crdb.ColDescTbl(tblName, "name", "unit_price", "quantity")
	fmt.Print(colDesc)
	// output: items.name, items.unit_price, items.quantity
}
