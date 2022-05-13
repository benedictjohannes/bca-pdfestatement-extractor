package extractPdf

import (
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ledongthuc/pdf"
)

var yearRegex = regexp.MustCompile("\\d\\d\\d\\d")
var months = []string{
	"JANUARI",
	"FEBRUARI",
	"MARET",
	"APRIL",
	"MEI",
	"JUNI",
	"JULI",
	"AGUSTUS",
	"NOVEMBER",
	"DESEMBER",
}

func ProcessPdfFromPath(path string) (transactions, error) {
	f, pdfR, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, err
	}
	return processPdf(pdfR)
}
func ProcessPdfFromReader(r io.Reader) (transactions, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ProcessPdfFromBytes(b)
}
func ProcessPdfFromBytes(b []byte) (transactions, error) {
	bytesR := bytes.NewReader(b)
	pdfR, err := pdf.NewReader(bytesR, bytesR.Size())
	if err != nil {
		return nil, err
	}
	return processPdf(pdfR)
}

func processPdf(pdfR *pdf.Reader) (transactions, error) {
	totalPage := pdfR.NumPage()
	transactions := make([]*transaction, 0)
	var currentTransaction *transaction = nil
	var isNew bool = false
	year := "1900"
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := pdfR.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		unsortedRows, _ := p.GetTextByRow()
		sortedRows := make(rowSortable, len(unsortedRows))
		for i, row := range unsortedRows {
			sortedRows[i] = row
		}
		sort.Sort(sortedRows)
		aftTanggal := false
		for _, row := range sortedRows {
			if aftTanggal {
				isNew, currentTransaction = IngestRow(currentTransaction, row, year)
				if isNew {
					transactions = append(transactions, currentTransaction)
				}
			} else {
				m := 0
				for wordIndex, word := range row.Content {
					if year == "1900" {
						for _, m := range months {
							if strings.Contains(word.S, m) {
								yearStr := strings.TrimPrefix(word.S, m+" ")
								if yearRegex.MatchString(yearStr) {
									_, err := strconv.ParseInt(yearStr, 10, 32)
									if err != nil {
										year = "1900"
									} else {
										year = yearStr
									}
								}
							}
						}
					}
					if strings.Contains("TANGGAL", word.S) && wordIndex == 0 {
						m++
					}
					if strings.Contains("KETERANGAN", word.S) && wordIndex == 1 {
						m++
					}
					if strings.Contains("CBG", word.S) && wordIndex == 2 {
						m++
					}
					if strings.Contains("MUTASI", word.S) && wordIndex == 3 {
						m++
					}
					if strings.Contains("SALDO", word.S) && wordIndex == 4 {
						m++
					}
					if m == 5 {
						aftTanggal = true
					}
				}
			}
		}
	}
	return transactions, nil
}
