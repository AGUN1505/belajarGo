package main

import "fmt"

var namaDepan string
var namaBelakang string
var umur int = 30

const pi float64 = 3.14

func main() {
	namaDepan = "John"
	namaBelakang = "Doe"

	fmt.Println(namaDepan + " " + namaBelakang)
	fmt.Println("Umur:", umur)
	fmt.Println("Nilai pi:", pi)
	fmt.Println("ini coba")
}
