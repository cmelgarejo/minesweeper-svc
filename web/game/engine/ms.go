package engine

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/cmelgarejo/minesweeper-svc/utils"
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
	Mine      bool     `json:"mine"`
	Clicked   bool     `json:"clicked"`   // indicated whether the field was clicked
	Flagged   bool     `json:"flagged"`   // red flag in the field
	AdjCount  int      `json:"adjMines"`  // count of adjacent mines
	Position  Position `json:"position"`  // position in the minefield
	ClickedBy string   `json:"clickedBy"` // who clicked this field
}

// Game contains the structure of the game
type Game struct {
	ID         string     `json:"id"`
	Rows       int        `json:"rows"`
	Cols       int        `json:"cols"`
	Mines      int        `json:"mines"`
	Status     GameStatus `json:"status"`
	MineField  [][]Field  `json:"mineField"`
	StartedAt  *time.Time `json:"startedAt,omitempty"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	CreatedBy  string     `json:"createdBy"` // who created this game
}

func (g *Game) Start() error {
	if g.Status == GameStatusCreated {
		g.Status = GameStatusStarted
		now := time.Now()
		g.StartedAt = &now
		return nil
	}
	return fmt.Errorf("Cannot start a game that is in status: %s", g.Status)
}

func (g *Game) Click(clickedBy string, clickType ClickType, row, col int) error {
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
	g.MineField[row][col].ClickedBy = clickedBy
	switch clickType {
	case GameClickTypeFlag:
		g.MineField[row][col].Flagged = true
	case GameClickTypeNormal:
		if !g.MineField[row][col].Flagged && g.MineField[row][col].Mine {
			g.Status = GameStatusDefeat
			now := time.Now()
			g.FinishedAt = &now
			return ErrDefeat
		}
		if g.MineField[row][col].AdjCount == 0 {
			_ = g.Click(clickedBy, GameClickTypeReveal, row, col)
		}
	case GameClickTypeReveal:
		g.autoReveal(clickedBy, row, col)
	}
	g.printMinefield()
	return nil
}

func NewGame(rows, cols, mines int, createdBy string) *Game {
	if rows < GameMinRows {
		rows = GameMinRows
	}
	if cols < GameMinCols {
		cols = GameMinCols
	}
	if mines < 1 || mines > rows*cols {
		mines = rows + cols // Make sure amount of mines is relative to a median of rows + cols
	}

	id, _ := utils.GenerateGUID()
	newGame := Game{
		ID:        id,
		Rows:      rows,
		Cols:      cols,
		Mines:     mines,
		Status:    GameStatusCreated,
		CreatedAt: time.Now(),
		// CreatedBy: User,
	}

	mineCount := 0
	newGame.MineField = make([][]Field, newGame.Rows)
	for i := 0; i < newGame.Rows; i++ {
		newGame.MineField[i] = make([]Field, newGame.Cols)
		for j := 0; j < newGame.Cols; j++ {
			newGame.MineField[i][j] = Field{
				Position: Position{i, j},
			}
		}
	}
	for mineCount < newGame.Mines {
		seed := rand.NewSource(time.Now().UnixNano())
		row := rand.New(seed).Intn(newGame.Rows)
		col := rand.New(seed).Intn(newGame.Cols)
		if !newGame.MineField[row][col].Mine {
			newGame.MineField[row][col].Mine = true
			mineCount++
			offRow := row
			if row == 0 {
				offRow++
			}
			offCol := col
			if col == 0 {
				offCol++
			}
			for i := offRow - 1; i < newGame.Rows && i <= row+1; i++ {
				for j := offCol - 1; j < newGame.Cols && j <= col+1; j++ {
					if !newGame.MineField[i][j].Mine {
						newGame.MineField[i][j].AdjCount++
					} else {
						newGame.MineField[i][j].AdjCount = 0
					}
				}
			}
		}
	}
	newGame.CreatedBy = createdBy
	newGame.CreatedAt = time.Now()

	return &newGame
}

func (g *Game) IsActive() bool {
	return g.Status == GameStatusStarted
}

func (g *Game) autoReveal(clickedBy string, row, col int) {
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
					_ = g.Click(clickedBy, GameClickTypeReveal, i, j)
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
	fmt.Println()
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
	fmt.Println()
}
