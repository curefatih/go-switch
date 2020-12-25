package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type Pattern struct {
	beginsAt int
	endsAt   int
	name     string
}

type ByBeginsAt []Pattern

func (a ByBeginsAt) Len() int           { return len(a) }
func (a ByBeginsAt) Less(i, j int) bool { return a[i].beginsAt < a[j].beginsAt }
func (a ByBeginsAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type Source struct {
	content string
}

func newPattern(name string, begins, ends int) *Pattern {
	p := Pattern{beginsAt: begins, endsAt: ends, name: name}
	return &p
}

func newSource(content string) *Source {
	s := Source{content: content}
	return &s
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func definePatterns(begins, ends string, data *[]byte) *ByBeginsAt {

	result := make(ByBeginsAt, 0)
	var inside bool = false
	var patternName string
	var patternBeginsAt int

	for i := 0; i < len(*data); i++ {
		if !inside {
			// if is not inside the pattern looks for the begining

			if (*data)[i] == begins[0] {
				for j := 0; j < len(begins); j++ { // should loop continue min of (i+len(begins)) and len(data)
					if begins[j] != (*data)[i+j] {
						break
					}
					if begins[j] == (*data)[i+j] && len(begins)-1 == j {
						fmt.Printf("Found begins in line: %d\n", i)
						inside = true
						patternBeginsAt = i
						i = i + j
					}
				}
			}
		} else {
			// if is inside the pattern looks for the ending
			if (*data)[i] == ends[0] {
				for j := 0; j < len(ends); j++ {
					if ends[j] != (*data)[i+j] {
						break
					}
					if ends[j] == (*data)[i+j] && len(ends)-1 == j {
						fmt.Printf("Found ends in line: %d, pattern: %s\n", i, patternName)

						// create pattern info
						pattern := newPattern(strings.Trim(patternName, " "), patternBeginsAt, i+j)
						result = append(result, *pattern)

						patternName = "" // reset patternName
						inside = false   // reset is inside info
					}
				}
			} else {
				patternName += string((*data)[i])
			}

		}

	}

	return &result
}

func defineSource(data *[]byte) *map[string]*Source {
	result := make(map[string]*Source)

	begins := "$$pattern:"
	ends := "$$end"
	var inside bool = false
	var patternName string
	var content string

	for i := 0; i < len(*data); i++ {
		if !inside {
			// if is not inside the pattern looks for the begining

			if (*data)[i] == begins[0] {
				for j := 0; j < len(begins); j++ { // should loop continue min of (i+len(begins)) and len(data)
					if begins[j] != (*data)[i+j] {
						break
					}
					if begins[j] == (*data)[i+j] && len(begins)-1 == j {
						j = j + 1
						for {
							if len(*data) <= i+j {
								break
							}
							if (*data)[i+j] != '\n' {
								patternName += string((*data)[i+j])
								j++
							} else {
								break
							}
						}
						fmt.Printf("Found source begins in line: %d, name: %s\n", i, patternName)

						inside = true
						i = i + j
					}
				}
			}
		} else {
			// if is inside the pattern looks for the ending
			if (*data)[i] == ends[0] {
				for j := 0; j < len(ends); j++ {
					if ends[j] != (*data)[i+j] {
						break
					}
					if ends[j] == (*data)[i+j] && len(ends)-1 == j {
						// create pattern info
						source := newSource(strings.TrimSpace(content))

						result[strings.TrimSpace(patternName)] = source

						fmt.Printf("Found source ends in line: %d, pattern: %s\n", i, patternName)
						content = ""     // reset content
						patternName = "" // reset patterName
						inside = false   // reset is inside info

					}
				}
			} else {
				content += string((*data)[i])
			}

		}

	}

	return &result
}

func changeWithSource(sources *map[string]*Source, patterns *ByBeginsAt, targetData *[]byte) *string {

	sort.Sort(ByBeginsAt(*patterns))
	var changedContent string
	currentIndex := 0
	for i := 0; i < len(*patterns); i++ {
		if val, ok := (*sources)[(*patterns)[i].name]; ok {
			changedContent += string((*targetData)[currentIndex:(*patterns)[i].beginsAt])
			changedContent += val.content
			currentIndex = (*patterns)[i].endsAt + 1
		} else {
			fmt.Println("->cant find!")
		}

	}

	return &changedContent
}

func main() {

	var sourceFileName string
	flag.StringVar(&sourceFileName, "sourceFile", "", "defines file name")

	var targetFileName string
	flag.StringVar(&targetFileName, "targetFile", "", "defines file name")

	patternBeginsPtr := flag.String("patternBegins", "${", "defines pattern begining")
	patternEndsPtr := flag.String("patternEnds", "}", "defines pattern begining")

	flag.Parse()

	targetData, err := ioutil.ReadFile(targetFileName)
	check(err)

	sourceData, err := ioutil.ReadFile(sourceFileName)
	check(err)

	sources := defineSource(&sourceData)
	patterns := definePatterns(*patternBeginsPtr, *patternEndsPtr, &targetData)
	res := changeWithSource(sources, patterns, &targetData)

	f, err := os.Create("new_" + targetFileName)
	check(err)

	defer f.Close()

	_, err = f.WriteString(*res)
	check(err)

	fmt.Println("Done!")
}
