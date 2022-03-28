package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"sync"
)

var pi []int

func printlner(i ...int) {
	pi = i
}

type slice struct {
	valuesGuard *sync.Mutex
	values      []int
}

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var s []int
	for i := 0; i < 99; i++ {
		s = append(s, i)
	}

	ms := newMySliceType(s)
	for i := 0; i < 10; i++ {
		slicerInBoundsChannels(ms)
		fmt.Println(pi)
	}

	return events.APIGatewayProxyResponse{Body: `{"status": "success"}`, StatusCode: 200}, nil
}

func (s slice) Get(idx int) int {
	s.valuesGuard.Lock()
	defer s.valuesGuard.Unlock()
	checkBuffer(s.values, idx)
	return s.values[idx]
}

func (s slice) GetCh(ch chan int, idx int) {
	s.valuesGuard.Lock()
	defer s.valuesGuard.Unlock()
	checkBuffer(s.values, idx)
	ch <- s.values[idx]
}

func newMySliceType(values []int) slice {
	return slice{
		valuesGuard: &sync.Mutex{},
		values:      values,
	}
}

func fillBuffer(slice []int) map[int]int {
	result := map[int]int{}
	for i := 0; i < 100; i++ {
		for j := 0; j < len(slice); j++ {
			result[i*len(slice)+j] = slice[j]
		}
	}
	return result
}

func checkBuffer(slice []int, idx int) {
	buffer := make(map[int]int, len(slice)*100)
	buffer = fillBuffer(slice)
	for i := range buffer {
		if i == idx {
			return
		}
	}
}

func slicerInBoundsChannels(slice slice) {
	ch := make(chan int, 8)
	for i := 0; i < 8; i++ {
		go slice.GetCh(ch, i*8+0)
		go slice.GetCh(ch, i*8+1)
		go slice.GetCh(ch, i*8+2)
		go slice.GetCh(ch, i*8+3)
		go slice.GetCh(ch, i*8+4)
		go slice.GetCh(ch, i*8+5)
		go slice.GetCh(ch, i*8+6)
		go slice.GetCh(ch, i*8+7)
		a0 := <-ch
		a1 := <-ch
		a2 := <-ch
		a3 := <-ch
		a4 := <-ch
		a5 := <-ch
		a6 := <-ch
		a7 := <-ch
		printlner(a0, a1, a2, a3, a4, a5, a6, a7)
	}
}
