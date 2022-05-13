package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

const bcaRegexStr = "(\\d+)([a-zA-Z]+)(\\d+)\\.pdf"

var renamedRegexStr = "\\d+-" + bcaRegexStr

var bcaRegex = regexp.MustCompile("^" + bcaRegexStr + "$")
var renamedRegex = regexp.MustCompile(renamedRegexStr)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Fatalln(err)
	}
	iForeign := 0
	iAlreadyMatched := 0
	iProcessed := 0
	for _, file := range files {
		if !file.IsDir() {
			fName := file.Name()
			if strings.Contains(fName, "pdf") {

				matches := bcaRegex.FindStringSubmatch(fName)
				if matches != nil && len(matches) == 4 {
					monthMmm := matches[2]
					year := matches[3]
					t, err := time.Parse("Jan2006", monthMmm+year)
					if err == nil {
						tF := t.Format("0601-")
						oldPath := wd + "/" + fName
						newPath := wd + "/" + tF + fName
						err = os.Rename(oldPath, newPath)
						if err != nil {
							log.Println("Failed:", fName)
						} else {
							log.Println("Renamed:", tF+fName)
							iProcessed++
						}
					} else {
						log.Println("Not time formatted:", fName)
					}
					continue
				}
				fNameSplit := strings.Split(fName, "-")
				_, err = time.Parse("0601", fNameSplit[0])
				if renamedRegex.MatchString(fName) && len(fNameSplit) == 2 && err == nil {
					iAlreadyMatched++
				} else {
					iForeign++
					log.Println("Foreign PDF:", fName)
				}
				continue
			}
		}
	}
	log.Println("Successed in processing PDFs:")
	fmt.Printf("    AlreadyFormatted: %d\n", iAlreadyMatched)
	fmt.Printf("    ForeignFormatted: %d\n", iForeign)
	fmt.Printf("    SuccessRenamed  : %d\n", iProcessed)
}
