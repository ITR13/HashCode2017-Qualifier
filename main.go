package main

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	system, videoSizes := ReadFile("current.in")

}

func ReadFile(path string) (*System, []int) {
	f, err := os.Open("current.in")
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
	fmt.Println("Video-Sizes: ", videoSizes)

	system := GetSystem(videoC, endpointC, cacheC, cacheSize)

	for i := 0; i < endpointC; i++ {
		var ILatency, connected int
		fmt.Fscanf(f, "%d %d\n", &ILatency, &connected)
		cacheN, latencies := make([]int, connected), make([]int, connected)
		for j := 0; j < connected; j++ {
			fmt.Fscanf(f, "%d %d\n", &(cacheN[j]), &(latencies[j]))
		}
		system.MakeLinks(system.endpoints[i], ILatency, cacheN, latencies)
	}

	for i := 0; i < rDescC; i++ {
		var videoN, endpointN, requests int
		fmt.Fscanf(f, "%d %d %d\n", &videoN, &endpointN, &requests)
		system.endpoints[endpointN].RegisterRequest(videoN, requests)
	}
	return system, videoSizes
}
