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

package service

import (
	"github.com/apache/servicecomb-service-center/pkg/util"
	"github.com/apache/servicecomb-service-center/scctl/pkg/model"
	"github.com/apache/servicecomb-service-center/scctl/pkg/writer"
)

const maxWidth = 35

var (
	longServiceTableHeader   = []string{"DOMAIN", "NAME", "APPID", "VERSIONS", "ENV", "FRAMEWORK", "ENDPOINTS", "AGE"}
	domainServiceTableHeader = []string{"DOMAIN", "NAME", "APPID", "VERSIONS", "ENV", "FRAMEWORK", "AGE"}
	shortServiceTableHeader  = []string{"NAME", "APPID", "VERSIONS", "ENV", "FRAMEWORK", "AGE"}
)

type Record struct {
	model.Service
}

func (s *Record) VersionsString() string {
	return util.StringJoin(s.Versions, "\n")
}

func (s *Record) FrameworksString() string {
	var arr []string
	for _, fw := range s.Frameworks {
		arr = append(arr, fw.Name)
	}
	return util.StringJoin(arr, "\n")
}

func (s *Record) EndpointsString() string {
	return util.StringJoin(s.Endpoints, "\n")
}

func (s *Record) AgeString() string {
	return writer.TimeFormat(s.Age())
}

func (s *Record) Domain() string {
	domain, _ := util.FromDomainProject(s.DomainProject)
	return domain
}

func (s *Record) PrintBody(fmt string, all bool) []string {
	switch {
	case fmt == "wide":
		return []string{s.Domain(), s.ServiceName, s.AppID, s.VersionsString(), s.Environment,
			s.FrameworksString(), s.EndpointsString(), s.AgeString()}
	case all:
		return writer.Reshape(maxWidth, []string{s.Domain(), s.ServiceName, s.AppID,
			s.VersionsString(), s.Environment, s.FrameworksString(), s.AgeString()})
	default:
		return writer.Reshape(maxWidth, []string{s.ServiceName, s.AppID,
			s.VersionsString(), s.Environment, s.FrameworksString(), s.AgeString()})
	}
}

type Printer struct {
	Records map[string]*Record
	flags   []interface{}
}

func (sp *Printer) SetOutputFormat(f string, all bool) {
	sp.Flags(f, all)
}

func (sp *Printer) Flags(flags ...interface{}) []interface{} {
	if len(flags) > 0 {
		sp.flags = flags
	}
	return sp.flags
}

func (sp *Printer) PrintBody() (slice [][]string) {
	for _, s := range sp.Records {
		slice = append(slice, s.PrintBody(sp.flags[0].(string), sp.flags[1].(bool)))
	}
	return
}

func (sp *Printer) PrintTitle() []string {
	switch {
	case sp.flags[0] == "wide":
		return longServiceTableHeader
	case sp.flags[1].(bool):
		return domainServiceTableHeader
	default:
		return shortServiceTableHeader
	}
}

func (sp *Printer) Sorter() *writer.RecordsSorter {
	return nil
}
