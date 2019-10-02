/*
w
==================================================================================
  Copyright (c) 2019 AT&T Intellectual Property.
  Copyright (c) 2019 Nokia

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
==================================================================================
*/
/*
  Mnemonic:	sbi.go
  Abstract:	Contains SBI (SouthBound Interface) module definitions and generic SBI components
  Date:		16 March 2019
*/

package sbi

import (
	"errors"
	"routing-manager/pkg/rtmgr"
	"strconv"
)

const DefaultNngPipelineSocketPrefix = "tcp://"
const DefaultNngPipelineSocketNumber = 4561
const PlatformType = "platform"

var (
	SupportedSbis = []*EngineConfig{
		{
			Name:        "nngpush",
			Version:     "v1",
			Protocol:    "nngpipeline",
			Instance:    NewNngPush(),
			IsAvailable: true,
		},
	}
)

func GetSbi(sbiName string) (Engine, error) {
	for _, sbi := range SupportedSbis {
		if sbi.Name == sbiName && sbi.IsAvailable {
			return sbi.Instance, nil
		}
	}
	return nil, errors.New("SBI:" + sbiName + " is not supported or still not available")
}

type Sbi struct {
}

func (s *Sbi) pruneEndpointList(sbi Engine) {
	for _, ep := range rtmgr.Eps {
		if !ep.Keepalive {
			rtmgr.Logger.Debug("deleting %v", ep)
			sbi.DeleteEndpoint(ep)
			delete(rtmgr.Eps, ep.Uuid)
		} else {
			rtmgr.Eps[ep.Uuid].Keepalive = false
		}
	}
}

func (s *Sbi) updateEndpoints(rcs *rtmgr.RicComponents, sbi Engine) {
	for _, xapp := range (*rcs).XApps {
		for _, instance := range xapp.Instances {
			uuid := instance.Ip + ":" + strconv.Itoa(int(instance.Port))
			if _, ok := rtmgr.Eps[uuid]; ok {
				rtmgr.Eps[uuid].Keepalive = true
			} else {
				ep := &rtmgr.Endpoint{
					Uuid:       uuid,
					Name:       instance.Name,
					XAppType:   xapp.Name,
					Ip:         instance.Ip,
					Port:       instance.Port,
					TxMessages: instance.TxMessages,
					RxMessages: instance.RxMessages,
					Socket:     nil,
					IsReady:    false,
					Keepalive:  true,
				}
				if err := sbi.AddEndpoint(ep); err != nil {
					rtmgr.Logger.Error("can't create socket for endpoint: " + ep.Name + " due to:" + err.Error())
					continue
				}
				rtmgr.Eps[uuid] = ep
			}
		}
	}
	s.updatePlatformEndpoints(&((*rcs).Pcs), sbi)
	s.pruneEndpointList(sbi)
}

func (s *Sbi) updatePlatformEndpoints(pcs *rtmgr.PlatformComponents, sbi Engine) {
	rtmgr.Logger.Debug("updatePlatformEndpoints invoked. PCS: %v", *pcs)
	for _, pc := range *pcs {
		uuid := pc.Fqdn + ":" + strconv.Itoa(int(pc.Port))
		if _, ok := rtmgr.Eps[uuid]; ok {
			rtmgr.Eps[uuid].Keepalive = true
		} else {
			ep := &rtmgr.Endpoint{
				Uuid:       uuid,
				Name:       pc.Name,
				XAppType:   PlatformType,
				Ip:         pc.Fqdn,
				Port:       pc.Port,
				TxMessages: rtmgr.PLATFORMMESSAGETYPES[pc.Name]["tx"],
				RxMessages: rtmgr.PLATFORMMESSAGETYPES[pc.Name]["rx"],
				Socket:     nil,
				IsReady:    false,
				Keepalive:  true,
			}
			rtmgr.Logger.Debug("ep created: %v", ep)
			if err := sbi.AddEndpoint(ep); err != nil {
				rtmgr.Logger.Error("can't create socket for endpoint: " + ep.Name + " due to:" + err.Error())
				continue
			}
			rtmgr.Eps[uuid] = ep
		}
	}
}
