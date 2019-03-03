package main

import (
	"fmt"
	"sort"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/sudomabider/scoop-ui/scoop"
)

type App struct {
	*scoop.App

	Index   int
	checked bool
}

type AppModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*App
}

func NewAppModel() *AppModel {
	m := new(AppModel)
	// m.Refresh()
	return m
}

// RowCount is called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *AppModel) RowCount() int {
	return len(m.items)
}

// Value is called by the TableView when it needs the text to display for a given cell.
func (m *AppModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index

	case 1:
		return item.Name

	case 2:
		return item.Version

	case 3:
		return item.Bucket
	}

	panic("unexpected col")
}

// Called by the TableView to retrieve if a given row is checked.
func (m *AppModel) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *AppModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked

	return nil
}

// Called by the TableView to sort the model.
func (m *AppModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)

		case 1:
			return c(a.Name < b.Name)

		case 2:
			return c(a.Version < b.Version)

		case 3:
			return c(a.Bucket < b.Bucket)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *AppModel) Refresh() {
	apps, _ := scoop.List()

	m.items = make([]*App, len(apps))

	for i := range m.items {
		m.items[i] = &App{
			App:   &apps[i],
			Index: i,
		}
	}

	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}

type ListPage struct {
	*walk.Composite
}

func newListPage(parent walk.Container) (Page, error) {
	p := new(ListPage)

	model := NewAppModel()

	var tv *walk.TableView

	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "listPage",
		Layout:   VBox{},
		Children: []Widget{
			PushButton{
				Text:      "Refresh Apps",
				OnClicked: model.Refresh,
			},
			TableView{
				AssignTo:              &tv,
				AlternatingRowBGColor: walk.RGB(239, 239, 239),
				CheckBoxes:            true,
				ColumnsOrderable:      true,
				MultiSelection:        true,
				Columns: []TableViewColumn{
					{Title: "#"},
					{Title: "Name"},
					{Title: "Version"},
					{Title: "Bucket", Alignment: AlignFar},
				},
				StyleCell: func(style *walk.CellStyle) {
					item := model.items[style.Row()]

					if item.checked {
						if style.Row()%2 == 0 {
							style.BackgroundColor = walk.RGB(159, 215, 255)
						} else {
							style.BackgroundColor = walk.RGB(143, 199, 239)
						}
					}
				},
				Model: model,
				OnSelectedIndexesChanged: func() {
					fmt.Printf("SelectedIndexes: %v\n", tv.SelectedIndexes())
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	if err := walk.InitWrapperWindow(p); err != nil {
		return nil, err
	}

	return p, nil
}
