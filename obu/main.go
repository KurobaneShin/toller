package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const sendInterval = time.Second

type OBUData struct {
	OBUID int     `json:"obuId,omitempty"`
	Lat   float64 `json:"lat,omitempty"`
	Long  float64 `json:"long,omitempty"`
}

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func genCoord() float64 {
	n, f := float64(rand.Intn(100)+1), rand.Float64()
	return n + f
}

func main() {
	obuIds := generateOBUIDS(20)
	for {
		for i := 0; i < len(obuIds); i++ {
			lat,long := genLatLong()
			data := OBUData{
				OBUID:obuIds[i],
				Lat:lat,
				Long:long,

			}
			fmt.Println(data)
			
		} 
		time.Sleep(sendInterval)
	}
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {

		ids[i] = rand.Intn(math.MaxInt)
	}

	return ids
}
