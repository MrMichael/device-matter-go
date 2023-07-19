// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2023 Jiangxing Intelligence Ltd
//
// SPDX-License-Identifier: Apache-2.0

// Package driver this package provides an GPIO implementation of
// ProtocolDriver interface.
//
package driver

import (
	// "errors"
	"fmt"
	// "time"

	sdkModels "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
)

type Driver struct {
	lc            logger.LoggingClient
	asyncCh       chan<- *sdkModels.AsyncValues

	cmdResp string
	parseSetupPayloadResp string
	commissionIntoWiFiOverBTResp string
	commissionWithQRCodeTResp string
	removeFromFabricResp string
	cluster_On_Off_Toggle_Resp string
	cluster_Read_On_Off_Resp string
}

// Initialize performs protocol-specific initialization for the device
// service.
func (s *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *sdkModels.AsyncValues, deviceCh chan<- []sdkModels.DiscoveredDevice) error {
	s.lc = lc
	s.asyncCh = asyncCh

	return nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *Driver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModels.CommandRequest) (res []*sdkModels.CommandValue, err error) {
	s.lc.Infof("Driver.HandleReadCommands: protocols: %v resource: %v attributes: %v", protocols, reqs[0].DeviceResourceName, reqs[0].Attributes)

	res = make([]*sdkModels.CommandValue, len(reqs))

	if len(reqs) == 1 {
		if reqs[0].DeviceResourceName == "cmdParams" {
			value := []string{s.cmdResp}
			cv, _ := sdkModels.NewCommandValue(reqs[0].DeviceResourceName, common.ValueTypeStringArray, value)
			res[0] = cv
		} else if reqs[0].DeviceResourceName == "qrcode_payload" {
			// ParseSetupPayload
			cv, _ := sdkModels.NewCommandValue(reqs[0].DeviceResourceName, common.ValueTypeString, s.parseSetupPayloadResp)
			res[0] = cv
		} else if reqs[0].DeviceResourceName == "node_id" {
			// RemoveFromFabric
			cv, _ := sdkModels.NewCommandValue(reqs[0].DeviceResourceName, common.ValueTypeString, s.removeFromFabricResp)
			res[0] = cv
		}
	} else if len(reqs) == 2 {
		if reqs[0].DeviceResourceName == "node_id" && reqs[1].DeviceResourceName == "qrcode_payload" {
			// CommissionWithQRCode
			cv, _ := sdkModels.NewCommandValue(reqs[0].DeviceResourceName, common.ValueTypeString, s.commissionWithQRCodeTResp)
			res[0] = cv
		} else if reqs[0].DeviceResourceName == "node_id" && reqs[1].DeviceResourceName == "endpoint_id" {
			// Cluster_Read_On_Off
			cv, _ := sdkModels.NewCommandValue(reqs[0].DeviceResourceName, common.ValueTypeString, s.cluster_Read_On_Off_Resp)
			res[0] = cv
		}
	} else if len(reqs) == 3 {
		if reqs[0].DeviceResourceName == "node_id" && reqs[1].DeviceResourceName == "endpoint_id" && reqs[2].DeviceResourceName == "on_off_toggle" {
			// Cluster_On_Off_Toggle
			cv, _ := sdkModels.NewCommandValue(reqs[0].DeviceResourceName, common.ValueTypeString, s.cluster_On_Off_Toggle_Resp)
			res[0] = cv
		} 
	} else if len(reqs) == 5 {
		if reqs[0].DeviceResourceName == "node_id" && reqs[1].DeviceResourceName == "ssid" && reqs[2].DeviceResourceName == "password" {
			// CommissionIntoWiFiOverBT
			cv, _ := sdkModels.NewCommandValue(reqs[0].DeviceResourceName, common.ValueTypeString, s.commissionIntoWiFiOverBTResp)
			res[0] = cv
		}
	}

	return res, nil
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource.
// Since the commands are actuation commands, params provide parameters for the individual
// command.
func (s *Driver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModels.CommandRequest,
	params []*sdkModels.CommandValue) error {

	var err error

	if len(reqs) == 1 {
		if	reqs[0].DeviceResourceName == "cmdParams" {
			var cmdParamsValue []string
			if cmdParamsValue, err = params[0].StringArrayValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			err = s.chipToolSendCMD(cmdParamsValue)
		} else if reqs[0].DeviceResourceName == "qrcode_payload" {
			// ParseSetupPayload
			var qrcode_payload string
			if qrcode_payload, err = params[0].StringValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			err = s.chipToolParseSetupPayload(qrcode_payload)
		} else if reqs[0].DeviceResourceName == "node_id" {
			// RemoveFromFabric
			var node_id string
			if node_id, err = params[0].StringValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			err = s.chipToolUnpair(node_id)
		}
	} else if len(reqs) == 2 {
		if reqs[0].DeviceResourceName == "node_id" && reqs[1].DeviceResourceName == "qrcode_payload" {
			// CommissionWithQRCode
			var node_id string
			var qrcode_payload string
			if node_id, err = params[0].StringValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			if qrcode_payload, err = params[1].StringValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			err = s.chipToolCommissionWithQRCode(node_id, qrcode_payload)
		} else if reqs[0].DeviceResourceName == "node_id" && reqs[1].DeviceResourceName == "endpoint_id" {
			// Cluster_Read_On_Off
			var node_id string
			var endpoint_id string
			if node_id, err = params[0].StringValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			if endpoint_id, err = params[1].StringValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			err = s.chipToolCluster_Read_On_Off(node_id, endpoint_id)
		}
	} else if len(reqs) == 3 {
		if reqs[0].DeviceResourceName == "node_id" && reqs[1].DeviceResourceName == "endpoint_id" && reqs[2].DeviceResourceName == "on_off_toggle" {
			// Cluster_On_Off_Toggle
			var node_id string
			var endpoint_id string
			var on_off_toggle string
			if node_id, err = params[0].StringValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			if endpoint_id, err = params[1].StringValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			if on_off_toggle, err = params[2].StringValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			err = s.chipToolCluster_On_Off_Toggle(node_id, endpoint_id, on_off_toggle)
		}
	} else if len(reqs) == 5 {
		if reqs[0].DeviceResourceName == "node_id" && reqs[1].DeviceResourceName == "ssid" && reqs[2].DeviceResourceName == "password" {
			// CommissionIntoWiFiOverBT
			fmt.Println("CommissionIntoWiFiOverBT")
			var node_id string
			var ssid string
			var password string
			var qrcode_payload string
			var timeout string
			if node_id, err = params[0].StringValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			if ssid, err = params[1].StringValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			if password, err = params[2].StringValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			if qrcode_payload, err = params[3].StringValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			if timeout, err = params[4].StringValue(); err != nil {
				err := fmt.Errorf("Driver.HandleWriteCommands; the data type of parameter should be string array, parameter: %s", params[0].String())
				return err
			}
			err = s.chipToolCommissionIntoWiFiOverBT(node_id, ssid, password, qrcode_payload, timeout)
		}
	}

	return err
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (s *Driver) Stop(force bool) error {
	// Then Logging Client might not be initialized
	if s.lc != nil {
		s.lc.Debugf("Driver.Stop called: force=%v", force)
	}
	return nil
}

// AddDevice is a callback function that is invoked
// when a new Device associated with this Device Service is added
func (s *Driver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	s.lc.Debugf("a new Device is added: %s", deviceName)
	return nil
}

// UpdateDevice is a callback function that is invoked
// when a Device associated with this Device Service is updated
func (s *Driver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	s.lc.Debugf("Device %s is updated", deviceName)
	return nil
}

// RemoveDevice is a callback function that is invoked
// when a Device associated with this Device Service is removed
func (s *Driver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	s.lc.Debugf("Device %s is removed", deviceName)
	return nil
}
