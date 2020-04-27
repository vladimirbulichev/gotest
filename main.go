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
	return &watchDog{avail: 5, total: 0}
}

func (w *watchDog) Add() {
	w.Lock()
	defer w.Unlock()
	w.avail--
}

func (w *watchDog) CheckAvail() bool {
	w.Lock()
	defer w.Unlock()
	return w.avail > 0
}

func (w *watchDog) Release(c int) {
	w.Lock()
	defer w.Unlock()
	w.avail++
	w.total += c
}

func main() {
	in := os.Stdin
	scanner := bufio.NewScanner(in)
	wd := newWatchDog()
	var wg sync.WaitGroup
	for scanner.Scan() {
		urlScan := scanner.Text()
		for !wd.CheckAvail() {
			fmt.Print("Wait\r")
		}
		wd.Add()
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			wd.Release(wordCount(url))
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
	//	time.Sleep(1*time.Second)
	return
}
