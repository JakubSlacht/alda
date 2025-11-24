package repl

import (
	"net"
	"strings"
	"testing"

	bencode "github.com/jackpal/bencode-go"
	"github.com/rs/zerolog"

	"alda.io/client/model"
)

func TestGenerateIdFormat(t *testing.T) {
	id := generateId()
	if len(id) != 3 {
		t.Fatalf("expected 3-char id, got %q", id)
	}
	for _, r := range id {
		if r < 'a' || r > 'z' {
			t.Fatalf("id contains invalid char: %q", r)
		}
	}
}

func TestResetStateClearsFields(t *testing.T) {
	s := &Server{
		input:        "abc",
		score:        model.NewScore(),
		eventIndex:   123,
		requestQueue: make(chan nREPLRequest),
	}
	err := s.resetState()
	if err != nil {
		t.Fatalf("resetState returned err: %v", err)
	}
	if s.input != "" {
		t.Errorf("expected cleared input, got %q", s.input)
	}
	if s.eventIndex != 0 {
		t.Errorf("expected eventIndex=0, got %d", s.eventIndex)
	}
	if s.score == nil {
		t.Error("expected new Score, got nil")
	}
}

func TestStateFilePath(t *testing.T) {
	s := &Server{id: "abc123"}
	got := s.stateFile()
	if !strings.Contains(got, "repl-servers/abc123.json") {
		t.Errorf("unexpected path: %s", got)
	}
}

func TestNewServerInitializes(t *testing.T) {
	s := NewServer(7777)
	if s.id == "" {
		t.Error("id not generated")
	}
	if s.Port != 7777 {
		t.Errorf("expected port 7777, got %d", s.Port)
	}
	if s.score == nil {
		t.Error("expected initialized score")
	}
	if s.requestQueue == nil {
		t.Error("expected requestQueue")
	}
}

func TestRespondWritesBencode(t *testing.T) {
	s := &Server{}
	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	req := nREPLRequest{conn: c1, msg: map[string]interface{}{"id": "1"}}

	go s.respondDone(req, map[string]interface{}{"foo": "bar"})

	v, err := bencode.Decode(c2)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	decoded, ok := v.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map[string]interface{}, got %T", v)
	}

	status, ok := decoded["status"].([]interface{})
	if !ok || len(status) != 1 || status[0] != "done" {
		t.Errorf("expected status=[\"done\"], got %#v", decoded["status"])
	}

	if decoded["id"] != "1" {
		t.Errorf("expected id=\"1\", got %#v", decoded["id"])
	}

	if decoded["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %#v", decoded["foo"])
	}
}

func TestUpdateScoreWithInputReturnsOptions(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	s := NewServer(0)
	opts, err := s.updateScoreWithInput("piano: c")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(opts) < 2 {
		t.Fatalf("expected at least 2 transmit options, got %d", len(opts))
	}
	if !strings.Contains(s.input, "piano") {
		t.Errorf("expected input to contain 'piano', got %q", s.input)
	}
}
