package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	names := make([]string, 0, 0)
	firstNames := make([]string, 0, 0)
	dates := make([]string, 0, 0)

	nameMap := make(map[string]int)
	dateMap := make(map[string]int)

	start := time.Now()
	scanner := bufio.NewScanner(file)

	common := ""
	commonCount := 0

	for scanner.Scan() {
		// get all the names
		text := scanner.Text()

		split := strings.SplitN(text, "|", 9) // 10.95
		name := strings.TrimSpace(split[7])
		names = append(names, name)

		// extract first names
		if name != "" {
			startOfName := strings.Index(name, ", ") + 2

			var firstName string
			if endOfName := strings.Index(name[startOfName:], " "); endOfName < 0 {
				firstName = name[startOfName:]
			} else {
				firstName = name[startOfName : startOfName+endOfName]
			}
			firstNames = append(firstNames, firstName)

			ccount := nameMap[firstName]
			nameMap[firstName] = ccount + 1
			if ccount+1 > commonCount {
				common = firstName
				commonCount = ccount + 1
			}
		}

		// extract dates
		chars := strings.TrimSpace(split[4])[:6]
		date := chars[:4] + "-" + chars[4:6]
		dates = append(dates, date)
		dateMap[date] += 1
	}

	fmt.Printf("Name: %s at index: %v\n", names[0], 0)
	fmt.Printf("Name: %s at index: %v\n", names[432], 432)
	fmt.Printf("Name: %s at index: %v\n", names[43243], 43243)

	fmt.Printf("Name time: %v\n", time.Since(start))
	fmt.Printf("Total file line count: %v\n", len(names))
	fmt.Printf("Line count time: : %v\n", time.Since(start))

	for k, v := range dateMap {
		fmt.Printf("Donations per month and year: %v and donation ncount: %v\n", k, v)
	}
	fmt.Printf("Donations time: : %v\n", time.Since(start))

	fmt.Printf("The most common first name is: %s and it occurs: %v times.\n", common, commonCount)
	fmt.Printf("Most common name time: %v\n", time.Since(start))
}
