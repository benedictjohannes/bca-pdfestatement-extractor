package main

import (
	"bytes"
	"encoding/base64"
	"strings"
	"syscall/js"

	"github.com/benedictjohannes/bca-pdfestatementindividual-extractor/extractpdf"
)

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js ./src

func exportExcel(this js.Value, inputs []js.Value) any {
	if len(inputs) != 2 {
		return "ERR: must have 2 argument: fileName and its arrayBuffer"
	}
	fileName := inputs[0].String()
	println("fileName")
	println(fileName)
	if !strings.HasSuffix(fileName, ".pdf") {
		return "ERR: fileName does not ends with .pdf"
	}
	inJsUint8Input := inputs[1]
	inJsUint8Arr := make([]uint8, inJsUint8Input.Get("byteLength").Int())
	js.CopyBytesToGo(inJsUint8Arr, inJsUint8Input)
	transactions, err := extractpdf.ProcessPdfFromBytes(inJsUint8Arr)
	if err != nil {
		return "ERR: extractTransactions errored: " + err.Error()
	}
	excelFile := transactions.ExportExcel()
	var excelBuf bytes.Buffer
	_, err = excelFile.WriteTo(&excelBuf)
	if err != nil {
		return "ERR: exportingExcel errored: " + err.Error()
	}
	excelBytes := excelBuf.Bytes()
	excelB64Str := base64.RawStdEncoding.EncodeToString(excelBytes)
	return excelB64Str
}

func main() {
	c := make(chan struct{})
	js.Global().Set("excelExport", js.FuncOf(exportExcel))
	<-c
}
