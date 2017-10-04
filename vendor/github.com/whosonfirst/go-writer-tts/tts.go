package tts

import (
	"errors"
	"github.com/whosonfirst/go-writer-tts/speakers"
)

func NewSpeakerForEngine(engine string, options ...interface{}) (speakers.Speaker, error) {

	if engine == "osx" {
		return speakers.NewOSXSpeaker()
	} else if engine == "polly" {
		return speakers.NewPollySpeaker()
	} else {
	}

	return nil, errors.New("Unknown or unsupported text to speech engine")
}
