package main

import (
    "fmt"
    "bufio"
    "log"
    "os"
    "sync"
)

// helper
func chunkArray(foo []string, n int) [][]int {
    chunks := make([][]int, 0)
    chunkSize := (len(foo) + n - 1) / n
    for start := 0; start < len(foo); start += chunkSize {
        end := start + chunkSize
        if end > len(foo) {
            end = len(foo)
        }
        chunks = append(chunks, []int{start,end})
    }
    return chunks
}

// part 2: struct Counter
type Counter struct {
    // part 3: concurrency
    mu sync.Mutex
    counts map[string]int
    topFreq int
    topVisitor string
}

func NewCounter() *Counter {
    counts := make(map[string]int)
    return &Counter{counts: counts, topFreq: 0, topVisitor: ""}
}

// part 2: count all frequencies
func (c *Counter) countIPs(ips []string) {
    for i, _ := range ips {
        c.countIP(ips[i])
    }
}

// part 2: add method for Counter
func (c *Counter) countIP(ip string) int {
    // part 3: concurrency
    c.mu.Lock()
    defer c.mu.Unlock()
    c.counts[ip]++
    if c.counts[ip] > c.topFreq {
        c.topFreq = c.counts[ip]
        c.topVisitor = ip
    }
    ipCount := c.counts[ip]
    return ipCount
}

// part 3: count all frequencies using goroutines
func (c *Counter) countIPsConcurrently(ips []string, nroutines int) {
    wg := &sync.WaitGroup{}
    for _, idxs := range chunkArray(ips, nroutines) {
        wg.Add(1)
        ipChunk := ips[idxs[0]:idxs[1]]
        go func() {
            defer wg.Done()
            c.countIPs(ipChunk)
        }()
    }
    wg.Wait()
}

var counter = NewCounter()
var ips []string

func main() {
    // part 0: reading data line by line
    file, err := os.Open("data.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        ips = append(ips, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    // part 2: call counter on IPs
    //counter.countIPs(ips)

    // part 3: call counter concurrently
    counter.countIPsConcurrently(ips, 10)

    // part 5
    fmt.Println(counter.topFreq)
    fmt.Println(counter.topVisitor)
}
