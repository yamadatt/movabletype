package main

import (
	"fmt"
	"os"

	"github.com/yamadatt/movabletype"
)

func main() {
	entries, _ := movabletype.Parse(os.Stdin)

	for _, e := range entries {
		fmt.Println(e.Date)
		fmt.Println(e.Title)
		fmt.Println(e.Image)
	}
}
