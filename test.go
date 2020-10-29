package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
)

func main() {
	f, err := os.Open("d:/gop/goai/data/iris.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	irisDF := dataframe.ReadCSV(f)
	fmt.Println(irisDF)
}
