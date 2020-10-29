package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("d:/data/iris.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Println("hello, world!")

}
