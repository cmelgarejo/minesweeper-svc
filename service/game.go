package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDefeat    = errors.New("Defeat")
	ErrNotActive = errors.New("Game not active")
)

// Some default game parameters, if the user does not provide those.
const (
	GameMinRows         = 3
	GameMinCols         = 3
	GameMaxRows         = 50
	GameMaxCols         = 50
	GameStatusCreated   = "created"
	GameStatusStarted   = "started"
	GameStatusVictory   = "victory"
	GameStatusDefeat    = "defeat"
	GameClickTypeNormal = 1
	GameClickTypeFlag   = 2
	GameClickTypeReveal = 3
)

type GameStatus string
type ClickType int

// Position stores the position of the field in the board
type Position struct {
	Row int `json:"row"` // row of the field position
	Col int `json:"col"` // col of the field position
}

//Field represents a square unit in the MineField
type Field struct {
	Mine     bool     `json:"-"`
	Clicked  bool     `json:"clicked"`  // indicated whether the field was clicked
	Flagged  bool     `json:"flagged"`  // red flag in the field
	AdjCount int      `json:"adjMines"` // count of adjacent mines
	Position Position `json:"position"` // position in the minefield
	// ClickedBy User   `json:"clickedBy"` // who clicked this field
}

// Game contains the structure of the game
type Game struct {
	ID         uuid.UUID  `json:"id"`
	Rows       int        `json:"rows"`
	Cols       int        `json:"cols"`
	Mines      int        `json:"mines"`
	Status     GameStatus `json:"status"`
	MineField  [][]Field  `json:"mineField"`
	CreatedAt  time.Time  `json:"createdAt"`
	StartedAt  *time.Time `json:"startedAt,omitempty"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
	// CreatedBy User   `json:"createdBy"` // who created this game
}

func (g *Game) Start() error {
	if g.Status == GameStatusCreated {
		g.Status = GameStatusStarted
		return nil
	}
	return fmt.Errorf("Cannot start a game that is in status: %s", g.Status)
}

func (g *Game) Click(clickType ClickType, row, col int) error {
	g.printMinefieldPos()
	if !g.IsActive() {
		return ErrNotActive
	}
	if row < 0 || row > g.Rows || col < 0 || col > g.Cols {
		return fmt.Errorf("Field [%d, %d] out of bounds", row, col)
	}
	if g.MineField[row][col].Clicked && clickType != GameClickTypeReveal {
		return nil
	}
	g.MineField[row][col].Clicked = true
	switch clickType {
	case GameClickTypeFlag:
		g.MineField[row][col].Flagged = true
	case GameClickTypeNormal:
		if !g.MineField[row][col].Flagged && g.MineField[row][col].Mine {
			g.Status = GameStatusDefeat
			return ErrDefeat
		}
		if g.MineField[row][col].AdjCount == 0 {
			_ = g.Click(GameClickTypeReveal, row, col)
		}
	case GameClickTypeReveal:
		g.autoReveal(row, col)
	}
	g.printMinefield()
	return nil
}

func (g *Game) InitializeMinefield() {
	mineCount := 0
	g.MineField = make([][]Field, g.Rows)
	for i := 0; i < g.Rows; i++ {
		g.MineField[i] = make([]Field, g.Cols)
		for j := 0; j < g.Cols; j++ {
			g.MineField[i][j] = Field{
				Position: Position{i, j},
			}
		}
	}
	for mineCount < g.Mines {
		seed := rand.NewSource(time.Now().UnixNano())
		row := rand.New(seed).Intn(g.Rows)
		col := rand.New(seed).Intn(g.Cols)
		if !g.MineField[row][col].Mine {
			g.MineField[row][col].Mine = true
			mineCount++
			offRow := row
			if row == 0 {
				offRow++
			}
			offCol := col
			if col == 0 {
				offCol++
			}
			for i := offRow - 1; i < g.Rows && i <= row+1; i++ {
				for j := offCol - 1; j < g.Cols && j <= col+1; j++ {
					if !g.MineField[i][j].Mine {
						g.MineField[i][j].AdjCount++
					} else {
						g.MineField[i][j].AdjCount = 0
					}
				}
			}
		}
	}
	g.Status = GameStatusCreated
}

func (g *Game) IsActive() bool {
	return g.Status == GameStatusStarted
}

func (g *Game) autoReveal(row, col int) {
	offRow := row
	if row == 0 {
		offRow++
	}
	offCol := col
	if col == 0 {
		offCol++
	}
	for i := offRow - 1; i < g.Rows && i <= row+1; i++ {
		for j := offCol - 1; j < g.Cols && j <= col+1; j++ {
			if !g.MineField[i][j].Mine && g.MineField[i][j].AdjCount == 0 {
				if !g.MineField[i][j].Clicked {
					_ = g.Click(GameClickTypeReveal, i, j)
				}
			}
		}
	}
}

// Functions that pretty print the mine field :)

func (g *Game) printMinefield() {
	fmt.Println()
	for i := 0; i < g.Rows; i++ {
		fmt.Println()
		for j := 0; j < g.Cols; j++ {
			fmt.Print("[")
			if g.MineField[i][j].Flagged {
				fmt.Print("F")
			} else if g.MineField[i][j].Clicked && g.MineField[i][j].AdjCount > 0 {
				fmt.Printf("%d", g.MineField[i][j].AdjCount)
			} else if g.MineField[i][j].Clicked {
				fmt.Print("C")
			} else if g.MineField[i][j].Mine {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
			fmt.Print("] ")
		}
	}
}

func (g *Game) printMinefieldPos() {
	fmt.Println()
	for i := 0; i < g.Rows; i++ {
		fmt.Println()
		for j := 0; j < g.Cols; j++ {
			fmt.Printf("%v", g.MineField[i][j].Position)
			// fmt.Printf("[%v", g.MineField[i][j].Position.Row)
			// fmt.Printf(", %v]", g.MineField[i][j].Position.Col)
		}
	}
}
