package test

import (
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"testing"
	"time"
)

func TestPerfBegin(t *testing.T) {

	Profile(
		func() {
			ropCase01(1)

		},
		"case01",
		true,
		true)

}

func Profile(toProfileF func(), fileSuffix string, addMem bool, addTrace bool) {

	cpuProfileFile, err := os.Create(fmt.Sprintf("./results/cpu_rop_%s.prof", fileSuffix))
	if err != nil {
		panic(err)
	}
	defer cpuProfileFile.Close()

	if err := pprof.StartCPUProfile(cpuProfileFile); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	if addTrace {

		// Start tracing
		traceFile, err := os.Create(fmt.Sprintf("./results/trace_rop_%s.prof", fileSuffix))
		if err != nil {
			panic(err)
		}
		defer traceFile.Close()

		if err := trace.Start(traceFile); err != nil {
			panic(err)
		}
		defer trace.Stop()
	}

	toProfileF()

	if addMem {

		memProfileFile, err := os.Create(fmt.Sprintf("./results/mem_rop_%s.prof", fileSuffix))
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
