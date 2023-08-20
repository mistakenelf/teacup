// package csv allows you to render a csv file as a table
package csv

import (
	"encoding/csv"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	headerOffset  = 2
	columnSpacing = 5
)

type renderCSVMsg [][]string
type errorMsg error

// Model represents the properties of a pdf bubble.
type Model struct {
	Table    table.Model
	Active   bool
	FileName string
}

// ReadCSV reads a PDF file given a name.
func ReadCSV(name string) ([][]string, error) {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// renderCSVCmd reads the content of a CSV and returns its content as a cmd.
func renderCSVCmd(filename string) tea.Cmd {
	return func() tea.Msg {
		csvContent, err := ReadCSV(filename)
		if err != nil {
			return errorMsg(err)
		}

		return renderCSVMsg(csvContent)
	}
}

// New creates a new instance of a CSV.
func New(active bool) Model {
	t := table.New(
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	return Model{
		Table: t,
	}
}

// Init initializes the CSV bubble.
func (m Model) Init() tea.Cmd {
	return nil
}

// SetFileName sets current file to render, this
// returns a cmd which will render the csv.
func (m *Model) SetFileName(filename string) tea.Cmd {
	m.FileName = filename

	return renderCSVCmd(filename)
}

// SetIsActive sets if the bubble is currently active.
func (m *Model) SetIsActive(active bool) {
	m.Active = active
}

// Update handles updating the UI of a code bubble.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Table.SetHeight(msg.Height - headerOffset)
		m.Table.SetWidth(msg.Width)
	case renderCSVMsg:
		var columns []table.Column
		var rows []table.Row

		for _, element := range msg[0] {
			columns = append(columns, table.Column{
				Title: element,
				Width: lipgloss.Width(element) + columnSpacing},
			)
		}

		for _, rowData := range msg[1:] {
			var row table.Row
			for _, colData := range rowData {
				row = append(row, colData)
			}

			rows = append(rows, row)
		}

		m.Table.SetColumns(columns)
		m.Table.SetRows(rows)

		return m, nil
	case errorMsg:
		m.FileName = ""

		return m, nil
	}

	if m.Active {
		m.Table, cmd = m.Table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the csv bubble.
func (m Model) View() string {
	return m.Table.View()
}
