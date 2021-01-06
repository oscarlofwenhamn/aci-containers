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
	"github.com/sirupsen/logrus"
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
	DstAddr           string
	DstPort           int
	Version           string
	ActiveFlowTimeOut int
	IdleFlowTimeOut   int
	SamplingRate      int
	Name              string
	gs                *gbpserver.Server
}

type NetflowWatcher struct {
	log *logrus.Entry
	gs  *gbpserver.Server
	idb *intentDB
	rc  restclient.Interface
}

func NewNetflowWatcher(gs *gbpserver.Server) (*NetflowWatcher, error) {
	level, err := logrus.ParseLevel(gs.Config().WatchLogLevel)
	if err != nil {
		panic(err.Error())
	}
	logger := logrus.New()
	logger.Level = level
	log := logger.WithField("mod", "K8S-W")
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

	NetflowLw := cache.NewListWatchFromClient(nfw.rc, "netflowpolicies", sysNs, fields.Everything())
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
		DstAddr: netflow.Spec.FlowSamplingPolicy.DstAddr,
		DstPort: netflow.Spec.FlowSamplingPolicy.DstPort,
		Version: netflow.Spec.FlowSamplingPolicy.Version,
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
		DstAddr: netflow.Spec.FlowSamplingPolicy.DstAddr,
		DstPort: netflow.Spec.FlowSamplingPolicy.DstPort,
		Version: netflow.Spec.FlowSamplingPolicy.Version,
	}
	nfw.gs.DelGBPCustomMo(netflowMO)
}

func (nf *netflowCRD) Subject() string {
	return netflowCRDSub
}

func (nf *netflowCRD) URI() string {
	return nf.gs.GetURIBySubject("NetflowExporterConfig")
}

func (nf *netflowCRD) ParentSub() string {
	return netflowCRDParentSub
}

func (nf *netflowCRD) ParentURI() string {
	return nf.gs.GetURIBySubject("PlatformConfig")
}

func (nf *netflowCRD) Properties() map[string]interface{} {
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

func (nf *netflowCRD) Children() []string {
	return []string{}
}
