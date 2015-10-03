package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/influxdb/influxdb/client"
	"github.com/ogier/pflag"

	"github.com/nesurion/temperature-go/service"
)

var (
	config = service.Config{}
)

func main() {
	pflag.Parse()

	glog.Info("=== Temperature-go ===")

	srv := service.New(serviceConfig(), serviceDeps())

	if srv.DryRun {
		glog.Info("running in dryRun mode")
	}

	go srv.Serve()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch

	glog.Info("Received termination signal. Closing ...")
	srv.Close()

}

func init() {
	pflag.DurationVar(&config.TickerTime, "ticker-time", 15*time.Minute, "Ticker time.")
	pflag.BoolVar(&config.DryRun, "dry-run", true, "Write to STDOUT instead of InfluxDB")
	pflag.IntVar(&config.OWMcityID, "owm-city-id", 0, "Open Weather Map city ID")

	pflag.IntVar(&config.InfluxPort, "influx-port", 8086, "InfluxDB Port")
	pflag.StringVar(&config.InfluxHost, "influx-host", "localhost", "InfluxDB Port")
	pflag.StringVar(&config.InfluxUser, "influx-user", "", "InfluxDB User")
	pflag.StringVar(&config.InfluxDB, "influx-db", "", "InfluxDB Database")
	pflag.StringVar(&config.InfluxPassword, "influx-password", "", "InfluxDB Password")
	pflag.StringVar(&config.InfluxRetention, "influx-retention", "default", "InfluxDB Retention")

	pflag.StringVar(&config.DhtType, "dht-type", "DHT22", "DHT Type (DHT11, DHT22)")
	pflag.IntVar(&config.DhtPin, "dht-pin", 4, "Pin Number DHT Data is connected to")
	pflag.BoolVar(&config.DHTPerf, "dht-perf", false, "Run DHT read in Boost Performance Mode - true will result in needing sudo")
	pflag.IntVar(&config.DhtRetries, "dht-retries", 15, "Number of reading data retries")

}

func serviceDeps() service.Deps {
	deps := service.Deps{
		InfluxClient: influx(),
	}
	return deps
}

func serviceConfig() service.Config {
	return config
}

func influx() client.Client {
	host, err := url.Parse(fmt.Sprintf("http://%s:%d", config.InfluxHost, config.InfluxPort))
	if err != nil {
		glog.Fatal(err)
	}
	influxConfig := client.Config{
		URL:      *host,
		Username: config.InfluxUser,
		Password: config.InfluxPassword,
	}
	influxClient, err := client.NewClient(influxConfig)
	if err != nil {
		glog.Fatal(err)
	}
	return *influxClient
}
