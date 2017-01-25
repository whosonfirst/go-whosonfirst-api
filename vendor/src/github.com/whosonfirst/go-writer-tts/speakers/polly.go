package speakers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"io"
	"io/ioutil"
	_ "log"
	"os"
	"strings"
)

type PollySpeaker struct {
	Speaker
	service      *polly.Polly
	OutputFormat string
	VoiceId      string
	Filename     string
	buffer       []string
}

func NewPollySpeaker() (*PollySpeaker, error) {

	// please fix me - this assumes shared credentials

	cfg := aws.NewConfig()
	cfg.WithRegion("us-east-1") // please fix me...

	sess, err := session.NewSession(cfg)

	if err != nil {
		return nil, err
	}

	buffer := make([]string, 0)

	svc := polly.New(sess)

	s := PollySpeaker{
		service:      svc,
		OutputFormat: "mp3",
		VoiceId:      "Russell",
		Filename:     "polly",
		buffer:       buffer,
	}

	return &s, nil
}

func (s *PollySpeaker) Read(reader io.Reader) error {

	tee := io.TeeReader(reader, s)
	_, err := ioutil.ReadAll(tee)
	return err
}

func (s *PollySpeaker) WriteString(text string) (int64, error) {
	r := strings.NewReader(text)
	return r.WriteTo(s)
}

func (s *PollySpeaker) Write(p []byte) (int, error) {

	s.buffer = append(s.buffer, string(p))
	return len(p), nil
}

func (s *PollySpeaker) Record() error {

	fname := fmt.Sprintf("%s.%s", s.Filename, s.OutputFormat)

	fh, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0644)
	defer fh.Close()

	if err != nil {
		return err
	}

	// https://docs.aws.amazon.com/polly/latest/dg/API_SynthesizeSpeech.html

	text := strings.Join(s.buffer, "\n\n")

	params := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String(s.OutputFormat),
		VoiceId:      aws.String(s.VoiceId),
		Text:         aws.String(text),
	}

	resp, err := s.service.SynthesizeSpeech(params)

	if err != nil {
		return err
	}

	_, err = io.Copy(fh, resp.AudioStream)

	if err != nil {
		return err
	}

	return nil
}

func (s *PollySpeaker) Close() error {
	return s.Record()
}
