package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/client"
	"github.com/whosonfirst/go-whosonfirst-api/endpoint"
	"github.com/whosonfirst/go-whosonfirst-api/writer"
	"log"
	"os"
	"strings"
)

// please rename me...

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var myFlags arrayFlags

func main() {

	flag.Var(&myFlags, "F", "Some description for this param.")

	var api_key = flag.String("api-key", "", "A valid Mapzen API key")

	var stdout = flag.Bool("stdout", false, "")
	var geojson = flag.Bool("geojson", false, "")
	var paginated = flag.Bool("paginated", false, "")

	var tts_speak = flag.Bool("tts", false, "Output integers to a text-to-speak engine.")
	var tts_engine = flag.String("tts-engine", "", "A valid go-writer-tts text-to-speak engine. Valid options are: osx.")

	flag.Parse()

	e, _ := endpoint.NewMapzenAPIEndpoint(*api_key)
	c, _ := client.NewHTTPClient(e)

	var method string

	args := c.DefaultArgs()

	for _, str_pair := range myFlags {

		pair := strings.Split(str_pair, "=")

		if pair[0] == "method" {
			method = pair[1]
			continue
		}

		args.Set(pair[0], pair[1])
	}

	if method == "" {
		log.Fatal("You forgot to specify a method")
	}

	writers := make([]api.APIResultWriter, 0)

	if *tts_speak {

		ts, err := writer.NewTTSWriter(*tts_engine)

		if err != nil {
			log.Fatal(err)
		}

		writers = append(writers, ts)
	}

	if *geojson {

		// please give me a better output source...
		wr, err := writer.NewGeoJSONWriter(os.Stdout)

		if err != nil {
			log.Fatal(err)
		}

		writers = append(writers, wr)
	}

	if *stdout {

		st, err := writer.NewStdoutWriter()

		if err != nil {
			log.Fatal(err)
		}

		writers = append(writers, st)
	}

	multi := writer.NewAPIResultMultiWriter(writers...)

	if len(writers) == 0 {
		log.Fatal("You forgot to specify an output source")
	}

	defer func() {
		multi.Close()
	}()

	cb := func(rsp api.APIResponse) error {

		results, err := rsp.Results()

		if err != nil {
			return err
		}

		for _, r := range results {
			multi.Write(r)
		}

		return nil
	}

	var err error

	if *paginated {

		err = c.ExecuteMethodPaginated(method, args, cb)

	} else {

		err = c.ExecuteMethodWithCallback(method, args, cb)
	}

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
