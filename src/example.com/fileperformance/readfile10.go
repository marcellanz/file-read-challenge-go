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

	namesC := make(chan string)
	firstNamesC := make(chan string)
	datesC := make(chan string)
	wg := sync.WaitGroup{}

	go func() {
		for {
			select {
			case fn, ok := <-firstNamesC:
				if ok {
					firstNames = append(firstNames, fn)
					wg.Done()
				}
			case n, ok := <-namesC:
				if ok {
					names = append(names, n)
					wg.Done()
				}
			case d, ok := <-datesC:
				if ok {
					dates = append(dates, d)
					wg.Done()
				}
			}
		}
	}()

	for scanner.Scan() {
		// get all the names
		text := scanner.Text()
		wg.Add(3)
		go func() {
			split := strings.SplitN(text, "|", 9) // 10.95
			name := strings.TrimSpace(split[7])
			namesC <- name

			// extract first names
			if matches := namePat.FindAllStringSubmatch(name, 1); len(matches) > 0 {
				firstNamesC <- matches[0][1]
			} else {
				wg.Add(-1) // bug
			}

			// extract dates
			chars := strings.TrimSpace(split[4])[:6]
			date := chars[:4] + "-" + chars[4:6]
			datesC <- date
		}()
	}
	fmt.Printf("%v\n", wg)
	wg.Wait()
	fmt.Printf("%v\n", wg)
	close(namesC)
	close(firstNamesC)
	close(datesC)

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
