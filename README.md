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
Dalam pencarian resep terdapat dua opsi algoritma ynag dapat dipilih yaitu **BFS dan DFS**. Algoritma BFS menggunakan pendekatan penelusuran graf secara menyebar sedangkan algoritma DFS menggunakan pendekatan penelusuran graf secara mendalam.
3. Pohon Solusi Resep  
Pohon solusi resep yang ditampilkan dapat diexpand untuk melihat resep secara detail dan juga dapat ditutup untuk melihat resep secara umum. Hal ini memiliki keunggulan memberikan tampilan yang lebih fleksibel sesuai kebutuhan pengguna.
4. Detail Informasi  
Di bagian bawah solusi terdapat detail infromasi terkait jumlah node yang dikunjungi serta durasi pencarian resep.

## Teknologi yang digunakan
Untuk bagian back end, digunakan bahasa go dan framework gin untuk tugas ini. Algoritma diimplementasikan secara terpisah di folder algorithm. Algoritma tersebut lalu digunakan di main.go dengan cara mengimport nya. Di dalam main.go terdapat mekanisme yang menghubungkan antara front end dan back end dengan framework gin.

## Cara Penggunaan

