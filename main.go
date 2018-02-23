package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var requests, connections safeCounter

type safeCounter struct {
	uint
	sync.RWMutex
}

func main() {
	go func() {
		for {
			<-time.NewTicker(time.Second).C
			requests.Lock()
			connections.RLock()
			log.Printf("connections: %03d, requests: %04d", connections.uint, requests.uint)
			requests.uint = 0
			requests.Unlock()
			connections.RUnlock()
		}
	}()

	http.HandleFunc("/current", bytesHandler(current))
	http.HandleFunc("/forecast", bytesHandler(forecast))
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func bytesHandler(bs []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		connections.Lock()
		connections.uint++
		connections.Unlock()
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(bs)
		if err != nil {
			log.Print(errors.Wrap(err, "writing body"))
		}
		requests.Lock()
		requests.uint++
		requests.Unlock()

		connections.Lock()
		connections.uint--
		connections.Unlock()
	}
}

var current = []byte(fmt.Sprintf(`{
  "coord": {
    "lon": -78.12,
    "lat": 28.46
  },
  "weather": [
    {
      "id": 800,
      "main": "Clear",
      "description": "clear sky",
      "icon": "01d"
    }
  ],
  "base": "stations",
  "main": {
    "temp": 22.94,
    "pressure": 1038.18,
    "humidity": 100,
    "temp_min": 22.94,
    "temp_max": 22.94,
    "sea_level": 1038.2,
    "grnd_level": 1038.18
  },
  "wind": {
    "speed": 0.56,
    "deg": 30.5012
  },
  "clouds": {
    "all": 0
  },
  "dt": %d,
  "sys": {
    "message": 0.0019,
    "sunrise": 1518090923,
    "sunset": 1518130699
  },
  "id": 0,
  "name": "",
  "cod": 200
}`, time.Now().Unix()))

var forecast = []byte(fmt.Sprintf(`{
	"cod":"200",
	"message":0.0022,
	"cnt":40,
	"list":
	[{
		"dt":%d,
		"main":{"temp":22.51,"temp_min":22.37,"temp_max":22.51,"pressure":1037.57,"sea_level":1037.57,"grnd_level":1037.57,"humidity":100,"temp_kf":0.14},
		"weather":[{"id":802,"main":"Clouds","description":"scattered clouds","icon":"03n"}],
		"clouds":{"all":48},
		"wind":{"speed":5.02,"deg":90.5008},
		"rain":{},"sys":{"pod":"n"},"dt_txt":"2018-02-09 09:00:00"
	},
	{
		"dt":%d,
		"main":{"temp":22.82,"temp_min":22.72,"temp_max":22.82,"pressure":1038.23,"sea_level":1038.28,"grnd_level":1038.23,"humidity":100,"temp_kf":0.1},
		"weather":[{"id":801,"main":"Clouds","description":"few clouds","icon":"02d"}],
		"clouds":{"all":20},
		"wind":{"speed":5.23,"deg":95},
		"rain":{},"sys":{"pod":"d"},"dt_txt":"2018-02-09 12:00:00"
	},
	{
		"dt":%d,
		"main":{"temp":23.1,"temp_min":23.03,"temp_max":23.1,"pressure":1039.45,"sea_level":1039.51,"grnd_level":1039.45,"humidity":98,"temp_kf":0.07},
		"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}],
		"clouds":{"all":0},
		"wind":{"speed":5.56,"deg":101.501},
		"rain":{},"sys":{"pod":"d"},"dt_txt":"2018-02-09 15:00:00"
	}],
		"city":{"coord":{}}
		}`, time.Now().Add(17*time.Minute).Unix(), time.Now().Add(197*time.Minute).Unix(), time.Now().Add(377*time.Minute).Unix()))