package main

import (
	"fmt"
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	mc := memcache.New("memcached.xzdlgx.0001.use2.cache.amazonaws.com:11211")
	mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})

	it, err := mc.Get("foo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(it)
}
