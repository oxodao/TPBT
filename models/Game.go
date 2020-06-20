package models

type Game struct {

	Title    string `json:"Title"`
	Artist   string `json:"Artist"`
	Running  bool   `json:"Running"`
	IsTimed  bool   `json:"IsTimed"`
	TimeLeft int64    `json:"TimeLeft"`
	TimeSet  int64    `json:"TimeSet"`

	Players []*BTPlayer `json:"-"`

}
