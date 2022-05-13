package extractpdf

import (
	"github.com/ledongthuc/pdf"
)

type rowSortable []*pdf.Row

func (rows rowSortable) Len() int {
	return len(rows)
}
func (a rowSortable) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a rowSortable) Less(i, j int) bool {
	return a[i].Position > a[j].Position
}
