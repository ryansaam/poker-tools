package model

import (
	"time"
)

type PlayerStatus string

const (
	PlayerStatusUnknown PlayerStatus = ""
	PlayerSittingIn     PlayerStatus = "sitting_in"
	PlayerSittingOut    PlayerStatus = "sitting_out"
)

type StreetKind string

const (
	StreetPreflop StreetKind = "preflop"
	StreetFlop    StreetKind = "flop"
	StreetTurn    StreetKind = "turn"
	StreetRiver   StreetKind = "river"
)

type ActionKind string

const (
	ActionPostBlind    ActionKind = "post"
	ActionCheck        ActionKind = "check"
	ActionBet          ActionKind = "bet"
	ActionCall         ActionKind = "call"
	ActionRaiseTo      ActionKind = "raise_to"
	ActionFold         ActionKind = "fold"
	ActionMuck         ActionKind = "muck"
	ActionShow         ActionKind = "show"
	ActionUncalledBack ActionKind = "uncalled_return"
	ActionCollect      ActionKind = "collect"
	ActionUnknown      ActionKind = "unknown"
)

type Card struct {
	Rank string `json:"rank"` // "2"-"9","T","J","Q","K","A"
	Suit string `json:"suit"` // "c","d","h","s"
	// Raw is the exact glyph form like "6♣" for debugging/provenance
	Raw string `json:"raw,omitempty"`
}

type Player struct {
	Seat   int          `json:"seat"`
	Name   string       `json:"name"`
	Chips  int64        `json:"chips"`
	Bot    bool         `json:"bot"`              // had [B] tag in Seat line
	Empty  bool         `json:"empty"`            // seat empty
	Status PlayerStatus `json:"status,omitempty"` // sitting_in | sitting_out | ""
	IsHero bool         `json:"is_hero,omitempty"`
}

type Action struct {
	Player string     `json:"player"`
	Kind   ActionKind `json:"kind"`
	// Amount is the chip amount added/attributed in this atomic action (bets, calls, posts).
	Amount int64 `json:"amount,omitempty"`
	// To is used for "raises to X".
	To int64 `json:"to,omitempty"`
	// Info can carry hand text for "shows" or pot name for collects, etc.
	Info string `json:"info,omitempty"`
}

type Street struct {
	Kind    StreetKind `json:"kind"`
	Board   []Card     `json:"board,omitempty"` // board as of this street
	Actions []Action   `json:"actions,omitempty"`
}

type Stakes struct {
	SB int64 `json:"sb"`
	BB int64 `json:"bb"`
}

type Hand struct {
	ID         string    `json:"id"`
	PlayedAt   time.Time `json:"played_at"`
	Stakes     Stakes    `json:"stakes,omitempty"` // derived from first "posts 1/2" seen
	ButtonSeat int       `json:"button_seat"`
	Players    []Player  `json:"players"`
	Hero       string    `json:"hero,omitempty"`
	HoleCards  []Card    `json:"hole_cards,omitempty"`
	Streets    []Street  `json:"streets"`
	Board      []Card    `json:"board,omitempty"` // final board from Summary/river
	Meta       Meta      `json:"meta,omitempty"`
}

type Meta struct {
	SourceFile string     `json:"source_file,omitempty"`
	Source     SourceMeta `json:"source,omitempty"`
	Ingest     IngestMeta `json:"ingest,omitempty"`
}

type SourceMeta struct {
	Site          string `json:"site,omitempty"`           // "cpokers", "pokerstars", ...
	Format        string `json:"format,omitempty"`         // "js", "txt", ...
	FormatVersion string `json:"format_version,omitempty"` // "2", "classic", ...
	RawHandID     string `json:"raw_hand_id,omitempty"`    // provider's native hand id
	TableName     string `json:"table_name,omitempty"`     // optional
	TournamentID  string `json:"tournament_id,omitempty"`  // optional
}

type IngestMeta struct {
	ParserName    string    `json:"parser_name,omitempty"` // e.g., "cpokersjs2"
	ParserVersion string    `json:"parser_version,omitempty"`
	IngestedAt    time.Time `json:"ingested_at,omitempty"`
}
