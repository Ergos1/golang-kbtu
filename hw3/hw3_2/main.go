package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type UsedUrl struct {
	url map[string]bool
	mu  sync.Mutex
}

func (uu *UsedUrl) checkUsed(url string) bool {
	uu.mu.Lock()
	defer uu.mu.Unlock()
	if _, ok := uu.url[url]; ok {
		return true
	}
	return false
}

func (uu *UsedUrl) addUsed(url string) {
	uu.mu.Lock()
	defer uu.mu.Unlock()
	uu.url[url] = true
}

var (
	usedUrl = UsedUrl{url: make(map[string]bool)}
	wg      = sync.WaitGroup{}
)

func Crawl(url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		wg.Done()
		return
	}

	body, urls, err := fetcher.Fetch(url)
	usedUrl.addUsed(url)
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		if usedUrl.checkUsed(u) {
			continue
		}
		wg.Add(1)
		go Crawl(u, depth-1, fetcher)
	}
	wg.Done()
	return
}

func main() {
	wg.Add(1)
	Crawl("https://golang.org/", 4, fetcher)
	wg.Wait()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
