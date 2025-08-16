package cpokersjs2

import "github.com/ryansaam/poker-tools/internal/model"

// parsePreflop continues from the "Preflop" marker onward.
// NOTE: early "posts ..." and "Dealt ..." are handled in parseHeader.
func (p *parser) parsePreflop(lineReader *LineReader, hand *model.Hand, preflopStreet *model.Street) error {
	for {
		line, ok := lineReader.Next()
		if !ok || line == "" {
			return nil
		}
		if p.regxFlopHeader.MatchString(line) || p.regxTurnHeader.MatchString(line) || p.regxRiverHeader.MatchString(line) {
			lineReader.Unread(line)
			return nil
		}
		if p.regxPreflopHeader.MatchString(line) {
			continue
		}
		// Later: checks, bets, calls, raises, folds, shows, uncalled, collects...
	}
}
