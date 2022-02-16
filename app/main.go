package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	itunesman "github.com/kuwa72/ituweak"
	intf "github.com/kuwa72/ituweak/interfaces"
	"github.com/rivo/tview"
)

var itunes intf.ITunes
var currentTrack intf.Track

var application *tview.Application
var tracksList *tview.List
var assignedList *tview.List
var logText *tview.TextView

func playlistHandler(pl intf.Playlist) func() {
	return func() {
		if tracksList == nil {
			return
		}

		tracksList.Clear()
		tl, err := pl.Tracks()
		if err != nil {
			logError(err)
			return
		}
		for _, t := range tl {
			name, err := t.Name()
			if err != nil {
				logError(err)
				return
			}
			tracksList.AddItem(name, "", 'x', trackHandler(t))
		}

		application.SetFocus(tracksList)
		application.Sync()
	}
}

func trackHandler(t intf.Track) func() {
	return func() {
		if assignedList == nil {
			return
		}

		currentTrack = t

		allPlaylists, err := itunes.Playlists()
		if err != nil {
			logError(err)
			return
		}

		assignedList.Clear()
		assigneds, err := t.AssignedPlaylists()
		if err != nil {
			logError(err)
			return
		}

		for i, ap := range allPlaylists {
			aname, err := ap.Name()
			if err != nil {
				logError(err)
				return
			}

			matched := false
			for _, p := range assigneds {
				name, err := p.Name()
				if err != nil {
					logError(err)
					return
				}
				if aname == name {
					matched = true
					break
				}
			}
			assignedList.AddItem(aname, "", matchedString(matched), assignedHandler(i, matched, t, ap))
		}

		application.SetFocus(assignedList)
		//application.Sync()
	}
}

func matchedString(b bool) rune {
	if b {
		return '*'
	}
	return '-'
}

func assignedHandler(index int, assigned bool, t intf.Track, ap intf.Playlist) func() {
	return func() {
		if assigned {
			if err := ap.Delete(t); err != nil {
				logText.SetText(fmt.Sprintf("%s, %#v\n", logText.GetText(true), err))
			}
		} else {
			if err := ap.Add(t); err != nil {
				logText.SetText(fmt.Sprintf("%s, %#v\n", logText.GetText(true), err))
			}
		}
		trackHandler(t)()
		assignedList.SetCurrentItem(index)
	}
}

func logError(err error) {
	logText.SetText(fmt.Sprintf("%s, %#v\n", logText.GetText(true), err))
}

func logString(str string) {
	logText.SetText(fmt.Sprintf("%s, %s\n", logText.GetText(true), str))
}

func main() {
	tmp, err := itunesman.NewITunes()
	if err != nil {
		panic(err)
	}
	itunes = tmp

	pls, err := itunes.Playlists()
	if err != nil {
		panic(err)
	}

	application = tview.NewApplication()
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	playlists := tview.NewList()
	playlists.ShowSecondaryText(false)
	tracksList = tview.NewList()
	tracksList.ShowSecondaryText(false)
	tracksList.SetDoneFunc(func() {
		application.SetFocus(playlists)
	})
	assignedList = tview.NewList()
	assignedList.ShowSecondaryText(false)
	assignedList.SetDoneFunc(func() {
		application.SetFocus(tracksList)
	})
	logText = tview.NewTextView()

	for _, l := range pls {
		name, err := l.Name()
		if err != nil {
			logError(err)
		}
		fmt.Println(name)
		playlists.AddItem(name, "", '_', playlistHandler(l))
	}

	playlistSearch := tview.NewInputField().SetLabel(">")
	trackSearch := tview.NewInputField().SetLabel(">")
	assignedPlaylistSearch := tview.NewInputField().SetLabel(">")

	grid := tview.NewGrid().
		SetRows(3, 1, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(logText, 3, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(playlistSearch, 0, 0, 0, 0, 0, 0, false).
		AddItem(playlists, 0, 0, 0, 0, 0, 0, false).
		AddItem(trackSearch, 1, 0, 1, 3, 0, 0, true).
		AddItem(tracksList, 2, 0, 1, 3, 0, 0, true).
		AddItem(assignedPlaylistSearch, 0, 0, 0, 0, 0, 0, false).
		AddItem(assignedList, 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(playlistSearch, 1, 0, 1, 1, 0, 100, false).
		AddItem(playlists, 2, 0, 1, 1, 0, 100, false).
		AddItem(trackSearch, 1, 1, 1, 1, 0, 0, false).
		AddItem(tracksList, 2, 1, 1, 1, 0, 100, false).
		AddItem(assignedPlaylistSearch, 1, 2, 1, 1, 0, 100, false).
		AddItem(assignedList, 2, 2, 1, 1, 0, 100, false)

	application.SetRoot(grid, true).EnableMouse(true)
	time.Sleep(3)
	application.SetFocus(playlists)

	application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// 入力できるコンポーネントでは操作しない
		p := application.GetFocus()
		switch p.(type) {
		case *tview.InputField:
			return event
		}

		// TODO: フォーカス移動ここで実装する？

		k := event.Key()
		switch k {
		case tcell.KeyF5:
			if currentTrack == nil {
				return event
			}
			currentTrack.Play()
			return nil
		}

		return event
	})

	if err := application.Run(); err != nil {
		panic(err)
	}

}
