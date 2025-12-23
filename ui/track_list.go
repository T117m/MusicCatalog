package ui

import (
	"strconv"

	"github.com/T117m/MusicCatalog/storage"
	"github.com/T117m/MusicCatalog/music"

	"github.com/charmbracelet/bubbles/table"
)

var columns = []table.Column{
	{Title: "ID", Width: 4},
	{Title: "Название", Width: 12},
	{Title: "Исполнитель", Width: 12},
	{Title: "Тип файла", Width: 10},
	{Title: "Жанр", Width: 10},
}

func defaultTrackList(store *storage.Storage) table.Model {
	tracks, _ := store.GetAllTracks()
	return newTrackList(tracks)
}

func newTrackList(tracks []music.Track) table.Model {
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
