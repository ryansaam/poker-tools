package cpokersjs2

import (
	"strconv"
	"strings"

	"github.com/ryansaam/poker-tools/internal/model"
)

func parseInt(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
func parseInt64(s string) int64 {
	n, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	return n
}
func parsePlayerStatus(m []string, idx int) model.PlayerStatus {
	if idx < len(m) {
		switch strings.TrimSpace(m[idx]) {
		case "Sitting In":
			return model.PlayerSittingIn
		case "Sitting Out":
			return model.PlayerSittingOut
		}
	}
	return model.PlayerStatusUnknown
}

func newEmptySeatArray() [9]model.Player {
	var arr [9]model.Player
	for i := 0; i < 9; i++ {
		arr[i] = model.Player{Seat: i, Empty: true}
	}
	return arr
}

// cleanName strips the optional "[B] " bot prefix and trims spaces.
func cleanName(str string) string {
	str = strings.TrimSpace(str)
	if strings.HasPrefix(str, "[B] ") {
		return strings.TrimSpace(str[4:])
	}
	return str
}

// isPreflopStart reports whether a line marks the beginning of preflop content.
// We treat any of: literal "Preflop", an early "Dealt ..." line, or a "posts ..."
// blind line as the start boundary. This uses the parser's compiled regexes. :contentReference[oaicite:1]{index=1}
func (p *parser) isPreflopStart(line string) bool {
	return p.regxPreflopHeader.MatchString(line) ||
		p.regxDealt.MatchString(line) ||
		p.regxPosts.MatchString(line)
}
