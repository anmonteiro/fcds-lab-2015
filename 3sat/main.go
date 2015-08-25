package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

var iVar []int64
var maxNumber int64
var workerCount int

// Solves the 3-SAT system using an exaustive search
// It finds all the possible values for the set of variables using
// a number. The binary representation of this number represents 
// the values of the variables. Since a long variable has 64 bits, 
// this implementations works with problems with up to 64 variables.
func clauseSolverWorker(clauses [][]int8, nClauses, startIdx int, ch chan<- int64) {
	for number := int64(startIdx); number < maxNumber; number += int64(workerCount) {
		var c int
		for c = 0; c < nClauses; c++ {
			variable := clauses[0][c]
			if variable > 0 && (number & iVar[variable - 1]) > 0 {
				continue // clause is true
			} else if variable < 0 && (number & iVar[-variable - 1]) == 0 {
				continue // clause is true
			}

			variable = clauses[1][c]
			if variable > 0 && (number & iVar[variable - 1]) > 0 {
				continue // clause is true
			} else if variable < 0 && (number & iVar[-variable - 1]) == 0 {
				continue // clause is true
			}

			variable = clauses[2][c]
			if variable > 0 && (number & iVar[variable - 1]) > 0 {
				continue // clause is true
			} else if variable < 0 && (number & iVar[-variable - 1]) == 0 {
				continue // clause is true
			}

			break
		}

		if c == nClauses {
			ch <- number
			break
		}
	}

	ch <- -1
}

func solveClauses(clauses [][]int8, nClauses, nVar int) int64 {
	iVar = make([]int64, nVar)
	for i := 0; i < nVar; i++ {
		iVar[i] = int64(math.Exp2(float64(i)))
	}
	maxNumber = int64(math.Exp2(float64(nVar)))

	ch := make(chan int64)
	for i := 0; i < workerCount; i++ {
		go clauseSolverWorker(clauses, nClauses, i, ch)
	}

	var solution int64 = -1
	for i := 0; i < workerCount; i++ {
		solution = <-ch
		if solution < 0 {
			continue
		} else {
			fmt.Println(solution)
			break
		}
	}

	return solution
}

// Read nClauses clauses of size 3. nVar represents the number of variables
// Clause[0][i], Clause[1][i] and Clause[2][i] contains the 3 elements of the i-esime clause.
// Each element of the caluse vector may contain values selected from:
// k = -nVar, ..., -2, -1, 1, 2, ..., nVar. The value of k represents the index of the variable.
// A negative value remains the negation of the variable.
func readClauses(nClauses, nVar int) [][]int8 {
	clauses := make([][]int8, 3)
	clauses[0] = make([]int8, nClauses)
	clauses[1] = make([]int8, nClauses)
	clauses[2] = make([]int8, nClauses)

	for i := 0; i < nClauses; i++ {
		fmt.Scanf("%d %d %d", &(clauses[0][i]), &(clauses[1][i]), &(clauses[2][i]))
	}

	return clauses
}

func main() {
	var err error

	workerCount, err = strconv.Atoi(os.Getenv("GOMAXPROCS"))
	if err != nil {
		workerCount = 1
	}

	var nClauses, nVar int

	fmt.Scanf("%d %d", &nClauses, &nVar);

	clauses := readClauses(nClauses, nVar)

	startTime := time.Now()

	solution := solveClauses(clauses, nClauses, nVar)

	if solution >= 0 {
		fmt.Printf("Solution found [%d]: ", solution)
		for i := 0; i < nVar; i++ {
			fmt.Printf("%d ", int64(float64(solution & int64(math.Exp2(float64(i)))) / math.Exp2(float64(i))));
		}
		fmt.Printf("\n");
	} else {
		fmt.Printf("Solution not found.\n")
	}

	duration := time.Since(startTime)
	fmt.Println("Duration: ", duration)
}