# BCA PDF Estatement Extractor

Duuuh.... ribet ga sih... kalau pengen hitung2 atau proses data dari eStatement Individual BCA?

Harus rapiin dulu satu-satu, copass satu-satu?

Apalagi kalau banyak transaksinya, harus diolah via spreadsheet?

Program ini solusinya! Unduh versi terkini di [releases](https://github.com/benedictjohannes/bca-pdfestatement-extractor/releases/). Atau, gunakan versi web di [https://bcapdfestatementtoexcel.web.app/](https://bcapdfestatementtoexcel.web.app/).

## Export Excel satu eStatement BCA

Tinggal drag and drop file PDF (misal `1271228193Jan2022.pdf`) ke program `excelkanBCAeStatement.exe`.  
Atau, filenya dipanggil sebagai argument pertama ke executable nya, `./excelkanBCAeStatement /path/to/folder`

Excel nya akan langsung jadi! Menyenangkan sekali bukan?

## Folder berisi eStatement BCA Individual: Urutkan dan Export Excel

Tiap eStatement BCA punya format `NomorRekeningMMMYYYY`. Contoh, untuk Januari 2022, nomor rekening 1271228193 akan punya nama file `1271228193Jan2022.pdf`.

Nah, tinggal download semua eStatement yang mau diproses.

Tunggu dulu, kalau pakai nama bulan, kan urutannya berantakan...?

Tenang, program ini akan mengenali format nama ini, dan satu demi satu:

-   nama file nya akan di rename jadi `YYMM-namaSebelumnya.pdf`. Contoh tadi, akan jadi `2201-1271228193Jan2022.pdf`. Jadi enak mengurutkan filenya berdasarkan nama toh...
-   file nya akan dibikinin excel nya dengan nama yang sama!

Tinggal drag and drop folder nya ke program `excelkanBCAeStatement.exe`.  
Atau, foldernya dipanggil sebagai argument pertama ke executable nya, `./excelkanBCAeStatement /path/to/folder`

## Folder berisi eStatement BCA Corporate: Urutkan dan Export Excel

Yang corporate (rekening Giro) format nama file nya sudah sesuai untuk di urutkan. Jadi program ini hanya akan berusaha untuk export Excel untuk tiap file (tanpa menimpa file excel yang sudah ada). 

Tinggal drag and drop folder nya ke program `excelkanBCAeStatement.exe`.  
Atau, foldernya dipanggil sebagai argument pertama ke executable nya, `./excelkanBCAeStatement /path/to/folder`

# Bagi rekan developer

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](https://pkg.go.dev/github.com/benedictjohannes/bca-pdfestatement-extractor)

Bila ingin menggunakan package ini dalam program Anda, monggo menggunakan exports di dalam package `extractpdf`.

# Bagi eksekutif perbankan yang membaca

Export dari web e banking BCA itu, masih harus dikenalin menggunakan Text-To-Column (dalam Excel). Dan, setelah di text-to-column formatnya cukup berantakan. Tidak semua pengguna puas hanya cek saldo di ATM. Tampilan perbankan nya harus lebih bagus, dengan format export yang dapat digunakan secara mudah (tanpa harus manual copass dan cleaning). Ini kan bisa dilakukan oleh mesin, mengapa harus manusia yang melakukannya? Agar tercipta lapangan kerja yang tidak efisien?

Program ini dibuat untuk mempermudah pengguna rekening BCA (sesuatu yang harusnya kalian sudah sadar untuk selalu lakukan).