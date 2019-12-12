package cache

import (
	"github.com/patrickmn/go-cache"
	"sync"
)

var shared *cache.Cache
var once sync.Once

func Init(){
	once.Do(func() {
		shared = cache.New(cache.NoExpiration, cache.NoExpiration)
	})
}

func GetShared() cache.Cache{
	return *shared
}
func SetValue(k string,v interface{}){
	shared.Set(k, v, cache.NoExpiration)
}

func GetString(k string) string{
	if x, found := shared.Get(k); found {
		return x.(string)
	}

	return ""
}

func GetRaw(k string) (interface{}, bool){
	return shared.Get(k)
}