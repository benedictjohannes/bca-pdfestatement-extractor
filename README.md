# BCA PDF Estatement Extractor

Duuuh.... ribet ga sih... kalau pengen hitung2 atau proses data dari eStatement Individual BCA?

Harus rapiin dulu satu-satu, copass satu-satu?

Apalagi kalau banyak transaksinya, harus diolah via spreadsheet?

Program ini solusinya! Unduh versi terkini di [releases](https://github.com/benedictjohannes/bca-pdfestatement-extractor/releases/).

## Baca satu eStatement BCA

Tinggal drag and drop file PDF (misal `1271228193Jan2022.pdf`) ke program `processSinglePdf.exe`.

Excel nya akan langsung jadi! Menyenangkan sekali bukan?

## Baca folder berisi eStatement BCA Individual

Tiap eStatement BCA punya format `NomorRekeningMMMYYYY` (individual) dan `ESTATEMENT_NomorRekening_YYYYMM`. Contoh, untuk Januari 2022, nomor rekening 1271228193 akan punya nama file `1271228193Jan2022.pdf` (individual) dan `ESTATEMENT_1271228193_202201.pdf`.

Nah, tinggal download semua eStatement yang mau diproses.

Tunggu dulu, kalau pakai nama bulan, kan urutannya berantakan...?

Tenang, program ini akan mengenali format nama ini, dan satu demi satu:

-   nama file nya akan di rename jadi `YYMM-namaSebelumnya.pdf`. Contoh tadi, akan jadi `2201-1271228193Jan2022.pdf`. Jadi enak mengurutkan filenya toh...
-   file nya akan dibikinin excel nya dengan nama yang sama!

Tinggal drag and drop folder nya ke program `processFolder.exe`

### Untuk Corporate

Yang corporate (rekening Giro) format file nya sudah enak untuk di sort. Jadi program ini hanya akan berusaha untuk export Excel untuk tiap file (dengan skip excel yang sudah ada). 

Tinggal drag and drop folder nya ke program `processCorporateFolder.exe`

# Bagi rekan developer

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](https://pkg.go.dev/github.com/benedictjohannes/bca-pdfestatement-extractor)

Bila ingin menggunakan package ini dalam program Anda, monggo menggunakan exports di dalam package `extractpdf`.

# Bagi eksekutif perbankan yang membaca

1. Zaman sudah canggih
2. Tidak saatnya lagi nih orang harus copass manual satu demi satu
3. PDF eStatement sejak awal adanya rekening itu bagus. Lebih bagus lagi dalam bentuk yang bisa diolah, seperti Excel!
