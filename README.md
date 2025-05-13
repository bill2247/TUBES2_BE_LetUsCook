# TUBES2_BE_LetUsCook
By Kelompok 51 Let Us Cook
Website untuk pencarian resep little alchemy 2
<br>

## Contributors
<div align="center">

| **NIM**  | **Nama** |
| ------------- |:-------------:|
| 13523066   | Muhammad Ghifary Komara Putra |
| 13523072   | Sabilul Huda |
| 13523080   | Diyah Susan Nugrahani |

</div>

## Konsep yang digunakan

Tugas pencarian resepp ini menggunakan konsep pencarian graf dengan algoritma BFS dan DFS
untuk mencari rute resep dari input query yang ingin dicari. Tujuan objektif dari tugas ini
adalah mencari rute untuk mencapai resep tertentu, dimana data resep diperoleh dari hasil 
scrapping website little alchemy 2

## Features
1. Mode Pencarian  
Dalam tugas ini terdapat dua mode pencarian yaitu **single recipe** dan **multiple recipes**. Untuk single recipes akan dicari jalur resep terpendek sesuai dengan wuery yang diberikan. Sedangkan untuk multiple recipes akan dicari alternatif cara mencapai resep tertentu dengan jumlah resepnya sesuai dengan inputan pengguna.
2. Algoritma Pencarian  
Dalam pencarian resep terdapat dua opsi algoritma ynag dapat dipilih yaitu **BFS, DFS, atau bidirectional**. Algoritma BFS menggunakan pendekatan penelusuran graf secara menyebar sedangkan algoritma DFS menggunakan pendekatan penelusuran graf secara mendalam. Algoritma bidirectional melakukan pencarian dari root dan leaf sehingga pencarian akan relatif lebih cepat.
3. Pohon Solusi Resep  
Pohon solusi resep yang ditampilkan dapat diexpand untuk melihat resep secara detail dan juga dapat ditutup untuk melihat resep secara umum. Hal ini memiliki keunggulan memberikan tampilan yang lebih fleksibel sesuai kebutuhan pengguna.
4. Detail Informasi  
Di bagian bawah solusi terdapat detail infromasi terkait jumlah node yang dikunjungi serta durasi pencarian resep.

## Teknologi yang digunakan
Untuk bagian back end, digunakan bahasa go dan framework gin untuk tugas ini. Algoritma diimplementasikan secara terpisah di folder algorithm. Algoritma tersebut lalu digunakan di main.go dengan cara mengimport nya. Di dalam main.go terdapat mekanisme yang menghubungkan antara front end dan back end dengan framework gin.

## Cara Penggunaan
1. Akses website LetUsCook pada tautan berikut
``` sh
https://let-us-cook-new.vercel.app
```
2. Masukkan query resep yang ingin dicari
3. Pilih mode single atau multiple
4. Pilih algoritma yang ingin digunakan, antara BFS, DFS, atau bidirectional
5. Tekan tombol search
6. Hasil pohon solusi akan ditampilkan beserta dnegan informasi detail node ynag dikunjungi dan durasi
![Imgur](https://imgur.com/n7tnAU8.jpg)

Catatan : jika mengalami kendala dalam mengakses website melalui internet, maka website dapat diakses dengan cara menjalankan front end dan back end di lokal.
Ubah terlebih dahulu tautan yang menghubungkan backend menjadi
``` sh
http://localhost:8080/api/search
```
Untuk menjalankan front end, masuk ke directory src lalu gunakan command
``` sh
npm install
npm run dev
```
Untuk menjalankan back end, gunakan command
``` sh
go run main.go
```
