package main

import (
	"flag"
	"fmt"
	"github.com/captain686/portScan/core"
	"log"
	"os"
	"time"
)

var (
	ip      string
	port    string
	num     int
	timeout int
)

func init() {
	flag.StringVar(&ip, "ip", "", "Target Ip 10.0.0.1, 10.0.0.5-10, 192.168.1.*, 192.168.10.0/24")
	flag.StringVar(&port, "p", "21,22,80,81,135,139,443,445,1433,3306,5432,6379,7001,8000,8080,8089,9000,9200,11211,27017", "Scan Ports")
	flag.IntVar(&timeout, "t", 3, "TimeOut Of Tcp Connect")
	flag.IntVar(&num, "n", 500, "Number of concurrent scans")
	flag.Usage = func() {
		_, err := fmt.Fprintf(os.Stderr, "Usage :\n")
		if err != nil {
			return
		}
		flag.PrintDefaults()
	}
}

func TimeTrack(start time.Time, logName string) {
	elapsed := time.Since(start)
	log.Printf("%s 耗时 %s", logName, elapsed)
}

func main() {
	flag.Parse()
	if flag.NFlag() == 0 || ip == "" {
		_, err := fmt.Fprintf(os.Stderr, "Usage :\n")
		if err != nil {
			return
		}
		flag.PrintDefaults()
	} else {
		startT := time.Now()
		host := core.CheckCidrIPs(ip)
		if host == nil {
			core.PortChecker(ip, port, num, time.Duration(timeout)*time.Second)
		} else {
			for ip := range host {
				core.PortChecker(*ip, port, num, time.Duration(timeout)*time.Second)
			}
		}
		defer TimeTrack(startT, "main")
	}

}
