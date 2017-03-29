package main

import "fmt"
import "net"
import "sync"
import "time"

func main() {
	var array []byte = make([]byte, 512)

	for i := 0; i < 512; i += 1 {
		array[i] = 'A'
	}

	var wg sync.WaitGroup
	service := "192.168.130.129:8888"
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", service)
	for i := 0; i < 50000; i += 1 {
		wg.Add(1)
		go func(id int) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
				}
				wg.Done()
			}()

			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				fmt.Printf("%d conn err: %s\n", id, err)
				return
			}

			defer func() { conn.Close() }()

			conn.SetReadDeadline(time.Now().Add(180 * time.Second))
			var result []byte = make([]byte, 512)
			count, err := conn.Read(result)
			if err != nil {
				fmt.Printf("%d read err: %s\n", id, err)
				return
			}

			if string(result) != string(array) {
				fmt.Printf("%d data err: %d:%s", id, count, result)
			}

			time.Sleep(1800 * time.Second)
		}(i)
	}

	wg.Wait()
}
