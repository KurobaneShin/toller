package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"

	"github.com/KurobaneShin/tolling/types"
)

const (
	sendInterval = time.Second * 5
	wsEndpoint   = "ws://127.0.0.1:30000/ws"
)

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func genCoord() float64 {
	n, f := float64(rand.Intn(100)+1), rand.Float64()
	return n + f
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	obuIds := generateOBUIDS(20)
	for {
		for i := 0; i < len(obuIds); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: obuIds[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}

		}
		time.Sleep(sendInterval)
	}
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(999999)
	}

	return ids
}
