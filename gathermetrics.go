package main

import (
	"reflect"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

func Cleanup() {
	if r := recover(); r != nil {
		log.Debug("recovered in cleanup: ", r)
	}
}

func CountThreads(c chan<- map[string]int) {
	defer Cleanup()
	m := make(map[string]int)
	threads, err := n.GetThreads()
	if err != nil {
		panic(err)
	}
	m[n.BaseURI] = threads.NumberOfThreads
	log.Debug("CountThreads done")
	c <- m
}

var rememberDate time.Time = time.Time{}
var value string = "2006-01-02 15:04:05"
var (
	rollbackmap = make(map[string]int)
)

func Rollbacks(c chan<- map[string]int) {
	defer Cleanup()
	rollbacks, err := n.GetRollbacks()
	if err != nil {
		panic(err)
	}
	for _, rollback := range rollbacks.File {
		date, _ := time.Parse(value, rollback.Date)
		if date.After(rememberDate) {
			rollbackmap[rollback.Creator] += 1
		}
	}
	rememberDate, _ = time.Parse(value, rollbacks.File[0].Date)
	log.Debug("Rollback done")
	c <- rollbackmap
}

type UtilisationResults struct {
	Id         int
	Success    bool
	InputLoad  int
	OutputLoad int
	Time       string
}

func CountServices(c chan<- map[string]int) {
	serviceData, err := n.GetService()
	if err != nil {
		panic(err)
	}
	return_map := make(map[string]int)
	services := cfg.NSO.Services
	for key, value := range serviceData {
		for _, service := range services {
			serviceName := strings.Split(key, ":")
			if len(serviceName) < 2 {
				continue
			}
			if serviceName[1] == service {
				v := reflect.ValueOf(value)
				if v.Kind() != reflect.Map {
					continue
				}
				for _, label := range v.MapKeys() {
					labelValuesReflect := v.MapIndex(label)
					labelValues := labelValuesReflect.Elem().Interface().([]interface{})
					return_map[label.String()] = len(labelValues)
				}
			}
		}
	}
	log.Debug("CountServices done")
	c <- return_map
}

func Neds() {
	defer Cleanup()
	neds, err := n.GetNeds()
	if err != nil {
		panic(err)
	}
	nedSyncMap = make(map[string]int)
	for _, devices := range neds.TailfNcsDevice {
		v := reflect.ValueOf(devices.NedSettings)
		if v.Kind() == reflect.Map {
			for _, key := range v.MapKeys() {
				keyValue := key.Interface().(string)
				nedSyncMap[keyValue] += 1
			}
		}
	}
	log.Debug("Neds done")
}

func DeviceSyncs() {
	defer Cleanup()
	devices, err := n.GetDeviceSyncs()
	if err != nil {
		panic(err)
	}
	for _, device := range devices.Devices {
		result := 0
		if device.Result == "in-sync" {
			result = 1
		}
		deviceSyncMap[device.Device] = result
	}
	log.Debug("Devicesync done")
}
