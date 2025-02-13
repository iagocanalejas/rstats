package types

type BoatConfig struct {
	Weight            float64      `json:"weight"`
	Length            int          `json:"length"`
	Rowlocks          []int        `json:"rowlocks"`
	RowlockGap        int16        `json:"rowlock_gap"`
	Ribs              []int        `json:"ribs"`
	Seats             []SeatConfig `json:"seats"`
	BowFloatingLine   int          `json:"bow_floating_line"`
	SternFloatingLine int          `json:"stern_floating_line"`
}

type SeatConfig struct {
	Position        Position `json:"position"`
	Weight          float64  `json:"weight"`
	RowlockPosition int16    `json:"rowlock_position"`
	BenchDistance   int      `json:"bench_distance"`
	Side            *Side    `json:"side"`
}

type Side string

const (
	STARBOARD Side = "STARBOARD"
	PORT           = "PORT"
)

type Position string

const (
	COXSWAIN Position = "COXSWAIN"
	STROKE            = "STROKE"
	TWO               = "2"
	THREE             = "3"
	FOUR              = "4"
	FIVE              = "5"
	SIX               = "6"
	BOW               = "BOW"
)

var VALID_POSITIONS = []Position{COXSWAIN, STROKE, TWO, THREE, FOUR, FIVE, SIX, BOW}
