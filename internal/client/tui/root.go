package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Screen interface {
	Init() tea.Cmd
	Update() (Screen, tea.Cmd)
	View() string
}

type RootModel struct {
	currentScreen Screen
}

func (this RootModel) Update(message tea.Msg) (Screen, tea.Cmd) {
	this.currentScreen, cmd := this.currentScreen.Update(message)
	return this, cmd
}

func (this RootModel) View() string {
	return this.currentScreen.View()
}
