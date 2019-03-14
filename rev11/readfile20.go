package rev11

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// why?
// https://stuartmarks.wordpress.com/2019/01/11/processing-large-files-in-java/

// original challenge in JS: https://itnext.io/using-node-js-to-read-really-really-large-files-pt-1-d2057fe76b33
// original challenge in JAVA: https://itnext.io/using-java-to-read-really-really-large-files-a6f8a3f44649

// mds off => sudo mdutil -i off /

// original on stuartmarks mac: 108s
// original on my mac: 86s

// Variation7 on stuartmarks mac: 32s
// Variation7 on my mac: 23s

// Variation7.go naive on my mac: 38s (+40% to Variation7), rerun later: 27-28s
// Variation7.go rev10 (stupid channel) on my mac: 12m6s :-o
// Variation7.go rev11 (entries channel) on my mac: 39s
// Variation7.go rev11 (1k line-chunks channel) on my mac: 39s
// Variation7.go rev11 (4k line-chunks channel) on my mac: 33s
// Variation7.go rev11 (8k line-chunks channel) on my mac: 32s
// Variation7.go rev11 (32k line-chunks channel) on my mac: 26s
// Variation7.go rev11 (64k line-chunks channel) on my mac: 23s

// Variation7.go rev12 (1k entry-chunks channel) on my mac: 24s
// Variation7.go rev12 (8k entry-chunks channel) on my mac: 23s
// Variation7.go rev12 (64k entry-chunks channel) on my mac: 24s

// Variation7.go rev13 (1k rev12 mutex) on my mac: 13s
// Variation7.go rev13 (64k rev12 mutex) on my mac: 11.5s

// Variation7.go rev15 (64k rev15 common-name in loop) on my mac: 10.5s
// Variation7.go rev15 (64k rev15 common-name + date in loop) on my mac: 9.8s (min) 10.5 - 11.0

// Variation7.go rev16 (64k rev15 + sync.Pool) on my mac: 9.8s (min)... 10.0 - 10.5

// Variation7.go rev17 (64k rev16 + regex bug) on my mac: 9.15 (min)
// Variation7.go rev17 (64k rev16 + simpler regex) on my mac: 8.32 (min) 8.7..9.2 with about 80% CPU

// Variation7.go rev18 (64k rev17 + no regex) on my mac: 6.46 (min) 6.8 - 7.1 with about 50% CPU

// Variation7.go rev0.1 port mac + no regex and date/name inner loop my mac: 16.6

// Variation7.go rev19 (64k rev18 + minor rearrangements) on my mac: 5.99

// Variation7.go rev20 (rev19 + use bytes instead of string) on my mac:

func main() {
	// go tool trace trace.pprof
	//
	//trace.Start(os.Stderr)
	//defer trace.Stop()

	// go tool pprof cpu.pprof
	//
	//pprof.StartCPUProfile(os.Stderr)
	//defer pprof.StopCPUProfile()

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	names := make([][12]byte, 0, 0)
	dates := make([][6]byte, 0, 0)

	start := time.Now()

	nameMap := make(map[string]int)
	dateMap := make(map[string]int)
	common := ""
	commonCount := 0

	type entry struct {
		firstName [12]byte
		name      [12]byte
		date      [6]byte
	}

	linesChunkLen := 64 * 1024
	linesChunkPoolAllocated := int64(0)
	linesPool := sync.Pool{New: func() interface{} {
		e := make([][512]byte, linesChunkLen, linesChunkLen)[:]
		atomic.AddInt64(&linesChunkPoolAllocated, 1)
		return e
	}}

	collectedPoolAllocated := int64(0)
	collectedPool := sync.Pool{New: func() interface{} {
		e := make([]entry, 0, linesChunkLen)
		atomic.AddInt64(&collectedPoolAllocated, 1)
		return e
	}}

	mutex := &sync.Mutex{}
	wg := sync.WaitGroup{}

	var lines = linesPool.Get().([][512]byte)
	current := 0

	maxLenName := 0
	maxLenFirstName := 0

	scanner.Scan()
	for {
		// get all the names
		i := scanner.Bytes()
		copy(lines[current][:], i)
		current++
		willScan := scanner.Scan()

		if current == linesChunkLen || !willScan {
			linesToProcess := lines // bug
			wg.Add(len(linesToProcess))
			current = 0
			go func() {
				collected := collectedPool.Get().([]entry)[:0]
				for _, byteLine := range linesToProcess {
					e := entry{}
					split := bytes.SplitN(byteLine[:], []byte("|"), 9) // 10.95

					if l := len(bytes.TrimSpace(split[7])); l > maxLenName {
						maxLenName = l
					}

					if copy(e.name[:], bytes.TrimSpace(split[7])) != 0 {
						startOfName := bytes.Index(e.name[:], []byte(", ")) + 2
						if endOfName := bytes.Index(e.name[startOfName:], []byte(" ")); endOfName < 0 {
							copy(e.firstName[:], e.name[startOfName:])
						} else {
							copy(e.firstName[:], e.name[startOfName:startOfName+endOfName])
						}
						if bytes.HasSuffix(e.firstName[:], []byte(",")) {
							copy(e.firstName[:], bytes.Replace(e.firstName[:], []byte(","), []byte(""), -1))
						}

					}

					// extract dates
					copy(e.date[:], bytes.TrimSpace(split[4])[:6]) // year-month das not needed but produces new object
					collected = append(collected, e)
				}

				linesPool.Put(linesToProcess)
				mutex.Lock()
				for _, e0 := range collected {
					if len(e0.firstName) != 0 {
						fns := string(e0.firstName[:])
						ncount := nameMap[fns] + 1
						nameMap[fns] = ncount
						if ncount > commonCount {
							commonCount = ncount
							common = fns
						}
					}
					names = append(names, e0.name)
					dates = append(dates, e0.date)
					dateMap[string(e0.date[:])]++
				}
				mutex.Unlock()
				collectedPool.Put(collected)
				wg.Add(-len(linesToProcess))
			}()
			lines = linesPool.Get().([][512]byte)
		}

		if !willScan {
			break
		}
	}
	wg.Wait()

	fmt.Printf("Name: %s at index: %v\n", names[0], 0)
	fmt.Printf("Name: %s at index: %v\n", names[432], 432)
	fmt.Printf("Name: %s at index: %v\n", names[43243], 43243)

	nameTime := time.Since(start)
	fmt.Printf("Name time: %v\n", nameTime)
	fmt.Printf("Total file line count: %v\n", len(names))
	lineCountTime := time.Since(start)
	fmt.Printf("Line count time: : %v\n", lineCountTime)

	for k, v := range dateMap {
		fmt.Printf("Donations per month and year: %v and donation ncount: %v\n", k, v)
	}
	donationsTime := time.Since(start)
	fmt.Printf("Donations time: : %v\n", donationsTime)

	fmt.Printf("The most common first name is: %s and it occurs: %v times.\n", common, commonCount)
	mostCommonTime := time.Since(start)
	fmt.Printf("Most common name time: %v\n", mostCommonTime)

	// other stats
	fmt.Printf("linesChunkPoolAllocated: %v, collectedPoolAllocated: %v\n", linesChunkPoolAllocated, collectedPoolAllocated)
	fmt.Printf("nameTime: %v, lineCountTime: %v, donationsTime: %v, mostCommonTime: %v\n", nameTime, lineCountTime, donationsTime, mostCommonTime)

	fmt.Printf("maxLenName: %v\n", maxLenName)
	fmt.Printf("maxLenFirstName: %v\n", maxLenFirstName)
}

func m() {
	b := make([]byte, 5, 5)
	copy(b, []byte("yes"))

	if len(b) == 0 {
		fmt.Printf("b is empty: %v", b)
	} else {
		fmt.Printf("b: %v", b)
	}
}
