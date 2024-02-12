package main

import (
	"context"
	"fmt"
	test3 "github.com/ib-77/rop/test"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

func main() {
	test()
}

func test() {
	Profile(
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
	go func() {
		defer close(inputs)
		for i := 0; i < amount; i++ {
			inputs <- fixedValue
		}
	}()
	return inputs
}

func Profile(toProfileF func(), fileSuffix string, addCpu bool, addTrace bool, addMem bool) {

	if addCpu {

		cpuProfileFile, err := os.Create(fmt.Sprintf("cpu_rop_%s.prof", fileSuffix))
		if err != nil {
			panic(err)
		}
		defer cpuProfileFile.Close()

		if err := pprof.StartCPUProfile(cpuProfileFile); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()

	}

	if addTrace {

		// Start tracing
		traceFile, err := os.Create(fmt.Sprintf("trace_rop_%s.prof", fileSuffix))
		if err != nil {
			panic(err)
		}
		defer traceFile.Close()

		if err := trace.Start(traceFile); err != nil {
			panic(err)
		}
		defer trace.Stop()
	}

	// run code
	toProfileF()

	if addMem {

		memProfileFile, err := os.Create(fmt.Sprintf("mem_rop_%s.prof", fileSuffix))
		if err != nil {
			panic(err)
		}
		defer memProfileFile.Close()

		if err := pprof.WriteHeapProfile(memProfileFile); err != nil {
			panic(err)
		}

		time.Sleep(5 * time.Second)
	}
}
