package main

import (
	"fmt"
	"math"
)

// CATATAN PRIBADI: File ini fokus pada struct dan interface di Go
// Struct adalah cara untuk mengelompokkan data yang terkait
// Interface adalah kontrak yang mendefinisikan method yang harus diimplementasikan

// CATATAN PRIBADI: Struct dasar untuk menyimpan informasi mahasiswa
type Mahasiswa struct {
	NIM     string
	Nama    string
	Jurusan string
	IPK     float64
	Aktif   bool
}

// CATATAN PRIBADI: Embedded struct (struct di dalam struct)
type Alamat struct {
	Jalan    string
	Kota     string
	KodePos  string
	Provinsi string
}

type MahasiswaLengkap struct {
	Mahasiswa // Embedded struct
	Alamat    Alamat
	Umur      int
}

// CATATAN PRIBADI: Interface untuk operasi matematika
type BentukGeometri interface {
	Luas() float64
	Keliling() float64
	Info() string
}

// CATATAN PRIBADI: Struct untuk persegi panjang
type PersegiPanjang struct {
	Panjang float64
	Lebar   float64
}

// CATATAN PRIBADI: Method untuk PersegiPanjang yang mengimplementasikan interface BentukGeometri
func (p PersegiPanjang) Luas() float64 {
	return p.Panjang * p.Lebar
}

func (p PersegiPanjang) Keliling() float64 {
	return 2 * (p.Panjang + p.Lebar)
}

func (p PersegiPanjang) Info() string {
	return fmt.Sprintf("Persegi Panjang dengan panjang %.2f dan lebar %.2f", p.Panjang, p.Lebar)
}

// CATATAN PRIBADI: Struct untuk lingkaran
type Lingkaran struct {
	JariJari float64
}

// CATATAN PRIBADI: Method untuk Lingkaran yang mengimplementasikan interface BentukGeometri
func (l Lingkaran) Luas() float64 {
	return math.Pi * l.JariJari * l.JariJari
}

func (l Lingkaran) Keliling() float64 {
	return 2 * math.Pi * l.JariJari
}

func (l Lingkaran) Info() string {
	return fmt.Sprintf("Lingkaran dengan jari-jari %.2f", l.JariJari)
}

// CATATAN PRIBADI: Fungsi yang menerima interface sebagai parameter
// Ini menunjukkan polymorphism - satu fungsi bisa bekerja dengan berbagai tipe
func hitungDanTampilkan(bentuk BentukGeometri) {
	fmt.Println(bentuk.Info())
	fmt.Printf("Luas: %.2f\n", bentuk.Luas())
	fmt.Printf("Keliling: %.2f\n", bentuk.Keliling())
	fmt.Println("---")
}

// CATATAN PRIBADI: Method untuk struct Mahasiswa
func (m Mahasiswa) TampilkanInfo() {
	fmt.Printf("NIM: %s\n", m.NIM)
	fmt.Printf("Nama: %s\n", m.Nama)
	fmt.Printf("Jurusan: %s\n", m.Jurusan)
	fmt.Printf("IPK: %.2f\n", m.IPK)
	fmt.Printf("Status: %s\n", m.statusKeaktifan())
}

// CATATAN PRIBADI: Method private (huruf kecil di awal)
func (m Mahasiswa) statusKeaktifan() string {
	if m.Aktif {
		return "Aktif"
	}
	return "Tidak Aktif"
}

// CATATAN PRIBADI: Method dengan pointer receiver
// Digunakan ketika ingin memodifikasi nilai struct
func (m *Mahasiswa) UpdateIPK(ipkBaru float64) {
	m.IPK = ipkBaru
	fmt.Printf("IPK %s berhasil diupdate menjadi %.2f\n", m.Nama, ipkBaru)
}

// CATATAN PRIBADI: Interface kosong (empty interface)
// Bisa menerima nilai dari tipe data apapun
func tampilkanTipeData(data interface{}) {
	fmt.Printf("Nilai: %v, Tipe: %T\n", data, data)
}

// CATATAN PRIBADI: Type assertion untuk mengecek tipe data dari interface{}
func cekTipeData(data interface{}) {
	switch v := data.(type) {
	case string:
		fmt.Printf("Ini adalah string dengan panjang %d karakter\n", len(v))
	case int:
		fmt.Printf("Ini adalah integer dengan nilai %d\n", v)
	case float64:
		fmt.Printf("Ini adalah float64 dengan nilai %.2f\n", v)
	case Mahasiswa:
		fmt.Printf("Ini adalah struct Mahasiswa dengan nama %s\n", v.Nama)
	default:
		fmt.Printf("Tipe data tidak dikenali: %T\n", v)
	}
}

func demoStruct() {
	fmt.Println("\n=== Demo Struct ===")

	// Membuat instance struct dengan berbagai cara
	mhs1 := Mahasiswa{
		NIM:     "12345678",
		Nama:    "Ibrahim Huda",
		Jurusan: "Teknik Informatika",
		IPK:     3.75,
		Aktif:   true,
	}

	// Cara lain membuat struct
	mhs2 := Mahasiswa{"87654321", "Siti Aminah", "Sistem Informasi", 3.85, true}

	// Membuat struct dengan zero values
	var mhs3 Mahasiswa
	mhs3.NIM = "11111111"
	mhs3.Nama = "Ahmad Fauzi"
	mhs3.Jurusan = "Teknik Komputer"
	mhs3.IPK = 3.50
	mhs3.Aktif = false

	// Menampilkan informasi mahasiswa
	fmt.Println("Informasi Mahasiswa 1:")
	mhs1.TampilkanInfo()

	fmt.Println("\nInformasi Mahasiswa 2:")
	mhs2.TampilkanInfo()

	fmt.Println("\nInformasi Mahasiswa 3:")
	mhs3.TampilkanInfo()

	// Update IPK menggunakan pointer receiver
	fmt.Println("\nUpdate IPK Mahasiswa 1:")
	mhs1.UpdateIPK(3.90)
	mhs1.TampilkanInfo()
}

func demoEmbeddedStruct() {
	fmt.Println("\n=== Demo Embedded Struct ===")

	mhsLengkap := MahasiswaLengkap{
		Mahasiswa: Mahasiswa{
			NIM:     "99999999",
			Nama:    "Dewi Sartika",
			Jurusan: "Teknik Elektro",
			IPK:     3.95,
			Aktif:   true,
		},
		Alamat: Alamat{
			Jalan:    "Jl. Merdeka No. 123",
			Kota:     "Jakarta",
			KodePos:  "12345",
			Provinsi: "DKI Jakarta",
		},
		Umur: 20,
	}

	fmt.Println("Informasi Mahasiswa Lengkap:")
	// Bisa mengakses method dari embedded struct
	mhsLengkap.TampilkanInfo()
	fmt.Printf("Umur: %d tahun\n", mhsLengkap.Umur)
	fmt.Printf("Alamat: %s, %s, %s, %s\n",
		mhsLengkap.Alamat.Jalan,
		mhsLengkap.Alamat.Kota,
		mhsLengkap.Alamat.KodePos,
		mhsLengkap.Alamat.Provinsi)
}

func demoInterface() {
	fmt.Println("\n=== Demo Interface ===")

	// Membuat berbagai bentuk geometri
	persegi := PersegiPanjang{Panjang: 10, Lebar: 5}
	lingkaran := Lingkaran{JariJari: 7}

	// Menggunakan interface untuk polymorphism
	bentukGeometri := []BentukGeometri{persegi, lingkaran}

	for i, bentuk := range bentukGeometri {
		fmt.Printf("Bentuk %d:\n", i+1)
		hitungDanTampilkan(bentuk)
	}
}

func demoEmptyInterface() {
	fmt.Println("\n=== Demo Empty Interface ===")

	// Empty interface bisa menerima tipe data apapun
	data := []interface{}{
		"Hello World",
		42,
		3.14,
		true,
		Mahasiswa{NIM: "12345", Nama: "Test", Jurusan: "TI", IPK: 3.5, Aktif: true},
	}

	fmt.Println("Menampilkan berbagai tipe data:")
	for i, item := range data {
		fmt.Printf("Data %d - ", i+1)
		tampilkanTipeData(item)
	}

	fmt.Println("\nType assertion:")
	for i, item := range data {
		fmt.Printf("Data %d - ", i+1)
		cekTipeData(item)
	}
}

func demoStructInterface() {
	fmt.Println("Belajar Struct dan Interface di Go")

	// Demo penggunaan struct
	demoStruct()

	// Demo embedded struct
	demoEmbeddedStruct()

	// Demo interface dan polymorphism
	demoInterface()

	// Demo empty interface
	demoEmptyInterface()

	fmt.Println("\nProgram selesai!")
}
