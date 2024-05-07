# Perjalanan Belajar Golang

Proyek ini mendokumentasikan perjalanan belajar saya dalam mempelajari bahasa pemrograman Go (Golang). Saya mengimplementasikan berbagai konsep dasar pemrograman dalam Go sesuai dengan diagram yang ada di `Challenge.puml`.

## Proses Belajar

### Tahap 1: Pengenalan Dasar
- Mempelajari sintaks dasar Go
- Memahami cara mendeklarasikan variabel dengan `var` dan `:=`
- Memahami tipe data dasar seperti `string`, `int`, `bool`
- Membuat program "Hello, World!" pertama

### Tahap 2: Variabel dan Tipe Data
- Mempelajari deklarasi variabel global dan lokal
- Memahami type inference di Go
- Mempelajari konstanta dan cara mendeklarasikannya
- Memahami pointer dan penggunaannya dengan `new()`

```go
// Contoh deklarasi variabel yang saya pelajari
var firstName string = "Ibrahim"
var lastName = "Huda" // type inferred as string
var Satu, Dua, Tiga string

// Penggunaan pointer
name := new(string)
fmt.Println(name)   // Mencetak alamat memori
fmt.Println(*name)  // Mencetak nilai (kosong)
```

### Tahap 3: Fungsi dan Kontrol Alur
- Mempelajari cara mendefinisikan dan memanggil fungsi
- Memahami parameter dan nilai kembalian fungsi
- Mempelajari struktur kontrol seperti if-else dan for loop
- Mengimplementasikan fungsi rekursif (factorial, power)

### Tahap 4: Struktur Data
- Mempelajari array dan slice
- Memahami perbedaan antara array (ukuran tetap) dan slice (dinamis)
- Mengimplementasikan algoritma pencarian (findLargest)

### Tahap 5: Input/Output dan File
- Mempelajari cara membaca input dari pengguna
- Memahami package `bufio` dan `strings`
- Mempelajari operasi file dasar (membuka, membaca, menulis, menutup)
- Menggunakan `defer` untuk memastikan file ditutup

## Catatan Pribadi dan Pembelajaran

Selama proses belajar, saya membuat beberapa catatan pribadi untuk membantu pemahaman:

```go
// CATATAN PRIBADI: Deklarasi variabel di Go bisa menggunakan var atau :=
// CATATAN PRIBADI: := hanya bisa digunakan di dalam fungsi, tidak bisa di level package
// CATATAN PRIBADI: Variabel yang tidak digunakan di Go akan menyebabkan error kompilasi
```

Beberapa hal penting yang saya pelajari:
1. Go memiliki garbage collector, jadi tidak perlu mengelola memori secara manual
2. Go mendorong penanganan error yang eksplisit dengan mengembalikan error sebagai nilai
3. Fungsi di Go dapat mengembalikan beberapa nilai sekaligus
4. Penggunaan `defer` sangat berguna untuk clean-up resources

## Tantangan yang Dihadapi

1. Memahami konsep pointer di Go
2. Mengatasi error "multiple main function" saat memiliki beberapa file dengan fungsi main
3. Memahami cara kerja package dan import
4. Menyesuaikan diri dengan penanganan error yang eksplisit

## Fitur yang Diimplementasikan

Program ini mencakup implementasi dari berbagai konsep dasar pemrograman:

1. Mencetak "Hello, World!" ke konsol
2. Deklarasi dan manipulasi variabel
3. Menukar nilai variabel
4. Fungsi untuk menjumlahkan dua angka
5. Input dari pengguna
6. Pengecekan bilangan genap atau ganjil
7. Mencetak 10 bilangan pertama dari deret Fibonacci
8. Menghitung faktorial dari sebuah angka
9. Mencari elemen terbesar dalam array
10. Mencetak segitiga bintang dengan tinggi yang ditentukan pengguna
11. Fungsi rekursif untuk menghitung pangkat dari sebuah angka
12. Membaca file teks dan menghitung jumlah baris

## Struktur Proyek

- `main.go`: Berisi program utama yang mengimplementasikan semua fitur
- `Type.go`: Berisi contoh penggunaan tipe data di Go
- `main_test.go`: Berisi pengujian untuk fungsi-fungsi di `main.go`
- `Type_test.go`: Berisi pengujian untuk variabel dan konstanta di `Type.go`
- `Challenge.puml`: Diagram PlantUML yang menjelaskan alur program

## Cara Menjalankan Program

```bash
go run main.go
```

Atau untuk menjalankan semua file Go dalam proyek:

```bash
go run .
```

## Cara Menjalankan Pengujian

Untuk menjalankan semua pengujian:

```bash
go test ./...
```

Untuk menjalankan pengujian tertentu, misalnya untuk `main_test.go`:

```bash
go test -v -run TestAdd
```

## 1. Konkurensi dengan Goroutines dan Channels

### Konsep Dasar

**Goroutines** adalah fungsi yang berjalan secara konkuren (bersamaan) dengan fungsi lain. Goroutines sangat ringan dan efisien dibandingkan dengan thread tradisional.

**Channels** adalah cara untuk berkomunikasi antar goroutines. Channels memungkinkan goroutines untuk mengirim dan menerima data dengan aman.

### File Pembelajaran
- `LanjutBelajar.go` - Contoh dasar goroutines dan channels
- `KonkurensiLanjutan.go` - Contoh lanjutan dengan worker pool, select, mutex

### Konsep yang Dipelajari

#### A. Goroutines Dasar
```go
// Menjalankan fungsi sebagai goroutine
go namaFungsi()

// Anonymous function sebagai goroutine
go func() {
    fmt.Println("Ini berjalan secara konkuren")
}()
```

#### B. Channels
```go
// Membuat channel
ch := make(chan string)

// Mengirim data ke channel
ch <- "Hello"

// Menerima data dari channel
pesan := <-ch
```

#### C. Buffered vs Unbuffered Channels
- **Unbuffered**: `make(chan int)` - Blocking sampai ada receiver
- **Buffered**: `make(chan int, 5)` - Bisa menyimpan 5 nilai sebelum blocking

#### D. Select Statement
```go
select {
case msg1 := <-ch1:
    fmt.Println("Dari channel 1:", msg1)
case msg2 := <-ch2:
    fmt.Println("Dari channel 2:", msg2)
default:
    fmt.Println("Tidak ada channel yang siap")
}
```

#### E. Worker Pool Pattern
Pattern untuk membatasi jumlah goroutine yang berjalan bersamaan:
- Membuat channel untuk jobs dan results
- Membuat sejumlah worker goroutines
- Mengirim tugas ke channel jobs
- Mengumpulkan hasil dari channel results

#### F. Mutex untuk Race Condition
```go
var mu sync.Mutex
mu.Lock()   // Mengunci akses
// kode yang perlu dilindungi
mu.Unlock() // Membuka kunci
```

### Tips Praktis
1. Gunakan `sync.WaitGroup` untuk menunggu goroutines selesai
2. Selalu tutup channels dengan `close(ch)` setelah selesai mengirim
3. Gunakan buffered channels untuk menghindari deadlock
4. Gunakan mutex untuk melindungi shared data
5. Gunakan select untuk menangani multiple channels

---

## 2. Penggunaan Struct dan Interface

### Konsep Dasar

**Struct** adalah cara untuk mengelompokkan data yang terkait dalam satu unit. Mirip dengan class di bahasa pemrograman lain, tetapi tanpa inheritance.

**Interface** adalah kontrak yang mendefinisikan method yang harus diimplementasikan oleh suatu tipe. Interface memungkinkan polymorphism di Go.

### File Pembelajaran
- `Type.go` - Contoh dasar tipe data
- `StructInterface.go` - Contoh lengkap struct dan interface

### Konsep yang Dipelajari

#### A. Struct Dasar
```go
type Mahasiswa struct {
    NIM     string
    Nama    string
    IPK     float64
}

// Membuat instance struct
mhs := Mahasiswa{
    NIM:  "12345",
    Nama: "Ibrahim",
    IPK:  3.75,
}
```

#### B. Method pada Struct
```go
// Method dengan value receiver
func (m Mahasiswa) TampilkanInfo() {
    fmt.Printf("Nama: %s, IPK: %.2f\n", m.Nama, m.IPK)
}

// Method dengan pointer receiver (untuk modifikasi)
func (m *Mahasiswa) UpdateIPK(ipkBaru float64) {
    m.IPK = ipkBaru
}
```

#### C. Embedded Struct
```go
type MahasiswaLengkap struct {
    Mahasiswa  // Embedded struct
    Alamat     string
    Umur       int
}
```

#### D. Interface
```go
type BentukGeometri interface {
    Luas() float64
    Keliling() float64
}

// Struct yang mengimplementasikan interface
type PersegiPanjang struct {
    Panjang, Lebar float64
}

func (p PersegiPanjang) Luas() float64 {
    return p.Panjang * p.Lebar
}

func (p PersegiPanjang) Keliling() float64 {
    return 2 * (p.Panjang + p.Lebar)
}
```

#### E. Empty Interface
```go
// interface{} bisa menerima tipe data apapun
func tampilkanData(data interface{}) {
    fmt.Printf("Nilai: %v, Tipe: %T\n", data, data)
}

// Type assertion
if str, ok := data.(string); ok {
    fmt.Println("Ini adalah string:", str)
}
```

### Tips Praktis
1. Gunakan pointer receiver untuk method yang memodifikasi struct
2. Gunakan value receiver untuk method yang hanya membaca data
3. Interface di Go diimplementasikan secara implisit
4. Gunakan empty interface `interface{}` dengan hati-hati
5. Gunakan type assertion untuk mengecek tipe dari interface{}

---

## 3. Pengelolaan Dependensi dengan Go Modules

### Konsep Dasar

**Go Modules** adalah sistem manajemen dependensi resmi untuk Go. Go Modules memungkinkan kita untuk mengelola versi package dan dependensi dengan mudah.

### File yang Terkait
- `go.mod` - File konfigurasi module
- `go.sum` - File checksum untuk verifikasi integritas

### Perintah Dasar Go Modules

#### A. Inisialisasi Module
```bash
# Membuat module baru
go mod init github.com/username/nama-proyek

# Contoh untuk proyek ini
go mod init github.com/username/BelGolang
```

#### B. Menambah Dependensi
```bash
# Menambah dependensi baru
go get github.com/nama-package/nama-module

# Menambah versi spesifik
go get github.com/nama-package/nama-module@v1.2.3

# Menambah versi terbaru
go get github.com/nama-package/nama-module@latest
```

#### C. Mengelola Dependensi
```bash
# Membersihkan dependensi yang tidak digunakan
go mod tidy

# Melihat dependensi
go list -m all

# Download dependensi tanpa build
go mod download
```

#### D. Update Dependensi
```bash
# Update semua dependensi ke versi minor/patch terbaru
go get -u

# Update dependensi spesifik
go get -u github.com/nama-package/nama-module
```

### Struktur File go.mod
```go
module github.com/username/BelGolang

go 1.21

require (
    github.com/gorilla/mux v1.8.0
    github.com/lib/pq v1.10.7
)

require (
    github.com/gorilla/context v1.1.1 // indirect
)
```

### Tips Praktis
1. Selalu jalankan `go mod tidy` setelah menambah/menghapus dependensi
2. Commit file `go.mod` dan `go.sum` ke version control
3. Gunakan `go mod vendor` jika ingin menyimpan dependensi lokal
4. Gunakan `replace` directive untuk development lokal
5. Perhatikan versi Go yang digunakan dalam file go.mod

---


# 4. Sesi belajar bareng konkurensi dan indexing postgresql

### Catatan Pembelajaran: Concurrency di Go & Indexing PostgreSQL

---

#### **A. Concurrency di Go**
1. **Goroutine**  
   - Lightweight thread yang dikelola Go runtime.
   - Mulai dengan `go func()`.
   - Contoh:  
     ```go
     go func() { 
         fmt.Println("Pesan dari goroutine") 
     }()
     ```

2. **Channel**  
   - Komunikasi antar goroutine (thread-safe).
   - Unbuffered (sync) vs Buffered (async):  
     ```go
     ch := make(chan int)       // Unbuffered
     chBuf := make(chan int, 3) // Buffered (kapasitas 3)
     ```

3. **Select**  
   - Handle multiple channel operations:  
     ```go
     select {
         case msg := <-ch1: 
             fmt.Println(msg)
         case ch2 <- data:
             // data dikirim
         default:
             // fallback
     }
     ```

4. **Sync Package**  
   - **Mutex**:  
     ```go
     var mu sync.Mutex
     mu.Lock()
     // akses shared resource
     mu.Unlock()
     ```
   - **WaitGroup**: Tunggu goroutine selesai:  
     ```go
     var wg sync.WaitGroup
     wg.Add(1)
     go func() {
         defer wg.Done()
         // tugas
     }()
     wg.Wait()
     ```

---

#### **B. Tantangan Concurrency**
Implementasikan **parallel file processor** dengan Go:  
- Baca 100 file `.txt` bersamaan.
- Hitung frekuensi kata tertentu (e.g., "bteee") di tiap file.
- Aggregasi hasil secara konkuren.

**Solusi Kerangka**:  
```go
package main

import (
  "fmt"
  "sync"
  "io/ioutil"
  "strings"
)

func processFile(path string, targetWord string, wg *sync.WaitGroup, resultChan chan<- int) {
  defer wg.Done()
  data, _ := ioutil.ReadFile(path)
  content := string(data)
  count := strings.Count(content, targetWord)
  resultChan <- count
}

func main() {
  var wg sync.WaitGroup
  resultChan := make(chan int, 100) // Buffered channel
  target := "bteee"

  files := []string{"file1.txt", "file2.txt", ...} // 100 file

  for _, file := range files {
    wg.Add(1)
    go processFile(file, target, &wg, resultChan)
  }

  go func() {
    wg.Wait()
    close(resultChan)
  }()

  total := 0
  for count := range resultChan {
    total += count
  }
  fmt.Printf("Total kemunculan '%s': %d\n", target, total)
}
```

---

#### **C. Indexing di PostgreSQL**
1. **B-Tree (Balanced Tree)**  
   - **Digunakan untuk**: `=`, `>`, `>=`, `<`, `<=`, `BETWEEN`, `LIKE 'foo%'`.
   - Optimal untuk data berurutan & rentang nilai.
   - Contoh:  
     ```sql
     CREATE INDEX idx_name ON users (name); -- Default B-Tree
     ```

2. **Hash Index**  
   - Hanya mendukung equality (`=`), lebih cepat untuk lookup exact.
   - Contoh:  
     ```sql
     CREATE INDEX idx_hash ON orders USING HASH (order_id);
     ```

3. **GIN (Generalized Inverted Index)**  
   - Untuk data terkomposisi: array, jsonb, full-text search.
   - Contoh:  
     ```sql
     CREATE INDEX idx_gin_tags ON products USING GIN (tags);
     ```

4. **BRIN (Block Range Index)**  
   - Efisien untuk data besar dengan nilai terurut (e.g., timestamp).
   - Contoh:  
     ```sql
     CREATE INDEX idx_brin_date ON sales USING BRIN (sale_date);
     ```

5. **GiST & SP-GiST**  
   - Untuk data spasial, geometri, atau custom tipe data.

---

#### **D. B-Tree vs Binary Tree**
| **Aspek**            | **B-Tree**                            | **Binary Tree (e.g., BST)**       |
|-----------------------|---------------------------------------|-----------------------------------|
| Struktur              | Multi-node per level                 | Max 2 child per node             |
| Keseimbangan          | Auto-balanced                        | Tidak terjamin (bisa unbalanced) |
| I/O Disk              | Optimal (blok/node besar)            | Tidak optimal                    |
| Penggunaan di Database| ✔️ Standard indexing                 | ❌ Tidak dipakai                 |
| Keunggulan            | Cache-friendly untuk sistem penyimpanan | Sederhana untuk memori internal |

---

#### **E. Event Sourcing & Golang**
- **Event Store**:  
  Simpan perubahan state sebagai sequence of events (immutable).
- **Pattern di Go**:  
  - Gunakan slice untuk menyimpan event stream.
  - Implementasi replay event untuk rebuild state.

```go
type Event interface{ Apply(*State) }

type State struct { /* ... */ }

func (s *State) Replay(events []Event) {
  for _, e := range events {
    e.Apply(s)
  }
}
```

---

**Tips**:
- Untuk PostgreSQL, pilih indeks berdasarkan pola akses data:  
  - B-Tree untuk kolom dengan filter rentang nilai.  
  - GIN untuk kolom array/jsonb.  
  - BRIN untuk data timeseries.  
- Di Go, hindari shared state; gunakan channels dan mutex untuk safety.