package main

import (
	"fmt"
)

func main() {

	row, column := input()
	if !validate(row, column) {
		return
	}
	calculate(row, column)

}

func input() (int, int) {
	var row int
	var column int
	fmt.Println("Введите первое значение")
	fmt.Scan(&row)
	fmt.Println("Введите второе значение")
	fmt.Scan(&column)
	return row, column
}

func validate(x, y int) bool {
	if x == 0 || y == 0 {
		fmt.Println("Ошибка: значение не указано")
		return false
	}

	if x < 0 || y < 0 {
		fmt.Println("Ошибка: число не может быть отрицательным")
		return false
	}
	return true
}

func calculate(row, column int) {
	row2 := row / 2
	for q := 0; q < column; q++ {
		if q%2 == 0 {
			for q := 0; q < row2; q++ {
				fmt.Print("# ")
			}
			if row%2 == 1 {
				fmt.Print("#")
			}
		} else {
			for q := 0; q < row2; q++ {
				fmt.Print(" #")
			}
			if row%2 == 1 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
