package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	firstNamePat := regexp.MustCompile(", \\s*([^, ]+)")
	names := make([]string, 0, 0)
	firstNames := make([]string, 0, 0)
	dates := make([]string, 0, 0)
	commonName := ""
	commonCount := 0

	scanner := bufio.NewScanner(file)

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
			case entry, ok := <-entries:
				if ok {
					if entry.firstName != "" {
						firstNames = append(firstNames, entry.firstName)
					}
					names = append(names, entry.name)
					dates = append(dates, entry.date)
					entry.wg.Done()
				}
			}
		}
	}()

	chunkLen := 64 * 1024
	lines := make([]string, 0, 0)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		if len(lines) == chunkLen {
			wg.Add(len(lines))
			process := lines
			go func() {
				for _, text := range process {
					// get all the names
					e := entry{wg: &wg}
					split := strings.SplitN(text, "|", 9)
					name := strings.TrimSpace(split[7])
					e.name = name

					// extract first names
					if matches := firstNamePat.FindAllStringSubmatch(name, 1); len(matches) > 0 {
						e.firstName = matches[0][1]
					}
					// extract dates
					chars := strings.TrimSpace(split[4])[:6]
					e.date = chars[:4] + "-" + chars[4:6]
					entries <- e
				}
			}()
			lines = make([]string, 0, chunkLen)
		}
	}
	wg.Wait()
	close(entries)

	// report c2: names at index
	fmt.Printf("Name: %s at index: %v\n", names[0], 0)
	fmt.Printf("Name: %s at index: %v\n", names[432], 432)
	fmt.Printf("Name: %s at index: %v\n", names[43243], 43243)
	fmt.Printf("Name time: %v\n", time.Since(start))

	// report c1: total number of lines
	fmt.Printf("Total file line count: %v\n", len(names))
	fmt.Printf("Line count time: : %v\n", time.Since(start))

	// report c3: donation frequency
	dateMap := make(map[string]int)
	for _, date := range dates {
		dateMap[date] += 1
	}
	for k, v := range dateMap {
		fmt.Printf("Donations per month and year: %v and donation count: %v\n", k, v)
	}
	fmt.Printf("Donations time: : %v\n", time.Since(start))

	// report c4: most common firstName
	nameMap := make(map[string]int)
	ncount := 0 // new count
	for _, name := range firstNames {
		ncount = nameMap[name] + 1
		nameMap[name] = ncount
		if ncount > commonCount {
			commonName = name
			commonCount = ncount
		}
	}

	fmt.Printf("The most common first name is: %s and it occurs: %v times.\n", commonName, commonCount)
	fmt.Printf("Most common name time: %v\n", time.Since(start))
	fmt.Fprintf(os.Stderr, "revision: %v, runtime: %v\n", filepath.Base(os.Args[0]), time.Since(start))
}
