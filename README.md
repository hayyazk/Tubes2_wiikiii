# Tubes2_wiikiii - PEMANFAATAN ALGORITMA IDS DAN BFS DALAM PERMAINAN WIKIRACE
## Algoritma IDS
Algoritma mengambil tautan dari artikel awal dan memasukkannya ke dalam list. List diproses secara LIFO. Jika sudah memproses semua artikel sampai kedalaman tertentu dan artikel tujuan belum ditemukan, pencarian diulang kembali dengan kedalaman maksimal lebih tinggi. 

## Algoritma BFS
Algoritma mengambil tautan dari artikel awal dan memasukkannya ke dalam list. List diproses secara FIFO. Untuk semua tautan yang didapatkan pada kedalaman yang sama, digunakan goroutine untuk melakukan scraping kepada beberapa artikel tersebut secara bersamaan. Jika artikel yang akan diproses adalah artikel tujuan, program diberhentikan.

## Requirements
- Go

## Compiling
Clone repository
```shell
git clone https://github.com/hayyazk/Tubes2_wiikiii.git
```
Open folder src
```shell
cd Tubes2_wiikiii/src
```
Build
```shell
go build
```
Run executable
```shell
./Tubes2_wiikiii.exe
```
Open browser
```shell
localhost:8080
```

## Authors
| NIM | Nama |
| --- | --- |
| 10023519 | Muhammad Fiqri |
| 13521154 | Naufal Baldemar Ardanni |
| 13522102 | Hayya Zuhailii Kinasih |
