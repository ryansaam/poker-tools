package cpokersjs2

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ryansaam/poker-tools/internal/model"
)

// parseCardGlyph parses a single glyph like "9♦" or "A♠" into a Card.
// It treats the glyph as two runes (rank + suit), which is the correct way
// to handle Unicode symbols in Go. If the suit is unknown, it errors. :contentReference[oaicite:0]{index=0}
func parseCardGlyph(glyph string) (model.Card, error) {
	rs := []rune(strings.TrimSpace(glyph))
	if len(rs) != 2 {
		return model.Card{}, fmt.Errorf("bad card glyph: %q", glyph)
	}
	rank := strings.ToUpper(string(rs[0]))
	suit := suitFromGlyph(string(rs[1]))
	if suit == "?" {
		return model.Card{}, fmt.Errorf("unknown suit glyph: %q", glyph)
	}
	return model.Card{
		Rank: rank,
		Suit: suit,       // "c","d","h","s"
		Raw:  string(rs), // preserve original glyphs
	}, nil
}

func suitFromGlyph(g string) string {
	switch g {
	case "♣":
		return "c"
	case "♦":
		return "d"
	case "♥":
		return "h"
	case "♠":
		return "s"
	default:
		return "?"
	}
}

// parseCardRun splits compact runs like "9♦2♠6♦" or "5♥" into []Card.
// It walks the string as runes (rank + suit per card) to safely handle Unicode suits.
// Example inputs come from lines like "Flop:9♦2♠6♦" / "Turn:5♥".
// Any malformed pair is skipped.
// See: Go’s guidance on strings vs. runes for Unicode handling. :contentReference[oaicite:1]{index=1}
func parseCardRun(run string) []model.Card {
	run = strings.TrimSpace(run)
	if run == "" {
		return nil
	}
	var cards []model.Card
	runes := []rune(run)
	for i := 0; i+1 < len(runes); i += 2 {
		// Recompose the two-rune glyph (rank + suit) for reuse by parseCardGlyph.
		var buf bytes.Buffer
		buf.WriteRune(runes[i])
		buf.WriteRune(runes[i+1])
		cardGlyph := buf.String()
		if c, err := parseCardGlyph(cardGlyph); err == nil {
			cards = append(cards, c)
		}
		// Note: In CPokers logs, Ten is "T", so two runes per card holds.
		// If a site ever uses "10", you’d need a smarter tokenizer.
	}
	return cards
}
