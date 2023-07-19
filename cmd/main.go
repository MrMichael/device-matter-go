// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2023 Jiangxing Intelligence Ltd
//
// SPDX-License-Identifier: Apache-2.0

// Package driver this package provides an GPIO implementation of
// ProtocolDriver interface.
//
package main

import (
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/startup"

	"github.com/edgexfoundry/device-matter-go"
	"github.com/edgexfoundry/device-matter-go/internal/driver"
)

const (
	serviceName string = "device-matter"
)

func main() {
	d := driver.Driver{}
	startup.Bootstrap(serviceName, device_matter.Version, &d)
}
