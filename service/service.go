package service

import (
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/influxdb/influxdb/client"
)

type Config struct {
	TickerTime      time.Duration
	DryRun          bool
	OWMcityID       int
	InfluxPort      int
	InfluxDB        string
	InfluxUser      string
	InfluxPassword  string
	InfluxHost      string
	InfluxRetention string
	DhtType         string
	DhtPin          int
	DHTPerf         bool
	DhtRetries      int
}

type Deps struct {
	InfluxClient client.Client
}

type Service struct {
	Config
	Deps
	ticker     *time.Ticker
	shutdownWG sync.WaitGroup
}

func New(cfg Config, deps Deps) *Service {
	return &Service{
		Config: cfg,
		ticker: time.NewTicker(cfg.TickerTime),
		Deps:   deps,
	}
}

func (srv *Service) Close() {
	srv.ticker.Stop()
	srv.shutdownWG.Wait()
}

func (srv *Service) Serve() {
	srv.Reconcile()
	for range srv.ticker.C {
		srv.Reconcile()
	}
}

func (srv *Service) Reconcile() {
	glog.Info("--- measuring ---")
	srv.shutdownWG.Add(1)
	defer srv.shutdownWG.Done()
	// get temp/humid data and write to influx
	si := measure(srv)
	so := outside(srv)
	writeData(si, so, srv)
}
