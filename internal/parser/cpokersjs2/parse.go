package cpokersjs2

import (
	"bufio"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/ryansaam/poker-tools/internal/model"
)

/*
preflop is out of order
i see something wrong with hand starting on line 749
we need a pot tracker (maybe because pot might be able to be derived from the JSON)
*/

// ParseCPokersJS2File opens a CPokers "Log Version js2" text file and parses all hands.
func ParseFile(path string) ([]model.Hand, error) {
	fileHandle, fileOpenErr := os.Open(path)
	if fileOpenErr != nil {
		return nil, fileOpenErr
	}
	defer fileHandle.Close()

	parsedHands, parseErr := Parse(fileHandle)
	if parseErr != nil {
		return nil, parseErr
	}

	for i := range parsedHands {
		parsedHands[i].Meta.SourceFile = path
	}
	return parsedHands, nil
}

func Parse(reader io.Reader) ([]model.Hand, error) {
	lineReader := NewLineReader(bufio.NewScanner(reader))
	p := newJS2Parser()
	var hands []model.Hand
	for {
		h, done, err := p.parseOneHand(lineReader)
		if err != nil {
			return nil, err
		}
		if done {
			break
		}
		hands = append(hands, h)
	}
	return hands, nil
}

func (p *parser) parseOneHand(lineReader *LineReader) (model.Hand, bool, error) {
	var hand model.Hand
	seats := newEmptySeatArray()

	for {
		line, ok := lineReader.Next()
		if !ok {
			return model.Hand{}, true, nil
		}
		if line == "" {
			continue
		}
		if match := p.regxHandStart.FindStringSubmatch(line); match != nil {
			hand.ID = match[1]
			// Use our own canonical UUID for top-level id:
			hand.ID = uuid.NewString()
			// Preserve the room's native hand id in meta.source:
			hand.Meta.Source = model.SourceMeta{
				Site:          "cpokers",
				Format:        "js",
				FormatVersion: "2",
				RawHandID:     match[1],
			}
			// Record ingest info (optional but handy)
			hand.Meta.Ingest = model.IngestMeta{ParserName: "cpokersjs2", ParserVersion: "0.1.0", IngestedAt: time.Now().UTC()}
			break
		}
	}
	// Allocate preflop street up front so header can record early posts/Dealt.
	preflopStreet := model.Street{Kind: model.StreetPreflop}

	if err := p.parseHeader(lineReader, &hand, &seats, &preflopStreet); err != nil {
		return model.Hand{}, false, err
	}
	// Continue preflop parsing after the "Preflop" marker.
	if err := p.parsePreflop(lineReader, &hand, &preflopStreet); err != nil {
		return model.Hand{}, false, err
	}

	// Streets in order (they're optional)
	if err := p.parseStreet(model.StreetFlop, p.regxFlopHeader, lineReader, &hand); err != nil {
		return model.Hand{}, false, err
	}
	if err := p.parseStreet(model.StreetTurn, p.regxTurnHeader, lineReader, &hand); err != nil {
		return model.Hand{}, false, err
	}
	if err := p.parseStreet(model.StreetRiver, p.regxRiverHeader, lineReader, &hand); err != nil {
		return model.Hand{}, false, err
	}
	// // After streets, parse Summary (optional)
	// if err := p.parseSummary(lineReader, &hand); err != nil {
	// 	return model.Hand{}, false, err
	// }

	hand.Players = make([]model.Player, 9)
	copy(hand.Players, seats[:])
	// Append preflop if we captured anything in header or preflop.
	if len(preflopStreet.Actions) > 0 {
		hand.Streets = append(hand.Streets, preflopStreet)
	}
	return hand, false, nil
}
