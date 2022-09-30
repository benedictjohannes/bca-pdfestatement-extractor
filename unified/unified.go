package main

import (
	"fmt"
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
		arg1 := os.Args[1]
		info, err := os.Stat(arg1)
		if err != nil {
			log.Fatalln("failed to stat file:", arg1)
		}
		if !info.IsDir() {
			processFile(arg1)
			return
		}
		wd = arg1
	}
	files, err := os.ReadDir(wd)
	if err != nil {
		log.Fatalln(err)
	}
	cIndividualAlreadyRenamed := 0
	cIndividualSuccessRenamed := 0
	cIndividualFailedRename := 0
	cIndividualSuccessExported := 0
	cIndividualAlreadyExported := 0
	cIndividualFailedExport := 0
	cForeign := 0
	cCorporateAlreadyExported := 0
	cCorporateSuccessExported := 0
	cCorporateFailedExport := 0
	for _, file := range files {
		if !file.IsDir() {
			fName := file.Name()
			if !strings.Contains(fName, "pdf") {
				continue
			}
			// BCA individual PDF - rename and tryExport Excel
			if matches := bcaRegex.FindStringSubmatch(fName); len(matches) == 4 {
				monthMmm := matches[2]
				year := matches[3]
				ts, err := time.Parse("Jan2006", monthMmm+year)
				if err != nil {
					log.Println("Foreign individual PDF (illegal time format):", fName)
					cForeign++
					continue
				}
				tF := ts.Format("0601-")
				newName := tF + fName
				newPath := wd + "/" + newName
				err = os.Rename(fName, newPath)
				if err != nil {
					log.Println("Failed rename:", fName)
					cIndividualFailedRename++
					continue
				}
				log.Println("Renamed:", newName)
				cIndividualSuccessRenamed++
				excelFileName := strings.TrimSuffix(newName, ".pdf") + ".xlsx"
				excelFilePath := wd + "/" + excelFileName
				_, err = os.Stat(excelFilePath)
				if err == nil {
					log.Println("Excel already exist:", excelFileName)
					cIndividualAlreadyExported++
					continue
				}
				t, err := extractpdf.ProcessPdfFromPath(newPath)
				if err != nil {
					log.Println("Failed to extract transactions from", fName, "; err:", err)
					cIndividualFailedExport++
					continue
				}
				excel := t.ExportExcel()
				err = excel.SaveAs(excelFilePath)
				if err != nil {
					log.Println("Failed to save excel for file", fName, "; err:", err)
					cIndividualFailedExport++
				} else {
					log.Println("Succeed to save excel for file", fName, "; err:", err)
					cIndividualSuccessExported++
				}
				continue
			}

			// BCA individual PDF - already renamed - tryExport Excel
			if fNameSplit := strings.Split(fName, "-"); len(fNameSplit) == 2 && renamedRegex.MatchString(fName) {
				_, err := time.Parse("0601", fNameSplit[0])
				if err != nil {
					log.Println("Foreign individual renamed PDF (illegal time format):", fName)
					cForeign++
					continue
				}
				cIndividualAlreadyRenamed++
				excelFileName := strings.TrimSuffix(fName, ".pdf") + ".xlsx"
				excelFilePath := wd + "/" + excelFileName
				_, err = os.Stat(excelFilePath)
				if err == nil {
					log.Println("Excel already exist:", excelFileName)
					cIndividualAlreadyExported++
					continue
				}
				t, err := extractpdf.ProcessPdfFromPath(wd + "/" + fName)
				if err != nil {
					log.Println("Failed to extract transactions from", fName, "; err:", err)
					cIndividualFailedExport++
					continue
				}
				e := t.ExportExcel()
				err = e.SaveAs(excelFilePath)
				if err != nil {
					log.Println("Failed to save excel for file", fName, "; err:", err)
					cIndividualFailedExport++
				} else {
					log.Println("Succeed to save excel for file", fName, "; err:", err)
					cIndividualSuccessExported++
				}
				continue
			}

			// BCA corporate PDF - tryExport Excel
			if matches := corporateBcaRegex.FindStringSubmatch(fName); len(matches) == 2 {
				_, err := time.Parse("200601", matches[1])
				if err != nil {
					log.Println("Foreign corporate PDF (illegal time format):", fName)
					cForeign++
					continue
				}
				excelFileName := strings.TrimSuffix(fName, ".pdf") + ".xlsx"
				excelFilePath := wd + "/" + excelFileName
				_, err = os.Stat(excelFilePath)
				if err == nil {
					log.Println("Excel already exist:", excelFileName)
					cCorporateAlreadyExported++
					continue
				}
				t, err := extractpdf.ProcessPdfFromPath(wd + "/" + fName)
				if err != nil {
					log.Println("Failed to extract transactions from", err)
					cCorporateFailedExport++
					continue
				}
				e := t.ExportExcel()
				err = e.SaveAs(excelFilePath)
				if err != nil {
					log.Println("Failed to save excel file", err)
					cCorporateFailedExport++
				} else {
					log.Println("Succeed to save excel for file", fName, "; err:", err)
					cCorporateSuccessExported++
				}
				continue
			}
			log.Println("Foreign PDF:", fName)
			cForeign++
		}
	}
	cIndividuals := cIndividualAlreadyRenamed +
		cIndividualSuccessRenamed +
		cIndividualFailedRename +
		cIndividualSuccessExported +
		cIndividualAlreadyExported +
		cIndividualFailedExport
	if cIndividuals > 0 {
		log.Println("Successed in processing individual PDFs:")
		fmt.Printf("    Already Renamed  : %d\n", cIndividualAlreadyRenamed)
		fmt.Printf("    Success Rename   : %d\n", cIndividualSuccessRenamed)
		fmt.Printf("    Failed  Rename   : %d\n", cIndividualFailedRename)
		fmt.Printf("    Already Exported : %d\n", cIndividualAlreadyExported)
		fmt.Printf("    Success Exported : %d\n", cIndividualSuccessExported)
		fmt.Printf("    Fail To Export   : %d\n", cIndividualFailedExport)
	}
	cCorporates := cCorporateAlreadyExported +
		cCorporateSuccessExported +
		cCorporateFailedExport
	if cCorporates > 0 {
		log.Println("Successed in processing corporate PDFs:")
		fmt.Printf("    Already Exported : %d\n", cCorporateAlreadyExported)
		fmt.Printf("    Success Exported : %d\n", cCorporateSuccessExported)
		fmt.Printf("    Fail To Export   : %d\n", cCorporateFailedExport)
	}
	log.Printf("Files with foreign PDF names: %d\n", cCorporateFailedExport)
}

func processFile(fileName string) {
	if !strings.HasSuffix(fileName, ".pdf") {
		log.Fatalln("file should has .pdf extension")
	}
	transactions, err := extractpdf.ProcessPdfFromPath(fileName)
	if err != nil {
		log.Fatalln("failed to process PDF file:", err)
	}
	excelFile := transactions.ExportExcel()
	excelFileName := strings.TrimSuffix(fileName, ".pdf") + ".xlsx"
	err = excelFile.SaveAs(excelFileName)
	if err != nil {
		log.Fatalln("failed to write excel file", excelFileName)
	}
	log.Println("Successfully saved", excelFileName)
}
