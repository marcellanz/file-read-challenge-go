package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

// Variation7.go rev19 (64k rev18 + rearrangements: do not format date, name at index 'in-loop') on my mac: 4.10

/*
nameTime: 7.665251714s, lineCountTime: 7.668233675s, donationsTime: 7.668921714s, mostCommonTime: 7.66892908s
nameTime: 6.128489519s, lineCountTime: 6.128532829s, donationsTime: 6.128597896s, mostCommonTime: 6.128602122s
nameTime: 6.19426646s, lineCountTime: 6.194322791s, donationsTime: 6.194388264s, mostCommonTime: 6.194391954s
nameTime: 5.991776806s, lineCountTime: 5.99181482s, donationsTime: 5.99187407s, mostCommonTime: 5.991878205s
nameTime: 6.968188955s, lineCountTime: 6.968230746s, donationsTime: 6.968301294s, mostCommonTime: 6.968305262s
nameTime: 6.650723052s, lineCountTime: 6.65078616s, donationsTime: 6.65085454s, mostCommonTime: 6.650859126s
nameTime: 6.350935963s, lineCountTime: 6.350976192s, donationsTime: 6.351034922s, mostCommonTime: 6.351043701s

date as string

nameTime: 4.401865629s, lineCountTime: 4.401904082s, donationsTime: 4.401960196s, mostCommonTime: 4.401964586s
nameTime: 4.389938167s, lineCountTime: 4.389975238s, donationsTime: 4.390037839s, mostCommonTime: 4.390042675s
nameTime: 4.337873304s, lineCountTime: 4.337906089s, donationsTime: 4.337961556s, mostCommonTime: 4.337966203s
nameTime: 4.456026815s, lineCountTime: 4.456064193s, donationsTime: 4.456128439s, mostCommonTime: 4.456133392s
nameTime: 4.431424339s, lineCountTime: 4.431460791s, donationsTime: 4.431520898s, mostCommonTime: 4.431525321s
nameTime: 4.431707128s, lineCountTime: 4.431751947s, donationsTime: 4.431809144s, mostCommonTime: 4.431813877s
nameTime: 4.420236883s, lineCountTime: 4.420279724s, donationsTime: 4.420337929s, mostCommonTime: 4.420342382s
nameTime: 4.457460186s, lineCountTime: 4.457492907s, donationsTime: 4.457547151s, mostCommonTime: 4.457551674s
nameTime: 4.435049851s, lineCountTime: 4.43508394s, donationsTime: 4.43513404s, mostCommonTime: 4.435138693s

*/

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

	start := time.Now()

	nameMap := make(map[string]int)
	dateMap := make(map[int]int)
	common := ""
	commonCount := 0

	type entry struct {
		firstName string
		name      string
		date      int
	}

	linesChunkLen := 64 * 1024
	linesChunkPoolAllocated := int64(0)
	linesPool := sync.Pool{New: func() interface{} {
		e := make([]string, 0, linesChunkLen)
		atomic.AddInt64(&linesChunkPoolAllocated, 1)
		return e
	}}

	collectedPoolAllocated := int64(0)
	collectedPool := sync.Pool{New: func() interface{} {
		e := make([]entry, 0, linesChunkLen)
		atomic.AddInt64(&collectedPoolAllocated, 1)
		return e
	}}

	var lines = linesPool.Get().([]string)
	mutex := &sync.Mutex{}
	wg := sync.WaitGroup{}

	namesCounted := false
	namesCount := 0
	fileLineCount := int64(0)

	scanner.Scan()
	for {
		// get all the names
		lines = append(lines, scanner.Text())

		willScan := scanner.Scan()
		if len(lines) == linesChunkLen || !willScan {

			wg.Add(len(lines))
			linesToProcess := lines // bug
			go func() {
				atomic.AddInt64(&fileLineCount, int64(len(linesToProcess)))
				collected := collectedPool.Get().([]entry)[:0]

				for _, text := range linesToProcess {
					e := entry{}
					split := strings.SplitN(text, "|", 9) // 10.95
					e.name = strings.TrimSpace(split[7])

					if len(e.name) != 0 {
						startOfName := strings.Index(e.name, ", ") + 2
						if endOfName := strings.Index(e.name[startOfName:], " "); endOfName < 0 {
							e.firstName = e.name[startOfName:]
						} else {
							e.firstName = e.name[startOfName : startOfName+endOfName]
						}
						if cs := strings.Index(e.firstName, ","); cs > 0 {
							e.firstName = e.firstName[:cs]
						}
					}

					// extract dates
					e.date, _ = strconv.Atoi(split[4][:6])
					collected = append(collected, e)
				}
				linesPool.Put(linesToProcess)

				mutex.Lock()
				for _, e0 := range collected {
					if len(e0.firstName) != 0 {
						ncount := nameMap[e0.firstName] + 1
						nameMap[e0.firstName] = ncount
						if ncount > commonCount {
							commonCount = ncount
							common = e0.firstName
						}
					}
					if namesCounted == false {
						if namesCount == 0 {
							fmt.Printf("Name: %s at index: %v\n", e0.name, 0)
						} else if namesCount == 432 {
							fmt.Printf("Name: %s at index: %v\n", e0.name, 432)
						} else if namesCount == 43243 {
							fmt.Printf("Name: %s at index: %v\n", e0.name, 43243)
							namesCounted = true
						}
						namesCount++
					}
					dateMap[e0.date]++
				}
				mutex.Unlock()

				collectedPool.Put(collected)
				wg.Add(-len(collected))
			}()
			lines = linesPool.Get().([]string)[:0]
		}
		if !willScan {
			break
		}
	}
	wg.Wait()

	fmt.Printf("Name time: %v\n", time.Since(start))
	fmt.Printf("Total file line count: %v\n", fileLineCount)
	fmt.Printf("Line count time: : %v\n", time.Since(start))

	for k, v := range dateMap {
		fmt.Printf("Donations per month and year: %v and donation ncount: %v\n", k, v)
	}
	fmt.Printf("Donations time: : %v\n", time.Since(start))

	fmt.Printf("The most common first name is: %s and it occurs: %v times.\n", common, commonCount)
	fmt.Printf("Most common name time: %v\n", time.Since(start))

	// other stats
	fmt.Printf("linesChunkPoolAllocated: %v, collectedPoolAllocated: %v\n", linesChunkPoolAllocated, collectedPoolAllocated)
	fmt.Printf("nameTime: %v, lineCountTime: %v, donationsTime: %v, mostCommonTime: %v\n", time.Since(start), time.Since(start), time.Since(start), time.Since(start))
}
