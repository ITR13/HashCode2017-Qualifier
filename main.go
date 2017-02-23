package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func (ts TimeSaved) GetHighestUnder(max int, videoSizes []int) int {
	h := -1
	id := -1
	for i := 0; i < len(ts); i++ {
		if videoSizes[i] < max {
			if ts[i] > h {
				h = ts[i]
				id = i
			}
		}
	}
	return id
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	BadSolve("me_at_the_zoo")
	BadSolve("trending_today")
	BadSolve("videos_worth_spreading")
}

func BadSolve(name string) {
	fmt.Println("Solving " + name)
	system, videoSizes := ReadFile(name + ".in")

	for i := 0; i < len(system.caches); i++ {
		cache := system.caches[i]
		for highest :=
			cache.timeSaved.GetHighestUnder(cache.size,
				videoSizes); highest != -1; highest =
			cache.timeSaved.GetHighestUnder(cache.size, videoSizes) {
			cache.RegisterVideo(highest, videoSizes[highest])
		}
	}

	system.WriteFile(name + ".out")
}

func (system *System) WriteFile(path string) {
	f, err := os.Create(path)
	check(err)
	defer f.Close()

	writer := bufio.NewWriter(f)
	writer.WriteString(strconv.Itoa(len(system.caches)) + "\n")

	for i := 0; i < len(system.caches); i++ {
		writer.WriteString(strconv.Itoa(i))
		for vid, _ := range system.caches[i].videos {
			writer.WriteString(" " + strconv.Itoa(vid))
		}
		writer.WriteString("\n")

	}
}

func ReadFile(path string) (*System, []int) {
	f, err := os.Open(path)
	check(err)
	defer f.Close()

	var videoC, endpointC, rDescC, cacheC, cacheSize int
	fmt.Fscanf(f, "%d %d %d %d %d\n",
		&videoC, &endpointC, &rDescC, &cacheC, &cacheSize)

	fmt.Printf("Videos: %d\nEndpoints: %d\nDescriptions: %d\nCaches: %d\nSize: %d\n",
		videoC, endpointC, rDescC, cacheC, cacheSize)

	videoSizes := make([]int, videoC)
	for i := 0; i < videoC-1; i++ {
		fmt.Fscanf(f, "%d", &(videoSizes[i]))
	}
	fmt.Fscanf(f, "%d\n", &(videoSizes[len(videoSizes)-1]))

	system := GetSystem(videoC, endpointC, cacheC, cacheSize)

	for i := 0; i < endpointC; i++ {
		var ILatency, connected int
		fmt.Fscanf(f, "%d %d\n", &ILatency, &connected)
		cacheN, latencies := make([]int, connected), make([]int, connected)
		for j := 0; j < connected; j++ {
			fmt.Fscanf(f, "%d %d\n", &(cacheN[j]), &(latencies[j]))
		}
		system.MakeLinks(i, ILatency, cacheN, latencies)
	}

	for i := 0; i < rDescC; i++ {
		var videoN, endpointN, requests int
		fmt.Fscanf(f, "%d %d %d\n", &videoN, &endpointN, &requests)
		system.endpoints[endpointN].RegisterRequest(videoN, requests)
	}
	return system, videoSizes
}
