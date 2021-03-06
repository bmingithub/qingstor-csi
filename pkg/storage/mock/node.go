/*
Copyright (C) 2018 Yunify, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this work except in compliance with the License.
You may obtain a copy of the License in the LICENSE file, or at:

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mock

import (
	"errors"
	"github.com/yunify/qingstor-csi/pkg/common"
	"time"
)

//var deviceNo = 50

func (p *mockStorageProvider) NodeAttachVolume(volumeID string) error {
	_, ok := p.attachedVolumes[volumeID]
	if ok {
		return errors.New("volume already attached")
	}
	vol, err := p.FindVolume(volumeID)
	if err != nil {
		return err
	}
	//deviceNo ++
	p.attachedVolumes[volumeID] = &attachVolume{
		vol:    vol,
		device: common.GenerateHashInEightBytes(time.Now().UTC().String()),
	}
	return nil
}

func (p *mockStorageProvider) NodeDetachVolume(volumeID string) error {
	_, ok := p.attachedVolumes[volumeID]
	if !ok {
		return errors.New("volume not attached")
	}
	delete(p.attachedVolumes, volumeID)
	return nil
}

func (p *mockStorageProvider) NodeGetDevice(volumeID string) (string, error) {
	attachVol, ok := p.attachedVolumes[volumeID]
	if ok && attachVol != nil {
		return attachVol.device, nil
	}
	return "", errors.New("volume not found")
}
