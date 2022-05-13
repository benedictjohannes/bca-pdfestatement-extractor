package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/benedictjohannes/bca-pdfestatementindividual-extractor/extractpdf"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("This program expects exactly one argument: the eStatement PDF file to process")
	}
	pdfFileName := os.Args[1]
	if !strings.HasSuffix(pdfFileName, ".pdf") {
		panic("file should has .pdf extension")
	}
	transactions, err := extractpdf.ProcessPdfFromPath(pdfFileName)
	if err != nil {
		panic(err)
	}
	excelFile := transactions.ExportExcel()
	excelFileName := strings.TrimSuffix(pdfFileName, ".pdf") + ".xlsx"
	err = excelFile.SaveAs(excelFileName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully saved", excelFileName, "at", time.Now().Format("15:04:05"))
	return
}
