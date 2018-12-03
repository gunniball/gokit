package cache

type CacheService interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(k string)
	Flush()
}
