package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

type keymap struct {
	next, prev key.Binding
	quit       key.Binding
	action     key.Binding
}

type subKeymap struct {
	add, edit, delete key.Binding
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
			key.WithKeys("esc", "q"),
			key.WithHelp("Esc/q", "Закрыть программу"),
		),
		action: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("Enter", "Включить/Остановить трек"),
		),
	}

	trackListSubKeyMap = subKeymap{
		add: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("C-a", "Добавить трек"),
		),
		edit: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("C-r", "Редактировать трек"),
		),
		delete: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "Удалить трек"),
		),
	}

	addTrackKeyMap = keymap{
		next: key.NewBinding(
			key.WithKeys("tab", "ctrl+n", "enter"),
			key.WithHelp("Tab/Ctrl+n", "Следующее поле"),
		),
		prev: key.NewBinding(
			key.WithKeys("shift+tab", "ctrl+p", "shift+enter"),
			key.WithHelp("S-Tab/C-n", "Предыдущее поле"),
		),
		quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("Esc", "Отмена"),
		),
		action: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("C-s", "Сохранить"),
		),
	}

	editTrackKeyMap = keymap{
		next: key.NewBinding(
			key.WithKeys("tab", "ctrl+n", "enter"),
			key.WithHelp("Tab/Ctrl+n", "Следующее поле"),
		),
		prev: key.NewBinding(
			key.WithKeys("shift+tab", "ctrl+p", "shift+enter"),
			key.WithHelp("S-Tab/C-n", "Предыдущее поле"),
		),
		quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("Esc", "Отмена"),
		),
		action: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("C-s", "Сохранить изменения"),
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
			key.WithKeys("n", "esc", "q"),
			key.WithHelp("n/Esc/q", "Отмена"),
		),
		action: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "Удалить трек"),
		),
	}
)
