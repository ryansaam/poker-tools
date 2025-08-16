package cpokersjs2

import "regexp"

type parser struct {
	regxHandStart, regxButtonSeat, regxPlayedAt, regxLogVersion *regexp.Regexp
	regxSeatEmpty, regxSeatBot, regxSeat                        *regexp.Regexp
	regxPosts, regxDealt, regxPreflopHeader                     *regexp.Regexp
	regxFlopHeader, regxTurnHeader, regxRiverHeader             *regexp.Regexp
	regxCheck, regxBet, regxCall, regxRaiseTo                   *regexp.Regexp
	regxFold, regxMuck, regxShow                                *regexp.Regexp
	regxUncalled, regxCollects                                  *regexp.Regexp
	regxSummaryHeader                                           *regexp.Regexp
	regxTotalPotRake                                            *regexp.Regexp
	regxBoardSummary                                            *regexp.Regexp
}

func newJS2Parser() *parser {
	return &parser{
		regxHandStart:  regexp.MustCompile(`^CPokers Hand #(\d+)`),
		regxButtonSeat: regexp.MustCompile(`^Button is in Seat (\d+)`),
		regxPlayedAt:   regexp.MustCompile(`^Played at ([0-9T:\.\-+Z]+)$`),
		regxLogVersion: regexp.MustCompile(`^Log Version (\S+)`),

		regxSeatEmpty: regexp.MustCompile(`^Seat (\d+): empty$`),
		regxSeatBot:   regexp.MustCompile(`^Seat (\d+): \[B\]\s+([^(]+)\s+\(Chips:\s+(\d+)\)(?:\s+\(([^)]+)\))?$`),
		regxSeat:      regexp.MustCompile(`^Seat (\d+):\s+(?:\[B\]\s+)?([^(]+)\s+\(Chips:\s+(\d+)\)(?:\s+\(([^)]+)\))?$`),

		regxPosts:         regexp.MustCompile(`^(?:\[B\]\s+)?([^:]+): posts (\d+)`),
		regxDealt:         regexp.MustCompile(`^Dealt ([2-9TJQKA][♣♦♥♠])([2-9TJQKA][♣♦♥♠])$`),
		regxPreflopHeader: regexp.MustCompile(`^Preflop$`),

		regxFlopHeader:  regexp.MustCompile(`^Flop:(.+)$`),
		regxTurnHeader:  regexp.MustCompile(`^Turn:(.+)$`),
		regxRiverHeader: regexp.MustCompile(`^River:(.+)$`),

		regxCheck:   regexp.MustCompile(`^(?:\[B\]\s+)?([^:]+): checks$`),
		regxBet:     regexp.MustCompile(`^(?:\[B\]\s+)?([^:]+): bets (\d+)$`),
		regxCall:    regexp.MustCompile(`^(?:\[B\]\s+)?([^:]+): calls (\d+)$`),
		regxRaiseTo: regexp.MustCompile(`^(?:\[B\]\s+)?([^:]+): raises to (\d+)$`),

		regxFold: regexp.MustCompile(`^(?:\[B\]\s+)?([^:]+): folds$`),
		regxMuck: regexp.MustCompile(`^(?:\[B\]\s+)?([^:]+): ?mucks$`),
		regxShow: regexp.MustCompile(`^(?:\[B\]\s+)?([^:]+):shows ([^\s]+) \((.+)\)$`),

		regxUncalled: regexp.MustCompile(`^Uncalled bet of (\d+) returned to (.+)$`),
		regxCollects: regexp.MustCompile(`^(?:\[B\]\s+)?([^:]+): collects (\d+) from (?:the )?(Main pot|Side pot \d+)$`),

		regxSummaryHeader: regexp.MustCompile(`^Summary$`),

		regxTotalPotRake: regexp.MustCompile(`^Total pot (\d+) \| Rake (\d+)`),

		regxBoardSummary: regexp.MustCompile(`^Board (.+)$`),
	}
}
