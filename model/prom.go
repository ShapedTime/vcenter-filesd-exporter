package model

import (
	"log"
	vcenterhelper "vcenter-prom-filesd-go/vcenter-helper"
)

type SDOutput struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

type PromModel struct {
	vc         *vcenterhelper.VCenterHelper
	datacenter string
}

var (
	//promPs = [...]string{"name", "guest.ipAddress", "guest.guestFamily", "guest.hostName", "guest.guestFullName",
	//	"summary.config", "guest.disk"}

	promPs = [...]string{"name", "guest.ipAddress", "guest.guestFamily", "guest.hostName", "guest.guestFullName"}
)

func NewPromModel(vc *vcenterhelper.VCenterHelper, datacenter string) *PromModel {
	return &PromModel{vc: vc, datacenter: datacenter}
}

func (s PromModel) GetVMs(path string) ([]SDOutput, error) {
	f, err := s.vc.FindFolder(s.datacenter, path)
	if err != nil {
		return nil, err
	}

	vms, err := s.vc.GetVMs(f, promPs[:])
	if err != nil {
		return nil, err
	}

	r := make([]SDOutput, len(vms))

	for i, vmRef := range vms {
		if vmRef.Guest == nil {
			log.Printf("vmRef.Guest is nil for %s", vmRef.Name)
			continue
		}

		//log.Println("name:", vmRef.Name, "ip:", vmRef.Guest.IpAddress,
		//	"guestfamily:", vmRef.Guest.GuestFamily, "hostname:", vmRef.Guest.HostName,
		//	"fullname", vmRef.Guest.GuestFullName, "memory", vmRef.Summary.Config.MemorySizeMB, "cpu", vmRef.Summary.Config.CpuReservation,
		//	"disk", vmRef.Guest.Disk[0].Capacity)

		log.Println("name:", vmRef.Name, "ip:", vmRef.Guest.IpAddress,
			"guestfamily:", vmRef.Guest.GuestFamily, "hostname:", vmRef.Guest.HostName,
			"fullname", vmRef.Guest.GuestFullName)

		r[i] = SDOutput{Targets: []string{vmRef.Guest.IpAddress},
			Labels: map[string]string{
				"name":          vmRef.Name,
				"id":            vmRef.Reference().String(),
				"guestfamily":   vmRef.Guest.GuestFamily,
				"hostname":      vmRef.Guest.HostName,
				"guestfullname": vmRef.Guest.GuestFullName,
			}}
	}

	return r, err
}
