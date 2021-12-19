package board

import (
	"fmt"
	"math"
	"sync"

	"github.com/jgarrone/battleship-game/server/enum"
)

type BattleshipBoard interface {
	Attack(cell *Cell) enum.AttackOutcome
	Exists(cell *Cell) bool
	Restart() BattleshipBoard
}

type battleshipBoard struct {
	mu       sync.Mutex
	xLength  int
	yLength  int
	ships    map[Cell]bool
	hits     map[Cell]bool
	selector CellSelector
}

type IntGenerator func() int

func NewBattleshipBoard(xLength, yLength int, selector CellSelector) (BattleshipBoard, error) {
	switch {
	case xLength < 1:
		return nil, fmt.Errorf("x-axis length must be greater than zero")
	case yLength < 1:
		return nil, fmt.Errorf("y-axis length must be greater than zero")
	}

	return (&battleshipBoard{
		xLength:  xLength,
		yLength:  yLength,
		selector: selector,
	}).Restart(), nil
}

func (b *battleshipBoard) Attack(cell *Cell) enum.AttackOutcome {
	if cell == nil || !b.isCellInsideBoard(cell) {
		return enum.AttackOutcomeInvalid
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if shipExists := b.ships[*cell]; !shipExists {
		return enum.AttackOutcomeMiss
	}

	if alreadyHit := b.hits[*cell]; alreadyHit {
		return enum.AttackOutcomeAlreadyHit
	}

	b.hits[*cell] = true

	if len(b.hits) == len(b.ships) {
		return enum.AttackOutcomeHitAndWin
	}

	return enum.AttackOutcomeHit
}

func (b *battleshipBoard) Exists(cell *Cell) bool {
	return b.isCellInsideBoard(cell)
}

func (b *battleshipBoard) isCellInsideBoard(cell *Cell) bool {
	return 0 <= cell.X && cell.X < b.xLength && 0 <= cell.Y && cell.Y < b.yLength
}

func (b *battleshipBoard) Restart() BattleshipBoard {
	cellCount := b.xLength * b.yLength
	availableCells := make([]Cell, 0, cellCount)
	for i := 0; i < b.xLength; i++ {
		for j := 0; j < b.yLength; j++ {
			availableCells = append(availableCells, Cell{
				X: i,
				Y: j,
			})
		}
	}

	// Occupy 25% of the board with ships. Just because.
	requiredShips := int(math.Ceil(float64(cellCount) / 4))
	ships := make(map[Cell]bool)
	for i := 0; i < requiredShips; i++ {
		index := b.selector.ChooseFrom(availableCells)
		cell := availableCells[index]
		availableCells = append(availableCells[:index], availableCells[index+1:]...)
		ships[cell] = true
	}

	b.mu.Lock()
	b.ships = ships
	b.hits = make(map[Cell]bool)
	b.mu.Unlock()

	return b
}
