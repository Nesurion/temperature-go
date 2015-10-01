package service

import (
	"github.com/d2r2/go-dht"
	"github.com/influxdb/influxdb/client"
)

type sensorData struct {
	temperature float32
	humidity    float32
}

func measure(srv *Service) sensorData {
	// TODO: add ReadDHT parameters to config
	temp, hum, _, err := ReadDHTxxWithRetry(dht.DHT22, 4, true, 15)
	if err != nil {
		glog.Fatal(err)
	}
	s := sensorData{
		temperature: temp,
		humidity:    hum,
	}
	return s
}

func writeData(s sensorData, srv *Service) {
	sensorDataTypes := 2
	pts = make([]client.Point, sensorDataTypes)

	// TODO: make tags configurable
	pts[0] = client.Point{
		Measurement: "temperature",
		Tags: map[string]string{
			"room": "maik",
		},
		Fields: map[string]interface{}{
			"value": s.temperature,
		},
		Time:      time.Now(),
		Precision: "s",
	}
	pts[1] = client.Point{
		Measurement: "humidity",
		Tags: map[string]string{
			"room": "maik",
		},
		Fields: map[string]interface{}{
			"value": s.humidity,
		},
		Time:      time.Now(),
		Precision: "s",
	}

	bps := client.BatchPoints{
		Points:   pts,
		Database: srv.InfluxDB,
	}
	_, err := srv.InfluxClient.Write(bps)
	if err != nil {
		log.Fatal(err)
	}
}
