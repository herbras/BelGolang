package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// CATATAN PRIBADI: File ini fokus pada konkurensi di Go
// Konkurensi adalah kemampuan menjalankan beberapa tugas secara bersamaan
// Go memiliki dua fitur utama untuk konkurensi: goroutines dan channels

// CATATAN PRIBADI: Worker Pool Pattern
// Pattern ini berguna untuk membatasi jumlah goroutine yang berjalan bersamaan
func demoWorkerPool() {
	fmt.Println("\n=== Demo Worker Pool Pattern ===")

	// Jumlah tugas dan worker
	const numJobs = 10
	const numWorkers = 3

	// Membuat channel untuk tugas dan hasil
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Membuat worker pool
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// Mengirim tugas ke channel jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Menutup channel jobs setelah semua tugas dikirim

	// Mengumpulkan hasil
	for a := 1; a <= numJobs; a++ {
		<-results // Hanya menunggu semua hasil selesai
	}
}

// CATATAN PRIBADI: Fungsi worker untuk worker pool
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d memulai tugas %d\n", id, j)
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond) // Simulasi pekerjaan
		fmt.Printf("Worker %d selesai tugas %d\n", id, j)
		results <- j * 2 // Mengirim hasil (misalnya: j * 2)
	}
}

// CATATAN PRIBADI: Select Statement untuk menangani multiple channels
func demoSelect() {
	fmt.Println("\n=== Demo Select Statement ===")

	c1 := make(chan string)
	c2 := make(chan string)

	// Goroutine 1
	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "satu"
	}()

	// Goroutine 2
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "dua"
	}()

	// Menggunakan select untuk menunggu channel yang siap
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("Menerima dari c1:", msg1)
		case msg2 := <-c2:
			fmt.Println("Menerima dari c2:", msg2)
		}
	}
}

// CATATAN PRIBADI: Mutex untuk menghindari race condition
func demoMutex() {
	fmt.Println("\n=== Demo Mutex untuk Menghindari Race Condition ===")

	// Counter yang akan diakses oleh banyak goroutine
	var counter = 0

	// Mutex untuk mengunci akses ke counter
	var mu sync.Mutex

	// WaitGroup untuk menunggu semua goroutine selesai
	var wg sync.WaitGroup

	// Jumlah goroutine dan iterasi
	const numGoroutines = 5
	const numIterations = 1000

	wg.Add(numGoroutines)

	// Membuat beberapa goroutine yang mengakses counter
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numIterations; j++ {
				// Lock mutex sebelum mengakses counter
				mu.Lock()
				counter++
				mu.Unlock()
			}

			fmt.Printf("Goroutine %d selesai\n", id)
		}(i)
	}

	// Menunggu semua goroutine selesai
	wg.Wait()

	// Nilai akhir counter
	fmt.Printf("Nilai akhir counter: %d (seharusnya: %d)\n",
		counter, numGoroutines*numIterations)
}

// CATATAN PRIBADI: Buffered vs Unbuffered Channels
func demoChannelTypes() {
	fmt.Println("\n=== Demo Buffered vs Unbuffered Channels ===")

	// Unbuffered channel (kapasitas 0)
	unbuffered := make(chan int)

	// Buffered channel (kapasitas 2)
	buffered := make(chan int, 2)

	// Demonstrasi unbuffered channel
	fmt.Println("Demonstrasi Unbuffered Channel:")
	go func() {
		fmt.Println("Goroutine: Mengirim ke unbuffered channel")
		unbuffered <- 42
		fmt.Println("Goroutine: Berhasil mengirim ke unbuffered channel")
	}()

	time.Sleep(500 * time.Millisecond)
	fmt.Println("Main: Menerima dari unbuffered channel")
	fmt.Println("Nilai dari unbuffered channel:", <-unbuffered)

	// Demonstrasi buffered channel
	fmt.Println("\nDemonstrasi Buffered Channel:")
	fmt.Println("Main: Mengirim nilai 1 ke buffered channel")
	buffered <- 1
	fmt.Println("Main: Mengirim nilai 2 ke buffered channel")
	buffered <- 2

	fmt.Println("Main: Buffered channel sudah terisi 2 nilai")

	// Membaca dari buffered channel
	fmt.Println("Nilai 1 dari buffered channel:", <-buffered)
	fmt.Println("Nilai 2 dari buffered channel:", <-buffered)
}

func demoKonkurensiLanjutan() {
	// Seed random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Println("Belajar Konkurensi Lanjutan di Go")

	// Demo worker pool pattern
	demoWorkerPool()

	// Demo select statement
	demoSelect()

	// Demo mutex untuk menghindari race condition
	demoMutex()

	// Demo channel types
	demoChannelTypes()

	fmt.Println("\nProgram selesai!")
}
