package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/tidwall/pretty"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/client"
	"github.com/whosonfirst/go-whosonfirst-api/endpoint"
	"github.com/whosonfirst/go-whosonfirst-api/flags"
	"github.com/whosonfirst/go-whosonfirst-api/throttle"
	"github.com/whosonfirst/go-whosonfirst-api/writer"
	"io"
	"log"
	"os"
	"time"
)

func main() {

	var csv_flags flags.FileHandleFlags
	var filelist_flags flags.FileHandleFlags
	var geojson_flags flags.FileHandleFlags
	var geojsonls_flags flags.FileHandleFlags
	var stdout_flags flags.StdoutFlags

	var api_params api.APIParams

	flag.Var(&api_params, "param", "One or more Who's On First API query=value parameters.")

	flag.Var(&csv_flags, "csv", "")
	flag.Var(&filelist_flags, "filelist", "")
	flag.Var(&geojson_flags, "geojson", "")
	flag.Var(&geojsonls_flags, "geojson-ls", "")
	flag.Var(&stdout_flags, "stdout", "")

	// output/formatting

	var raw = flag.Bool("raw", false, "Dump raw Who's On First API responses.")
	var pretty_json = flag.Bool("pretty", false, "Pretty-print JSON results.")

	// pagination

	var paginated = flag.Bool("paginated", false, "Automatically paginate API results.")
	var async = flag.Bool("async", false, "Process API results asynchronously. If true then any errors processing a response are reported by will not stop execution.")

	// output formats

	// var filelist_prefix = flag.String("filelist-prefix", "", "Prepend each WOF \"file list\" result with this prefix.")
	// var filelist_output = flag.String("filelist-output", "", "The path to a file where WOF \"file list\"  output should be written. Output is written to STDOUT if empty.")

	// advanced

	var custom_endpoint = flag.String("endpoint", "", "Define a custom endpoint for the Who's On First API.")
	var oauth2 = flag.Bool("oauth2", false, "")

	// misc

	var timings = flag.Bool("timings", false, "Track and report total time to invoke an API method. Timings are printed to STDOUT.")

	flag.Parse()

	args := api_params.ToArgs()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var ep api.APIEndpoint

	if *oauth2 {

		access_token := args.Get("access_token")
		e, err := endpoint.NewOAuth2APIEndpoint(access_token)

		if err != nil {
			log.Fatal(err)
		}

		ep = e

	} else {

		api_key := args.Get("api_key")
		e, err := endpoint.NewMapzenAPIEndpoint(api_key)

		if err != nil {
			log.Fatal(err)
		}

		ep = e
	}

	if *custom_endpoint != "" {

		err := ep.SetEndpoint(*custom_endpoint)

		if err != nil {
			log.Fatal(err)
		}
	}

	th, err := throttle.NewDefaultThrottle(ctx)

	if err != nil {
		log.Fatal(err)
	}

	cl, err := client.NewHTTPClient(ep, th)

	if err != nil {
		log.Fatal(err)
	}

	writers := make([]api.APIResultWriter, 0)

	csv_func := func(fh io.Writer) (api.APIResultWriter, error) {
		return writer.NewCSVWriter(fh)
	}

	filelist_func := func(fh io.Writer) (api.APIResultWriter, error) {
		return writer.NewFileListWriter(fh)
	}

	geojson_func := func(fh io.Writer) (api.APIResultWriter, error) {
		return writer.NewGeoJSONWriter(fh)
	}

	geojsonls_func := func(fh io.Writer) (api.APIResultWriter, error) {
		return writer.NewGeoJSONLSWriter(fh)
	}

	stdout_func := func(fh io.Writer) (api.APIResultWriter, error) {
		return writer.NewStdoutWriter()
	}

	append_flags := func(fl flags.ResultWriterFlags, f flags.ResultWriterFunc, wr []api.APIResultWriter) ([]api.APIResultWriter, error) {

		filehandles, err := fl.FileHandles()

		if err != nil {
			return nil, err
		}

		for _, fh := range filehandles {

			wr, err := f(fh)

			if err != nil {
				return nil, err
			}

			writers = append(writers, wr)
		}

		return writers, nil
	}

	writers, err = append_flags(csv_flags, csv_func, writers)

	if err != nil {
		log.Fatal(err)
	}

	writers, err = append_flags(filelist_flags, filelist_func, writers)

	if err != nil {
		log.Fatal(err)
	}

	writers, err = append_flags(geojson_flags, geojson_func, writers)

	if err != nil {
		log.Fatal(err)
	}

	writers, err = append_flags(geojsonls_flags, geojsonls_func, writers)

	if err != nil {
		log.Fatal(err)
	}

	writers, err = append_flags(stdout_flags, stdout_func, writers)

	if err != nil {
		log.Fatal(err)
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
			_, err := dest.Write(raw)
			return err
		}
	}

	var t1 time.Time

	if *timings {
		t1 = time.Now()
	}

	method := args.Get("method")

	if *paginated {
		err = cl.ExecuteMethodPaginated(ctx, method, args, cb)

	} else {
		err = cl.ExecuteMethodWithCallback(ctx, method, args, cb)
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
