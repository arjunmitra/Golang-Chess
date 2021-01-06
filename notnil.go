package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/notnil/chess"
)

func main() {

	game := chess.NewGame()
	for game.Outcome() == chess.NoOutcome {

		printLine()
		fmt.Println(game.Position().Board().Draw())
		reader := bufio.NewReader(os.Stdin)
		moves := game.ValidMoves()
		if game.Position().Turn().String() == "w" {
			fmt.Println("Your Turn. Please make a move")
			for true {
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				input = strings.ToLower(input)
				isValid := validate(moves, input)
				if isValid == -1 {
					fmt.Println("please input a valid move")
					continue
				}
				game.Move(moves[isValid])
				break
			}

		} else {

			fmt.Println("Opponent Turn")
			index := 0
			maxim := -100000
			for i := 0; i < len(moves); i++ {
				gameCopy := game.Clone()
				gameCopy.Move(moves[i])
				score := AI_move(gameCopy, 3, false)

				if score > maxim {
					index = i
					maxim = score

				}
			}
			game.Move(moves[index])
		}
	}
	winner(game)

}

func printLine() {
	for i := 1; i < 100; i++ {
		fmt.Print("_")
	}
	fmt.Println("")
}
func winner(game *chess.Game) {
	fmt.Println("")
	result := game.Outcome().String()
	if result == "*" {
		fmt.Println("No outcome was possible")
	} else if result == "1-0" {
		fmt.Println("You won the game")
	} else if result == "0-1" {
		fmt.Println("You lost the game")
	} else {
		fmt.Println("The game ended in a draw")
	}
}

func validate(moves []*chess.Move, input string) int {
	for i := 0; i < len(moves); i++ {
		if moves[i].String() == strings.ToLower(input) {
			return i
		}
	}
	return -1
}

func AI_move(game *chess.Game, depth int, isMaximizing bool) int {

	if depth == 0 {

		value := 0
		mapping := game.Position().Board().SquareMap()

		for key := range mapping {

			row := getRow(key.String())
			col := getCol(key.String())

			element := mapping[key]

			curr := 1

			Type := element.Type().String()
			color := element.Color().String()

			if color == "w" {
				curr = curr * -1

			}
			if Type == "k" {
				score := king(row, col)
				curr = curr * 1000 * score

			} else if Type == "q" {
				score := queen(row, col)
				curr = curr * 500 * score

			} else if Type == "r" {
				score := rook(row, col)
				curr = curr * 300 * score
			} else if Type == "n" {
				score := knight(row, col)
				curr = curr * 100 * score
			} else if Type == "" {
				score := pawn(row, col)
				curr = curr * 1 * score

			} else if Type == "b" {
				score := bishop(row, col)
				curr = curr * 150 * score
			}
			value = value + curr

		}
		return value
	}

	if isMaximizing == true {
		moves := game.ValidMoves()
		value := -100000
		for i := 0; i < len(moves); i++ {
			copy := game.Clone()
			copy.Move(moves[i])
			eval := AI_move(copy, depth-1, false)
			if eval > value {
				value = eval
			}
		}
		return value
	} else {
		moves := game.ValidMoves()
		value := 100000
		for i := 0; i < len(moves); i++ {
			copy := game.Clone()
			copy.Move(moves[i])
			eval := AI_move(copy, depth-1, true)
			if eval < value {
				value = eval
			}

		}
		return value
	}

	return 0
}

func pawn(row int, col int) int {
	values := [8][8]int{

		{0, 0, 0, 0, 0, 0, 0, 0},
		{5, 10, 10, -20, -20, 10, 10, 5},
		{5, -5, -10, 0, 0, -10, -5, 5},
		{0, 0, 0, 20, 20, 0, 0, 0},
		{5, 5, 10, 25, 25, 10, 5, 5},
		{10, 10, 20, 30, 30, 20, 10, 10},
		{50, 50, 50, 50, 50, 50, 50, 50},
		{0, 0, 0, 0, 0, 0, 0, 0},
	}
	return values[row][col]

}

func knight(row int, col int) int {
	values := [8][8]int{
		{-50, -40, -30, -30, -30, -30, -40, -50},
		{-40, -20, 0, 5, 5, 0, -20, -40},
		{-30, 5, 10, 15, 15, 10, 5, -30},
		{-30, 0, 15, 20, 20, 15, 0, -30},
		{-30, 5, 15, 20, 20, 15, 5, -30},
		{-30, 0, 10, 15, 15, 10, 0, -30},
		{-40, -20, 0, 0, 0, 0, -20, -40},
		{-50, -40, -30, -30, -30, -30, -40, -50},
	}
	return values[row][col]
}

func bishop(row int, col int) int {
	values := [8][8]int{
		{-20, -10, -10, -10, -10, -10, -10, -20},
		{-10, 5, 0, 0, 0, 0, 5, -10},
		{-10, 10, 10, 10, 10, 10, 10, -10},
		{-10, 0, 10, 10, 10, 10, 0, -10},
		{-10, 5, 5, 10, 10, 5, 5, -10},
		{-10, 0, 5, 10, 10, 5, 0, -10},
		{-10, 0, 0, 0, 0, 0, 0, -10},
		{-20, -10, -10, -10, -10, -10, -10, -20},
	}
	return values[row][col]
}

func rook(row int, col int) int {
	values := [8][8]int{
		{0, 0, 0, 5, 5, 0, 0, 0},
		{-5, 0, 0, 0, 0, 0, 0, -5},
		{-5, 0, 0, 0, 0, 0, 0, -5},
		{-5, 0, 0, 0, 0, 0, 0, -5},
		{-5, 0, 0, 0, 0, 0, 0, -5},
		{-5, 0, 0, 0, 0, 0, 0, -5},
		{5, 10, 10, 10, 10, 10, 10, 5},
		{0, 0, 0, 0, 0, 0, 0, 0},
	}
	return values[row][col]
}

func queen(row int, col int) int {
	values := [8][8]int{
		{-20, -10, -10, -5, -5, -10, -10, -20},
		{-10, 0, 5, 0, 0, 0, 0, -10},
		{-10, 5, 5, 5, 5, 5, 0, -10},
		{0, 0, 5, 5, 5, 5, 0, -5},
		{-5, 0, 5, 5, 5, 5, 0, -5},
		{-10, 0, 5, 5, 5, 5, 0, -10},
		{-10, 0, 0, 0, 0, 0, 0, -10},
		{-20, -10, -10, -5, -5, -10, -10, -20},
	}
	return values[row][col]
}

func king(row int, col int) int {
	values := [8][8]int{
		{20, 30, 10, 0, 0, 10, 30, 20},
		{20, 20, 0, 0, 0, 0, 20, 20},
		{-10, -20, -20, -20, -20, -20, -20, -10},
		{-20, -30, -30, -40, -40, -30, -30, -20},
		{-30, -40, -40, -50, -50, -40, -40, -30},
		{-30, -40, -40, -50, -50, -40, -40, -30},
		{-30, -40, -40, -50, -50, -40, -40, -30},
		{-30, -40, -40, -50, -50, -40, -40, -30},
	}
	return values[row][col]
}

func getCol(sqr string) int {
	var letter string
	for i, c := range sqr {
		if i == 0 {
			letter = string(c)
			break
		}
	}
	if letter == "a" {
		return 0
	} else if letter == "b" {
		return 1
	} else if letter == "c" {
		return 2
	} else if letter == "d" {
		return 3
	} else if letter == "e" {
		return 4
	} else if letter == "f" {
		return 5
	} else if letter == "g" {
		return 6
	} else {
		return 7
	}

}
func getRow(sqr string) int {
	row := sqr[1:]
	if row == "8" {
		return 0
	} else if row == "7" {
		return 1
	} else if row == "6" {
		return 2
	} else if row == "5" {
		return 3
	} else if row == "4" {
		return 4
	} else if row == "3" {
		return 5
	} else if row == "2" {
		return 6
	} else {
		return 7
	}

}
