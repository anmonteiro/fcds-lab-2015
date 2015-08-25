package main

import (
	"fmt"
	"os"
	"encoding/binary"
	"time"
	"math"
	"sync"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("usage: %s [input_file] [output_file]\n", os.Args[0]);
		os.Exit(1)
	}

	fin, err := os.Open(os.Args[1])
	check(err)
	defer fin.Close()

	fout, err := os.Create(os.Args[2])
	check(err)
	defer fout.Close()

	var size uint64

    err = binary.Read(fin, binary.LittleEndian, &size)
    check(err)

    err = binary.Write(fout, binary.LittleEndian, size)
	check(err)

	pixels := make([]uint32, size * size)

	start := time.Now()

	err = binary.Read(fin, binary.LittleEndian, pixels)

	duration := time.Since(start)
	fmt.Printf("time to read: %v\n", duration);

	var SQRT_2 float64 = math.Sqrt(2)
	var mid uint64

	start = time.Now()

	for s := size; s > 1; s /= 2 {
		var rwg, cwg sync.WaitGroup
		mid = s / 2

		// row-transformation
		rwg.Add(int(mid))
		for y := uint64(0); y < mid; y++ {
			go func(y uint64) {
				defer rwg.Done()
				for x := uint64(0); x < mid; x++ {
					a := pixels[y * size + x]
					d := a
					a = uint32(float64(a + pixels[y * size + (x + mid)]) / SQRT_2)
					d = uint32(float64(d - pixels[y * size + (x + mid)]) / SQRT_2)
					pixels[y * size + x] = a
					pixels[y * size + (x + mid)] = d
				}
			}(y)
		}
		rwg.Wait()

		// column-transformation
		cwg.Add(int(mid))
		for y := uint64(0); y < mid; y++ {
			go func(y uint64) {
				defer cwg.Done()
				for x := uint64(0); x < mid; x++ {
					a := pixels[y * size + x]
					d := a
					a = uint32(float64(a + pixels[(y + mid) * size + x]) / SQRT_2)
					d = uint32(float64(d - pixels[(y + mid) * size + x]) / SQRT_2)
					pixels[y * size + x] = a;
					pixels[(y + mid) * size + x] = d;
				}
			}(y)
		}
		cwg.Wait()
	}

	duration = time.Since(start)
	fmt.Printf("time to haar: %v\n", duration)

	start = time.Now()

	err = binary.Write(fout, binary.LittleEndian, pixels)
	check(err)

	duration = time.Since(start)
	fmt.Printf("time to write: %v\n", duration)

}