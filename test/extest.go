package main

import (
	"fmt"
)

type person struct {
}

func (this *person) F() {
	fmt.Println("this is person")
}

type student struct {
	person
}

func (this *student) F() {
	fmt.Println("this is a student")
}

func main() {
	stu := student{}
	stu.F()
}
