package board

// A move is a function that takes two parameters:
// 	- The position of the piece, a number in between 0 and 63 being 0 the top left corner
// 		and 63 the bottom right corner
// 	- The number of times the movement has to be done
//
// This provides a general  and easy way to move pieces. For example, if we want to define
// a rook move, we would do it like this:
// 	rook.move(UP(5))
// Meanig that this rook can move up to 5 times up.
//
// And for pieces that do not have a basic movement, we can create our own movements just by
// combining them using a sum.
// For instance, a knight moves one diagonally one cell, and one cell to any other non diagonal
// cells, which would look like this:
// 	knight.move(UP(2)+RIGHT(1))
type move func(int) int

var (
	// the four basic movement functions
	UP move = func(n int) int {
		return (-8) * n
	}
	DOWN move = func(n int) int {
		return (+8) * n
	}
	LEFT move = func(n int) int {
		return (-1) * n
	}
	RIGHT move = func(n int) int {
		return (+1) * n
	}
)

type Movement struct {
	// definition of the movement
	Move move
	// max amount of cells that that movement can be performed
	Limit int
}
