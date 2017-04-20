package writer

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-writer-tts"
	"github.com/whosonfirst/go-writer-tts/speakers"
	_ "log"
)

type TTSWriter struct {
	api.APIResultWriter
	writer speakers.Speaker
}

func NewTTSWriter(engine string) (*TTSWriter, error) {

	writer, err := tts.NewSpeakerForEngine(engine)

	if err != nil {
		return nil, err
	}

	tw := TTSWriter{
		writer: writer,
	}

	return &tw, nil
}

func (w *TTSWriter) WriteResult(r api.APIPlacesResult) (int, error) {

	text := fmt.Sprintf("%s is a %s with Who's On First ID %d", r.WOFName(), r.WOFPlacetype(), r.WOFId())
	return w.Write([]byte(text))
}

func (w *TTSWriter) Write(p []byte) (int, error) {

	return w.writer.Write(p)
}

func (w *TTSWriter) Close() error {

	return w.writer.Close()
}
