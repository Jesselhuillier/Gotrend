package main

import (
	"fmt"
	"stringutil"
)

func main() {
	fmt.Println(reverse("Golang"))
}
func reverse(chaine string)string{
	return stringutil.Reverse(chaine)
}