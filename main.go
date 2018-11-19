package main

import (
	"net/http"
	"sync"
	"time"

	"git.is.comhem.com/ersa20/nso_exporter/config"
	"git.is.comhem.com/ersa20/nso_exporter/nsorest"
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	namespace = "nso"
	cfg       = config.GetConfig()
	n         = nsorest.NSO{
		Username: cfg.NSO.Username,
		Password: cfg.NSO.Password,
		BaseURI:  cfg.NSO.BaseURI,
	}
	deviceSyncMap = make(map[string]int)
	nedSyncMap    = make(map[string]int)
)

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type NsoCollector struct {
	mutex                 sync.Mutex
	ServicesCountMetric   *prometheus.Desc
	ThreadCountMetric     *prometheus.Desc
	DeviceCountMetric     *prometheus.Desc
	RollbackCounterMetric *prometheus.Desc
	DeviceSyncs           *prometheus.Desc
	NedMetric             *prometheus.Desc
}

func newNsoCollector() *NsoCollector {
	return &NsoCollector{
		ServicesCountMetric: prometheus.NewDesc("nso_exporter_services_count",
			"Tracks the amount of services deployed",
			[]string{"service_name"}, nil,
		),
		ThreadCountMetric: prometheus.NewDesc("nso_exporter_thread_count",
			"Tracks the amount of threads NSO uses",
			[]string{"nso_instance"}, nil,
		),
		// Rollback metric is named commit count, but we get the data from
		// gathering the rollback posibilities in nso
		RollbackCounterMetric: prometheus.NewDesc("nso_exporter_commit_count",
			"Tracks the amount of commits",
			[]string{"username"}, nil,
		),
		DeviceSyncs: prometheus.NewDesc("nso_exporter_device_sync",
			"0 if device is out of sync, 1 if it is in sync. This is only polled once every 30 minutes in reality",
			[]string{"device"}, nil,
		),
		NedMetric: prometheus.NewDesc("nso_exporter_neds",
			"0 if device is out of sync, 1 if it is in sync",
			[]string{"ned_name"}, nil,
		),
		DeviceCountMetric: prometheus.NewDesc("nso_exporter_device_count",
			"Tracks the amount of devices in nso",
			nil, nil,
		),
	}
}

func (collector *NsoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.ServicesCountMetric
	ch <- collector.DeviceSyncs
	ch <- collector.ThreadCountMetric
	ch <- collector.RollbackCounterMetric
	ch <- collector.DeviceCountMetric
	ch <- collector.NedMetric
}

func (collector *NsoCollector) Collect(ch chan<- prometheus.Metric) {
	collector.mutex.Lock() // To protect metrics from concurrent collects
	defer collector.mutex.Unlock()

	stats := collect_stats()

	for k, v := range stats.Rollbacks {
		ch <- prometheus.MustNewConstMetric(
			collector.RollbackCounterMetric, prometheus.CounterValue, float64(v), k,
		)
	}
	for k, v := range nedSyncMap {
		ch <- prometheus.MustNewConstMetric(
			collector.NedMetric, prometheus.GaugeValue, float64(v), k,
		)
	}
	for k, v := range deviceSyncMap {
		ch <- prometheus.MustNewConstMetric(
			collector.DeviceSyncs, prometheus.GaugeValue, float64(v), k,
		)
	}
	for k, v := range stats.CountThreads {
		ch <- prometheus.MustNewConstMetric(
			collector.ThreadCountMetric, prometheus.GaugeValue, float64(v), k,
		)
	}
	for k, v := range stats.CountServices {
		ch <- prometheus.MustNewConstMetric(
			collector.ServicesCountMetric, prometheus.GaugeValue, float64(v), k,
		)
	}
}

type Stats struct {
	Rollbacks     map[string]int
	CountThreads  map[string]int
	CountServices map[string]int
	DeviceSyncs   map[string]int
	Neds          map[string]int
}

func collect_stats() Stats {
	// We will maybe need a waitgroup in the future
	s := Stats{}
	rollbacks := make(chan map[string]int)
	countThreads := make(chan map[string]int)
	countServices := make(chan map[string]int)
	log.Debug("Running Rollbacks")
	go Rollbacks(rollbacks)
	log.Debug("Running CountThreads")
	go CountThreads(countThreads)
	log.Debug("Running CountServices")
	go CountServices(countServices)
	s.Rollbacks = <-rollbacks
	s.CountThreads = <-countThreads
	s.CountServices = <-countServices
	return s
}

func init() {
	// Set loglevels
	switch loglevel := cfg.Log.Loglevel; loglevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	prometheus.MustRegister(newNsoCollector())
}

func main() {
	go func() {
		for {
			DeviceSyncs()
			Neds()
			time.Sleep(time.Minute * 30)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :", cfg.HTTP.Port)
	log.Fatal(http.ListenAndServe(cfg.HTTP.Ip+":"+cfg.HTTP.Port, nil))
}
