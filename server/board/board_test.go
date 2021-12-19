package board

import (
	"fmt"
	"testing"

	"github.com/jgarrone/battleship-game/server/enum"
	"github.com/stretchr/testify/assert"
)

func TestNewBattleshipBoard(t *testing.T) {
	tests := []struct {
		name    string
		xLength int
		yLength int
		wantErr error
	}{
		{
			name:    "invalid x length",
			xLength: 0,
			yLength: 10,
			wantErr: fmt.Errorf("x-axis length must be greater than zero"),
		},

		{
			name:    "invalid y length",
			xLength: 10,
			yLength: 0,
			wantErr: fmt.Errorf("y-axis length must be greater than zero"),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			board, err := NewBattleshipBoard(test.xLength, test.yLength, NewDummySelector())
			assert.Equal(t, test.wantErr, err)
			if test.wantErr == nil {
				assert.NotNil(t, board)
			} else {
				assert.Nil(t, board)
			}
		})
	}
}

func TestBattleshipBoard_Attack(t *testing.T) {
	tests := []struct {
		name            string
		xLength         int
		yLength         int
		attackPositions []*Cell
		wantOutcomes    []enum.AttackOutcome
	}{
		{
			name:    "invalid attacks",
			xLength: 10,
			yLength: 10,
			attackPositions: []*Cell{
				{
					X: 0,
					Y: 20,
				},
				{
					X: 20,
					Y: 0,
				},
				{
					X: 0,
					Y: -1,
				},
				{
					X: -1,
					Y: 0,
				},
				{
					X: 0,
					Y: 10,
				},
				{
					X: 10,
					Y: 0,
				},
			},
			wantOutcomes: []enum.AttackOutcome{
				enum.AttackOutcomeInvalid,
				enum.AttackOutcomeInvalid,
				enum.AttackOutcomeInvalid,
				enum.AttackOutcomeInvalid,
				enum.AttackOutcomeInvalid,
				enum.AttackOutcomeInvalid,
			},
		},
		{
			name:    "hit same cell multiple times",
			xLength: 10,
			yLength: 10,
			attackPositions: []*Cell{
				{
					X: 0,
					Y: 0,
				},
				{
					X: 0,
					Y: 0,
				},
				{
					X: 0,
					Y: 0,
				},
			},
			wantOutcomes: []enum.AttackOutcome{
				enum.AttackOutcomeHit,
				enum.AttackOutcomeAlreadyHit,
				enum.AttackOutcomeAlreadyHit,
			},
		},
		{
			name:    "hit different cells",
			xLength: 10,
			yLength: 10,
			attackPositions: []*Cell{
				{
					X: 0,
					Y: 0,
				},
				{
					X: 0,
					Y: 1,
				},
				{
					X: 5,
					Y: 5,
				},
			},
			wantOutcomes: []enum.AttackOutcome{
				enum.AttackOutcomeHit,
				enum.AttackOutcomeHit,
				enum.AttackOutcomeMiss,
			},
		},
		{
			name:    "game won",
			xLength: 2,
			yLength: 4,
			attackPositions: []*Cell{
				{
					X: 0,
					Y: 0,
				},
				{
					X: 0,
					Y: 1,
				},
			},
			wantOutcomes: []enum.AttackOutcome{
				enum.AttackOutcomeHit,
				enum.AttackOutcomeHitAndWin,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			board, err := NewBattleshipBoard(test.xLength, test.yLength, NewDummySelector())
			if err != nil {
				t.Fatalf("error initializing board: %v", err)
			}
			for i, pos := range test.attackPositions {
				assert.Equal(t, test.wantOutcomes[i], board.Attack(pos))
			}
		})
	}
}
