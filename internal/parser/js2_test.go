package parser

// import (
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/ryansaam/poker-tools/internal/model"
// )

// const sampleHeader = `
// CPokers Hand #23594389
// Button is in Seat 3
// Played at 2025-08-15T07:10:06.360Z
// Log Version js2
// Seat 0: [B] Bean (Chips: 267)
// Seat 1: [B] Rudick (Chips: 200) (Sitting In)
// Seat 2: [B] Doherty (Chips: 413)
// Seat 3: [B] Murphy (Chips: 338)
// Seat 4: [B] Reed (Chips: 339)
// Seat 5: Guest6340 (Chips: 198)
// Seat 6: ryansamu5 (Chips: 200) (Sitting In)
// Seat 7: empty
// Seat 8: empty
// `

// func TestParseCPokersJS2(t *testing.T) {
// 	t.Parallel()

// 	hands, err := ParseCPokersJS2(strings.NewReader(sampleHeader))
// 	if err != nil {
// 		t.Fatalf("ParseCPokersJS2 error: %v", err)
// 	}
// 	if len(hands) != 1 {
// 		t.Fatalf("expected 1 hand, got %d", len(hands))
// 	}
// 	h := hands[0]

// 	// Hand header fields
// 	if got, want := h.ID, "23594389"; got != want {
// 		t.Errorf("ID = %q, want %q", got, want)
// 	}
// 	if got, want := h.ButtonSeat, 3; got != want {
// 		t.Errorf("ButtonSeat = %d, want %d", got, want)
// 	}
// 	if got, want := h.LogVersion, "js2"; got != want {
// 		t.Errorf("LogVersion = %q, want %q", got, want)
// 	}
// 	wantTime, _ := time.Parse(time.RFC3339Nano, "2025-08-15T07:10:06.360Z")
// 	if !h.PlayedAt.Equal(wantTime) {
// 		t.Errorf("PlayedAt = %s, want %s", h.PlayedAt.Format(time.RFC3339Nano), wantTime.Format(time.RFC3339Nano))
// 	}

// 	// Seats: exactly 9
// 	if len(h.Players) != 9 {
// 		t.Fatalf("expected 9 players, got %d", len(h.Players))
// 	}

// 	// Seat-by-seat checks
// 	checkSeat := func(seat int, name string, chips int64, bot, empty bool, status model.PlayerStatus) {
// 		p := h.Players[seat]
// 		if p.Seat != seat {
// 			t.Errorf("seat[%d].Seat = %d, want %d", seat, p.Seat, seat)
// 		}
// 		if p.Name != name {
// 			t.Errorf("seat[%d].Name = %q, want %q", seat, p.Name, name)
// 		}
// 		if p.Chips != chips {
// 			t.Errorf("seat[%d].Chips = %d, want %d", seat, p.Chips, chips)
// 		}
// 		if p.Bot != bot {
// 			t.Errorf("seat[%d].Bot = %v, want %v", seat, p.Bot, bot)
// 		}
// 		if p.Empty != empty {
// 			t.Errorf("seat[%d].Empty = %v, want %v", seat, p.Empty, empty)
// 		}
// 		if p.Status != status {
// 			t.Errorf("seat[%d].Status = %q, want %q", seat, p.Status, status)
// 		}
// 	}

// 	checkSeat(0, "Bean", 267, true, false, model.PlayerSittingOut) // no status suffix -> unknown; change to Unknown if you prefer
// 	checkSeat(1, "Rudick", 200, true, false, model.PlayerSittingIn)
// 	checkSeat(2, "Doherty", 413, true, false, model.PlayerSittingOut)
// 	checkSeat(3, "Murphy", 338, true, false, model.PlayerSittingOut)
// 	checkSeat(4, "Reed", 339, true, false, model.PlayerSittingOut)
// 	checkSeat(5, "Guest6340", 198, false, false, model.PlayerSittingOut)
// 	checkSeat(6, "ryansamu5", 200, false, false, model.PlayerSittingIn)
// 	checkSeat(7, "", 0, false, true, model.PlayerSittingOut)
// 	checkSeat(8, "", 0, false, true, model.PlayerSittingOut)
// }
