// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"encoding/json"
	"encoding/xml"
	"sort"

	"github.com/syncthing/syncthing/lib/util"
)

// VersioningConfiguration is used in the code and for JSON serialization
type VersioningConfiguration struct {
	Type             string            `json:"type"`
	Params           map[string]string `json:"params"`
	CleanupIntervalS int               `json:"cleanupIntervalS" default:"3600"`
}

// internalVersioningConfiguration is used in XML serialization
type internalVersioningConfiguration struct {
	Type             string          `xml:"type,attr,omitempty"`
	Params           []internalParam `xml:"param"`
	CleanupIntervalS int             `xml:"cleanupIntervalS" default:"3600"`
}

type internalParam struct {
	Key string `xml:"key,attr"`
	Val string `xml:"val,attr"`
}

func (c VersioningConfiguration) Copy() VersioningConfiguration {
	cp := c
	cp.Params = make(map[string]string, len(c.Params))
	for k, v := range c.Params {
		cp.Params[k] = v
	}
	return cp
}

func (c *VersioningConfiguration) UnmarshalJSON(data []byte) error {
	util.SetDefaults(c)
	type noCustomUnmarshal VersioningConfiguration
	ptr := (*noCustomUnmarshal)(c)
	return json.Unmarshal(data, ptr)
}

func (c *VersioningConfiguration) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var intCfg internalVersioningConfiguration
	util.SetDefaults(&intCfg)
	if err := d.DecodeElement(&intCfg, &start); err != nil {
		return err
	}
	c.fromInternal(intCfg)
	return nil
}

func (c *VersioningConfiguration) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(c.toInternal())
}

func (c *VersioningConfiguration) toInternal() internalVersioningConfiguration {
	var tmp internalVersioningConfiguration
	tmp.Type = c.Type
	tmp.CleanupIntervalS = c.CleanupIntervalS
	for k, v := range c.Params {
		tmp.Params = append(tmp.Params, internalParam{k, v})
	}
	sort.Slice(tmp.Params, func(a, b int) bool {
		return tmp.Params[a].Key < tmp.Params[b].Key
	})
	return tmp
}

func (c *VersioningConfiguration) fromInternal(intCfg internalVersioningConfiguration) {
	c.Type = intCfg.Type
	c.CleanupIntervalS = intCfg.CleanupIntervalS
	c.Params = make(map[string]string, len(intCfg.Params))
	for _, p := range intCfg.Params {
		c.Params[p.Key] = p.Val
	}
}
