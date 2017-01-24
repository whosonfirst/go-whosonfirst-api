package main

import (
       "flag"
       "github.com/whosonfirst/go-whosonfirst-api"
       "github.com/whosonfirst/go-whosonfirst-api/client"
       "github.com/whosonfirst/go-whosonfirst-api/endpoint"
       "log"
)

func main () {

     var api_key = flag.String("api-key", "", "A valid Mapzen API key")
     var field = flag.String("field", "q", "...")
     var query = flag.String("query", "", "...")     

     flag.Parse()
     
     e, _ := endpoint.NewMapzenAPIEndpoint(*api_key)
     c, _ := client.NewHTTPClient(e)

     method := "whosonfirst.places.search"
     
     args := c.DefaultArgs()
     args.Set(*field, *query)

     cb := func(rsp api.APIResponse) error {

     	results, err := rsp.Results()

	if err != nil {
	   return err
	}

	for _, r := range results {
		log.Println(r.WOFId(), r.WOFName())
		log.Println(r.URI())
	}
	
	return nil		
     }

     err := c.ExecuteMethodPaginated(method, args, cb)

     if err != nil {
     	log.Fatal(err)
     }     
}
