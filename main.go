package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
)

const CHAN_SIZE = 100

func main() {
	// CPU PROFILING
	var err error
	fCPU, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(fCPU)
	defer pprof.StopCPUProfile()

	ch := make(chan [][]byte, CHAN_SIZE)

	go ReadFromFile("LARGE_FILE.txt", ch)

	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	wg.Add(numCPU)
	for range numCPU {
		go func() {
			defer wg.Done()

			for range ch {
				fmt.Println("buffer processed")
			}
		}()
	}
	wg.Wait()

	// MEM PROFILING
	fMEM, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer fMEM.Close()
	if err := pprof.WriteHeapProfile(fMEM); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}
