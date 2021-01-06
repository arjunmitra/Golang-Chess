package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/notnil/chess"
)

func main() {

	game := chess.NewGame() //creating new game
	for game.Outcome() == chess.NoOutcome {
		//printing a line to segregate between moves
		printLine()
		fmt.Println(game.Position().Board().Draw()) //printing the current chess board
		reader := bufio.NewReader(os.Stdin)
		moves := game.ValidMoves() //retrieving valid moves
		if game.Position().Turn().String() == "w" {
			fmt.Println("Your Turn. Please make a move")
			for true {
				input, _ := reader.ReadString('\n') // getting user input for what move they want to play
				input = strings.TrimSpace(input)
				input = strings.ToLower(input)
				isValid := validate(moves, input) // validating if this a legal move
				if isValid == -1 {
					fmt.Println("please input a valid move")
					continue
				}
				game.Move(moves[isValid]) //updating chess board with move
				break
			}

		} else {

			fmt.Println("Opponent Turn")
			index := 0
			maxim := -100000 // using Minimax algorithm to get best move for AI
			for i := 0; i < len(moves); i++ {
				gameCopy := game.Clone() // creating copy of the game
				gameCopy.Move(moves[i])  // going through all possible moves and seeing which move leaves us in the best position to win
				score := AI_move(gameCopy, 3, false)

				if score > maxim { //updating best move
					index = i
					maxim = score

				}
			}
			game.Move(moves[index])
		}
	}
	winner(game) // printing the result of the game

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
		mapping := game.Position().Board().SquareMap() //retrieving map of peices and where these peices are located

		for key := range mapping {

			row := getRow(key.String()) //getting the location of each  peice
			col := getCol(key.String())

			element := mapping[key]

			curr := 1

			Type := element.Type().String()
			color := element.Color().String()

			if color == "w" { // if w, these peices will like to minimize the AI's score, thus we make them - ve of the black peices score
				curr = curr * -1

			} // each peice is given a weight (1000 for king) and we get a score based on where the peice is on the board
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
			value = value + curr // updating the score at this position

		}
		return value
	}
	// it is the AI's turn
	if isMaximizing == true {
		moves := game.ValidMoves()
		value := -100000
		for i := 0; i < len(moves); i++ {
			copy := game.Clone()
			copy.Move(moves[i])
			eval := AI_move(copy, depth-1, false) //searching all possible positions upto a certain depth
			if eval > value {
				value = eval
			}
		}
		return value
	} else { // it is the person's turn
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

}

//these table give a score to each peice depending on their location
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
