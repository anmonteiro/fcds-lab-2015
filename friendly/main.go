package main

import (
	"fmt"
	"sync"
	"math"
	"os"
	"strconv"
	"time"
)

var OSTHREADS int

func gcd(u, v int) int {
	if v == 0 {
		return u
	}
	return gcd(v, u % v)
}

func friendlyNumbers(start, end int) {
	last := end - start + 1

	the_num, num, den := make([]int, last), make([]int, last), make([]int, last)

	
	var wg sync.WaitGroup

	wg.Add(OSTHREADS)
	for i := 0; i < OSTHREADS; i++ {
		go func(idx int) {
			defer wg.Done()

			nrsPerThread := int(math.Ceil(float64(last) / float64(OSTHREADS)))
			startIdx := start + idx * nrsPerThread

			endIdx := (startIdx + nrsPerThread - 1)
			if endIdx > end {
				endIdx = end
			}

			for k := startIdx; k <= endIdx; k++ {
				var factor, ii, sum, done, n int
				ii = k - start
				sum = 1 + k
				the_num[ii] = k
				done = k
				factor = 2
				for factor < done {
					if k % factor == 0 {
						sum += (factor + (k / factor))
						if done = k / factor; done == factor {
							sum -= factor
						}
					}
					factor++
				}
				num[ii] = sum
				den[ii] = k
				n = gcd(num[ii], den[ii])
				num[ii] /= n
				den[ii] /= n
			}
		}(i)
	}
	
	wg.Wait()

	wg.Add(last)
	for i := 0; i < last; i++ {
		go func(idx int) {
			defer wg.Done()
			for j := idx+1; j < last; j++ {
				if num[idx] == num[j] && den[idx] == den[j] {
					fmt.Printf("%d and %d are FRIENDLY\n", the_num[idx], the_num[j])
				}
			}
		}(i)
	}
	wg.Wait()
}

func main() {
	var err error
	OSTHREADS, err = strconv.Atoi(os.Getenv("GOMAXPROCS"))
	if err != nil {
		OSTHREADS = 1
	}

	startTime := time.Now()

	for {
		var start, end int
		fmt.Scanf("%d %d", &start, &end)
		if start == 0 && end == 0 {
			break
		}
		fmt.Printf("Number %d to %d\n", start, end)
		
		friendlyNumbers(start, end)
	}
	duration := time.Since(startTime)
	fmt.Println("Duration: ", duration)
}