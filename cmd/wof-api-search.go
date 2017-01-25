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

	flag.Parse()

	e, _ := endpoint.NewMapzenAPIEndpoint(*api_key)
	c, _ := client.NewHTTPClient(e)

	method := "whosonfirst.places.search"

	args := c.DefaultArgs()
	args.Set(*field, *query)

	ts, err := writer.NewTTSWriter("polly")

	if err != nil {
		log.Fatal(err)
	}
	
	wr, err := writer.NewGeoJSONWriter(os.Stdout)

	if err != nil {
		log.Fatal(err)
	}

	defer wr.Close()

	writers := []api.APIResultWriter{
		wr,
		ts,			
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

	err = c.ExecuteMethodPaginated(method, args, cb)

	if err != nil {
		log.Fatal(err)
	}
}
