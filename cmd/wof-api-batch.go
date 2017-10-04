package main

// https://github.com/whosonfirst/whosonfirst-www-api/issues/99#issuecomment-333960724

import (
       "net/http"
)

func process() {

}

func handler() (http.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		query := req.Url.Query()

		
		
		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		// rsp.Write(enc)
	}

 	h := http.HandlerFunc(fn)
	return h, nil
}

func main () {


     batch_handler, err := handler()
}
