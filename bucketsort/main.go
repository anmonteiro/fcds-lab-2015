package main

import (
	"os"
	"fmt"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func parseArgs(args []string) {
	var usage = func() {
        fmt.Fprintf(os.Stderr, "usage: %s [inputFile] [outputFile]\n", os.Args[0])
	}

	if len(args) != 2 { 
		usage()
		os.Exit(1)
	}
}

func main() {
	args := os.Args[1:]
	parseArgs(args)

	// Open input & output files
	fin, err := os.Open(args[0])
	check(err)
	defer fin.Close()

	fout, err := os.Create(args[1])
	check(err)
	defer fout.Close()

	// Read number of lines
	var nLines uint32
	_, err = fmt.Fscanf(fin, "%d", &nLines)
	check(err)

	// Perform the sort
	r := StreamBucketSort(fin, 8, int(nLines))
	
	// Write to output file
	for _, v := range r {
		fmt.Fprintln(fout, v)
	}
}