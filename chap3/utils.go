package chap3

import (
	"fmt"
	"time"
)

func Tee(done <-chan interface{}, inputCh <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)

		for {
			select {
			case <-done:
				return
			case val, ok := <-inputCh:
				if !ok {
					return
				}

				for i := 0; i < 2; i++ {
					tmp1, tmp2 := out1, out2
					select {
					case <-done:
						return
					case tmp1 <- val:
						tmp1 = nil
					case tmp2 <- val:
						tmp2 = nil
					}
				}
			}
		}
	}()

	return out1, out2
}

func Ordone(done <-chan interface{}, inputCh <-chan interface{}) <-chan interface{} {
	outCh := make(chan interface{})
	go func() {
		defer close(outCh)
		for {
			select {
			case <-done:
				fmt.Println("Done received outside")
				return
			case val, ok := <-inputCh:
				if !ok {
					fmt.Println("Channel closed!")
					return
				}
				select {
				case outCh <- val:
				case <-done:
					fmt.Println("Done received inside")
					return
				}
			}
		}
	}()

	return outCh
}

func Repeat(done <-chan interface{}, values ...int) <-chan int {
	repCh := make(chan int)
	go func() {
		defer close(repCh)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case repCh <- v:
				}
			}
		}
	}()
	return repCh
}

func Take(done <-chan interface{}, input <-chan int, num int) <-chan int {
	outCh := make(chan int)
	go func() {
		defer close(outCh)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case outCh <- <-input:
			}
		}
	}()

	return outCh
}

func Generator(done <-chan interface{}, ints ...interface{}) <-chan interface{} {
	genCh := make(chan interface{})
	go func() {
		defer close(genCh)
		for _, v := range ints {
			select {
			case <-done:
				fmt.Println("Done closed in gen")
				return
			case genCh <- v:
				time.Sleep(1 * time.Second)
			}
		}
		fmt.Println("Gen exited")
	}()

	return genCh
}

func Add(done <-chan interface{}, intStream <-chan int, val int) <-chan int {
	outCh := make(chan int)
	go func() {
		defer close(outCh)
		for i := range intStream {
			select {
			case <-done:
				fmt.Println("Done closed in Add")
				return
			case outCh <- i + val:
			}
		}
		fmt.Println("Add exited")
	}()

	return outCh
}

func Multiply(done <-chan interface{}, intStream <-chan int, val int) <-chan int {
	outCh := make(chan int)
	go func() {
		defer close(outCh)
		for i := range intStream {
			select {
			case <-done:
				fmt.Println("Done closed in multiply")
				return
			case outCh <- val * i:
			}
		}
		fmt.Println("Multiply exited")
	}()

	return outCh
}
