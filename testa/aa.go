package main

import (
	_ "embed"
	"fmt"

	"github.com/eoe2005/goutils"
)

//go:embed t1.sql
var src string

//go:embed t2.sql
var desc string

func main() {
	fmt.Println(goutils.SqlMergeTable(src, desc))
}
