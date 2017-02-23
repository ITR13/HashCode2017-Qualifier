package main

type CachedVideo struct {
	timeSaved int
	size      int
}

var totalTimeSaved int

type CacheServer struct {
	ID, size  int
	timeSaved TimeSaved
	videos    map[int]*CachedVideo
	endpoints map[int]*CacheEndpointLink
}

type TimeSaved []int

type CacheEndpointLink struct {
	timeSaved   int
	cache       *CacheServer
	endpoint    *Endpoint
	timeSavedPV TimeSaved
}

type CacheRank struct {
	timeSaved int
	cacheN    int
}

type Endpoint struct {
	caches     map[int]*CacheEndpointLink
	cacheRanks []CacheRank
}

type System struct {
	caches    []*CacheServer
	endpoints []*Endpoint
}

func GetSystem(videoC, endpointC, cacheC int, cacheSize int) *System {
	caches := make([]*CacheServer, cacheC)
	endpoints := make([]*Endpoint, endpointC)
	for i := 0; i < cacheC; i++ {
		caches[i] = &CacheServer{i, cacheSize, make(TimeSaved, videoC),
			make(map[int]*CachedVideo), make(map[int]*CacheEndpointLink)}
	}
	for i := 0; i < endpointC; i++ {
		endpoints[i] = &Endpoint{make(map[int]*CacheEndpointLink, cacheC),
			make([]CacheRank, videoC)}
	}

	return &System{caches, endpoints}
}

func (system *System) MakeLinks(epN int, ILatency int, cacheN []int, latencies []int) {
	endpoint := system.endpoints[epN]
	for i := 0; i < len(cacheN); i++ {
		cn := cacheN[i]
		cache := system.caches[cn]
		link := &CacheEndpointLink{ILatency - latencies[i],
			cache, endpoint, make(TimeSaved, len(cache.timeSaved))}
		endpoint.caches[cn] = link
		cache.endpoints[epN] = link
	}
}

func (endpoint *Endpoint) RegisterRequest(video, times int) {
	for _, value := range endpoint.caches {
		saved := times * value.timeSaved
		value.timeSavedPV[video] += saved
		value.cache.timeSaved[video] += saved
	}
}

//Note, cannot unregister yet
//TODO consider not using append and [:]
//TODO make it possible to remove video in any order
func (cache *CacheServer) RegisterVideo(vID, size int) {
	save := cache.timeSaved[vID]
	totalTimeSaved += save
	cache.size -= size
	cache.videos[vID] = &CachedVideo{save, size}
	for _, link := range cache.endpoints {
		ep := link.endpoint
		rank := ep.cacheRanks[vID]
		if rank.timeSaved < save {
			ep.ReplaceBest(vID, save, cache.ID)
		}
	}
}

func (cache *CacheServer) UnregisterVideo(vID int) {
	video := cache.videos[vID]
	delete(cache.videos, vID)
	cache.size += video.size
	if video.timeSaved != 0 {
		for _, link := range cache.endpoints {
			ep := link.endpoint
			rank := ep.cacheRanks[vID]
			if rank.cacheN == cache.ID {
				bestID, bestSave := -1, -1
				for cID, link2 := range ep.caches {
					video := link2.cache.videos[vID]
					if video != nil {
						if link2.timeSavedPV[vID] > bestSave {
							bestSave = link2.timeSavedPV[vID]
							bestID = cID
						}
					}
				}

				ep.ReplaceBest(vID, bestSave, bestID)
			}
		}
	}
}

func (ep *Endpoint) ReplaceBest(vID, save, cID int) {
	rank := ep.cacheRanks[vID]
	oldTS := rank.timeSaved
	rank.timeSaved = save
	rank.cacheN = cID

	for _, link2 := range ep.caches {
		linkSave := link2.timeSavedPV[vID] - save
		if linkSave < 0 {
			linkSave = 0
		}

		link2.cache.timeSaved[vID] -=
			(link2.timeSavedPV[vID] - linkSave) - oldTS

		otherVideo := link2.cache.videos[vID]
		if otherVideo != nil {
			otherVideo.timeSaved -= linkSave
		}
	}
}
