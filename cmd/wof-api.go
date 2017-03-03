package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/client"
	"github.com/whosonfirst/go-whosonfirst-api/endpoint"
	"github.com/whosonfirst/go-whosonfirst-api/writer"
	"log"
	"os"
	"time"
)

func main() {

	var api_params api.APIParams

	flag.Var(&api_params, "param", "One or more Who's On First API query=value parameters.")

	var stdout = flag.Bool("stdout", false, "Write API results to STDOUT")
	var geojson = flag.Bool("geojson", false, "Transform API results to source GeoJSON for each Who's On First place.")
	var raw = flag.Bool("raw", false, "Dump raw Who's On First API responses.")
	var async = flag.Bool("async", false, "...")
	var timings = flag.Bool("timings", false, "...")	
	var paginated = flag.Bool("paginated", false, "Automatically paginate API results.")

	var output = flag.String("output", "", "...")

	var tts_speak = flag.Bool("tts", false, "Output integers to a text-to-speak engine.")
	var tts_engine = flag.String("tts-engine", "", "A valid go-writer-tts text-to-speak engine. Valid options are: osx, polly.")

	var custom_endpoint = flag.String("endpoint", "", "...")

	flag.Parse()

	args := api_params.ToArgs()

	api_key := args.Get("api_key")
	method := args.Get("method")

	if method == "" {
		log.Fatal("You forgot to specify a method")
	}

	e, err := endpoint.NewMapzenAPIEndpoint(api_key)

	if err != nil {
		log.Fatal(err)
	}

	if *custom_endpoint != "" {

		err := e.SetEndpoint(*custom_endpoint)

		if err != nil {
			log.Fatal(err)
		}

	}

	c, _ := client.NewHTTPClient(e)

	writers := make([]api.APIResultWriter, 0)

	if *tts_speak {

		ts, err := writer.NewTTSWriter(*tts_engine)

		if err != nil {
			log.Fatal(err)
		}

		writers = append(writers, ts)
	}

	if *geojson {

		dest := os.Stdout

		if *output != "" {

			f, err := os.OpenFile(*output, os.O_RDWR|os.O_CREATE, 0644)

			if err != nil {
				log.Fatal(err)
			}

			dest = f
		}

		wr, err := writer.NewGeoJSONWriter(dest)

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

	if *async {
		// please fix me...
	}

	multi := writer.NewAPIResultAsyncWriter(writers...)

	//multi = writer.NewAPIResultMultiWriter(writers...)

	if len(writers) == 0 && !*raw {
		log.Fatal("You forgot to specify an output source")
	}

	// why doesn't this work? see below... (20170125/thisisaaronland)

	/*
		defer func() {
			multi.Close()
		}()
	*/

	var cb api.APIResponseCallback

	cb = func(rsp api.APIResponse) error {

		results, err := rsp.Results()

		if err != nil {
			return err
		}

		for _, r := range results {
			multi.Write(r)
		}

		return nil
	}

	if *raw {
		cb = func(rsp api.APIResponse) error {
			_, err = os.Stdout.Write(rsp.Raw())
			return err
		}
	}

	var t1 time.Time

	if *timings {
		t1 = time.Now()
	}
	
	if *paginated {
		err = c.ExecuteMethodPaginated(method, args, cb)

	} else {
		err = c.ExecuteMethodWithCallback(method, args, cb)
	}

	if err != nil {
		msg := fmt.Sprintf("Failed to call '%s' because %s", method, err)
		log.Fatal(msg)
	}

	// I don't really understand why the defer func() stuff above
	// to do this doesn't work... (20170125/thisisaaronland)

	multi.Close()

	if *timings {
	   t2 := time.Since(t1)
	   log.Println("time %t\n", t2)	      
	}
	
	os.Exit(0)
}
