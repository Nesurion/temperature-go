package service

import (
	// "time"

	owm "github.com/briandowns/openweathermap"
	"github.com/d2r2/go-dht"
	"github.com/golang/glog"
	"github.com/influxdb/influxdb/client"
)

type sensorData struct {
	temperature float64
	humidity    float64
}

func measure(srv *Service) sensorData {
	var sensorType dht.SensorType
	switch srv.Config.DhtType {
	case "DHT22":
		sensorType = dht.DHT22
	case "DHT11":
		sensorType = dht.DHT11
	}

	temp, hum, _, err := dht.ReadDHTxxWithRetry(sensorType, srv.Config.DhtPin, srv.Config.DHTPerf, srv.Config.DhtRetries)
	if err != nil {
		glog.Info(err)
	}
	s := sensorData{
		temperature: float64(temp),
		humidity:    float64(hum),
	}
	return s
}

func outside(srv *Service) sensorData {
	w, err := owm.NewCurrent("C", "EN")
	if err != nil {
		glog.Infoln(err)
	}
	w.CurrentByID(srv.OWMcityID)
	s := sensorData{
		temperature: w.Main.Temp,
		humidity:    float64(w.Main.Humidity),
	}
	return s
}

func writeData(si sensorData, so sensorData, srv *Service) {
	// now := time.Now()
	glog.Infof("[Inside] Temperature: %.1f°C | Humidity: %.1f%%", si.temperature, si.humidity)
	glog.Infof("[Outside] Temperature: %.1f°C | Humidity: %.1f%%", so.temperature, so.humidity)
	if srv.Config.DryRun {
		return
	}
	sensorPoints := 4
	var pts = make([]client.Point, sensorPoints)

	// TODO: - make tags configurable
	// 		 - Do this in a loop
	pts[0] = client.Point{
		Measurement: "temperature",
		Tags: map[string]string{
			"room": "maik",
		},
		Fields: map[string]interface{}{
			"value": si.temperature,
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
			"value": si.humidity,
		},
		// Time:      now,
		Precision: "s",
	}
	pts[2] = client.Point{
		Measurement: "temperature",
		Tags: map[string]string{
			"room": "outside",
		},
		Fields: map[string]interface{}{
			"value": so.temperature,
		},
		// Time:      now,
		Precision: "s",
	}
	pts[3] = client.Point{
		Measurement: "humidity",
		Tags: map[string]string{
			"room": "outside",
		},
		Fields: map[string]interface{}{
			"value": so.humidity,
		},
		// Time:      now,
		Precision: "s",
	}

	bps := client.BatchPoints{
		Points:          pts,
		Database:        srv.InfluxDB,
		RetentionPolicy: srv.Config.InfluxRetention,
	}
	_, err := srv.InfluxClient.Write(bps)
	if err != nil {
		glog.Info(err)
	}
}
