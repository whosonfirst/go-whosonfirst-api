# go-whosonfirst-api

Go package for working with Who's On First API.

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

## Interfaces

Interfaces are still a bit of a moving target. Or more specifically existing interfaces that have been defined should not change but there also aren't interfaces for many types of WOF API responses. That's why the default `APIResponse` interface defines a `Raw()` that returns plain-vanilla bytes and leaves it as an exercise to consumers to figure out what to do with them.

While all of the interfaces still need to be documented properly the most important ones are:

```
type APIResponse interface {
	Raw() []byte
	String() string
	Ok() (bool, APIError)
	Pagination() (APIPagination, error)
	Places() ([]APIPlacesResult, error)
}

type APIPlacesResult interface {
	WOFId() int64
	WOFParentId() int64
	WOFName() string
	WOFPlacetype() string
	WOFCountry() string
	WOFRepo() string
	Path() string
	URI() string
	String(...APIResultFlag) string
}

type APIError interface {
	String() string
	Code() int64
	Message() string
}

type APIPagination interface {
	String() string
	Pages() int
	Page() int
	PerPage() int
	Total() int
	Cursor() string
	NextQuery() string
}
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

* Fetch all the venues that are in [San Francisco](https://whosonfirst.mapzen.com/spelunker/id/85922583/) by passing the `-param placetype=venue` and `-param locality_id=85922583` flags, respectively.
* Do so in batches of 500 and handle pagination automatically by passing the `-param per_page=500` and `-paginated` flags respectively.
* For each result fetch the source GeoJSON file over the network, and do so asynchronously, creating a new `FeatureCollection` and save it as `venues.geojson` by passing the `-geojson`, `-async` and `-output venues.geojson` flags, respectively.
* Print how long the whole thing takes by passing the `-timings` flag.

You could also do the same by calling the `whosonfirst.places.getDescendants` API method, like this:

```
./bin/wof-api -param method=whosonfirst.places.getDescendants -param id=85922583 -param api_key=mapzen-XXXXXX -param per_page=500 -param placetype=venue -paginated -geojson -output descendants.geojson -timings -async
2017/03/03 17:56:14 Failed to retrieve https://whosonfirst.mapzen.com/data/110/880/049/1/1108800491.geojson because 404 Not Found
2017/03/03 17:56:14 Failed to retrieve https://whosonfirst.mapzen.com/data/110/880/049/3/1108800493.geojson because 404 Not Found
2017/03/03 17:56:14 Failed to retrieve https://whosonfirst.mapzen.com/data/110/882/755/7/1108827557.geojson because 404 Not Found
2017/03/03 17:56:15 time to 'whosonfirst.places.getDescendants': 5m16.811679531s
```

If you're wondering about the `-geojson` flag it's useful because the Who's On First API returns a minimum subset of a record's properties by default and does not return geometries at all (at least not yet). For example, here's what a default API response for a place looks like:

```
{
	"wof:id": 202863435,
	"wof:parent_id": "85887433",
	"wof:name": "18th Ave Photo",
	"wof:placetype": "venue",
	"wof:country": "US",
	"wof:repo": "whosonfirst-data-venue-us-ca"
}
```																					

The `-geojson` flag will instruct the `wof-api` tool to determine the fully qualified URL for a record – for example `202863435` becomes `https://whosonfirst.mapzen.com/data/202/863/435/202863435.geojson` – and then fetch [the contents of that file](https://whosonfirst.mapzen.com/data/202/863/435/202863435.geojson) and use that (rather than the default response above) in your final output.

## See also

* https://mapzen.com/documentation/wof/
