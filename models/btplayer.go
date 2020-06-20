package models

type BTPlayer struct {
	ID   string `json:"-" db:"tw_id"`
	Name string `json:"Username" db:"last_known_username"`

	Score int `json:"score"`
	Game *Game `json:"-"`

	TitleFoundAt int64 `json:"TitleFoundAt"`
	TitleFound  bool `json:"ArtistTitle"`

	ArtistFoundAt int64 `json:"TitleFoundAt"`
	ArtistFound bool `json:"ArtistFound"`
}

func (p *BTPlayer) SetFoundTitle() {
	p.TitleFound = true
	p.TitleFoundAt = p.Game.TimeLeft
}

func (p *BTPlayer) SetFoundArtist() {
	p.ArtistFound = true
	p.ArtistFoundAt = p.Game.TimeLeft
}
