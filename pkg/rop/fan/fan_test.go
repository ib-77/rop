package fan

import (
	"context"
	"github.com/ib-77/rop/pkg/rop"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func Test_OutNext(t *testing.T) {

	inputCh := make(chan rop.Result[int], 5)
	inputCh <- rop.Success(1)
	inputCh <- rop.Success(3)
	inputCh <- rop.Success(4)
	inputCh <- rop.Success(2)
	inputCh <- rop.Success(5)

	ctx := context.Background()
	outputChs := OutNext[int](ctx, inputCh, 4)

	wg := sync.WaitGroup{}
	wg.Add(len(outputChs))

	for index, ch := range outputChs {
		go func(c chan rop.Result[int], index int) {
			defer wg.Done()

			if index == 0 {
				assert.Equal(t, rop.Success(1), <-c)
				assert.Equal(t, rop.Success(5), <-c)
			}

			if index == 1 {
				assert.Equal(t, rop.Success(3), <-c)
			}

			if index == 2 {
				assert.Equal(t, rop.Success(4), <-c)
			}

			if index == 3 {
				assert.Equal(t, rop.Success(2), <-c)
			}

		}(ch, index)
	}

	wg.Wait()

}

func Test_OutNext_SliceToChs(t *testing.T) {

	inputCh := make(chan rop.Result[int], 5)
	inputCh <- rop.Success(1)
	inputCh <- rop.Success(3)
	inputCh <- rop.Success(4)
	inputCh <- rop.Success(2)
	inputCh <- rop.Success(5)

	ctx := context.Background()
	outputChs := SliceToChs(OutNext[int](ctx, inputCh, 4))

	wg := sync.WaitGroup{}
	wg.Add(len(outputChs))

	index := 0
	for ch := range outputChs {
		go func(c chan rop.Result[int], index int) {
			defer wg.Done()

			if index == 0 {
				assert.Equal(t, rop.Success(1), <-c)
				assert.Equal(t, rop.Success(5), <-c)
			}

			if index == 1 {
				assert.Equal(t, rop.Success(3), <-c)
			}

			if index == 2 {
				assert.Equal(t, rop.Success(4), <-c)
			}

			if index == 3 {
				assert.Equal(t, rop.Success(2), <-c)
			}

		}(ch, index)
		index++
	}

	wg.Wait()

}
