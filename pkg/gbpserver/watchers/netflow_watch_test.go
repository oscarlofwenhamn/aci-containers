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
	"github.com/davecgh/go-spew/spew"
	"github.com/noironetworks/aci-containers/pkg/gbpserver"
	netflowpolicy "github.com/noironetworks/aci-containers/pkg/netflowpolicy/apis/aci.netflow/v1alpha"
	log "github.com/sirupsen/logrus"
	"reflect"
	"testing"
	"time"
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

func (s *netflow_suite) setup() {
	gCfg := &gbpserver.GBPServerConfig{}
	gCfg.GRPCPort = 19999
	gCfg.ProxyListenPort = 8899
	gCfg.PodSubnet = "10.2.56.1/21"
	gCfg.NodeSubnet = "1.100.201.0/24"
	gCfg.AciPolicyTenant = "defaultTenant"
	log := log.WithField("mod", "test")
	s.s = gbpserver.NewServer(gCfg)
	//s.nf = NewNetflowWatcher(s.s)

	s.nf = &NetflowWatcher{
		log: log,
		gs:  s.s,
		idb: newIntentDB(s.s, log),
	}
}

func (s *netflow_suite) expectMsg(op int, msg interface{}) error {
	gotOp, gotMsg, err := s.s.UTReadMsg(200 * time.Millisecond)
	if err != nil {
		return err
	}

	if gotOp != op {
		return fmt.Errorf("Exp op: %d, got: %d", op, gotOp)
	}

	if !reflect.DeepEqual(msg, gotMsg) {
		spew.Dump(msg)
		spew.Dump(gotMsg)
		return fmt.Errorf("msgs don't match")
	}

	return nil
}

func (s *netflow_suite) expectOp(op int) error {
	gotOp, _, err := s.s.UTReadMsg(200 * time.Millisecond)
	if err != nil {
		return err
	}

	if gotOp != op {
		return fmt.Errorf("Exp op: %d, got: %d", op, gotOp)
	}

	return nil
}

var netflow_policy = netflowpolicy.NetflowType{
	DstAddr:           "1.1.1.1",
	DstPort:           2055,
	Version:           "netflow",
	ActiveFlowTimeOut: 60,
	IdleFlowTimeOut:   15,
	SamplingRate:      5,
}

func TestNetflowGBP(t *testing.T) {
	ns := &netflow_suite{}
	ns.setup()

	netflowMO := &netflowCRD{
		DstAddr:           "1.1.1.1",
		DstPort:           2055,
		Version:           "netflow",
		ActiveFlowTimeOut: 60,
		IdleFlowTimeOut:   15,
		SamplingRate:      5,
	}
	ns.s.AddGBPCustomMo(netflowMO)

	ns.s.DelGBPCustomMo(netflowMO)

	ns.nf.netflowAdded(&netflowpolicy.NetflowPolicySpec{FlowSamplingPolicy: netflow_policy})
	err := ns.expectOp(gbpserver.OpaddGBPCustomMo)
	if err != nil {
		t.Error(err)
	}
	ns.nf.netflowDeleted(&netflowpolicy.NetflowPolicySpec{FlowSamplingPolicy: netflow_policy})
	err = ns.expectOp(gbpserver.OpdelGBPCustomMo)
	if err != nil {
		t.Error(err)
	}

}
