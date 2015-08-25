package main

import (
	"os"
	"fmt"
	"math/rand"
	"time"
	"encoding/binary"
)

// change size to make image bigger/smaller
const IMGSIZE uint64 = 0x3000  // 0x6000

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s [create_file]\n", os.Args[0]);
		os.Exit(1)
	}

	var size uint64

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	fout, err := os.Create(os.Args[1])
	check(err)
	defer fout.Close()

	size = IMGSIZE
	fmt.Printf("%d\n", size);

	err = binary.Write(fout, binary.LittleEndian, size)
	check(err)

	pixels := make([]uint32, size * size)

	for y := uint64(0); y < size; y++ {
		for x := uint64(0); x < size; x++ {
			pixels[y * size + x] = r.Uint32() & 0xFFFF
		}
	}

	err = binary.Write(fout, binary.LittleEndian, pixels)
	check(err)
}