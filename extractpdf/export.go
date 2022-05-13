package extractpdf

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

const excelSheet = "transactions"

var last4DigitStrIsNumber = regexp.MustCompile("\\d\\d\\d\\d$")

var excelZeroTime time.Time


// ExportExcel exports the transactions into an *excelize.File
// 
// The resulting Excel file can be written to file, etc
//
// Reasonable column headers and basic formatting included
// to ensure the file is easily readable and workable
func (t Transactions) ExportExcel() (f *excelize.File) {
	excelZeroTime, _ = time.Parse("2006/01/02", "1900/01/01")
	f = excelize.NewFile()
	index := f.NewSheet(excelSheet)
	descriptionStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
	})
	numbersStyle, _ := f.NewStyle(&excelize.Style{
		NumFmt: 2,
	})
	dateStyle, _ := f.NewStyle(&excelize.Style{NumFmt: 14})
	f.SetActiveSheet(index)
	f.SetCellValue(excelSheet, "A1", "Date")
	f.SetColWidth(excelSheet, "A", "A", 12)
	f.SetColStyle(excelSheet, "A", dateStyle)
	f.SetCellValue(excelSheet, "B1", "Description1")
	f.SetColStyle(excelSheet, "B", descriptionStyle)
	f.SetCellValue(excelSheet, "C1", "Description2")
	f.SetColStyle(excelSheet, "C", descriptionStyle)
	f.SetColWidth(excelSheet, "B", "C", 45)
	f.SetCellValue(excelSheet, "D1", "Branch")
	f.SetColWidth(excelSheet, "D", "D", 7)
	f.SetCellValue(excelSheet, "E1", "Change")
	f.SetColStyle(excelSheet, "E", numbersStyle)
	f.SetColWidth(excelSheet, "E", "E", 20)
	f.SetCellValue(excelSheet, "F1", "Direction")
	f.SetColWidth(excelSheet, "F", "F", 10)
	f.SetCellValue(excelSheet, "G1", "Balance")
	f.SetColWidth(excelSheet, "G", "G", 20)
	f.SetColStyle(excelSheet, "G", numbersStyle)
	for i, transaction := range t {
		transaction.writeExcelCell(f, i)
	}
	return
}
func (t *Transaction) writeExcelCell(f *excelize.File, i int) {
	rowIdxStr := strconv.FormatInt(int64(i+2), 10)
	if !t.Date.IsZero() {
		d := t.Date.Sub(excelZeroTime).Hours()/24 + 2
		f.SetCellValue(excelSheet, "A"+rowIdxStr, d)
	}
	maxLines := 1

	if t.Description1 != "" {
		length := len(strings.Split(t.Description1, "\n"))
		if length > maxLines {
			maxLines = length
		}
		f.SetCellValue(excelSheet, "B"+rowIdxStr, t.Description1)
	}

	if t.Description2 != "" {
		length := len(strings.Split(t.Description2, "\n"))
		if length > maxLines {
			maxLines = length
		}
		f.SetCellValue(excelSheet, "C"+rowIdxStr, t.Description2)
	}
	f.SetRowHeight(excelSheet, i+2, float64(maxLines*16))
	if t.Branch != "" {
		f.SetCellValue(excelSheet, "D"+rowIdxStr, t.Branch)
	}
	if t.Change != 0 {
		f.SetCellValue(excelSheet, "E"+rowIdxStr, t.Change)
	}
	if t.DirectionCr != nil {
		if *t.DirectionCr {
			f.SetCellValue(excelSheet, "F"+rowIdxStr, "CR")
		} else {
			f.SetCellValue(excelSheet, "F"+rowIdxStr, "DB")
		}
	}
	if t.Balance != 0 {
		f.SetCellValue(excelSheet, "G"+rowIdxStr, t.Balance)
	}
}
