package main

import "fmt"

const (
	LEFT_PAREN rune = '('
	RIGHT_PAREN rune = ')'
)

func scanTokens(fileContents string){
	for _, current := range fileContents{
		switch current {
			case LEFT_PAREN:
				fmt.Println("LEFT_PAREN ( null")
			case RIGHT_PAREN:
				fmt.Println("RIGHT_PAREN ) null")
		}
	}
	fmt.Println("EOF  null")
}