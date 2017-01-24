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

     	   
     	log.Println(rsp.String())
	return nil		
     }

     err := c.ExecuteMethodPaginated(method, args, cb)

     if err != nil {
     	log.Fatal(err)
     }     
}
