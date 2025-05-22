package client

import (
	"fmt"
	"net"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	// PlayerGridSize is the size of the player grid
	GridSize = 10
)

type Client struct {
	app  *tview.Application
	conn *net.Conn
	// state        *GameState
	playerGrid   *tview.Table
	opponentGrid *tview.Table
	statusView   *tview.TextView
}

func NewClient() (*Client, error) {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		return nil, err
	}

	app := tview.NewApplication()
	client := &Client{
		app:  app,
		conn: &conn,
		// state:        &GameState{Player: player, Status: "Connecting..."},
		playerGrid:   tview.NewTable(),
		opponentGrid: tview.NewTable(),
		statusView:   tview.NewTextView(),
	}
	client.setupUI()
	return client, nil
}

func (c *Client) setupStatusView() {
	tv := c.statusView

	tv.SetDynamicColors(true)
	tv.SetRegions(true)

	tv.SetChangedFunc(func() {
		c.app.Draw()
	})

	tv.SetTextAlign(tview.AlignCenter)
	tv.SetBorder(true)
	tv.SetTitle("Status")

	fmt.Fprintf(tv, "\nSet up fleet\n")
}

func newTableCell() *tview.TableCell {
	cell := tview.NewTableCell("     ")
	cell.SetAlign(tview.AlignCenter)
	cell.SetBackgroundColor(tcell.ColorDimGray)
	// cell.SetExpansion(1)
	return cell
}

func setupGrid(t *tview.Table) {
	t.SetBorders(true)
	// t.SetBorder(true)
	t.SetFixed(GridSize, GridSize)

	for row := 0; row < GridSize; row++ {
		for col := 0; col < GridSize; col++ {
			cell := newTableCell()
			t.SetCell(row, col, cell)
		}
	}
}

func (c *Client) setupPlayerGrid() {
	setupGrid(c.playerGrid)
	// c.playerGrid.SetTitle("Player Grid")
}

func (c *Client) setupOpponentGrid() {
	setupGrid(c.opponentGrid)
	// c.opponentGrid.SetTitle("Opponent Grid")
}

func setupLetterLabels() *tview.Table {
	const cellWidth = 2
	// Create letter labels (A-J) for left table (above)
	letterLabels := tview.NewTable()
	col := 1
	for i := 0; i < 10; i++ {
		col = col + 5
		letterLabels.SetCell(0, col, tview.NewTableCell(string(rune('A'+i))).
			SetAlign(tview.AlignCenter).
			SetMaxWidth(cellWidth))
	}
	letterLabels.SetFixed(1, 10) // Fix 1 row, 10 columns
	// letterLabels.SetBorder(true)

	return letterLabels
}

func setupNumberLabels() *tview.Table {
	const cellWidth = 2
	// Create number labels (1-10) for top table (left)
	numberLabels := tview.NewTable()
	row := -1
	for i := 0; i < 10; i++ {
		row = row + 2
		numberLabels.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%d", i+1)).
			SetAlign(tview.AlignCenter).
			SetMaxWidth(cellWidth))
	}
	numberLabels.SetFixed(10, 1) // Fix 10 rows, 1 column
	// numberLabels.SetBorder(true)

	return numberLabels
}

func (c *Client) setupFirstRow() *tview.Flex {
	const tableWidth = 63
	const cellWidth = 4

	// come above
	letterLabels := setupLetterLabels()
	// come left
	numberLabels := setupNumberLabels()

	leftTableContainer := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 2, 1, false).
		AddItem(numberLabels, 3, 1, false).           // Empty space on the left
		AddItem(c.playerGrid, tableWidth, 20, false). // Table in the middle
		AddItem(nil, 0, 1, false)                     // Empty space on the right
	// leftTableContainer.SetBorder(true)

	leftFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 1, 1, false).
		AddItem(letterLabels, 1, 1, false). // Empty space on the top
		AddItem(leftTableContainer, 0, 1, false)
	leftFlex.SetTitle("Player Fleet")
	leftFlex.SetBorder(true)

	rightTableContainer := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 2, 1, false).
		AddItem(numberLabels, 3, 1, false).             // Empty space on the left
		AddItem(c.opponentGrid, tableWidth, 20, false). // Table in the middle
		AddItem(nil, 0, 1, false)                       // Empty space on the right
	// rightTableContainer.SetBorder(true)

	rightFlex := tview.NewFlex().
		AddItem(nil, 1, 1, false).
		SetDirection(tview.FlexRow).
		AddItem(letterLabels, 1, 1, false). // Empty space on the top
		AddItem(rightTableContainer, 0, 1, false)
	rightFlex.SetTitle("Opponent Fleet")
	rightFlex.SetBorder(true)

	firstRow := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(leftFlex, 0, 1, false).
		AddItem(rightFlex, 0, 1, false)

	firstRow.SetBorder(true).SetTitle("Player vs Opponent")
	// firstRow.SetBorder(false)
	return firstRow
}

func (c *Client) setupUI() {
	c.setupStatusView()
	c.setupPlayerGrid()
	c.setupOpponentGrid()
	firstRow := c.setupFirstRow()

	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(firstRow, 0, 1, false).      // First row with tables
		AddItem(c.statusView, 3+2, 1, false) // Status view
	mainFlex.SetBorder(true).SetTitle("Main Layout")

	c.app.SetRoot(mainFlex, true)
}

func (c *Client) Run() error {
	if err := c.app.Run(); err != nil {
		return err
	}
	return nil
}
