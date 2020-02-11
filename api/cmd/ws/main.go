package main

import (
	"fmt"
	"koffee/pkg/websocketkoffee"
	"log"
	"os/exec"
	"sync"
)

const addr = ":3333"

// R error fast
func R(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

var wg sync.WaitGroup

func jobify(jobs <-chan func(), id int) {
	for {
		select {
		case job := <-jobs:
			if job == nil {
				continue
			}
			fmt.Println("Jober number ", id)
			job()
		}
	}

}

func work() {

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	defer close(jobStream)
	// 	for i := 0; i < 1000000; i++ {
	// 		// buf := []byte("")
	// 		// s := strconv.AppendInt(buf, int64(i), 10)
	// 		jobStream <- func() {
	// 			for j := 0; j < 10000; j++ {
	// 			}
	// 			// fmt.Println("Finished ")
	// 		}
	// 	}
	// }()
	// wg.Wait()
}

func iserr(e error) bool {
	if e != nil {
		return true
	}
	return false
}

type intmutex struct {
	mu sync.Mutex
	n  uint
}

var concConn intmutex = intmutex{mu: sync.Mutex{}, n: 0}

func addConcurrent() {
	concConn.mu.Lock()
	concConn.n++
	fmt.Println("Concurrent connections: ", concConn.n)
	concConn.mu.Unlock()
}

func substractConcurrent() {
	concConn.mu.Lock()
	concConn.n--
	fmt.Println("Concurrent connections: ", concConn.n)
	concConn.mu.Unlock()
}

func worker(jobs <-chan func()) {
	for {
		select {
		case job := <-jobs:
			job()
		}
	}
}

func main() {
	const port = ":3333"
	if wserr := websocketkoffee.WSKoffeeHandle(websocketkoffee.Configuration{Port: port}, func(c websocketkoffee.Connection) {
		fmt.Printf("Listening on port %s\n", port)
		c.OnClose(func(cc *websocketkoffee.ConnectionData) {
			fmt.Printf("Closed %d\n", cc.ID)
			substractConcurrent()
		})
		c.OnOpen(func(cc *websocketkoffee.ConnectionData) {
			fmt.Printf("Opened %d\n", cc.ID)
			addConcurrent()
		})
		c.OnMessage(func(m *websocketkoffee.Message) error {
			fmt.Printf("Message %d: %s\n", m.ID, m.Data)
			websocketkoffee.Write(m.Conn, []byte("Yes/...."))
			e := exec.Command("clear")
			e.Run()
			return nil
		})
	}); wserr != nil {
		log.Fatal(wserr)
	}

	for {
		select {}
	}

}
