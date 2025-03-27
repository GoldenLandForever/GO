package InitLib1

import (
	_ "awesomeProject/InitLib2"
	"fmt"
)

func init() {
	fmt.Println("lib1")
}
