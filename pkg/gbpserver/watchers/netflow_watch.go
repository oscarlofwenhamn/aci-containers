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
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"github.com/noironetworks/aci-containers/pkg/gbpserver"
	netflowpolicy "github.com/noironetworks/aci-containers/pkg/netflowpolicy/apis/aci.netflow/v1alpha"
	netflowclientset "github.com/noironetworks/aci-containers/pkg/netflowpolicy/clientset/versioned"
)

const (
	netflowCRDSub       = "NetflowExporterConfig"
	netflowCRDParentSub = "PlatformConfig"
)

type netflowCRD struct {
	DstAddr           string `json:"destIp"`
	DstPort           int    `json:"destPort"`
	Version           string `json:"flowType,omitempty"`
	ActiveFlowTimeOut int    `json:"activeFlowTimeOut,omitempty"`
	IdleFlowTimeOut   int    `json:"idleFlowTimeOut,omitempty"`
	SamplingRate      int    `json:"samplingRate,omitempty"`
	Name              string `json:"name,omitempty"`
}

type NetflowWatcher struct {
	log *log.Entry
	gs  *gbpserver.Server
	idb *intentDB
	rc  restclient.Interface
}

func NewNetflowWatcher(gs *gbpserver.Server) (*NetflowWatcher, error) {
	gcfg := gs.Config()
	level, err := log.ParseLevel(gcfg.WatchLogLevel)
	if err != nil {
		panic(err.Error())
	}
	logger := log.New()
	logger.Level = level
	log := logger.WithField("mod", "NETFLOW-W")
	cfg, err := restclient.InClusterConfig()
	if err != nil {
		return nil, err
	}

	netflowclient, err := netflowclientset.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	restClient := netflowclient.AciV1alpha().RESTClient()
	return &NetflowWatcher{
		log: log,
		rc:  restClient,
		gs:  gs,
		idb: newIntentDB(gs, log),
	}, nil
}

func (nfw *NetflowWatcher) InitNetflowInformer(stopCh <-chan struct{}) error {
	nfw.watchNetflow(stopCh)
	return nil
}

func (nfw *NetflowWatcher) watchNetflow(stopCh <-chan struct{}) {

	NetflowLw := cache.NewListWatchFromClient(nfw.rc, "netflowpolicies", metav1.NamespaceAll, fields.Everything())
	_, netflowInformer := cache.NewInformer(NetflowLw, &netflowpolicy.NetflowPolicy{}, 0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				nfw.log.Infof("netflow added")
				nfw.netflowAdded(obj)
			},
			UpdateFunc: func(oldobj interface{}, newobj interface{}) {
				nfw.log.Infof("netflow updated")
				nfw.netflowAdded(newobj)
			},
			DeleteFunc: func(obj interface{}) {
				nfw.log.Infof("netflow deleted")
				nfw.netflowDeleted(obj)
			},
		})
	go netflowInformer.Run(stopCh)
}

func (nfw *NetflowWatcher) netflowAdded(obj interface{}) {
	netflow, ok := obj.(*netflowpolicy.NetflowPolicy)
	if !ok {
		nfw.log.Errorf("netflowAdded: Bad object type")
		return
	}

	nfw.log.Infof("netflowAdded - %s", netflow.ObjectMeta.Name)
	netflowMO := &netflowCRD{
		DstAddr:           netflow.Spec.FlowSamplingPolicy.DstAddr,
		DstPort:           netflow.Spec.FlowSamplingPolicy.DstPort,
		Version:           netflow.Spec.FlowSamplingPolicy.Version,
		ActiveFlowTimeOut: netflow.Spec.FlowSamplingPolicy.ActiveFlowTimeOut,
		IdleFlowTimeOut:   netflow.Spec.FlowSamplingPolicy.IdleFlowTimeOut,
		SamplingRate:      netflow.Spec.FlowSamplingPolicy.SamplingRate,
		Name:              netflow.ObjectMeta.Name,
	}
	nfw.gs.AddGBPCustomMo(netflowMO)
}

func (nfw *NetflowWatcher) netflowDeleted(obj interface{}) {
	netflow, ok := obj.(*netflowpolicy.NetflowPolicy)
	if !ok {
		nfw.log.Errorf("netflowDeleted: Bad object type")
		return
	}

	nfw.log.Infof("netflowDeleted - %s", netflow.ObjectMeta.Name)
	netflowMO := &netflowCRD{
		DstAddr:           netflow.Spec.FlowSamplingPolicy.DstAddr,
		DstPort:           netflow.Spec.FlowSamplingPolicy.DstPort,
		Version:           netflow.Spec.FlowSamplingPolicy.Version,
		ActiveFlowTimeOut: netflow.Spec.FlowSamplingPolicy.ActiveFlowTimeOut,
		IdleFlowTimeOut:   netflow.Spec.FlowSamplingPolicy.IdleFlowTimeOut,
		SamplingRate:      netflow.Spec.FlowSamplingPolicy.SamplingRate,
		Name:              netflow.ObjectMeta.Name,
	}
	nfw.gs.DelGBPCustomMo(netflowMO)
}

func (nf *netflowCRD) fillNetflowDefaults() {

	// setting default values if no values present for netflow attributes.
	if nf.Version == "netflow" {
		nf.Version = "v5"
	} else if nf.Version == "ipfix" {
		nf.Version = "v9"
	} else {
		nf.Version = "v5"
	}
	if nf.DstPort == 0 {
		nf.DstPort = 2055
	}
	if nf.ActiveFlowTimeOut == 0 {
		nf.ActiveFlowTimeOut = 60
	}
	if nf.IdleFlowTimeOut == 0 {
		nf.IdleFlowTimeOut = 15
	}
	if nf.SamplingRate == 0 {
		nf.SamplingRate = 0
	}
}

func (nf *netflowCRD) Subject() string {
	return netflowCRDSub
}

func (nf *netflowCRD) URI() string {
	gServer := &gbpserver.Server{}
	platformURI := gServer.GetPlatformURI()
	return fmt.Sprintf("%s%s/%s/", platformURI, netflowCRDSub, nf.Name)
}

func (nf *netflowCRD) Properties() map[string]interface{} {

	nf.fillNetflowDefaults()
	return map[string]interface{}{
		"dstAddr":           nf.DstAddr,
		"dstPort":           nf.DstPort,
		"version":           nf.Version,
		"activeFlowTimeOut": nf.ActiveFlowTimeOut,
		"idleFlowTimeOut":   nf.IdleFlowTimeOut,
		"samplingRate":      nf.SamplingRate,
		"name":              nf.Name,
	}
}

func (nf *netflowCRD) ParentSub() string {
	return netflowCRDParentSub
}

func (nf *netflowCRD) ParentURI() string {
	gServer := &gbpserver.Server{}
	nfParentURI := gServer.GetPlatformURI()
	return nfParentURI
}

func (nf *netflowCRD) Children() []string {
	return []string{}
}
