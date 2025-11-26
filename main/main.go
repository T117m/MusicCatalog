package main

import (
	"log"

	"github.com/T117m/MusicCatalog/player"
	"github.com/T117m/MusicCatalog/storage"
	"github.com/T117m/MusicCatalog/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	store, err := storage.New()
	if err != nil {
		log.Fatalf("Ошибка инициализации хранилища: %v\n", err)
	}
	defer store.Close()

	p := player.New()

	m := ui.New(store, p)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		log.Fatalf("error running tui: %v", err)
	}
}
