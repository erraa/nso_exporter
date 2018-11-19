package nsorest

import (
	"encoding/json"
	"fmt"
)

type TailfNcsMonitoringSmp struct {
	NumberOfThreads int `json:"number-of-threads"`
}

func (nso *NSO) GetThreads() (TailfNcsMonitoringSmp, error) {
	m := make(map[string]TailfNcsMonitoringSmp)
	err, resp := nso.get("operational/ncs-state/smp")
	defer resp.Body.Close()
	if err != nil {
		return TailfNcsMonitoringSmp{}, err
	}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return TailfNcsMonitoringSmp{}, err
	}
	return m["tailf-ncs-monitoring:smp"], nil
}

type Rollbacks struct {
	File []struct {
		Name    string `json:"name"`
		Creator string `json:"creator"`
		Date    string `json:"date"`
		Via     string `json:"via"`
		Label   string `json:"label"`
		Comment string `json:"comment"`
	} `json:"file"`
}

func (nso *NSO) GetRollbacks() (Rollbacks, error) {
	m := make(map[string]Rollbacks)
	err, resp := nso.get("rollbacks/")
	defer resp.Body.Close()
	if err != nil {
		return Rollbacks{}, err
	}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return Rollbacks{}, err
	}
	return m["rollbacks"], nil
}

// Note we are umarshaling sync-result into devices, cus thats what is rly is
type Syncs struct {
	Devices []struct {
		Device string `json:"device"`
		Result string `json:"result"`
		Info   string `json:"info,omitempty"`
	} `json:"sync-result"`
}

// This func really does a post to check all the syncs
// NSO really has to check every device if it is in sync, this takes a bit of time
func (nso *NSO) GetDeviceSyncs() (Syncs, error) {
	m := make(map[string]Syncs)
	err, resp := nso.send("POST", "running/devices/_operations/check-sync", nil)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		return Syncs{}, err
	}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return Syncs{}, err
	}
	// probably needs to change
	return m["tailf-ncs:output"], nil
}

type Neds struct {
	TailfNcsDevice []struct {
		Name        string      `json:"name"`
		NedSettings interface{} `json:"ned-settings"`
	} `json:"tailf-ncs:device"`
}

func (nso *NSO) GetNeds() (Neds, error) {
	m := make(map[string]Neds)
	err, resp := nso.get("operational/devices/device?select=name;ned-settings(*)")
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		return Neds{}, err
	}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return Neds{}, err
	}
	// probably needs to change
	return m["collection"], nil
}

func (nso *NSO) GetService() (map[string]interface{}, error) {
	err, resp := nso.get("running/")
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, err
	}
	// TEMP unmarshaling for testing only

	return m, nil
}
