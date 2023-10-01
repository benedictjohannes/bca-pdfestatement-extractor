package extractpdf

import (
	"strconv"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
)

type Transaction struct {
	Date         time.Time `json:"date,omitempty"`
	Description1 string    `json:"description1,omitempty"`
	Description2 string    `json:"description2,omitempty"`
	Branch       string    `json:"branch,omitempty"`
	Change       float64   `json:"change,omitempty"`
	DirectionCr  *bool     `json:"directionCr,omitempty"`
	Balance      float64   `json:"balance,omitempty"`
}

type Transactions []*Transaction

type leftCol float64
type rightCol float64

func (c leftCol) Is(x float64) bool {
	diff := c - leftCol(x)
	if diff < 0 {
		diff = diff * -1
	}
	return diff < 5.0
}
func (c rightCol) Is(x float64) bool {
	return rightCol(x) > c
}

const dateCol leftCol = 46.04
const balanceCol rightCol = 470.0
const changeAmountCol rightCol = 340.0
const description1Col leftCol = 92.61
const description2Col leftCol = 196.71
const summaryFirstCol leftCol = 180.18

// a row with a new date signifies a new transaction
//
// as a PDF text row might not be a new transaction,
// but adds detail to the previous transaction, prevT
// is added as argument to add transaction detail
//
// the returned isNew tells whether the returned
// *transaction is a new transaction that should
// be added to a transaction slice
func IngestRow(prevT *Transaction, row *pdf.Row, year string) (isNew bool, t *Transaction, shouldStopProcessing bool) {
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
	hasDate := dateErr == nil && dateCol.Is(firstWord.X)
	if !hasDate {
		if prevT == nil {
			return
		}
		t = prevT
	} else {
		isNew = true
		words = words[1:]
		t = &Transaction{
			Date: date,
		}
	}

	lastWord := words[len(words)-1]
	balance, balanceErr := strconv.ParseFloat(
		strings.Replace(lastWord.S, ",", "", -1),
		32)
	hasBalance := balanceErr == nil && balanceCol.Is(lastWord.X)
	if hasBalance {
		t.Balance = balance
		words = words[:len(words)-1]
	}

	shouldStopProcessing = readSupplementary(t, words)
	return
}

// readSupplementary try to read words in a row that are apart of
// date and balance information
func readSupplementary(t *Transaction, words []pdf.Text) (stopProcessingNext bool) {
	for i, word := range words {
		if i==0 && word.S == "SALDO AWAL" && summaryFirstCol.Is(word.X) {
			return true
		}
		if changeAmountCol.Is(word.X) {
			amount, amountErr := strconv.ParseFloat(
				strings.Replace(word.S, ",", "", -1),
				32)
			if amountErr == nil {
				isCr := len(words) == i+1 || !(words[i+1].S == "DB")
				t.DirectionCr = &isCr
				t.Change = amount
			}
		}
		if description1Col.Is(word.X) {
			if t.Description1 == "" {
				t.Description1 = word.S
			} else {
				t.Description1 = t.Description1 + "\n" + word.S
			}
		}
		if description2Col.Is(word.X) {
			if t.Description2 == "" {
				t.Description2 = word.S
			} else {
				t.Description2 = t.Description2 + "\n" + word.S
			}
		}
	}
	return
}
