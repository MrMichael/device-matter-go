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
	"fmt"
	"os/exec"
	"context"
	"time"
	"bytes"
	"strconv"
)


func (s *Driver) chipToolSendCMD(params []string) error {

	// 要执行的二进制文件和命令行参数
	executable := "./chip-tool"
	args := params
	fmt.Println("args:", args)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable, args...)
	out, err := cmd.CombinedOutput()

	if ctx.Err() == nil {
		if err == nil {
			// 命令能执行，且结果正常
			cmd2 := exec.Command("grep", "-o", "CHIP:.*")
			cmd2.Stdin = bytes.NewReader(out)
			out2, _ := cmd2.Output()
			cmd3 := exec.Command("sed", "s/ //g")
			cmd3.Stdin = bytes.NewReader(out2)
			out3, _ := cmd3.Output()

			s.cmdResp = string(out3)
		} else {
			// 命令能执行，但结果异常（参数不对等）
			cmd3 := exec.Command("sed", "s/ //g")
			cmd3.Stdin = bytes.NewReader(out)
			out3, _ := cmd3.Output()
			s.cmdResp = string(out3)
		}
	} else {
		s.cmdResp = "Exec failure"
		return ctx.Err()
	}
	return nil
}


func (s *Driver) chipToolParseSetupPayload(payload string) error {

	// 要执行的二进制文件和命令行参数
	executable := "./chip-tool"
	args := []string{"payload","parse-setup-payload", payload}
	fmt.Println("args:", args)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable, args...)
	out, err := cmd.CombinedOutput()

	if ctx.Err() == nil {
		if err == nil {
			// 命令能执行，且结果正常
			cmd2 := exec.Command("grep", "-oE", "VendorID.*|ProductID.*|Discovery.*|Long.*|Passcode.*")
			cmd2.Stdin = bytes.NewReader(out)
			out2, _ := cmd2.Output()

			s.parseSetupPayloadResp = string(out2)
			fmt.Println("response:", string(out2))
		} else {
			// 命令能执行，但结果异常（参数不对等）
			s.parseSetupPayloadResp = "Result failure"
			return err
		}
	} else {
		s.parseSetupPayloadResp = "Exec failure"
		return ctx.Err()
	}
	return nil
}


func (s *Driver) chipToolCommissionIntoWiFiOverBT(node_id string, ssid string, password string, payload string, timeout string) error {

	// 要执行的二进制文件和命令行参数
	executable := "./chip-tool"
	args := []string{"pairing","code-wifi", node_id, ssid, password, payload, "--paa-trust-store-path", "credentials/production/paa-root-certs/"}
	fmt.Println("args:", args)

	tmp, _ := strconv.Atoi(timeout)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(tmp)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable, args...)
	out, err := cmd.CombinedOutput()

	if ctx.Err() == nil {
		if err == nil {
			// 命令能执行，且结果正常
			cmd2 := exec.Command("grep", "-o", "Received Command Response Data.*")
			cmd2.Stdin = bytes.NewReader(out)
			out2, _ := cmd2.Output()

			s.commissionIntoWiFiOverBTResp = string(out2)
			fmt.Println("response:", string(out2))
		} else {
			// 命令能执行，但结果异常（参数不对等）
			s.commissionIntoWiFiOverBTResp = "result failure"
			return err
		}
	} else {
		s.commissionIntoWiFiOverBTResp = "Already networked or exec failure"
		return ctx.Err()
	}
	return nil
}


func (s *Driver) chipToolCommissionWithQRCode(node_id string, payload string) error {

	// 要执行的二进制文件和命令行参数
	executable := "./chip-tool"
	args := []string{"pairing","code", node_id, payload}
	fmt.Println("args:", args)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable, args...)
	out, err := cmd.CombinedOutput()

	if ctx.Err() == nil {
		if err == nil {
			// 命令能执行，且结果正常
			cmd2 := exec.Command("grep", "-o", "Received Command Response Data.*")
			cmd2.Stdin = bytes.NewReader(out)
			out2, _ := cmd2.Output()

			s.commissionWithQRCodeTResp = string(out2)
			fmt.Println("response:", string(out2))
		} else {
			// 命令能执行，但结果异常（参数不对等）
			s.commissionWithQRCodeTResp = "result failure"
			return err
		}
	} else {
		s.commissionWithQRCodeTResp = "Already networked or exec failure"
		return ctx.Err()
	}
	return nil
}

func (s *Driver) chipToolUnpair(node_id string) error {

	// 要执行的二进制文件和命令行参数
	executable := "./chip-tool"
	args := []string{"pairing","unpair", node_id}
	fmt.Println("args:", args)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable, args...)
	out, err := cmd.CombinedOutput()

	if ctx.Err() == nil {
		if err == nil {
			// 命令能执行，且结果正常
			cmd2 := exec.Command("grep", "-o", "Received Command Response Data.*")
			cmd2.Stdin = bytes.NewReader(out)
			out2, _ := cmd2.Output()

			s.removeFromFabricResp = string(out2)
			fmt.Println("response:", string(out2))
		} else {
			// 命令能执行，但结果异常（参数不对等）
			s.removeFromFabricResp = "Result failure"
			return err
		}
	} else {
		s.removeFromFabricResp = "Not connected or exec failure"
		return ctx.Err()
	}
	return nil
}


func (s *Driver) chipToolCluster_On_Off_Toggle(node_id string, endpoint_id string, on_off_toggle string) error {

	// 要执行的二进制文件和命令行参数
	executable := "./chip-tool"
	var action string = "on"
	if on_off_toggle == "off" {
		action = "off"
	} else if on_off_toggle == "on" {
		action = "on"
	} else if on_off_toggle == "toggle" {
		action = "toggle"
	}
	args := []string{"onoff", action, node_id, endpoint_id}
	fmt.Println("args:", args)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable, args...)
	out, err := cmd.CombinedOutput()

	if ctx.Err() == nil {
		if err == nil {
			// 命令能执行，且结果正常
			cmd2 := exec.Command("grep", "-o", "CHIP:DMG.*")
			cmd2.Stdin = bytes.NewReader(out)
			out2, _ := cmd2.Output()
			cmd3 := exec.Command("grep", "-oE", "EndpointId =.*|ClusterId =.*|CommandId =.*|status =.*")
			cmd3.Stdin = bytes.NewReader(out2)
			out3, _ := cmd3.Output()
			fmt.Println("response:", string(out3))

			s.cluster_On_Off_Toggle_Resp = string(out3)
		} else {
			// 命令能执行，但结果异常（参数不对等）
			s.cluster_On_Off_Toggle_Resp = "Result failure"
			return err
		}
	} else {
		s.cluster_On_Off_Toggle_Resp = "Exec failure"
		return ctx.Err()
	}
	return nil
}

func (s *Driver) chipToolCluster_Read_On_Off(node_id string, endpoint_id string) error {

	// 要执行的二进制文件和命令行参数
	executable := "./chip-tool"
	args := []string{"onoff","read", "on-off", node_id, endpoint_id}
	fmt.Println("args:", args)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable, args...)
	out, err := cmd.CombinedOutput()

	if ctx.Err() == nil {
		if err == nil {
			// 命令能执行，且结果正常
			cmd2 := exec.Command("grep", "-o", "CHIP:DMG.*")
			cmd2.Stdin = bytes.NewReader(out)
			out2, _ := cmd2.Output()
			cmd3 := exec.Command("grep", "-oE", "Endpoint =.*|Cluster =.*|Attribute =.*|Data =.*")
			cmd3.Stdin = bytes.NewReader(out2)
			out3, _ := cmd3.Output()
			fmt.Println("response:", string(out3))

			s.cluster_Read_On_Off_Resp = string(out3)
		} else {
			// 命令能执行，但结果异常（参数不对等）
			s.cluster_Read_On_Off_Resp = "Result failure"
			return err
		}
	} else {
		s.cluster_Read_On_Off_Resp = "Exec failure"
		return ctx.Err()
	}
	return nil
}