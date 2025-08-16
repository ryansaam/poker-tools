package parser

// import (
// 	"bufio"
// 	"io"
// 	"os"
// 	"regexp"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/ryansaam/poker-tools/internal/model"
// )

// // ParseCPokersJS2File opens a CPokers "Log Version js2" text file and parses all hands.
// func ParseCPokersJS2File(path string) ([]model.Hand, error) {
// 	fileHandle, fileOpenErr := os.Open(path)
// 	if fileOpenErr != nil {
// 		return nil, fileOpenErr
// 	}
// 	defer fileHandle.Close()

// 	parsedHands, parseErr := ParseCPokersJS2(fileHandle)
// 	if parseErr != nil {
// 		return nil, parseErr
// 	}

// 	for i := range parsedHands {
// 		parsedHands[i].Meta.SourceFile = path
// 	}
// 	return parsedHands, nil
// }

// func ParseCPokersJS2(reader io.Reader) ([]model.Hand, error) {
// 	lineReader := NewLineReader(bufio.NewScanner(reader))
// 	p := NewJS2Parser()

// 	var hands []model.Hand
// 	for {
// 		hand, done, err := p.parseOneHand(lineReader)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if done {
// 			break
// 		}
// 		hands = append(hands, hand)
// 	}
// 	return hands, nil
// }

// type Parser struct {
// 	// header
// 	regxHandStart  *regexp.Regexp
// 	regxButtonSeat *regexp.Regexp
// 	regxPlayedAt   *regexp.Regexp
// 	regxLogVersion *regexp.Regexp
// 	regxSeatEmpty  *regexp.Regexp
// 	regxSeatBot    *regexp.Regexp
// 	regxSeat       *regexp.Regexp

// 	// preflop/actions (you’ll add more as you go)
// 	regxPosts         *regexp.Regexp
// 	regxDealt         *regexp.Regexp
// 	regxPreflopHeader *regexp.Regexp

// 	// street headers
// 	regxFlopHeader  *regexp.Regexp
// 	regxTurnHeader  *regexp.Regexp
// 	regxRiverHeader *regexp.Regexp
// }

// func NewJS2Parser() *Parser {
// 	return &Parser{
// 		regxHandStart:  regexp.MustCompile(`^CPokers Hand #(\d+)`),
// 		regxButtonSeat: regexp.MustCompile(`^Button is in Seat (\d+)`),
// 		regxPlayedAt:   regexp.MustCompile(`^Played at ([0-9T:\.\-+Z]+)$`),
// 		regxLogVersion: regexp.MustCompile(`^Log Version (\S+)`),

// 		// Seats: capture optional "(Sitting In|Sitting Out)" in group 4
// 		regxSeatEmpty: regexp.MustCompile(`^Seat (\d+): empty$`),
// 		regxSeatBot:   regexp.MustCompile(`^Seat (\d+): \[B\]\s+([^(]+)\s+\(Chips:\s+(\d+)\)(?:\s+\(([^)]+)\))?$`),
// 		regxSeat:      regexp.MustCompile(`^Seat (\d+):\s+(?:\[B\]\s+)?([^(]+)\s+\(Chips:\s+(\d+)\)(?:\s+\(([^)]+)\))?$`),

// 		// Early preflop
// 		regxPosts:         regexp.MustCompile(`^(?:\[B\]\s+)?([^:]+): posts (\d+)`),
// 		regxDealt:         regexp.MustCompile(`^Dealt ([2-9TJQKA][♣♦♥♠])([2-9TJQKA][♣♦♥♠])$`),
// 		regxPreflopHeader: regexp.MustCompile(`^Preflop$`),

// 		// Next streets (just for boundary detection right now)
// 		regxFlopHeader:  regexp.MustCompile(`^Flop:(.+)$`),
// 		regxTurnHeader:  regexp.MustCompile(`^Turn:(.+)$`),
// 		regxRiverHeader: regexp.MustCompile(`^River:(.+)$`),
// 	}
// }

// // parseOneHand reads from the current position to build 1 hand.
// // Returns (hand, done=true) at EOF with no more hands.
// func (p *Parser) parseOneHand(lineReader *LineReader) (model.Hand, bool, error) {
// 	var hand model.Hand
// 	seats := newEmptySeatArray()

// 	// 1) Seek to a hand start (skip blank lines)
// 	for {
// 		line, ok := lineReader.Next()
// 		if !ok {
// 			return model.Hand{}, true, nil // EOF, no more hands
// 		}
// 		if line == "" {
// 			continue
// 		}
// 		if m := p.regxHandStart.FindStringSubmatch(line); m != nil {
// 			hand.ID = m[1]
// 			break
// 		}
// 		// Ignore any preamble garbage; keep scanning
// 	}

// 	// 2) Header
// 	if err := p.parseHeader(lineReader, &hand, &seats); err != nil {
// 		return model.Hand{}, false, err
// 	}

// 	// We stop here for now if you want strictly header-only. Otherwise,
// 	// start consuming the preflop section (will return at street boundary).
// 	if err := p.parsePreflop(lineReader, &hand); err != nil {
// 		return model.Hand{}, false, err
// 	}

// 	// Optional: parse flop/turn/river using p.parseStreet(...)
// 	// For now we leave them for later.

// 	// finalize players (exactly 9 seats snapshot)
// 	hand.Players = make([]model.Player, 9)
// 	copy(hand.Players, seats[:])
// 	return hand, false, nil
// }

// func (p *Parser) parseHeader(lineReader *LineReader, hand *model.Hand, seats *[9]model.Player) error {
// 	for {
// 		line, ok := lineReader.Next()
// 		if !ok || line == "" {
// 			// Boundary: end of header block
// 			return nil
// 		}

// 		if m := p.regxButtonSeat.FindStringSubmatch(line); m != nil {
// 			hand.ButtonSeat = parseInt(m[1])
// 			continue
// 		}
// 		if m := p.regxPlayedAt.FindStringSubmatch(line); m != nil {
// 			// RFC3339 with fractions; fallback to plain RFC3339
// 			if t, err := time.Parse(time.RFC3339Nano, m[1]); err == nil {
// 				hand.PlayedAt = t
// 			} else if t2, err2 := time.Parse(time.RFC3339, m[1]); err2 == nil {
// 				hand.PlayedAt = t2
// 			}
// 			continue
// 		}
// 		if m := p.regxLogVersion.FindStringSubmatch(line); m != nil {
// 			hand.LogVersion = m[1]
// 			continue
// 		}

// 		// Seats
// 		if m := p.regxSeatEmpty.FindStringSubmatch(line); m != nil {
// 			if s := parseInt(m[1]); s >= 0 && s <= 8 {
// 				(*seats)[s] = model.Player{Seat: s, Empty: true}
// 			}
// 			continue
// 		}
// 		if m := p.regxSeatBot.FindStringSubmatch(line); m != nil {
// 			if s := parseInt(m[1]); s >= 0 && s <= 8 {
// 				(*seats)[s] = model.Player{
// 					Seat:   s,
// 					Name:   strings.TrimSpace(m[2]),
// 					Chips:  parseInt64(m[3]),
// 					Bot:    true,
// 					Empty:  false,
// 					Status: parsePlayerStatus(m, 4),
// 				}
// 			}
// 			continue
// 		}
// 		if m := p.regxSeat.FindStringSubmatch(line); m != nil {
// 			if s := parseInt(m[1]); s >= 0 && s <= 8 {
// 				(*seats)[s] = model.Player{
// 					Seat:   s,
// 					Name:   strings.TrimSpace(m[2]),
// 					Chips:  parseInt64(m[3]),
// 					Bot:    false,
// 					Empty:  false,
// 					Status: parsePlayerStatus(m, 4),
// 				}
// 			}
// 			continue
// 		}

// 		// First non-header token → push back for next phase (preflop)
// 		if p.isPreflopStart(line) {
// 			lineReader.Unread(line)
// 			return nil
// 		}
// 		// Also guard for street headers if the site ever omits "Preflop"
// 		if p.regxFlopHeader.MatchString(line) || p.regxTurnHeader.MatchString(line) || p.regxRiverHeader.MatchString(line) {
// 			lineReader.Unread(line)
// 			return nil
// 		}

// 		// Unknown header-adjacent lines? Push back and let next phase decide.
// 		lineReader.Unread(line)
// 		return nil
// 	}
// }

// func (p *Parser) parsePreflop(lineReader *LineReader, hand *model.Hand) error {
// 	// Optional: create street if you already have model.Street
// 	// st := model.Street{Kind: model.StreetPreflop}
// 	// defer func() { if len(st.Actions) > 0 { hand.Streets = append(hand.Streets, st) } }()

// 	for {
// 		line, ok := lineReader.Next()
// 		if !ok || line == "" {
// 			// boundary (blank line or EOF)
// 			return nil
// 		}
// 		// stop if we hit another section
// 		if p.regxFlopHeader.MatchString(line) || p.regxTurnHeader.MatchString(line) || p.regxRiverHeader.MatchString(line) {
// 			lineReader.Unread(line)
// 			return nil
// 		}

// 		// Record SB/BB posts now (useful for inferring stakes)
// 		if m := p.regxPosts.FindStringSubmatch(line); m != nil {
// 			// Example: name := cleanName(m[1]); amt := parseInt64(m[2])
// 			// TODO: push into st.Actions; derive hand.Stakes if needed
// 			continue
// 		}

// 		// Capture "Dealt XXYY" for the hero later (we’ll decide hero upstream)
// 		if m := p.regxDealt.FindStringSubmatch(line); m != nil {
// 			// TODO: if hero is known, set hand.HoleCards = those cards
// 			continue
// 		}

// 		// "Preflop" header line (some logs include it). Safe to ignore.
// 		if p.regxPreflopHeader.MatchString(line) {
// 			continue
// 		}

// 		// For now, ignore other actions until you wire Action regexes.
// 		// You can lr.Unread(line) and switch on more regex as you implement.
// 	}
// }

// // --- helpers ---

// func (p *Parser) isPreflopStart(line string) bool {
// 	return p.regxPreflopHeader.MatchString(line) ||
// 		p.regxDealt.MatchString(line) ||
// 		p.regxPosts.MatchString(line)
// }

// func parsePlayerStatus(submatches []string, idx int) model.PlayerStatus {
// 	if idx < len(submatches) {
// 		switch strings.TrimSpace(submatches[idx]) {
// 		case "Sitting In":
// 			return model.PlayerSittingIn
// 		case "Sitting Out":
// 			return model.PlayerSittingOut
// 		}
// 	}
// 	return model.PlayerStatusUnknown
// }

// func newEmptySeatArray() [9]model.Player {
// 	var arr [9]model.Player
// 	for i := 0; i < 9; i++ {
// 		arr[i] = model.Player{Seat: i, Empty: true}
// 	}
// 	return arr
// }

// func parseInt(s string) int {
// 	n, _ := strconv.Atoi(strings.TrimSpace(s))
// 	return n
// }
// func parseInt64(s string) int64 {
// 	n, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
// 	return n
// }

// type LineReader struct {
// 	scanner    *bufio.Scanner
// 	pushedBack *string
// }

// func NewLineReader(scanner *bufio.Scanner) *LineReader {
// 	// If you ever see "token too long", increase the buffer (bufio.Scanner docs).
// 	// scanner.Buffer(make([]byte, 0, 64<<10), 1<<20) // example tweak if needed.
// 	return &LineReader{scanner: scanner}
// }

// func (lineReader *LineReader) Next() (string, bool) {
// 	if lineReader.pushedBack != nil {
// 		line := *lineReader.pushedBack
// 		lineReader.pushedBack = nil
// 		return line, true
// 	}
// 	if !lineReader.scanner.Scan() {
// 		return "", false
// 	}
// 	return strings.TrimSpace(lineReader.scanner.Text()), true
// }

// func (lr *LineReader) Unread(line string) {
// 	l := line
// 	lr.pushedBack = &l
// }
