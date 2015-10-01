package service

import (
	// "time"

	"github.com/d2r2/go-dht"
	"github.com/golang/glog"
	"github.com/influxdb/influxdb/client"
)

type sensorData struct {
	temperature float32
	humidity    float32
}

func measure(srv *Service) sensorData {
	// TODO: add ReadDHT parameters to config
	var sensorType dht.SensorType
	switch srv.Config.DhtType {
	case "DHT22":
		sensorType = dht.DHT22
	case "DHT11":
		sensorType = dht.DHT11
	}

	temp, hum, _, err := dht.ReadDHTxxWithRetry(sensorType, srv.Config.DhtPin, srv.Config.DHTPerf, srv.Config.DhtRetries)
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
	// now := time.Now()
	glog.Infof("Temperature: %v°C | Humidity: %v%%", s.temperature, s.humidity)
	if srv.Config.DryRun {
		return
	}
	sensorDataTypes := 2
	var pts = make([]client.Point, sensorDataTypes)

	// TODO: make tags configurable
	pts[0] = client.Point{
		Measurement: "temperature",
		Tags: map[string]string{
			"room": "maik",
		},
		Fields: map[string]interface{}{
			"value": s.temperature,
		},
		// Time:      now,
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
		// Time:      now,
		Precision: "s",
	}

	bps := client.BatchPoints{
		Points:   pts,
		Database: srv.InfluxDB,
	}
	_, err := srv.InfluxClient.Write(bps)
	if err != nil {
		glog.Fatal(err)
	}
}
