package cpokersjs2

import (
	"regexp"

	"github.com/ryansaam/poker-tools/internal/model"
)

func (p *parser) parseStreet(kind model.StreetKind, header *regexp.Regexp, lineReader *LineReader, hand *model.Hand /*, pt *PotTracker*/) error {
	// Consume header line and parse cards (e.g., "Flop:9♦2♠6♦" or "Turn:5♥")
	line, ok := lineReader.Next()
	if !ok {
		return nil
	}
	headerMatch := header.FindStringSubmatch(line)
	if headerMatch == nil {
		lineReader.Unread(line)
		return nil
	}

	// Update board and start street
	newCards := parseCardRun(headerMatch[1])
	switch kind {
	case model.StreetFlop:
		hand.Board = append([]model.Card{}, newCards...)
	case model.StreetTurn, model.StreetRiver:
		hand.Board = append(hand.Board, newCards...)
	}
	street := model.Street{Kind: kind, Board: append([]model.Card{}, hand.Board...)}

	// Actions until next header / Summary / blank
	for {
		line, ok := lineReader.Next()
		if !ok || line == "" || p.regxFlopHeader.MatchString(line) || p.regxTurnHeader.MatchString(line) || p.regxRiverHeader.MatchString(line) || p.regxSummaryHeader.MatchString(line) {
			if ok && line != "" {
				lineReader.Unread(line)
			}
			if len(street.Actions) > 0 {
				hand.Streets = append(hand.Streets, street)
			}
			return nil
		}
		// match actions (check/bet/call/raise/fold/muck/show/uncalled/collects)
		if match := p.regxCheck.FindStringSubmatch(line); match != nil {
			street.Actions = append(street.Actions, model.Action{Player: cleanName(match[1]), Kind: model.ActionCheck})
			continue
		}
		if match := p.regxBet.FindStringSubmatch(line); match != nil {
			street.Actions = append(street.Actions, model.Action{Player: cleanName(match[1]), Kind: model.ActionBet, Amount: parseInt64(match[2])})
			continue
		}
		if match := p.regxCall.FindStringSubmatch(line); match != nil {
			street.Actions = append(street.Actions, model.Action{Player: cleanName(match[1]), Kind: model.ActionCall, Amount: parseInt64(match[2])})
			continue
		}
		if match := p.regxRaiseTo.FindStringSubmatch(line); match != nil {
			street.Actions = append(street.Actions, model.Action{Player: cleanName(match[1]), Kind: model.ActionRaiseTo, To: parseInt64(match[2])})
			continue
		}
		if match := p.regxFold.FindStringSubmatch(line); match != nil {
			street.Actions = append(street.Actions, model.Action{Player: cleanName(match[1]), Kind: model.ActionFold})
			continue
		}
		if match := p.regxMuck.FindStringSubmatch(line); match != nil {
			street.Actions = append(street.Actions, model.Action{Player: cleanName(match[1]), Kind: model.ActionMuck})
			continue
		}
		if match := p.regxShow.FindStringSubmatch(line); match != nil {
			street.Actions = append(street.Actions, model.Action{Player: cleanName(match[1]), Kind: model.ActionShow, Info: match[2]})
			continue
		}
		if match := p.regxUncalled.FindStringSubmatch(line); match != nil {
			street.Actions = append(street.Actions, model.Action{Player: cleanName(match[2]), Kind: model.ActionUncalledBack, Amount: parseInt64(match[1])})
			continue
		}
		if match := p.regxCollects.FindStringSubmatch(line); match != nil {
			street.Actions = append(street.Actions, model.Action{Player: cleanName(match[1]), Kind: model.ActionCollect, Amount: parseInt64(match[2]), Info: match[3]})
			continue
		}

		// If nothing matched, push back and let next phase handle (e.g., Summary)
		lineReader.Unread(line)
		if len(street.Actions) > 0 {
			hand.Streets = append(hand.Streets, street)
		}
		return nil
	}
}
