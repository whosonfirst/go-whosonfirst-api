package main

import (
	"flag"
	"fmt"
	"github.com/tidwall/pretty"
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

	// output/formatting

	var stdout = flag.Bool("stdout", false, "Write API results to STDOUT")
	var raw = flag.Bool("raw", false, "Dump raw Who's On First API responses.")
	var pretty_json = flag.Bool("pretty", false, "Pretty-print JSON results.")

	// pagination

	var paginated = flag.Bool("paginated", false, "Automatically paginate API results.")
	var async = flag.Bool("async", false, "Process API results asynchronously. If true then any errors processing a response are reported by will not stop execution.")

	// output formats

	var csv = flag.Bool("csv", false, "Transform API results to source CSV for each API result.")
	var csv_output = flag.String("csv-output", "", "The path to a file where CSV output should be written. Output is written to STDOUT if empty.")

	var filelist = flag.Bool("filelist", false, "Transform API results to a WOF \"file list\".")
	var filelist_prefix = flag.String("filelist-prefix", "", "Prepend each WOF \"file list\" result with this prefix.")
	var filelist_output = flag.String("filelist-output", "", "The path to a file where WOF \"file list\"  output should be written. Output is written to STDOUT if empty.")

	var geojson = flag.Bool("geojson", false, "Transform API results to source GeoJSON for each API result, collating everything in to a single GeoJSON Feature Collection.")
	var geojson_output = flag.String("geojson-output", "", "The path to a file where GeoJSON output should be written. Output is written to STDOUT if empty.")

	var geojson_ls = flag.Bool("geojson-ls", false, "Transform API results to line-separated source GeoJSON for each API result, with one GeoJSON Feature per line.")
	var geojson_ls_output = flag.String("geojson-ls-output", "", "The path to a file where line-separated GeoJSON output should be written. Output is written to STDOUT if empty.")

	// silly

	var tts_speak = flag.Bool("tts", false, "Output integers to a text-to-speak engine.")
	var tts_engine = flag.String("tts-engine", "", "A valid go-writer-tts text-to-speak engine. Valid options are: osx, polly.")

	// advanced

	var custom_endpoint = flag.String("endpoint", "", "Define a custom endpoint for the Who's On First API.")

	// misc

	var timings = flag.Bool("timings", false, "Track and report total time to invoke an API method. Timings are printed to STDOUT.")

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

		if *geojson_output != "" {

			f, err := os.OpenFile(*geojson_output, os.O_RDWR|os.O_CREATE, 0644)

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

	if *geojson_ls {

		dest := os.Stdout

		if *geojson_ls_output != "" {

			f, err := os.OpenFile(*geojson_ls_output, os.O_RDWR|os.O_CREATE, 0644)

			if err != nil {
				log.Fatal(err)
			}

			dest = f
		}

		wr, err := writer.NewGeoJSONLSWriter(dest)

		if err != nil {
			log.Fatal(err)
		}

		writers = append(writers, wr)
	}

	if *csv {

		dest := os.Stdout

		if *csv_output != "" {

			fh, err := os.OpenFile(*csv_output, os.O_RDWR|os.O_CREATE, 0644)

			if err != nil {
				log.Fatal(err)
			}

			dest = fh
		}

		wr, err := writer.NewCSVWriter(dest)

		if err != nil {
			log.Fatal(err)
		}

		writers = append(writers, wr)

	}

	if *filelist {

		dest := os.Stdout

		if *filelist_output != "" {

			fh, err := os.OpenFile(*filelist_output, os.O_RDWR|os.O_CREATE, 0644)

			if err != nil {
				log.Fatal(err)
			}

			dest = fh
		}

		wr, err := writer.NewFileListWriter(dest)

		if err != nil {
			log.Fatal(err)
		}

		wr.Prefix = *filelist_prefix

		writers = append(writers, wr)
	}

	if *stdout || len(writers) == 0 {

		st, err := writer.NewStdoutWriter()

		if err != nil {
			log.Fatal(err)
		}

		writers = append(writers, st)
	}

	var multi api.APIResultMultiWriter

	multi = writer.NewAPIResultMultiWriterSync(writers...)

	if *async {
		multi = writer.NewAPIResultMultiWriterAsync(writers...)
	}

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

		results, err := rsp.Places()

		if err != nil {
			log.Println(err)
			return err
		}

		for _, r := range results {
			multi.Write(r)
		}

		return nil
	}

	if *raw {

		dest := os.Stdout

		cb = func(rsp api.APIResponse) error {

			raw := rsp.Raw()

			if *pretty_json {
				raw = pretty.Pretty(raw)
			}
			_, err = dest.Write(raw)
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
		log.Printf("time to '%s': %v\n", method, t2)
	}

	os.Exit(0)
}
