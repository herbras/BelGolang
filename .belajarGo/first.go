package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Variabel global yang diminta pengguna
var firstName string = "Ibrahim"
var lastName = "Huda" // type inferred as string
var Satu, Dua, Tiga string

// CATATAN PRIBADI: Deklarasi variabel di Go bisa menggunakan var atau :=
// CATATAN PRIBADI: := hanya bisa digunakan di dalam fungsi, tidak bisa di level package
// CATATAN PRIBADI: Variabel yang tidak digunakan di Go akan menyebabkan error kompilasi

// Function to add two numbers
func add(a, b int) int {
	// CATATAN PRIBADI: Fungsi di Go bisa mengembalikan beberapa nilai sekaligus
	return a + b
}

// Recursive function to calculate power of a number
func power(base, exponent int) int {
	// CATATAN PRIBADI: Rekursi di Go sama seperti bahasa lain, tapi hati-hati dengan stack overflow
	if exponent == 0 {
		return 1
	}
	return base * power(base, exponent-1)
}

// Function to calculate factorial of a number
func factorial(n int) int {
	// CATATAN PRIBADI: Faktorial bisa juga diimplementasikan dengan loop untuk performa lebih baik
	if n == 0 || n == 1 {
		return 1
	}
	return n * factorial(n-1)
}

// Function to find the largest element in an array
func findLargest(arr []int) int {
	// CATATAN PRIBADI: Di Go, array punya ukuran tetap, sedangkan slice lebih fleksibel
	largest := arr[0]
	for _, value := range arr {
		if value > largest {
			largest = value
		}
	}
	return largest
}

// Function to print a star triangle
func printStarTriangle(height int) {
	// CATATAN PRIBADI: Loop di Go hanya ada for, tidak ada while atau do-while
	for i := 1; i <= height; i++ {
		for j := 1; j <= i; j++ {
			fmt.Print("*")
		}
		fmt.Println()
	}
}

// Function to print Fibonacci numbers
func printFibonacci(n int) {
	// CATATAN PRIBADI: Multiple assignment di Go sangat berguna untuk swap dan kasus seperti ini
	a, b := 0, 1
	fmt.Print("Fibonacci Series: ")
	for i := 0; i < n; i++ {
		fmt.Printf("%d ", a)
		a, b = b, a+b
	}
	fmt.Println()
}

// Function to count lines in a file
func countLines(filePath string) (int, error) {
	// CATATAN PRIBADI: Go memiliki error handling yang eksplisit, bukan exception seperti bahasa lain
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// CATATAN PRIBADI: defer akan menjalankan fungsi setelah fungsi yang melingkupinya selesai
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

func main() {
	// CATATAN PRIBADI: Menaruh type inference := hanya bisa digunakan di dalam fungsi
	// CATATAN PRIBADI: Atau pemanggilan berulang juga harus dalam fungsi

	// Print "Hello, World!" to console
	fmt.Println("Hello, World!")

	// Contoh penggunaan variabel yang dideklarasikan di level package
	fmt.Printf("Nama lengkap: %s %s\n", firstName, lastName)

	// Menggunakan pointer seperti di kode asli pengguna
	name := new(string)
	fmt.Println("Pointer address:", name) // Akan mencetak alamat memori
	fmt.Println("Pointer value:", *name)  // Akan mencetak nilai kosong ("")

	// Menggunakan variabel middleName seperti di kode asli pengguna
	middleName := "Nurul"
	fmt.Printf("halo, %s %s %s ! \n", firstName, middleName, lastName)

	// Menggunakan variabel Satu, Dua, Tiga seperti di kode asli pengguna
	Satu, Dua, Tiga = "Satu", "Dua", "Tiga"
	fmt.Println("halo", Satu, Dua, Tiga)

	// Declare two variables a and b
	var a, b int

	// Set their values
	a = 10
	b = 20
	fmt.Printf("Before swap: a = %d, b = %d\n", a, b)

	// Swap their values
	a, b = b, a
	fmt.Printf("After swap: a = %d, b = %d\n", a, b)

	// Define a function that takes two numbers and returns their sum
	result := add(a, b)
	fmt.Printf("Sum of %d and %d is %d\n", a, b, result)

	// Ask for an integer input
	fmt.Print("Please enter an integer: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input. Using default value 10.")
		num = 10
	}

	// Check if the number is even or odd
	if num%2 == 0 {
		fmt.Printf("%d is even\n", num)
	} else {
		fmt.Printf("%d is odd\n", num)
	}

	// Print 10 first Fibonacci numbers
	printFibonacci(10)

	// Calculate factorial of a number
	factResult := factorial(num)
	fmt.Printf("Factorial of %d is %d\n", num, factResult)

	// Find the largest element in an array
	arr := []int{3, 7, 2, 9, 1, 5}
	largest := findLargest(arr)
	fmt.Printf("The largest element in the array is %d\n", largest)

	// Ask for the height of a star triangle
	fmt.Print("Enter the height of the star triangle: ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	height, err := strconv.Atoi(input)
	if err != nil || height <= 0 {
		fmt.Println("Invalid input. Using default height 5.")
		height = 5
	}

	// Print a star triangle with the provided height
	printStarTriangle(height)

	// Define a recursive function to calculate power of a number
	base := 2
	exponent := 3
	powerResult := power(base, exponent)
	fmt.Printf("%d raised to the power %d is %d\n", base, exponent, powerResult)

	// Create a sample text file
	filePath := "sample.txt"
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
	} else {
		defer file.Close()
		file.WriteString("Line 1\nLine 2\nLine 3\nLine 4\nLine 5\n")
		fmt.Println("Sample file created successfully.")
	}

	// Read a text file and count the number of lines
	lineCount, err := countLines(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	} else {
		fmt.Printf("Number of lines in the file: %d\n", lineCount)
	}

	// CATATAN PRIBADI: Memanggil fungsi dari Type.go untuk demonstrasi tipe data
	fmt.Println("\n=== Demonstrasi Tipe Data dari Type.go ===")
	typeExample()

	// CATATAN PRIBADI: Memanggil fungsi pembelajaran lanjutan
	fmt.Println("\n=== Pembelajaran Lanjutan Golang ===")
	demoLanjutBelajar()
	demoKonkurensiLanjutan()
	demoStructInterface()
	demoGoModules()

	fmt.Println("\nEnd Program")
}
