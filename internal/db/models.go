package db

// import (
// 	"github.com/jackc/pgx/v5/pgtype"
// )
//
// type Entity struct {
// 	ID             int64
// 	FromDate       pgtype.Timestamptz
// 	ToDate         pgtype.Timestamptz
// 	IsActive       bool
// 	Name           string
// 	KnownNames     []string
// 	Type           string
// 	Symbol         pgtype.Text
// 	Metadata       []byte
// 	IsPartnership  bool
// 	NormalizedName string
// 	ParentID       pgtype.Int8
// }
//
// type EntityPartnership struct {
// 	ID       int64
// 	IsActive bool
// 	PartID   int64
// 	TargetID int64
// }
//
// type Flag struct {
// 	ID             int64
// 	CreationDate   pgtype.Timestamptz
// 	Name           string
// 	Tokens         []string
// 	Verified       bool
// 	QualifiesForID pgtype.Int8
// }
//
// type League struct {
// 	ID       int64
// 	FromDate pgtype.Timestamptz
// 	ToDate   pgtype.Timestamptz
// 	IsActive bool
// 	Name     string
// 	Symbol   string
// 	Gender   pgtype.Text
// 	ParentID pgtype.Int8
// }
//
// type Participant struct {
// 	ID       int64
// 	ClubName pgtype.Text
// 	Distance pgtype.Int4
// 	Laps     []pgtype.Time
// 	Lane     pgtype.Int2
// 	Series   pgtype.Int2
// 	Gender   string
// 	Category string
// 	ClubID   int64
// 	RaceID   int64
// 	Handicap pgtype.Time
// }
//
// type Penalty struct {
// 	ID               int64
// 	Penalty          int32
// 	Disqualification bool
// 	Reason           pgtype.Text
// 	ParticipantID    int64
// }
//
// type Race struct {
// 	ID                  int64
// 	CreationDate        pgtype.Timestamptz
// 	Laps                pgtype.Int2
// 	Lanes               pgtype.Int2
// 	Type                string
// 	Date                pgtype.Date
// 	Day                 int16
// 	Cancelled           bool
// 	CancellationReasons []string
// 	RaceName            pgtype.Text
// 	Sponsor             pgtype.Text
// 	TrophyEdition       pgtype.Int2
// 	FlagEdition         pgtype.Int2
// 	Modality            string
// 	Metadata            []byte
// 	FlagID              pgtype.Int8
// 	LeagueID            pgtype.Int8
// 	OrganizerID         pgtype.Int8
// 	TrophyID            pgtype.Int8
// 	AssociatedID        pgtype.Int8
// 	Gender              string
// }
//
// type Trophy struct {
// 	ID             int64
// 	CreationDate   pgtype.Timestamptz
// 	Name           string
// 	Tokens         []string
// 	Verified       bool
// 	QualifiesForID pgtype.Int8
// }
