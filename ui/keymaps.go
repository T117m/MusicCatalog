package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

type keymap struct {
	next, prev             key.Binding
	quit                   key.Binding
	insert, update, remove key.Binding
}

var (
	trackListKeyMap = keymap{
		next: key.NewBinding(
			key.WithKeys("tab", "ctrl+n"),
			key.WithHelp("Tab/Ctrl+n", "Следующий трек"),
		),
		prev: key.NewBinding(
			key.WithKeys("shift+tab", "ctrl+p"),
			key.WithHelp("S-Tab/C-n", "Предыдущий трек"),
		),
		quit: key.NewBinding(
			key.WithKeys("esc", "ctrl-c"),
			key.WithHelp("Esc/C-c", "Закрыть программу"),
		),
		insert: key.NewBinding(
			key.WithKeys("ctrl-a"),
			key.WithHelp("C-a", "Добавить трек"),
		),
		update: key.NewBinding(
			key.WithKeys("ctrl-r"),
			key.WithHelp("C-r", "Редактировать трек"),
		),
		remove: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("C-r", "Редактировать трек"),
		),
	}

	addTrackKeyMap = keymap{
		next: key.NewBinding(
			key.WithKeys("tab", "ctrl+n"),
			key.WithHelp("Tab/Ctrl+n", "Следующее поле"),
		),
		prev: key.NewBinding(
			key.WithKeys("shift+tab", "ctrl+p"),
			key.WithHelp("S-Tab/C-n", "Предыдущее поле"),
		),
		quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("Esc", "Отмена"),
		),
		insert: key.NewBinding(
			key.WithKeys("ctrl-s"),
			key.WithHelp("C-s", "Сохранить"),
		),
		update: key.NewBinding(
			key.WithDisabled(),
		),
		remove: key.NewBinding(
			key.WithDisabled(),
		),
	}

	editTrackKeyMap = keymap{
		next: key.NewBinding(
			key.WithKeys("tab", "ctrl+n"),
			key.WithHelp("Tab/Ctrl+n", "Следующее поле"),
		),
		prev: key.NewBinding(
			key.WithKeys("shift+tab", "ctrl+p"),
			key.WithHelp("S-Tab/C-n", "Предыдущее поле"),
		),
		quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("Esc", "Отмена"),
		),
		insert: key.NewBinding(
			key.WithDisabled(),
		),
		update: key.NewBinding(
			key.WithKeys("ctrl-s"),
			key.WithHelp("C-s", "Сохранить изменения"),
		),
		remove: key.NewBinding(
			key.WithDisabled(),
		),
	}

	deleteTrackKeyMap = keymap{
		next: key.NewBinding(
			key.WithDisabled(),
		),
		prev: key.NewBinding(
			key.WithDisabled(),
		),
		quit: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "Отмена"),
		),
		insert: key.NewBinding(
			key.WithDisabled(),
		),
		update: key.NewBinding(
			key.WithDisabled(),
		),
		remove: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "Удалить трек"),
		),
	}
)
