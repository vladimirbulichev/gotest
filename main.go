package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

type watchDog struct {
	avail int
	total int
	sync.Mutex
}

func newWatchDog() *watchDog {
	return &watchDog{
		avail: 5,
		total: 0,
	}
}

func (w *watchDog) Add() {
	w.Lock()
	defer w.Unlock()
	w.avail--
}

func (w *watchDog) CheckAvail() int {
	w.Lock()
	defer w.Unlock()
	return w.avail
}

func (w *watchDog) Release(c int) {
	w.Lock()
	defer w.Unlock()
	w.avail++
	w.total += c
}

func main() {
	in := os.Stdin
	lock := make(chan int)
	scanner := bufio.NewScanner(in)
	wd := newWatchDog()
	var wg sync.WaitGroup
	for scanner.Scan() {
		urlScan := scanner.Text()
		if wd.CheckAvail() == 0 {
			select {
			case <-lock:
			}
		}
		wd.Add()
		wg.Add(1)
		go func(url string) {
			wd.Release(wordCount(url))
			wg.Done()
			lock <- 1
		}(urlScan)
	}
	wg.Wait()
	wd.Lock()
	fmt.Println("Total", wd.total)
	wd.Unlock()
}

func wordCount(url string) (c int) {
	resp, err := http.Get(url)
	if err != nil {
		return 0
	}
	body, _ := ioutil.ReadAll(resp.Body)
	text := string(body)
	c = strings.Count(text, "Go")
	fmt.Printf("Count at %s : %d\n", url, c)
	return
}
