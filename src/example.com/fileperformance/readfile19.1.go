package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

// Variation7.go rev19 (64k rev18 + minor rearrangements) on my mac:

func main() {
	// go tool trace trace.pprof
	//
	//trace.Start(os.Stderr)
	//defer trace.Stop()

	// go tool pprof cpu.pprof
	//
	//pprof.StartCPUProfile(os.Stderr)
	//defer pprof.StopCPUProfile()

	//runtime.GOMAXPROCS(5)

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	start := time.Now()

	type entry struct {
		firstName string
		name      string
		date      string
	}

	//scanner.Buffer(make([]byte, 1e6), 1e6)
	scanner.Scan()
	for {
		// get all the names
		text := scanner.Bytes()

		if len(text) == 0 {
		}

		//index := bytes.Index(text, []byte(","))
		//fmt.Println(index)

		//s := "C00401224|A|M6|P|201804059101532000|15|IND|LIEBERMAN, MATTHEW|DEVON|PA|19333|STREAMLIGHT, INC.|PRODUCTION ASSEMBLER|05152017|0||SA11AI_81294473|1217152||CONTRIBUTION TO ACTBLUE|4050820181542985770"

		willScan := scanner.Scan()
		if !willScan {
			break
		}
	}

	nameTime := time.Since(start)
	fmt.Printf("lines time: %v\n", nameTime)
}
