package main

import (
	"context"
	"fmt"
	"net/http"
	"github.com/redis/go-redis/v9"
// 	"html"
	"errors"
)

var pluginName = "rate_limiting"

var HandlerRegisterer = registerer(pluginName)

type registerer string

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(_ context.Context, extra map[string]interface{}, h http.Handler) (http.Handler, error) {
	config, ok := extra[pluginName].(map[string]interface{})
	if !ok {
		return h, errors.New("configuration not found")
	}

	path, _ := config["path"].(string)
	fmt.Printf("The plugin is now hijacking the path %s\n", path)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		if req.URL.Path != path {
// 			h.ServeHTTP(w, req)
// 			return
// 		}

// 		fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))

		apiKey := req.Header.Get("API_KEY")

		if apiKey == "" {
			fmt.Println("Error in API_KEY")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
        fmt.Println(apiKey)
		client := redis.NewClient(&redis.Options{
			Addr:     "host.docker.internal:6379",
			Password: "", // No password set
			DB:       0,  // Use default DB
		})

		ctx := context.Background()

// 		err := client.Set(ctx, "apikey1", "gold", 0).Err()
//         fmt.Println(err)
// 		if err != nil {
// 			panic(err)
// 		}
        all_records,err:=client.Keys(ctx,"*").Result()
        if err!=nil{
            panic(err)
            return
        }
        fmt.Println(all_records)
		val, err := client.Get(ctx, apiKey).Result()

		if err != nil {
// 		    panic(err)
		    fmt.Println("Invalid apikey")
		    w.WriteHeader(http.StatusUnauthorized)
		    fmt.Fprintf(w,"Invalid apikey")
// 		    return w.Header(http.StatusUnauthorized)
return
		}

// 		fmt.Println(val==apiKey)
//         fmt.Println(val)
//         fmt.Fprintf(w,"Rate limiting to be applied...")
        headers:=req.Header
        headers.Set("Tier",val)

                /* in zebra , the tenantID is probably embedded in the apikey itself and has to be decoded from the apikey to apply the
                rate limiting to each tenant.
                */
        queryValues := req.URL.Query()
        tenantId := queryValues.Get("tenantId")
//         fmt.Println(tenantId)
        if tenantId==""{
            fmt.Println("TenantID not found")
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Fprintf(w,"TenantID not found")
               return
        }
        fmt.Println(tenantId)
        headers.Set("Tenant-Id",tenantId)
        fmt.Println(req)
//         fmt.Println(w)

		h.ServeHTTP(w, req)
	}), nil
}

func main() {}
