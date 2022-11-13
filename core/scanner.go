package core

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func PortChecker(address, ports string, num int, timeout time.Duration) {
	var portList []string
	var wg sync.WaitGroup
	n := make(chan int, num)
	success := make(chan string)
	if strings.Contains(ports, ",") {
		portList = strings.Split(ports, ",")
		success = make(chan string, len(portList))
		for i, port := range portList {
			iPort, err := strconv.Atoi(port)
			if err != nil {
				//log.Println(err)
				return
			}
			if iPort > 65535 || iPort < 1 {
				fmt.Println("[X] Port Range Is 1-65535 !")
				continue
			}
			wg.Add(1)
			target := address + ":" + port
			n <- i
			go CoonTester(&wg, target, success, n, timeout)
		}
	} else if strings.Contains(ports, "-") {
		startToEnd := strings.Split(ports, "-")
		start, end := startToEnd[0], startToEnd[1]
		iEnd, err := strconv.Atoi(end)
		if err != nil {
			//log.Println(err)
			return
		}
		iStart, err := strconv.Atoi(start)
		if err != nil {
			//log.Println(err)
			return
		}
		if iEnd > 65535 || iStart < 1 {
			fmt.Println("[X] Port Range Is 1-65535 !")
			iEnd = 65535
		}
		lenOfSuccess := iEnd - iStart
		success = make(chan string, lenOfSuccess)
		for i := iStart; i < iEnd+1; i++ {
			wg.Add(1)
			target := address + ":" + strconv.Itoa(i)
			n <- i
			go CoonTester(&wg, target, success, n, timeout)
		}
	} else {
		iPort, err := strconv.Atoi(ports)
		if err != nil {
			//log.Println(err)
			return
		}
		if iPort > 65535 || iPort < 1 {
			fmt.Println("[X] Port Input Error !")
			return
		}
		success = make(chan string, 1)
		wg.Add(1)
		target := address + ":" + ports
		n <- 1
		go CoonTester(&wg, target, success, n, timeout)
	}
	wg.Wait()
	close(success)
}

func CoonTester(wg *sync.WaitGroup, address string, success chan string, n chan int, timeout time.Duration) {
	//fmt.Println(address)
	defer func(wg *sync.WaitGroup) {
		wg.Done()
		<-n
	}(wg)
	coon, err := net.DialTimeout("tcp", address, timeout)
	//coon, err := net.Dial("tcp", address)
	if err != nil {
		//log.Println(err)
		return
	}
	fmt.Println(address + " Is Open")
	defer func(coon net.Conn) {
		err := coon.Close()
		if err != nil {
			//log.Println(err)
			return
		}
	}(coon)
	success <- address
}
