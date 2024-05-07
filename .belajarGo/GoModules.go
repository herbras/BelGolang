package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CATATAN PRIBADI: File ini fokus pada Go Modules
// Go Modules adalah sistem manajemen dependensi resmi untuk Go
// File ini menjelaskan konsep dasar Go Modules dan cara menggunakannya

// CATATAN PRIBADI: Fungsi untuk menampilkan informasi tentang Go Modules
func tampilkanInfoGoModules() {
	fmt.Println("\n=== Informasi Go Modules ===")
	fmt.Println("Go Modules adalah sistem manajemen dependensi resmi untuk Go.")
	fmt.Println("Diperkenalkan pada Go 1.11 dan menjadi default pada Go 1.13.")
	fmt.Println("\nKeuntungan Go Modules:")
	fmt.Println("1. Versioning yang lebih baik")
	fmt.Println("2. Reproducible builds")
	fmt.Println("3. Dependency management yang lebih mudah")
	fmt.Println("4. Tidak perlu lagi GOPATH")
}

// CATATAN PRIBADI: Fungsi untuk menampilkan perintah dasar Go Modules
func tampilkanPerintahDasarGoModules() {
	fmt.Println("\n=== Perintah Dasar Go Modules ===")
	fmt.Println("1. Inisialisasi module baru:")
	fmt.Println("   go mod init github.com/username/nama-proyek")

	fmt.Println("\n2. Menambah dependensi:")
	fmt.Println("   go get github.com/nama-package/nama-module")
	fmt.Println("   go get github.com/nama-package/nama-module@v1.2.3  # versi spesifik")
	fmt.Println("   go get github.com/nama-package/nama-module@latest # versi terbaru")

	fmt.Println("\n3. Membersihkan dependensi yang tidak digunakan:")
	fmt.Println("   go mod tidy")

	fmt.Println("\n4. Melihat dependensi:")
	fmt.Println("   go list -m all")

	fmt.Println("\n5. Download dependensi tanpa build:")
	fmt.Println("   go mod download")

	fmt.Println("\n6. Vendor dependensi (menyimpan lokal):")
	fmt.Println("   go mod vendor")
}

// CATATAN PRIBADI: Fungsi untuk menampilkan struktur file go.mod
func tampilkanStrukturGoMod() {
	fmt.Println("\n=== Struktur File go.mod ===")
	fmt.Println("File go.mod adalah file konfigurasi utama untuk Go Modules.")
	fmt.Println("Contoh struktur file go.mod:")

	fmt.Println("```")
	fmt.Println("module github.com/username/nama-proyek")
	fmt.Println("")
	fmt.Println("go 1.21")
	fmt.Println("")
	fmt.Println("require (")
	fmt.Println("    github.com/gorilla/mux v1.8.0")
	fmt.Println("    github.com/lib/pq v1.10.7")
	fmt.Println(")")
	fmt.Println("")
	fmt.Println("require (")
	fmt.Println("    github.com/gorilla/context v1.1.1 // indirect")
	fmt.Println(")")
	fmt.Println("```")

	fmt.Println("\nPenjelasan:")
	fmt.Println("- module: Nama module (biasanya URL repository)")
	fmt.Println("- go: Versi Go minimum yang dibutuhkan")
	fmt.Println("- require: Dependensi yang dibutuhkan")
	fmt.Println("- // indirect: Dependensi tidak langsung (dibutuhkan oleh dependensi lain)")
}

// CATATAN PRIBADI: Fungsi untuk menampilkan informasi tentang Semantic Versioning
func tampilkanSemanticVersioning() {
	fmt.Println("\n=== Semantic Versioning di Go Modules ===")
	fmt.Println("Go Modules menggunakan Semantic Versioning (SemVer).")
	fmt.Println("Format: vMAJOR.MINOR.PATCH")
	fmt.Println("- MAJOR: Perubahan yang tidak kompatibel dengan versi sebelumnya")
	fmt.Println("- MINOR: Penambahan fitur yang kompatibel dengan versi sebelumnya")
	fmt.Println("- PATCH: Bug fixes yang kompatibel dengan versi sebelumnya")

	fmt.Println("\nContoh:")
	fmt.Println("v1.2.3 -> v1.2.4: Patch update (aman)")
	fmt.Println("v1.2.3 -> v1.3.0: Minor update (biasanya aman)")
	fmt.Println("v1.2.3 -> v2.0.0: Major update (mungkin ada breaking changes)")

	fmt.Println("\nDi Go, major version v2+ memerlukan suffix di import path:")
	fmt.Println("import \"github.com/username/module/v2\"")
}

// CATATAN PRIBADI: Fungsi untuk memeriksa file go.mod di proyek ini
func periksaGoModProyek() {
	fmt.Println("\n=== Memeriksa go.mod di Proyek Ini ===")

	// Cek apakah file go.mod ada
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		fmt.Println("File go.mod tidak ditemukan di direktori ini.")
		return
	}

	// Baca isi file go.mod
	data, err := os.ReadFile("go.mod")
	if err != nil {
		fmt.Printf("Error membaca file go.mod: %v\n", err)
		return
	}

	fmt.Println("Isi file go.mod:")
	fmt.Println("```")
	fmt.Println(string(data))
	fmt.Println("```")

	// Parse informasi dasar
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			fmt.Printf("Nama module: %s\n", strings.TrimPrefix(line, "module "))
		} else if strings.HasPrefix(line, "go ") {
			fmt.Printf("Versi Go: %s\n", strings.TrimPrefix(line, "go "))
		}
	}
}

// CATATAN PRIBADI: Fungsi untuk menampilkan tips penggunaan Go Modules
func tampilkanTipsGoModules() {
	fmt.Println("\n=== Tips Penggunaan Go Modules ===")
	fmt.Println("1. Selalu jalankan 'go mod tidy' setelah menambah/menghapus dependensi")
	fmt.Println("2. Commit file go.mod dan go.sum ke version control")
	fmt.Println("3. Gunakan 'go mod vendor' jika ingin menyimpan dependensi lokal")
	fmt.Println("4. Gunakan 'replace' directive untuk development lokal")
	fmt.Println("5. Perhatikan versi Go yang digunakan dalam file go.mod")
	fmt.Println("6. Gunakan 'go mod why' untuk melihat kenapa suatu package dibutuhkan")
	fmt.Println("7. Gunakan 'go mod graph' untuk melihat dependency graph")
}

// CATATAN PRIBADI: Fungsi untuk menampilkan contoh penggunaan replace directive
func tampilkanContohReplace() {
	fmt.Println("\n=== Contoh Penggunaan Replace Directive ===")
	fmt.Println("Replace directive berguna untuk:")
	fmt.Println("1. Development lokal package")
	fmt.Println("2. Mengganti dependensi dengan fork")
	fmt.Println("3. Mengatasi masalah dengan dependensi")

	fmt.Println("\nContoh dalam go.mod:")
	fmt.Println("```")
	fmt.Println("module github.com/username/myproject")
	fmt.Println("")
	fmt.Println("go 1.21")
	fmt.Println("")
	fmt.Println("require github.com/username/library v1.0.0")
	fmt.Println("")
	fmt.Println("replace github.com/username/library => ../library")
	fmt.Println("```")

	fmt.Println("\nIni akan menggunakan versi lokal library dari ../library")
	fmt.Println("daripada mengunduhnya dari GitHub.")
}

// CATATAN PRIBADI: Fungsi untuk menampilkan contoh workspace dengan multiple modules
func tampilkanWorkspaceMultipleModules() {
	fmt.Println("\n=== Workspace dengan Multiple Modules (Go 1.18+) ===")
	fmt.Println("Go 1.18 memperkenalkan fitur workspace untuk bekerja dengan multiple modules.")

	fmt.Println("\nStruktur direktori:")
	fmt.Println("```")
	fmt.Println("myworkspace/")
	fmt.Println("├── go.work")
	fmt.Println("├── module1/")
	fmt.Println("│   ├── go.mod")
	fmt.Println("│   └── main.go")
	fmt.Println("└── module2/")
	fmt.Println("    ├── go.mod")
	fmt.Println("    └── lib.go")
	fmt.Println("```")

	fmt.Println("\nContoh file go.work:")
	fmt.Println("```")
	fmt.Println("go 1.18")
	fmt.Println("")
	fmt.Println("use (")
	fmt.Println("    ./module1")
	fmt.Println("    ./module2")
	fmt.Println(")")
	fmt.Println("```")

	fmt.Println("\nPerintah untuk membuat workspace:")
	fmt.Println("go work init ./module1 ./module2")

	fmt.Println("\nPerintah untuk menambah module ke workspace:")
	fmt.Println("go work use ./module3")
}

// CATATAN PRIBADI: Fungsi untuk menampilkan langkah-langkah membuat dan publish module
func tampilkanLangkahMembuatModule() {
	fmt.Println("\n=== Langkah-langkah Membuat dan Publish Module ===")
	fmt.Println("1. Buat repository di GitHub/GitLab/dll")
	fmt.Println("2. Clone repository ke lokal")
	fmt.Println("3. Inisialisasi Go module:")
	fmt.Println("   go mod init github.com/username/nama-module")
	fmt.Println("4. Buat kode Go Anda")
	fmt.Println("5. Commit dan push ke repository")
	fmt.Println("6. Tag versi dengan Git:")
	fmt.Println("   git tag v1.0.0")
	fmt.Println("   git push origin v1.0.0")
	fmt.Println("7. Sekarang orang lain bisa menggunakan module Anda:")
	fmt.Println("   go get github.com/username/nama-module")
}

func demoGoModules() {
	fmt.Println("Belajar Go Modules")

	// Tampilkan informasi dasar tentang Go Modules
	tampilkanInfoGoModules()

	// Tampilkan perintah dasar Go Modules
	tampilkanPerintahDasarGoModules()

	// Tampilkan struktur file go.mod
	tampilkanStrukturGoMod()

	// Tampilkan informasi tentang Semantic Versioning
	tampilkanSemanticVersioning()

	// Periksa file go.mod di proyek ini
	periksaGoModProyek()

	// Tampilkan tips penggunaan Go Modules
	tampilkanTipsGoModules()

	// Tampilkan contoh penggunaan replace directive
	tampilkanContohReplace()

	// Tampilkan informasi tentang workspace dengan multiple modules
	tampilkanWorkspaceMultipleModules()

	// Tampilkan langkah-langkah membuat dan publish module
	tampilkanLangkahMembuatModule()

	// Tampilkan direktori cache Go Modules
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			gopath = filepath.Join(home, "go")
		}
	}

	if gopath != "" {
		modCache := filepath.Join(gopath, "pkg", "mod")
		fmt.Println("\n=== Direktori Cache Go Modules ===")
		fmt.Printf("Cache Go Modules berada di: %s\n", modCache)
		fmt.Println("Direktori ini menyimpan semua dependensi yang diunduh.")
	}

	fmt.Println("\nProgram selesai!")
}
