package unusual_generics_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xakep666/unusual_generics"
)

func ExampleSingleFlightGroup() {
	const (
		concurrency = 5
		wait        = 2 * time.Second
	)

	var (
		wgProduce, wgConsume sync.WaitGroup
		sfg                  unusual_generics.SingleFlightGroup[string]
	)

	var (
		calls   int32
		results = make(chan unusual_generics.SingleFlightResult[string])
	)

	wgConsume.Add(1)
	go func() {
		defer wgConsume.Done()
		for result := range results {
			fmt.Printf("Val: %q, Shared: %t, Error: %v\n", result.Val, result.Shared, result.Err)
		}
	}()

	for i := 0; i < concurrency; i++ {
		wgProduce.Add(1)
		go func() {
			defer wgProduce.Done()

			results <- <-sfg.DoChan("key", func() (string, error) {
				// do something heavy
				time.Sleep(wait)
				atomic.AddInt32(&calls, 1)
				return "test", nil
			})
		}()
	}

	wgProduce.Wait()
	close(results)
	wgConsume.Wait()

	fmt.Println("Calls:", calls)

	// Output:
	// Val: "test", Shared: true, Error: <nil>
	// Val: "test", Shared: true, Error: <nil>
	// Val: "test", Shared: true, Error: <nil>
	// Val: "test", Shared: true, Error: <nil>
	// Val: "test", Shared: true, Error: <nil>
	// Calls: 1
}
