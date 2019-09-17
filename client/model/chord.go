package model

import "math"

// A Chord is a collection of notes and rests starting at the same point in
// time.
//
// Certain other types of events are allowed to occur between the notes and
// rests, e.g. octave and other attribute changes.
type Chord struct {
	Events []ScoreUpdate
}

func (chord Chord) updateScore(score *Score) error {
	score.chordMode = true
	for _, event := range chord.Events {
		if err := event.updateScore(score); err != nil {
			return err
		}
	}
	score.chordMode = false

	for _, part := range score.CurrentParts {
		// Notes/rests in a chord can have different durations. Following a chord, the
		// next note/rest is placed after the shortest note/rest in the chord.
		shortestDurationMs := math.MaxFloat64
		for _, event := range chord.Events {
			var specifiedDuration Duration
			switch event.(type) {
			case Note:
				specifiedDuration = event.(Note).Duration
			case Rest:
				specifiedDuration = event.(Rest).Duration
			}

			var duration Duration

			// If no duration is specified, use the part's default duration.
			if specifiedDuration.Components == nil {
				duration = part.Duration
			} else {
				duration = specifiedDuration
			}

			durationMs := float64(duration.Ms(part.Tempo))
			shortestDurationMs = math.Min(shortestDurationMs, durationMs)
		}

		part.LastOffset = part.CurrentOffset
		part.CurrentOffset += float32(shortestDurationMs)
	}

	return nil
}
