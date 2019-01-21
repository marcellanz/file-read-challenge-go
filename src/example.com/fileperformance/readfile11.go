package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

func main() {

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	namePat := regexp.MustCompile(", \\s*([^, ]+)")
	names := make([]string, 0, 0)
	firstNames := make([]string, 0, 0)
	dates := make([]string, 0, 0)

	start := time.Now()

	scanner := bufio.NewScanner(file)

	nameMap := make(map[string]int)
	common := ""
	commonCount := 0

	type entry struct {
		firstName string
		name      string
		date      string
		wg        *sync.WaitGroup
	}
	entries := make(chan entry)
	wg := sync.WaitGroup{}

	go func() {
		for {
			select {
			case e, ok := <-entries:
				if ok {
					if e.firstName != "" {
						firstNames = append(firstNames, e.firstName)
					}
					names = append(names, e.name)
					dates = append(dates, e.date)
					e.wg.Done()
				}
			}
		}
	}()

	ii := 0
	chunkLen := 256 * 1024 //4 * 1024 // 1024 => 5.9 vs 4k => 5.1
	chunks := make([]string, 0, 0)
	for scanner.Scan() {
		// get all the names
		line := scanner.Text()
		chunks = append(chunks, line)

		if len(chunks) == chunkLen {
			wg.Add(len(chunks))
			process := chunks // bug
			go func() {
				for _, text := range process {
					e := entry{wg: &wg}
					split := strings.SplitN(text, "|", 9) // 10.95
					name := strings.TrimSpace(split[7])
					e.name = name

					// extract first names
					if matches := namePat.FindAllStringSubmatch(name, 1); len(matches) > 0 {
						e.firstName = matches[0][1]
					}
					// extract dates
					chars := strings.TrimSpace(split[4])[:6]
					e.date = chars[:4] + "-" + chars[4:6]
					entries <- e
				}
			}()
			chunks = make([]string, 0, 0) // bug
		}
		ii++
	}
	wg.Wait()
	close(entries)

	fmt.Printf("Name: %s at index: %v\n", names[0], 0)
	fmt.Printf("Name: %s at index: %v\n", names[432], 432)
	fmt.Printf("Name: %s at index: %v\n", names[43243], 43243)

	fmt.Printf("Name time: %v\n", time.Since(start))
	fmt.Printf("Total file line count: %v\n", len(names))
	fmt.Printf("Line count time: : %v\n", time.Since(start))

	dateMap := make(map[string]int)
	for _, date := range dates {
		dateMap[date] += 1
	}
	for k, v := range dateMap {
		fmt.Printf("Donations per month and year: %v and donation ncount: %v\n", k, v)
	}
	fmt.Printf("Donations time: : %v\n", time.Since(start))

	// this takes about 7-10 seconds
	//
	ccount := 0 // current count
	ncount := 0 // new count
	for _, name := range firstNames {
		ccount = nameMap[name]

		ncount = ccount + 1
		nameMap[name] = ncount
		if ncount > commonCount {
			common = name
			commonCount = ncount
		}
	}

	fmt.Printf("The most common first name is: %s and it occurs: %v times.\n", common, commonCount)
	fmt.Printf("Most common name time: %v\n", time.Since(start))
}
