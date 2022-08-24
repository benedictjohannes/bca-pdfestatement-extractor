package main

import (
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
const corporateBcaRegexStr = `^ESTATEMENT_\d+_(\d{6,6})\.pdf$`

var renamedRegexStr = "\\d+-" + bcaRegexStr

var bcaRegex = regexp.MustCompile("^" + bcaRegexStr + "$")
var renamedRegex = regexp.MustCompile(renamedRegexStr)
var corporateBcaRegex = regexp.MustCompile(corporateBcaRegexStr)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	argsLen := len(os.Args)
	if argsLen > 1 {
		if argsLen > 2 {
			log.Fatalln("this program can only have folder path as first argument")
		}
		wd = os.Args[1]
	}
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Fatalln(err)
	}
	iForeign := 0
	iAlreadyMatched := 0
	iExported := 0
	iFailExport := 0
	for _, file := range files {
		if !file.IsDir() {
			fName := file.Name()
			if strings.Contains(fName, "pdf") {
				matches := corporateBcaRegex.FindStringSubmatch(fName)
				if matches != nil && len(matches) == 2 {
					_, err := time.Parse("200601", matches[1])
					excelFileName := strings.TrimSuffix(fName, ".pdf") + ".xlsx"
					excelFilePath := wd + "/" + excelFileName
					_, err = os.Stat(excelFilePath)
					fmt.Println(excelFilePath)
					if err == nil {
						iAlreadyMatched++
						continue
					}
					f, err := os.Create(excelFilePath)
					if err != nil {
						log.Println("Failed to create excel file", excelFileName)
						iFailExport++
						continue
					}
					defer f.Close()
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
						} else {
							iExported++
						}
					}
				} else {
					iForeign++
				}
			}
		}
	}
	log.Println("Successed in processing PDFs:")
	fmt.Printf("    ForeignFormatted: %d\n", iForeign)
	fmt.Printf("    AlreadyExported : %d\n", iAlreadyMatched)
	fmt.Printf("    SuccessExported : %d\n", iExported)
	if iFailExport > 0 {
		fmt.Printf("    FailToExport    : %d\n", iExported)
	}
}
