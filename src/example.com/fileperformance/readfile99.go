package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	file, err := os.Open("indiv18/itcont.txt")
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
	//dateMap := make(map[string]int)
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

			//ncount := nameMap[firstName] + 1
			//nameMap[firstName] = ncount
			//if ncount > commonCount {
			//	common = firstName
			//	commonCount = ncount
			//}
		}

		// extract dates
		chars := strings.TrimSpace(split[4])[:6]
		date := chars[:4] + "-" + chars[4:6]
		//dateMap[date] += 1
		dates = append(dates, date)
	}

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
	//fnlen := len(firstNames)
	for _, name := range firstNames {
		ccount = nameMap[name]
		//if ccount < 0 {
		//	continue
		//}
		//// check if we can beat the current max count with the rest of names
		//if (ccount + (fnlen - i)) < commonCount {
		//	nameMap[name] = -1
		//	continue
		//}

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

func genSplit(s, sep string, n int) string {
	n--
	i := 0
	r := ""
	for i < n {
		m := strings.Index(s, sep)
		if m < 0 {
			break
		}
		r = s[:m]
		s = s[m+len(sep):]
		i++
	}
	return r
}
