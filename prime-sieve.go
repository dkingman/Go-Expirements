package main

import (
  "fmt"
  "runtime"
  "time"
  "math"
)

type Result struct {
	terms int
	prime int
}

func CreateIntegers(n int) map[int]int {
  result := make(map[int]int)
  for i := 2; i < n; i++ {
    result[i] = i
  }
  return result
}

func filterMultiple(primes map[int]int, lowest int, highest int) map[int]int {
	// Remove all the multiples of 'lowest' up to 'highest'
	for i := lowest; i < highest; i = lowest + i {
		// If our multiple is in the list, delete it
		_, ok := primes[i]
		if ok == true {
			delete(primes, i)
		}
	}
	return primes
}

func checkForPrimes(m map[int]int, lowest int, highest int) int {
	for i := lowest; i < highest; i++ {
		_, ok := m[i]
		if ok == true {
			return i
		}
	}
	return -1
}

func FindConsecPrimes(primes []int, prime int, ch chan Result) {
	// finished := false
	for i, value := range primes {
		totalTerms := 0
		total := 0
		if value >= prime { 
			break
		}
		for _, currPrime := range primes[i:] {
			// Optimization, stop checking when we are approaching the prime we are looking for
			total += currPrime
			totalTerms += 1
			if total >= prime {
				if total == prime {
					r := Result{totalTerms, prime}
					// fin = true
					ch <- r
				}
				break
			} 
		}
	}

	if primes[len(primes)-1] == prime {
		ch <- Result{-1,-1}
	}
	// fmt.Println(prime, "finished")

}

func FindPrimes(highest int) []int {
  	m := CreateIntegers(highest)
  	result := make([]int, 0)
  	result = append(result, 2)
  	multiple := 2
  	// Check for multiples only up to the the sqrt of the number limit, ie. check 100 entries if we are looking up to 10,000 for primes
  	for ; math.Sqrt(float64(multiple)) < float64(highest); {
  		updatedPrimes := filterMultiple(m,multiple, highest)
  		prime := checkForPrimes(updatedPrimes, multiple, highest)
  		multiple = prime
  		if prime == -1 { break }
  		result = append(result, prime)
  	}
  	
  	return result
}

func spawnChildren(primes []int, ch chan Result) {
	for _, prime := range primes {
		go FindConsecPrimes(primes, prime, ch)
	}
}

func main() {
	numCpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numCpu)
	// fmt.Println(numCpu)

	// Time how long it takes to find the primes
	s := time.Now()
	primes := FindPrimes(100000)
	e := time.Now()
	fmt.Println("Finding all primes took: ", e.Sub(s))

	curr := Result{}
	largest := Result{}
  	ch := make(chan Result, len(primes))
  	start := time.Now()
  	go spawnChildren(primes, ch)
  	for {
  		// Blocks until a case can proceed
  		select {
  			case curr = <- ch: {
  				if curr.terms == -1 {
  					end := time.Now()
  					fmt.Println("Time: ", end.Sub(start))
  					return
  				}
  				if curr.terms > largest.terms {
  					largest = curr
  					fmt.Println("largest = ",largest)
  				}
  				
  			}
  		}
  	}
  	fmt.Println("DONE")
}