package main

type CachedVideo struct {
	timeSaved int
	size      int
}

var totalTimeSaved int

type CacheServer struct {
	size      int
	timeSaved TimeSaved
	videos    []*CachedVideo
	endpoints []*CacheEndpointLink
}

type TimeSaved []int

type CacheEndpointLink struct {
	timeSaved   int
	cache       *CacheServer
	endpoint    *Endpoint
	timeSavedPV TimeSaved
}

type Endpoint []*CacheEndpointLink

type System struct {
	caches    []*CacheServer
	endpoints []*Endpoint
}

func GetSystem(videoC, endpointC, cacheC int, cacheSize int) *System {
	caches := make([]*CacheServer, cacheC)
	endpoints := make([]*Endpoint, endpointC)
	for i := 0; i < cacheC; i++ {
		caches[i] = &CacheServer{cacheSize, make(TimeSaved, videoC),
			make([]*CachedVideo, 0), make([]*CacheEndpointLink, 0)}
	}
	return &System{caches, endpoints}
}

func (system *System) MakeLinks(endpoint *Endpoint, ILatency int, cacheN []int, latencies []int) {
	ep := Endpoint(make([]*CacheEndpointLink, len(cacheN)))
	endpoint = &ep
	for i := 0; i < len(ep); i++ {
		cache := system.caches[cacheN[i]]
		ep[i] = &CacheEndpointLink{ILatency - latencies[i],
			cache, endpoint, make(TimeSaved, len(cache.timeSaved))}
		cache.endpoints = append(cache.endpoints, ep[i])
	}
}

func (endpoint *Endpoint) RegisterRequest(video, times int) {
	for i := 0; i < len(*endpoint); i++ {
		saved := times * (*endpoint)[i].timeSaved
		(*endpoint)[i].timeSavedPV[video] += saved
		(*endpoint)[i].cache.timeSaved[video] += saved
	}
}

//Note, cannot unregister yet
//TODO consider not using append and [:]
//TODO make it possible to remove video in any order
func (cache *CacheServer) RegisterVideo(video, size int) {
	save := cache.timeSaved[video]
	totalTimeSaved += save
	cache.size -= size
	cache.videos = append(cache.videos,
		&CachedVideo{save, size})
	for i := 0; i < len(cache.endpoints); i++ {
		ep := cache.endpoints[i].endpoint
		for j := 0; j < len(*ep); j++ {
			link := (*ep)[j]
			link.cache.timeSaved[video] = 0
		}
	}

}
