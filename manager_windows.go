package itunesman

import (
	"fmt"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"

	intf "github.com/kuwa72/ituweak/interfaces"
)

var _ intf.ITunes = &ITunes{}
var _ intf.Playlist = &Playlist{}
var _ intf.Track = &Track{}

type ITunes struct {
	instance *ole.IDispatch
}

type Playlist struct {
	instance *ole.IDispatch
}

type Track struct {
	instance *ole.IDispatch
}

func NewITunes() (intf.ITunes, error) {
	if err := ole.CoInitializeEx(0, 0); err != nil {
		return nil, err
	}
	unknown, err := oleutil.CreateObject("iTunes.Application")
	if err != nil {
		return nil, err
	}
	itunes, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, err
	}
	return &ITunes{itunes}, nil
}

func (it *ITunes) Playlists() ([]intf.Playlist, error) {
	ls, err := oleutil.GetProperty(it.instance, "LibrarySource")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	pls, err := oleutil.GetProperty(ls.ToIDispatch(), "Playlists")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	ret := []intf.Playlist{}
	oleutil.ForEach(pls.ToIDispatch(), func(v *ole.VARIANT) error {
		ret = append(ret, &Playlist{v.ToIDispatch()})
		return nil
	})

	return ret, nil
}
func (it *ITunes) Library() (intf.Playlist, error) {
	l, err := oleutil.GetProperty(it.instance, "LibraryPlaylist")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &Playlist{l.ToIDispatch()}, nil
}
func (it *ITunes) Play() error {
	_, err := oleutil.CallMethod(it.instance, "Play")
	return err
}
func (it *ITunes) Stop() error {
	_, err := oleutil.CallMethod(it.instance, "Stop")
	return err
}
func (it *ITunes) CurrentTrack() (intf.Track, error) {
	return nil, nil
}
func (it *ITunes) CurrentPlaylist() (intf.Playlist, error) {
	return nil, nil
}

//------------

func (p *Playlist) ID() (int, error) {
	v, err := oleutil.GetProperty(p.instance, "PlaylistID")
	if err != nil {
		return 0, err
	}

	return int(v.Val), nil
}
func (p *Playlist) Tracks() ([]intf.Track, error) {
	v, err := oleutil.GetProperty(p.instance, "Tracks")
	if err != nil {
		return nil, err
	}

	ret := []intf.Track{}
	oleutil.ForEach(v.ToIDispatch(), func(v *ole.VARIANT) error {
		ret = append(ret, &Track{v.ToIDispatch()})
		return nil
	})

	return ret, nil
}
func (p *Playlist) Add(t intf.Track) error {
	_, err := oleutil.CallMethod(p.instance, "AddTrack", t.(*Track).instance)
	return err
}
func (p *Playlist) Delete(t intf.Track) error {
	tl, err := p.Tracks()
	if err != nil {
		return err
	}
	tNumber, err := t.TrackNumber()
	if err != nil {
		return err
	}
	for _, at := range tl {
		atNumber, err := at.TrackNumber()
		if err != nil {
			return err
		}
		if tNumber == atNumber {
			if err := at.Delete(); err != nil {
				return err
			}
			break
		}
	}
	return nil
}
func (p *Playlist) Show() error { return nil }
func (p *Playlist) Name() (string, error) {
	v, err := oleutil.GetProperty(p.instance, "Name")
	return (string)(v.ToString()), err
}
func (p *Playlist) Index() (int, error) {
	v, err := oleutil.GetProperty(p.instance, "Index")
	return (int)(v.Val), err
}

// -----------

func (t *Track) TrackNumber() (int, error) {
	v, err := oleutil.GetProperty(t.instance, "TrackNumber")
	return (int)(v.Val), err
}
func (t *Track) Name() (string, error) {
	v, err := oleutil.GetProperty(t.instance, "Name")
	return (string)(v.ToString()), err
}
func (t *Track) AssignedPlaylists() ([]intf.Playlist, error) {
	v, err := oleutil.GetProperty(t.instance, "Playlists")
	if err != nil {
		return nil, err
	}

	ret := []intf.Playlist{}
	oleutil.ForEach(v.ToIDispatch(), func(v *ole.VARIANT) error {
		ret = append(ret, &Playlist{v.ToIDispatch()})
		return nil
	})

	return ret, nil
}
func (t *Track) Play() error {
	_, err := oleutil.CallMethod(t.instance, "Play")
	return err
}
func (t *Track) Stop() error {
	_, err := oleutil.CallMethod(t.instance, "Stop")
	return err
}
func (t *Track) Delete() error {
	_, err := oleutil.CallMethod(t.instance, "Delete")
	return err
}
