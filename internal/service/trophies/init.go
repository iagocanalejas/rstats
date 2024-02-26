package trophies

type Trophy struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Edition int64  `json:"edition,omitempty"`
}
