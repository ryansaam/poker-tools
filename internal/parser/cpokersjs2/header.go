package cpokersjs2

import (
	"strings"
	"time"

	"github.com/ryansaam/poker-tools/internal/model"
)

// parseHeader reads the hand header & Seat lines, and ALSO consumes early
// "posts ..." and "Dealt ..." lines that precede the literal "Preflop" marker.
func (p *parser) parseHeader(lineReader *LineReader, hand *model.Hand, seats *[9]model.Player, preflopStreet *model.Street) error {
	for {
		line, ok := lineReader.Next()
		if !ok || line == "" {
			return nil
		}

		if match := p.regxButtonSeat.FindStringSubmatch(line); match != nil {
			hand.ButtonSeat = parseInt(match[1])
			continue
		}
		if match := p.regxPlayedAt.FindStringSubmatch(line); match != nil {
			if t, err := time.Parse(time.RFC3339Nano, match[1]); err == nil {
				hand.PlayedAt = t
			} else if t2, err2 := time.Parse(time.RFC3339, match[1]); err2 == nil {
				hand.PlayedAt = t2
			}
			continue
		}
		if match := p.regxLogVersion.FindStringSubmatch(line); match != nil {
			// Keep top-level JSON room-agnostic; store room format in meta.source
			// If parseOneHand already set Format/Version, only update version here.
			if hand.Meta.Source.Site == "" {
				hand.Meta.Source.Site = "cpokers"
				hand.Meta.Source.Format = "js"
			}
			hand.Meta.Source.FormatVersion = match[1]
			continue
		}

		if match := p.regxSeatEmpty.FindStringSubmatch(line); match != nil {
			if seat := parseInt(match[1]); seat >= 0 && seat <= 8 {
				(*seats)[seat] = model.Player{Seat: seat, Empty: true}
			}
			continue
		}
		if match := p.regxSeatBot.FindStringSubmatch(line); match != nil {
			if seat := parseInt(match[1]); seat >= 0 && seat <= 8 {
				(*seats)[seat] = model.Player{
					Seat: seat, Name: strings.TrimSpace(match[2]), Chips: parseInt64(match[3]),
					Bot: true, Empty: false, Status: parsePlayerStatus(match, 4),
				}
			}
			continue
		}
		if match := p.regxSeat.FindStringSubmatch(line); match != nil {
			if seat := parseInt(match[1]); seat >= 0 && seat <= 8 {
				(*seats)[seat] = model.Player{
					Seat: seat, Name: strings.TrimSpace(match[2]), Chips: parseInt64(match[3]),
					Bot: false, Empty: false, Status: parsePlayerStatus(match, 4),
				}
			}
			continue
		}

		// ---- Early preflop content before the "Preflop" marker ----
		// Posts: derive stakes (first two distinct amounts => SB, BB) and
		// append a post-blind action into the preflop street.
		if match := p.regxPosts.FindStringSubmatch(line); match != nil {
			playerName := cleanName(strings.TrimSpace(match[1]))
			amount := parseInt64(match[2])
			if hand.Stakes.SB == 0 {
				hand.Stakes.SB = amount
			} else if hand.Stakes.BB == 0 && amount != hand.Stakes.SB {
				hand.Stakes.BB = amount
			}
			preflopStreet.Actions = append(preflopStreet.Actions, model.Action{
				Player: playerName,
				Kind:   model.ActionPostBlind,
				Amount: amount,
			})
			continue
		}
		// Dealt: capture hero hole cards here (hero resolution can be added later).
		if match := p.regxDealt.FindStringSubmatch(line); match != nil {
			card1, err1 := parseCardGlyph(match[1])
			card2, err2 := parseCardGlyph(match[2])
			if err1 == nil && err2 == nil {
				hand.HoleCards = []model.Card{card1, card2}
			}
			continue
		}

		if p.isPreflopStart(line) || p.regxFlopHeader.MatchString(line) || p.regxTurnHeader.MatchString(line) || p.regxRiverHeader.MatchString(line) {
			lineReader.Unread(line)
			return nil
		}
		lineReader.Unread(line)
		return nil
	}
}
