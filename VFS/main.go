package main

import (
	"flag"
	"fmt"
	"time"
)

// main() В данном коде основная функция  main
// использует пакет  flag  для определения строковых флагов
// командной строки. Далее вызывается функция  StartScan(Path)
// которая выполняет дальнейшую логику.
func main() {
	starttime := time.Now()
	var Path string
	// Define string flags
	flag.StringVar(&Path, "ROOT", "", "Dir for scan")

	// Parse the command-line arguments
	flag.Parse()

	// Access the string flag values
	fmt.Println("Dir for scan: ", Path)

	errMessage := StartScan(Path)
	if errMessage != nil {
		fmt.Println(errMessage)
	}

	workTime := time.Since(starttime)
	fmt.Println("The program has worked: ", workTime, " seconds.")
}

//go run . --ROOT=/home/anton/go
