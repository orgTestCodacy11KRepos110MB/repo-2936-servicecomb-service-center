/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import (
	"strconv"
	"time"

	"github.com/apache/servicecomb-service-center/datasource/etcd/path"
	"github.com/apache/servicecomb-service-center/pkg/dump"
	"github.com/apache/servicecomb-service-center/pkg/util"
	"github.com/go-chassis/cari/discovery"
)

func GetDomainProject(resource interface{}) (domainProject string) {
	switch r := resource.(type) {
	case *dump.Microservice:
		_, domainProject = path.GetInfoFromSvcKV(
			util.StringToBytesWithNoCopy(r.Key))
	case *dump.Instance:
		_, _, domainProject = path.GetInfoFromInstKV(
			util.StringToBytesWithNoCopy(r.Key))
	}
	return
}

type Service struct {
	DomainProject string
	Environment   string
	AppID         string
	ServiceName   string
	Versions      []string
	Frameworks    []*discovery.FrameWork
	Endpoints     []string
	Timestamp     int64 // the seconds from 0 to now
}

func (s *Service) AppendVersion(v string) {
	s.Versions = append(s.Versions, v)
}

func (s *Service) AppendFramework(property *discovery.FrameWork) {
	if property == nil || property.Name == "" {
		return
	}
	for _, fw := range s.Frameworks {
		if fw.Name == property.Name && fw.Version == property.Version {
			return
		}
	}
	s.Frameworks = append(s.Frameworks, property)
}

func (s *Service) AppendEndpoints(endpoints []string) {
	s.Endpoints = append(s.Endpoints, endpoints...)
}

func (s *Service) UpdateTimestamp(t string) {
	d, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return
	}
	if s.Timestamp == 0 || s.Timestamp > d {
		s.Timestamp = d
	}
}

func (s *Service) Age() time.Duration {
	return time.Since(time.Unix(s.Timestamp, 0).Local())
}

type Instance struct {
	DomainProject string
	Host          string
	Endpoints     []string
	Environment   string
	AppID         string
	ServiceName   string
	Version       string
	Framework     *discovery.FrameWork
	Lease         int64 // seconds
	Timestamp     int64 // the seconds from 0 to now
}

func (s *Instance) SetLease(hc *discovery.HealthCheck) {
	if hc == nil {
		s.Lease = -1
		return
	}
	if hc.Mode == discovery.CHECK_BY_PLATFORM {
		s.Lease = 0
		return
	}
	s.Lease = int64(hc.Interval * (hc.Times + 1))
}

func (s *Instance) UpdateTimestamp(t string) {
	d, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return
	}
	if s.Timestamp == 0 || s.Timestamp > d {
		s.Timestamp = d
	}
}

func (s *Instance) Age() time.Duration {
	return time.Since(time.Unix(s.Timestamp, 0).Local())
}
