package main

import (
	"github.com/n0rad/go-erlog/logs"
	"os"
)

func (aci *Aci) Clean() {
	logs.WithF(aci.fields).Debug("Cleaning")

	aci.checkCompatibilityVersions()
	aci.checkLatestVersions()

	if err := os.RemoveAll(aci.target + "/"); err != nil {
		logs.WithEF(err, aci.fields).WithField("dir", aci.target).Warn("Cannot remove directory")
	}
}
