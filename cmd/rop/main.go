package main

import (
	"context"
	"fmt"
	test3 "github.com/ib-77/rop/test"
	"log"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

func main() {
	test()
}

func test() {
	_ = Profile(
		RunMassCase01,
		"mass-case01",
		false,
		true,
		false)
}

func test2() {
	Profile(
		RunCase01,
		"case01",
		false,
		false,
		true)
}

func RunMassCase01() {
	ctx := context.Background()
	test3.MassRopCase01(ctx, generateFixedValueUnbufferedChan(5, 1))
}

func RunCase01() {
	test3.RopCase01(1)
}

func generateFixedValueUnbufferedChan(amount int, fixedValue int) chan int {
	inputs := make(chan int)
	go func(am int, fv int) {
		defer close(inputs)
		for i := 0; i < am; i++ {
			inputs <- fv
		}
	}(amount, fixedValue)
	return inputs
}

func Profile(toProfileF func(), fileSuffix string, addCpu bool, addTrace bool, addMem bool) (err error) {

	if addCpu {

		cpuProfileFile, errFile := os.Create(fmt.Sprintf("cpu_rop_%s.prof", fileSuffix))
		if errFile != nil {
			err = errFile
			return
		}
		defer func() {
			errClose := cpuProfileFile.Close()
			if err != nil {
				if errClose != nil {
					log.Printf("failed to close cpuProfileFile %v", errClose)
				}
				return
			}
			err = errClose
		}()

		if errCpu := pprof.StartCPUProfile(cpuProfileFile); errCpu != nil {
			err = errCpu
			return
		}
		defer pprof.StopCPUProfile()

	}

	if addTrace {

		// Start tracing
		traceFile, errTraceFile := os.Create(fmt.Sprintf("trace_rop_%s.prof", fileSuffix))
		if errTraceFile != nil {
			err = errTraceFile
			return
		}
		defer func() {
			errClose := traceFile.Close()
			if err != nil {
				if errClose != nil {
					log.Printf("failed to close traceFile %v", errClose)
				}
				return
			}
			err = errClose
		}()

		if errTrace := trace.Start(traceFile); errTrace != nil {
			err = errTrace
			return
		}
		defer trace.Stop()
	}

	// run code
	toProfileF()

	if addMem {

		memProfileFile, errMem := os.Create(fmt.Sprintf("mem_rop_%s.prof", fileSuffix))
		if errMem != nil {
			err = errMem
			return
		}
		defer func() {
			errClose := memProfileFile.Close()
			if err != nil {
				if errClose != nil {
					log.Printf("failed to close memProfileFile %v", errClose)
				}
				return
			}
			err = errClose
		}()

		if errHeap := pprof.WriteHeapProfile(memProfileFile); errHeap != nil {
			err = errHeap
			return
		}

		time.Sleep(5 * time.Second)
	}

	return
}
