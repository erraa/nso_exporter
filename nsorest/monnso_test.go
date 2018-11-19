package nsorest

import (
	"testing"

	"git.is.comhem.com/ersa20/nso_exporter/config"
)

var (
	cfg = config.GetConfig()
	n   = NSO{
		Username: cfg.NSO.Username,
		Password: cfg.NSO.Password,
		BaseURI:  cfg.NSO.BaseURI,
	}
)

func TestServices(t *testing.T) {
	_, err := n.GetService()
	if err != nil {
		t.Errorf("Error:", err)
	}
}
