package core

import (
	"github.com/malfunkt/iprange"
)

func CheckCidrIPs(data string) <-chan *string {
	ch := make(chan *string)
	list, err := iprange.ParseList(data)
	if err != nil {
		return nil
	}
	go func() {
		for _, rng := range list.Expand() {
			ing := rng.String()
			ch <- &ing
		}
		close(ch)
	}()
	return ch
}
