# go-whosonfirst-api

Go package for working with Who's On First API.

## Important

Too soon, move along. Probably.

## Install

You will need to have both `Go` and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Usage

_Note that all error handling in the examples below has been removed for the sake of brevity._

### Simple

```
import (
	"github.com/whosonfirst/go-whosonfirst-api/client"
	"github.com/whosonfirst/go-whosonfirst-api/endpoint"
	"os"
)

api_key := "mapzen-xxxxxxx"
	
api_endpoint, _ := endpoint.NewMapzenAPIEndpoint(api_key)
api_client, _ := client.NewHTTPClient(api_endpoint)

method := "whosonfirst.places.search"
	
args := api_client.DefaultArgs()
args.Set("query", "poutine")
args.Set("placetype", "venue")	
	
rsp, _ := c.ExecuteMethod(method, args)
os.Stdout.Write(rsp.Raw())
```

### Paginated

```
import (
	"github.com/whosonfirst/go-whosonfirst-api/client"
	"github.com/whosonfirst/go-whosonfirst-api/endpoint"
	"os"
)

api_key := "mapzen-xxxxxxx"
	
api_endpoint, _ := endpoint.NewMapzenAPIEndpoint(api_key)
api_client, _ := client.NewHTTPClient(api_endpoint)

method := "whosonfirst.places.search"
	
args := api_client.DefaultArgs()
args.Set("query", "beer")
args.Set("placetype", "venue")
args.Set("locality_id", "101748417")

cb := func(rsp api.APIResponse) error {
	_, err := os.Stdout.Write(rsp.Raw())
	return err
}

c.ExecuteMethodPaginated(method, args, cb)
```

## Tools

### wof-api

`wof-api` is a command line tool for calling the Who's On First API.

```
./bin/wof-api -h
Usage of ./bin/wof-api:
  -async
    	Process API results asynchronously. If true then any errors processing a response are reported by will not stop execution.
  -endpoint string
    	Define a custom endpoint for the Who's On First API.
  -geojson
    	Transform API results to source GeoJSON for each Who's On First place.
  -output string
    	The path to a file where output should be written.
  -paginated
    	Automatically paginate API results.
  -param value
    	One or more Who's On First API query=value parameters.
  -raw
    	Dump raw Who's On First API responses.
  -stdout
    	Write API results to STDOUT
  -timings
    	Track and report total time to invoke an API method. Timings are printed to STDOUT.
  -tts
    	Output integers to a text-to-speak engine.
  -tts-engine string
    	A valid go-writer-tts text-to-speak engine. Valid options are: osx, polly.
```

#### Example

![](images/sf-venues.png)

Fetch all 63, 387 venues in [San Francisco](https://whosonfirst.mapzen.com/spelunker/id/85922583/) as a single GeoJSON `FeatureCollection` by calling the `whosonfirst.places.search` API method, like this:

```
./bin/wof-api -param method=whosonfirst.places.search -param locality_id=85922583 -param api_key=mapzen-XXXXXXX -param per_page=500 -param placetype=venue -paginated -geojson -output venues.geojson -timings -async
2017/03/03 17:29:11 Failed to retrieve https://whosonfirst.mapzen.com/data/110/880/049/3/1108800493.geojson because 404 Not Found
2017/03/03 17:29:11 Failed to retrieve https://whosonfirst.mapzen.com/data/110/880/049/1/1108800491.geojson because 404 Not Found
2017/03/03 17:29:11 Failed to retrieve https://whosonfirst.mapzen.com/data/110/882/755/7/1108827557.geojson because 404 Not Found
2017/03/03 17:30:09 Failed to retrieve https://whosonfirst.mapzen.com/data/236/676/137/236676137.geojson because 500 Internal Server Error
2017/03/03 17:31:17 time to 'whosonfirst.places.search': 5m22.656896289s
```

Here's what's going on:

* Fetch all the venues that are in [San Francisco](https://whosonfirst.mapzen.com/spelunker/id/85922583/) _by passing the `-param placetype=venue` and `-param locality_id=85922583` flags, respectively_.
* Do so in batches of 500 and handle pagination automatically _by passing the `-param per_page=500` and `-paginated` flags respectively_.
* For each result fetch the source GeoJSON file over the network, and do so asynchronously, creating a new `FeatureCollection` and save it as `venues.geojson` _by passing the `-geojson`, `-async` and `-output venues.geojson` flags, respectively_.
* Print how long the whole thing takes _by passing the `-timings` flag_.

You could also do the same by calling the `whosonfirst.places.getDescendants` API method, like this:

```
./bin/wof-api -param method=whosonfirst.places.getDescendants -param id=85922583 -param api_key=mapzen-XXXXXX -param per_page=500 -param placetype=venue -paginated -geojson -output descendants.geojson -timings -async
2017/03/03 17:56:14 Failed to retrieve https://whosonfirst.mapzen.com/data/110/880/049/1/1108800491.geojson because 404 Not Found
2017/03/03 17:56:14 Failed to retrieve https://whosonfirst.mapzen.com/data/110/880/049/3/1108800493.geojson because 404 Not Found
2017/03/03 17:56:14 Failed to retrieve https://whosonfirst.mapzen.com/data/110/882/755/7/1108827557.geojson because 404 Not Found
2017/03/03 17:56:15 time to 'whosonfirst.places.getDescendants': 5m16.811679531s
```

## See also


