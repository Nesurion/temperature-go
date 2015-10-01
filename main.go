package main

import (
	"github.com/influxdb/influxdb/client"
	"github.com/ogier/pflag"

	"github.com/nesurion/temperature-go/service"
)

var (
	config = service.Config{}
)

func main() {
	srv := service.New(serviceConfig(), serviceDeps())

	if srv.DryRun {
		glog.Info("running in dryRun mode")
	}

	go srv.Serve()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	glog.Info("Received termination signal. Closing ...")
	srv.Close()

}

func init() {
	pflag.DurationVar(&config.TickerTime, "ticker-time", 10*time.Minutes, "Ticker time.")
	pflag.IntVar(&config.InfluxPort, "influx-port", 8086, "InfluxDB Port")
	pflag.StringVar(&config.InfluxHost, "influx-host", "localhost", "InfluxDB Port")
	pflag.StringVar(&config.InfluxUser, "influx-user", "", "InfluxDB User")
	pflag.StringVar(&config.InfluxPassword, "influx-password", "", "InfluxDB Password")
	pflag.BoolVar(&config.DryRun, "dry-run", true, "Write to STDOUT instead of InfluxDB")

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

func influx() (client.Client, error) {
	host, err := url.Parse(fmt.Sprintf("http://%s:%d", config.InfluxHost, config.InfluxPort))
	if err != nil {
		log.Fatal(err)
	}
	influxConfig := client.Config{
		URL:      *host,
		Username: config.InfluxUser,
		Password: config.InfluxPassword,
	}
	con, err := client.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}
	return con
}
