package ui

import (
	"strconv"

	"github.com/T117m/MusicCatalog/storage"

	"github.com/charmbracelet/bubbles/table"
)

func newTrackList(store *storage.Storage) table.Model {
	tracks, _ := store.GetAllTracks()

	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Название", Width: 12},
		{Title: "Исполнитель", Width: 12},
		{Title: "Тип файла", Width: 10},
		{Title: "Жанр", Width: 10},
	}

	var rows []table.Row
	for _, track := range tracks {
		row := []string{
			strconv.Itoa(track.ID),
			track.Title,
			track.Artist,
			track.FileType,
			track.Genre,
		}

		rows = append(rows, row)
	}

	return newStyledTable(columns, rows)
}
