/***
Copyright 2021 Cisco Systems Inc. All rights reserved.

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

package watchers

import (
	"fmt"
	"github.com/noironetworks/aci-containers/pkg/gbpserver"
	log "github.com/sirupsen/logrus"
	"testing"
)

var uriTenant string

const (
	testnetflowCRDSub       = "NetflowExporterConfig"
	testnetflowCRDParentSub = ""
	testCRDParentUri        = ""
)

type netflow_suite struct {
	s  *gbpserver.Server
	nf *NetflowWatcher
}

type testnetflowCRD struct {
	DstAddr           string `json:"destIp"`
	DstPort           int    `json:"destPort"`
	Version           string `json:"flowType,omitempty"`
	ActiveFlowTimeOut int    `json:"activeFlowTimeOut,omitempty"`
	IdleFlowTimeOut   int    `json:"idleFlowTimeOut,omitempty"`
	SamplingRate      int    `json:"samplingRate,omitempty"`
	Name              string `json:"name,omitempty"`
}

func (s *netflow_suite) setup() {
	gCfg := &gbpserver.GBPServerConfig{}
	gCfg.GRPCPort = 19999
	gCfg.ProxyListenPort = 8899
	gCfg.PodSubnet = "10.2.56.1/21"
	gCfg.NodeSubnet = "1.100.201.0/24"
	gCfg.AciPolicyTenant = "defaultTenant"
	log := log.WithField("mod", "test")
	s.s = gbpserver.NewServer(gCfg)

	s.nf = &NetflowWatcher{
		log: log,
		gs:  s.s,
		idb: newIntentDB(s.s, log),
	}
}

func (nt *testnetflowCRD) Subject() string {
	return testnetflowCRDSub
}

func (nt *testnetflowCRD) URI() string {
	return fmt.Sprintf("%s%s/", uriTenant, testnetflowCRDSub)
}

func (nt *testnetflowCRD) Properties() map[string]interface{} {
	return map[string]interface{}{
		"dstAddr":           nt.DstAddr,
		"dstPort":           nt.DstPort,
		"version":           nt.Version,
		"activeFlowTimeOut": nt.ActiveFlowTimeOut,
		"idleFlowTimeOut":   nt.IdleFlowTimeOut,
		"samplingRate":      nt.SamplingRate,
		"name":              nt.Name,
	}
}

func (nt *testnetflowCRD) ParentSub() string {
	return testnetflowCRDParentSub
}

func (nt *testnetflowCRD) ParentURI() string {
	return testCRDParentUri
}

func (nt *testnetflowCRD) Children() []string {
	return []string{}
}

func TestNetflowGBP(t *testing.T) {
	ns := &netflow_suite{}
	ns.setup()

	netflowcrd := &testnetflowCRD{
		DstAddr: "1.1.1.1",
		DstPort: 2055,
		Version: "netflow",
	}
	ns.s.AddGBPCustomMo(netflowcrd)

	ns.s.DelGBPCustomMo(netflowcrd)

}
