/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package csi_common

import (
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

type CSIDriver struct {
	name    string
	nodeID  string
	version *csi.Version
	supVers []*csi.Version
	cap     []*csi.ControllerServiceCapability
	vc      []*csi.VolumeCapability_AccessMode
}

func NewCSIDriver(name string, v *csi.Version, vers []*csi.Version, nodeID string) *CSIDriver {
	driver := CSIDriver{
		name:    name,
		version: v,
		supVers: vers,
		nodeID:  nodeID,
	}
	return &driver
}

func (d *CSIDriver) CheckVersion(v *csi.Version) error {
	// Assumes always backward compatible
	for _, sv := range d.supVers {
		if v.Major == sv.Major && v.Minor <= sv.Minor {
			return nil
		}
	}

	return status.Error(codes.InvalidArgument, "Unsupported version: %s"+GetVersionString(v))
}

func (d *CSIDriver) ValidateRequest(v *csi.Version, c csi.ControllerServiceCapability_RPC_Type) error {
	if v == nil {
		return status.Error(codes.InvalidArgument, "Version not specified")
	}

	if err := d.CheckVersion(v); err != nil {
		return status.Error(codes.InvalidArgument, "Unsupported version")
	}

	if c == csi.ControllerServiceCapability_RPC_UNKNOWN {
		return nil
	}

	for _, cap := range d.cap {
		if c == cap.GetRpc().GetType() {
			return nil
		}
	}

	return status.Error(codes.InvalidArgument, "Unsupported version: %s"+GetVersionString(v))
}

func (d *CSIDriver) AddControllerServiceCapabilities(cl []csi.ControllerServiceCapability_RPC_Type) {
	var csc []*csi.ControllerServiceCapability

	for _, c := range cl {
		glog.Infof("Enabling controller service capability: %v", c.String())
		csc = append(csc, NewControllerServiceCapability(c))
	}

	d.cap = csc

	return
}

func (d *CSIDriver) AddVolumeCapabilityAccessModes(vc []csi.VolumeCapability_AccessMode_Mode) []*csi.VolumeCapability_AccessMode {
	var vca []*csi.VolumeCapability_AccessMode
	for _, c := range vc {
		glog.Infof("Enabling volume access mode: %v", c.String())
		vca = append(vca, NewVolumeCapabilityAccessMode(c))
	}
	return vca
}
