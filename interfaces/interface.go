package interfaces

type ITunes interface {
	Playlists() ([]Playlist, error)
	Library() (Playlist, error) // as main playlist, contains all tracks
	Play() error
	Stop() error
	CurrentTrack() (Track, error)
	CurrentPlaylist() (Playlist, error)
}

type Track interface {
	TrackNumber() (int, error)
	Name() (string, error)
	AssignedPlaylists() ([]Playlist, error)
	Play() error
	Stop() error
	Delete() error
}

type Playlist interface {
	ID() (int, error)
	Index() (int, error)
	Name() (string, error)
	Tracks() ([]Track, error)
	Add(Track) error
	Delete(Track) error
	Show() error
}
