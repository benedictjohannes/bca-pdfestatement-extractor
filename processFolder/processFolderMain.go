package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/benedictjohannes/bca-pdfestatementindividual-extractor/extractpdf"
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
	argsLen := len(os.Args)
	if argsLen > 1 {
		if argsLen > 2 {
			log.Fatalln("this program expects none or exactly 1 argument: the folder to process")
		}
		wd = os.Args[1]
	}
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Fatalln(err)
	}
	iForeign := 0
	iAlreadyMatched := 0
	iRenamed := 0
	iExported := 0
	iFailExport := 0
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
						newName := tF + fName
						newPath := wd + "/" + newName
						err = os.Rename(fName, newPath)
						if err != nil {
							log.Println("Failed:", fName)
						} else {
							log.Println("Renamed:", newName)
							iRenamed++
							t, err := extractpdf.ProcessPdfFromPath(newPath)
							if err != nil {
								log.Println("Failed to extract transactions from", fName)
							}
							excel := t.ExportExcel()
							excelFileName := strings.TrimSuffix(newName, ".pdf") + ".xlsx"
							err = excel.SaveAs(excelFileName)
							if err != nil {
								iFailExport++
								log.Println("Failed to save excel for file", fName)
							} else {
								iExported++
							}

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
					excelFileName := strings.TrimSuffix(fName, ".pdf") + ".xlsx"
					excelFilePath := wd + "/" + excelFileName
					f, err := os.Create(excelFilePath)
					if err != nil {
						if !errors.Is(err, os.ErrExist) {
							log.Println("Failed to create excel file", excelFileName)
							iFailExport++
						}
					} else {
						defer f.Close()
					}
					t, err := extractpdf.ProcessPdfFromPath(wd + "/" + fName)
					if err != nil {
						log.Println("Failed to extract transactions from", err)
						iFailExport++
					} else {
						e := t.ExportExcel()
						_, err := e.WriteTo(f)
						if err != nil {
							log.Println("Failed to save excel file", err)
							iFailExport++
						}
					}
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
	fmt.Printf("    SuccessRenamed  : %d\n", iRenamed)
	fmt.Printf("    SuccessExported : %d\n", iExported)
	if iFailExport > 0 {
		fmt.Printf("    FailToExport    : %d\n", iExported)
	}
}
