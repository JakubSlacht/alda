package repl

import (
	"testing"

	"alda.io/client/model"
	"alda.io/client/system"
	"alda.io/client/transmitter"
)

// --- Stubs ---

var stubPlayer = system.PlayerState{ID: "player1", Port: 1234}

// --- Tests ---

func TestHasPlayerAndUnsetPlayer(t *testing.T) {
	s := &Server{}
	if s.hasPlayer() {
		t.Fatal("Expected hasPlayer to be false for zero-value player")
	}

	s.player = stubPlayer
	if !s.hasPlayer() {
		t.Fatal("Expected hasPlayer to be true after setting player")
	}

	s.unsetPlayer()
	if s.hasPlayer() {
		t.Fatal("Expected hasPlayer to be false after unsetPlayer")
	}
}

func TestWithTransmitterNoPlayer(t *testing.T) {
	s := &Server{} // no player
	err := s.withTransmitter(func(t transmitter.OSCTransmitter) error { return nil })
	if err == nil {
		t.Fatal("Expected error when no player is available")
	}
}

func TestResetStateWithoutPlayer(t *testing.T) {
	s := &Server{
		score: model.NewScore(),
	}
	err := s.resetState()
	if err != nil {
		t.Fatal(err)
	}
	if s.input != "" || s.eventIndex != 0 || s.score == nil {
		t.Fatal("Expected server state to be reset")
	}
}

func TestUpdateScoreWithInput(t *testing.T) {
	s := &Server{
		score: model.NewScore(),
	}
	opts, err := s.updateScoreWithInput("piano: c d e f")
	if err != nil {
		t.Fatal(err)
	}
	if len(opts) == 0 {
		t.Fatal("Expected transmission options to be returned")
	}
	if s.input == "" {
		t.Fatal("Expected server input to be updated")
	}
	if s.eventIndex != len(s.score.Events) {
		t.Fatal("Expected eventIndex to be updated")
	}
}

func TestEvalAndPlayWithNoPlayer(t *testing.T) {
	s := &Server{} // no player
	err := s.evalAndPlay("piano: c d e f")
	if err == nil {
		t.Fatal("Expected error when no player is available")
	}
}

func TestShutdownPlayerNoPlayer(t *testing.T) {
	s := &Server{} // no player
	err := s.shutdownPlayer()
	if err == nil {
		t.Fatal("Expected error when no player is available")
	}
}

func TestWithTransmitterSuccess(t *testing.T) {
	// Prepare a server with a mock player
	s := &Server{player: stubPlayer}

	// Wrap the original transmitter constructor temporarily
	executeCalled := false

	err := s.withTransmitter(func(t transmitter.OSCTransmitter) error {
		executeCalled = true
		return nil
	})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if !executeCalled {
		t.Fatal("Expected execute function to be called")
	}
}

func BenchmarkUpdateScoreWithInput(b *testing.B) {
	for i := 0; i < 100000; i++ {
		s := &Server{
			score: model.NewScore(),
		}
		s.updateScoreWithInput("piano: c d e f c d c d c d o5 d c d c d c d c")
	}
}
