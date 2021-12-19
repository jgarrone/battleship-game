package board

import (
	"math/rand"
	"time"
)

// CellSelector provides an interface for choosing a cell from a slice using different strategies.
type CellSelector interface {
	// ChooseFrom returns the index of the selected cell.
	ChooseFrom(cells []Cell) int
}

// dummySelector always chooses the first cell.
type dummySelector struct {
}

func NewDummySelector() CellSelector {
	return &dummySelector{}
}

func (f *dummySelector) ChooseFrom(_ []Cell) int {
	return 0
}

// randomSelector chooses items at random.
type randomSelector struct {
}

func NewRandomSelector() CellSelector {
	// Seed random generator with current time.
	now := time.Now().UTC().UnixNano()
	rand.Seed(now)

	return &randomSelector{}
}

func (f *randomSelector) ChooseFrom(cells []Cell) int {
	return rand.Intn(len(cells))
}
