package main

import (
	"fmt"
	"sync"
	"time"
)

// CATATAN PRIBADI: Struct adalah kumpulan field dengan tipe data tertentu
// Struct mirip dengan class di bahasa pemrograman lain
type Orang struct {
	Nama   string
	Umur   int
	Alamat string
	Hobi   []string
}

// CATATAN PRIBADI: Interface adalah kumpulan method signature
// Interface digunakan untuk polymorphism di Go
type Sapa interface {
	SayHello() string
}

// CATATAN PRIBADI: Method adalah fungsi yang terikat dengan tipe data tertentu
// Method ini terikat dengan struct Orang
func (o Orang) SayHello() string {
	return fmt.Sprintf("Halo, nama saya %s dan umur saya %d tahun", o.Nama, o.Umur)
}

// CATATAN PRIBADI: Struct lain yang juga mengimplementasikan interface Sapa
type Robot struct {
	Model string
	Tahun int
}

func (r Robot) SayHello() string {
	return fmt.Sprintf("Beep boop! Saya robot model %s keluaran tahun %d", r.Model, r.Tahun)
}

// CATATAN PRIBADI: Fungsi yang menerima interface sebagai parameter
// Ini menunjukkan polymorphism di Go
func PerkenalkanDiri(s Sapa) {
	fmt.Println(s.SayHello())
}

// CATATAN PRIBADI: Goroutine adalah fungsi yang berjalan secara konkuren
// Channel digunakan untuk komunikasi antar goroutine
func kirimPesan(ch chan string, pesan string) {
	fmt.Printf("Mengirim: %s\n", pesan)
	ch <- pesan // Mengirim pesan ke channel
}

func terimaPesan(ch chan string, wg *sync.WaitGroup) {
	pesan := <-ch // Menerima pesan dari channel
	fmt.Printf("Menerima: %s\n", pesan)
	wg.Done()
}

// CATATAN PRIBADI: Contoh konkurensi dengan multiple goroutines
func demoKonkurensi() {
	fmt.Println("\n=== Demo Konkurensi dengan Goroutines dan Channels ===")

	// Membuat channel untuk komunikasi
	ch := make(chan string)

	// WaitGroup untuk menunggu semua goroutine selesai
	var wg sync.WaitGroup

	// Menjalankan beberapa goroutine
	pesan := []string{"Halo", "Apa kabar", "Selamat belajar Go"}

	wg.Add(len(pesan))

	// Goroutine untuk menerima pesan
	for i := 0; i < len(pesan); i++ {
		go terimaPesan(ch, &wg)
	}

	// Mengirim pesan ke channel
	for _, msg := range pesan {
		go kirimPesan(ch, msg)
	}

	// Menunggu semua goroutine selesai
	wg.Wait()
}

// CATATAN PRIBADI: Contoh penggunaan struct dan interface
func demoStructInterfaceDasar() {
	fmt.Println("\n=== Demo Struct dan Interface ===")

	// Membuat instance dari struct Orang
	ibrahim := Orang{
		Nama:   "Ibrahim",
		Umur:   25,
		Alamat: "Jakarta",
		Hobi:   []string{"Coding", "Membaca", "Traveling"},
	}

	// Mengakses field dari struct
	fmt.Printf("Nama: %s\n", ibrahim.Nama)
	fmt.Printf("Umur: %d\n", ibrahim.Umur)
	fmt.Printf("Alamat: %s\n", ibrahim.Alamat)
	fmt.Printf("Hobi: %v\n", ibrahim.Hobi)

	// Membuat instance dari struct Robot
	r2d2 := Robot{
		Model: "R2D2",
		Tahun: 1977,
	}

	// Menggunakan interface untuk polymorphism
	fmt.Println("\nDemo Interface (Polymorphism):")
	PerkenalkanDiri(ibrahim) // Orang.SayHello()
	PerkenalkanDiri(r2d2)    // Robot.SayHello()
}

// CATATAN PRIBADI: Contoh penggunaan Go Modules
// Untuk Go Modules, kita perlu menjalankan beberapa perintah di terminal:
// 1. go mod init github.com/username/projectname
// 2. go get github.com/nama-package/nama-module
// 3. go mod tidy
func demoGoModulesDasar() {
	fmt.Println("\n=== Demo Go Modules ===")
	fmt.Println("Go Modules adalah sistem manajemen dependensi resmi untuk Go.")
	fmt.Println("File go.mod sudah ada di proyek ini, berisi:")
	fmt.Println("module github.com/username/BelGolang")
	fmt.Println("go 1.x")
	fmt.Println("\nUntuk menambahkan dependensi eksternal, gunakan perintah:")
	fmt.Println("go get github.com/nama-package/nama-module")
}

func demoLanjutBelajar() {
	fmt.Println("Belajar Lanjutan Golang: Konkurensi, Struct/Interface, dan Go Modules")

	// Demo penggunaan struct dan interface
	demoStructInterfaceDasar()

	// Demo konkurensi dengan goroutines dan channels
	demoKonkurensi()

	// Demo Go Modules
	demoGoModulesDasar()

	// Contoh tambahan: goroutine dengan timer
	fmt.Println("\n=== Demo Goroutine dengan Timer ===")
	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Printf("Goroutine: hitungan ke-%d\n", i)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Program utama juga berjalan
	for i := 1; i <= 3; i++ {
		fmt.Printf("Main: hitungan ke-%d\n", i)
		time.Sleep(1 * time.Second)
	}

	// Beri waktu agar goroutine selesai
	fmt.Println("\nMenunggu goroutine selesai...")
	time.Sleep(1 * time.Second)
	fmt.Println("Program selesai!")
}
