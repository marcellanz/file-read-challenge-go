package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	readfile0(os.Args[1], log.New(os.Stdout, "", log.LstdFlags))
}

func readfile0(f string, log *log.Logger) {
	file, err := os.Open(f)
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

	common := ""
	commonCount := 0

	for scanner.Scan() {
		// get all the names
		text := scanner.Text()

		split := strings.SplitN(text, "|", 9) // 10.95
		name := strings.TrimSpace(split[7])
		names = append(names, name)

		// extract first names
		if matches := namePat.FindAllStringSubmatch(name, 1); len(matches) > 0 {
			firstNames = append(firstNames, matches[0][1])
		}

		// extract dates
		chars := strings.TrimSpace(split[4])[:6]
		date := chars[:4] + "-" + chars[4:6]
		dates = append(dates, date)
	}

	log.Printf("Name: %s at index: %v\n", names[0], 0)
	log.Printf("Name: %s at index: %v\n", names[432], 432)
	log.Printf("Name: %s at index: %v\n", names[43243], 43243)

	log.Printf("Name time: %v\n", time.Since(start))
	log.Printf("Total file line count: %v\n", len(names))
	log.Printf("Line count time: : %v\n", time.Since(start))

	dateMap := make(map[string]int)
	for _, date := range dates {
		dateMap[date] += 1
	}
	for k, v := range dateMap {
		log.Printf("Donations per month and year: %v and donation ncount: %v\n", k, v)
	}
	log.Printf("Donations time: : %v\n", time.Since(start))

	// this takes about 7-10 seconds
	//
	nameMap := make(map[string]int)
	ncount := 0 // new count
	for _, name := range firstNames {
		ncount = nameMap[name] + 1
		nameMap[name] = ncount
		if ncount > commonCount {
			common = name
			commonCount = ncount
		}
	}

	log.Printf("The most common first name is: %s and it occurs: %v times.\n", common, commonCount)
	log.Printf("Most common name time: %v\n", time.Since(start))
}
