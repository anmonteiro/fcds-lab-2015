package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
	"io"
)

const N_BUCKETS int = 94

type bucket struct {
	Data []string
	total int
}

func sortBuckets(buckets []bucket) {
	var wg sync.WaitGroup

	wg.Add(len(buckets))
	for _, b := range buckets {
		go func(b bucket) {
			defer wg.Done()
			sort.Strings(b.Data)
		}(b)
	}
	wg.Wait()
}

func StreamBucketSort(stream io.Reader, length int, size int) []string {
	var b *bucket
	buckets := make([]bucket, N_BUCKETS)
	returns := make([]string, size)

	slice := func (i int) int {
		return (i * size / N_BUCKETS)
	}

	for i := 0; i < N_BUCKETS; i++ {
		buckets[i].Data = returns[slice(i):slice(i+1)]
		// for completion; Go already initializes ints as 0
		buckets[i].total = 0
	}

	readBuffer := make([]byte, length)
	for i := 0; i < size; i++ {
		_, err := stream.Read(readBuffer)
		check(err)
		str := string(readBuffer[:(len(readBuffer)-1)])
		b = &buckets[str[0] - 0x21]

		b.Data[b.total] = str
		b.total++
	}
	start := time.Now()
	sortBuckets(buckets)
	duration := time.Since(start)
	fmt.Println("Sort time: ", duration)	

	return returns
}




