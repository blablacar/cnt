package main

import "github.com/n0rad/go-erlog/logs"

func (p *Pod) Test() error {
	logs.WithF(p.fields).Info("Testing")

	for _, e := range p.manifest.Pod.Apps {
		aci, err := NewAciWithManifest(p.path+"/"+e.Name, p.args, p.toAciManifest(e))
		if err != nil {
			logs.WithEF(err, p.fields).WithField("name", e.Name).Fatal("Cannot prepare aci")
		}
		aci.podName = &p.manifest.Name
		if err := aci.Test(); err != nil {
			return err
		}
	}
	return nil
}
