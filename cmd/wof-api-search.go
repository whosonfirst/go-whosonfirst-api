package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/client"
	"github.com/whosonfirst/go-whosonfirst-api/endpoint"
	"github.com/whosonfirst/go-whosonfirst-api/writer"
	"log"
	"os"
)

func main() {

	var api_key = flag.String("api-key", "", "A valid Mapzen API key")
	var field = flag.String("field", "q", "...")
	var query = flag.String("query", "", "...")

	var geojson = flag.Bool("geojson", false, "")

	var tts_speak = flag.Bool("tts", false, "Output integers to a text-to-speak engine.")
	var tts_engine = flag.String("tts-engine", "", "A valid go-writer-tts text-to-speak engine. Valid options are: osx.")

	flag.Parse()

	e, _ := endpoint.NewMapzenAPIEndpoint(*api_key)
	c, _ := client.NewHTTPClient(e)

	method := "whosonfirst.places.search"

	args := c.DefaultArgs()
	args.Set(*field, *query)

	writers := make([]api.APIResultWriter, 0)

	if *tts_speak {

		ts, err := writer.NewTTSWriter(*tts_engine)

		if err != nil {
			log.Fatal(err)
		}

		writers = append(writers, ts)

		defer func() { ts.Close() }()
	}

	if *geojson {
		wr, err := writer.NewGeoJSONWriter(os.Stdout)

		if err != nil {
			log.Fatal(err)
		}

		writers = append(writers, wr)

		defer func() { wr.Close() }()
	}

	multi := writer.NewAPIResultMultiWriter(writers...)

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

	err := c.ExecuteMethodPaginated(method, args, cb)

	if err != nil {
		log.Fatal(err)
	}

}
