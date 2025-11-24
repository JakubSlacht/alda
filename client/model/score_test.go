package model

import (
	"fmt"
	"strings"
	"testing"

	"alda.io/client/json"

	_ "alda.io/client/testing"
)

// ScoreUpdate stub
type stubUpdate struct {
	err error
}

func (m stubUpdate) UpdateScore(score *Score) error                  { return m.err }
func (m stubUpdate) DurationMs(part *Part) float64                   { return 0 }
func (m stubUpdate) VariableValue(score *Score) (ScoreUpdate, error) { return nil, nil }
func (m stubUpdate) GetSourceContext() AldaSourceContext             { return AldaSourceContext{"stub", 0, 0} }
func (m stubUpdate) JSON() *json.Container                           { return json.ToJSON("stub") }

func TestNewScoreInitializesProperly(t *testing.T) {
	score := NewScore()
	if score.Parts == nil || score.Aliases == nil || score.Markers == nil || score.Variables == nil || score.GlobalAttributes == nil {
		t.Error("Expected all score maps/slices to be initialized")
	}
}

func TestScoreUpdate_UpdatePaths(t *testing.T) {
	executeScoreUpdateTestCases(
		t,
		scoreUpdateTestCase{
			label: "update success",
			updates: []ScoreUpdate{
				stubUpdate{err: nil},
			},
			expectations: []scoreUpdateExpectation{
				func(s *Score) error {
					// score should be unchanged, just not nil
					if s == nil {
						return fmt.Errorf("score was nil after update")
					}
					return nil
				},
			},
		},
		scoreUpdateTestCase{
			label: "update fail",
			updates: []ScoreUpdate{
				stubUpdate{err: fmt.Errorf("bad update")},
			},
			errorExpectations: []scoreUpdateErrorExpectation{
				func(err error) error {
					if err == nil {
						return fmt.Errorf("expected an error but got nil")
					}
					got := err.Error()
					if got == "" {
						return fmt.Errorf("empty error received")
					}
					if !strings.Contains(got, "bad update") {
						return fmt.Errorf("expected 'bad update' in error, got: %s", got)
					}
					return nil
				},
			},
		},
	)
}

func TestPartOffsets_Basic(t *testing.T) {
	score := NewScore()
	part, err := score.NewPart("piano")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	part.CurrentOffset = 123.0
	score.Parts = []*Part{part}

	offsets := score.PartOffsets()
	got := offsets[part.origin]
	if !equalish(got, 123.0) {
		t.Errorf("expected offset ~123, got %f", got)
	}
}

func TestPartOffsets_SingleVoice(t *testing.T) {
	score := NewScore()
	part, err := score.NewPart("piano")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	score.Parts = []*Part{part}
	score.CurrentParts = []*Part{part}

	vm := VoiceMarker{AldaSourceContext{"stub", 0, 0}, 1}
	vm.UpdateScore(score)

	voice1 := part.GetVoice(1)
	voice1.CurrentOffset = 50.0

	offsets := score.PartOffsets()
	got := offsets[part.origin]
	if !equalish(got, 0.0) {
		t.Errorf("expected offset ~0, got %f. Voices length was %d and voiceChange was %d. Voice1 offset was %f", got, len(part.voices.voices), score.voiceChange, part.voices.voices[1].CurrentOffset)
	}
	offsets = score.PartOffsets()
	got = offsets[part.origin]
	if !equalish(got, 50.0) {
		t.Errorf("expected offset ~50, got %f. Voices length was %d and voiceChange was %d. Voice1 offset was %f", got, len(part.voices.voices), score.voiceChange, part.voices.voices[1].CurrentOffset)
	}

}

func TestPartOffsets_MultiVoice(t *testing.T) {
	score := NewScore()
	part, err := score.NewPart("piano")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	score.Parts = []*Part{part}
	score.CurrentParts = []*Part{part}

	vm := VoiceMarker{AldaSourceContext{"stub", 0, 0}, 1}
	vm.UpdateScore(score)

	voice1 := part.GetVoice(1)
	voice1.CurrentOffset = 50.0

	offsets := score.PartOffsets()
	got := offsets[part.origin]
	if !equalish(got, 0.0) {
		t.Errorf("expected offset ~0, got %f. Voices length was %d and voiceChange was %d. Voice1 offset was %f", got, len(part.voices.voices), score.voiceChange, part.voices.voices[1].CurrentOffset)
	}
	offsets = score.PartOffsets()
	got = offsets[part.origin]
	if !equalish(got, 50.0) {
		t.Errorf("expected offset ~50, got %f. Voices length was %d and voiceChange was %d. Voice1 offset was %f", got, len(part.voices.voices), score.voiceChange, part.voices.voices[1].CurrentOffset)
	}

	vm = VoiceMarker{AldaSourceContext{"stub", 0, 0}, 2}
	vm.UpdateScore(score)

	voice2 := part.GetVoice(2)
	voice2.CurrentOffset = 25.0

	offsets = score.PartOffsets()
	got = offsets[part.origin]
	if !equalish(got, 0.0) {
		t.Errorf("expected offset ~0, got %f. Voices length was %d and voiceChange was %d. Voice1 offset was %f", got, len(part.voices.voices), score.voiceChange, part.voices.voices[2].CurrentOffset)
	}
	offsets = score.PartOffsets()
	got = offsets[part.origin]
	if !equalish(got, 25.0) {
		t.Errorf("expected offset ~25, got %f. Voices length was %d and voiceChange was %d. Voice1 offset was %f", got, len(part.voices.voices), score.voiceChange, part.voices.voices[2].CurrentOffset)
	}

}

func BenchmarkPlayScore(b *testing.B) {
	score := NewScore()
	part, _ := score.NewPart("piano")
	score.Parts = []*Part{part}
	score.CurrentParts = []*Part{part}

	vm := VoiceMarker{AldaSourceContext{"stub", 0, 0}, 1}
	vm.UpdateScore(score)

	voice1 := part.GetVoice(1)
	voice1.CurrentOffset = 50.0

}
