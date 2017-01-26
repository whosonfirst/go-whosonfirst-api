# go-whosonfirst-api

Go package for working with Who's On First API.

## Important

Too soon, move along.

## Install

You will need to have both `Go` and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Usage

### Simple

```
import (
	"github.com/whosonfirst/go-whosonfirst-api/client"
	"github.com/whosonfirst/go-whosonfirst-api/endpoint"
	"log"
)

api_key := "mapzen-xxxxx"
	
api_endpoint, _ := endpoint.NewMapzenAPIEndpoint(api_key)
api_client, _ := client.NewHTTPClient(api_endpoint)

method := "whosonfirst.places.search"
	
args := api_client.DefaultArgs()
args.Set("query", "Poutine")
args.Set("placetype", "venue")	
	
rsp, _ := c.ExecuteMethod(method, args)
log.Println(rsp.Raw())
```

_Note that error handling has been left out for the sake of brevity._

## Tools

### wof-api

```
./bin/wof-api -h
Usage of ./bin/wof-api:
  -endpoint string
    	TBW
  -geojson
    	Transform API results to source GeoJSON for each Who's On First place.
  -output string
    	TBW
  -paginated
    	Automatically paginate API results.
  -param value
    	One or more Who's On First API query=value parameters.
  -raw
    	Dump raw Who's On First API responses.
  -stdout
    	Write API results to STDOUT
  -tts
    	Output integers to a text-to-speak engine.
  -tts-engine string
    	A valid go-writer-tts text-to-speak engine. Valid options are: osx, polly.
```

## See also


