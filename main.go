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
	ReadFile("current.in")
}

func ReadFile(path string) (*System, []int) {
	f, err := os.Open("current.in")
	check(err)
	defer f.Close()

	var videoC, endpointC, rDescC, cacheC, cacheSize int
	fmt.Fscanf(f, "%d %d %d %d %d\n",
		&videoC, &endpointC, &rDescC, &cacheC, &cacheSize)
	videoSizes := make([]int, videoC)
	for i := 0; i < videoC-1; i++ {
		fmt.Fscanf(f, "%d ", &(videoSizes[i]))
	}
	fmt.Fscanf(f, "%d\n", &(videoSizes[len(videoSizes)-1]))

	system := GetSystem(videoC, endpointC, cacheC, cacheSize)

}
