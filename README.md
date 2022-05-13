# BCA Individual PDF Estatement Extractor

Duuuh.... ribet ga sih... kalau pengen hitung2 atau proses data dari eStatement Individual BCA?

Harus rapiin dulu satu-satu, copass satu-satu?

Apalagi kalau banyak transaksinya, harus diolah via spreadsheet?

Program ini solusinya!

## Baca satu eStatement BCA

Tinggal drag and drop file PDF (misal `1271228193Jan2022.pdf`) ke program `processSinglePdf.exe`.

Excel nya akan langsung jadi! Menyenangkan sekali bukan?

## Baca folder berisi eStatement BCA

Tiap eStatement BCA punya format `NomorRekeningMMMYYYY`. Contoh, untuk Januari 2022, nomor rekening 1271228193 akan punya nama file `1271228193Jan2022.pdf`

Nah, tinggal download semua eStatement yang mau diproses.

Tunggu dulu, kalau pakai nama bulan, kan urutannya berantakan...?

Tenang, program ini akan mengenali format nama ini, dan satu demi satu:

-   nama file nya akan di rename jadi `YYMM-namaSebelumnya.pdf`. Contoh tadi, akan jadi `2201-1271228193Jan2022.pdf`. Jadi enak mengurutkan filenya toh...
-   file nya akan dibikinin excel nya dengan nama yang sama!

Tinggal drag and drop folder nya ke program `processFolder.exe`

# Bagi eksekutif perbankan yang membaca

1. Zaman sudah canggih
2. Tidak saatnya lagi nih orang harus copass manual satu demi satu
3. PDF eStatement sejak awal adanya rekening itu bagus. Lebih bagus lagi dalam bentuk yang bisa diolah, seperti Excel!
