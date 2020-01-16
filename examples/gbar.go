package main

import (
	"math/rand"
	"time"
)

func main() {
	names := []string{
		"MySQL Write",
		"MySQL Read",
		"Redis Write",
		"Redis Read",
		"Channal Push",
		"File I/O",
		"Logger Write",
		"HTTP GET",
		"Database Align",
		"Analysis Status",
	}
	for i := int64(1); i <= 9; i++ {
		go func(i int64) {
			step := int64(0)
			for {
				step++
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				time.Sleep(time.Duration(r.Int63n(int64(time.Second))))
				if i%2 == 1 {
					//Status(names[i])
					Info(names[i], ".....")
				} else {
					if Progress(names[i], step) {
						return
					}
				}
			}
		}(i)
	}
	s := make(chan int)
	<-s
}
