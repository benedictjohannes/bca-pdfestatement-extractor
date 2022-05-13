package extractpdf

import (
	"strconv"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
)

type transaction struct {
	Date         time.Time
	Description1 string
	Description2 string
	Branch       string
	Change       float64
	DirectionCr  *bool
	Balance      float64
}

type transactions []*transaction

type leftCol float64
type rightCol float64

func (c leftCol) Is(x float64) bool {
	diff := c - leftCol(x)
	if diff < 0 {
		diff = diff * -1
	}
	return diff < 1.0
}
func (c rightCol) Is(x float64) bool {
	return rightCol(x) > c
}

const DateCol leftCol = 46.04
const BalanceCol rightCol = 470.0
const ChangeAmountCol rightCol = 340.0
const DbCrCol leftCol = 447.22
const Description1Col leftCol = 92.61
const Description2Col leftCol = 196.71

func IngestRow(prevT *transaction, row *pdf.Row, year string) (isNew bool, t *transaction) {
	words := make([]pdf.Text, len(row.Content))
	for i, w := range row.Content {
		words[i] = w
	}
	if len(words) < 2 {
		if prevT == nil {
			return
		}
		t = prevT
		readSupplementary(t, words)
		return

	}
	firstWord := words[0]
	date, dateErr := time.Parse("02/01/2006", firstWord.S+"/"+year)
	hasDate := dateErr == nil && DateCol.Is(firstWord.X)
	if !hasDate {
		if prevT == nil {
			return
		}
		t = prevT
	} else {
		isNew = true
		words = words[1:]
		t = &transaction{
			Date: date,
		}
	}

	lastWord := words[len(words)-1]
	balance, balanceErr := strconv.ParseFloat(
		strings.Replace(lastWord.S, ",", "", -1),
		32)
	hasBalance := balanceErr == nil && BalanceCol.Is(lastWord.X)
	if hasBalance {
		t.Balance = balance
		words = words[:len(words)-1]
	}

	readSupplementary(t, words)
	return
}

func readSupplementary(t *transaction, words []pdf.Text) {
	for i, word := range words {
		if ChangeAmountCol.Is(word.X) {
			amount, amountErr := strconv.ParseFloat(
				strings.Replace(word.S, ",", "", -1),
				32)
				if amountErr == nil {
				isCr := len(words) == i+1 || !(words[i+1].S == "DB")
				t.DirectionCr = &isCr
				t.Change = amount
			}
		}
		if Description1Col.Is(word.X) {
			if t.Description1 == "" {
				t.Description1 = word.S
			} else {
				t.Description1 = t.Description1 + "\n" + word.S
			}
		}
		if Description2Col.Is(word.X) {
			if t.Description2 == "" {
				t.Description2 = word.S
			} else {
				t.Description2 = t.Description2 + "\n" + word.S
			}
		}
	}
}
