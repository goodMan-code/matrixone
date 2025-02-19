// Copyright 2022 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fileservice

import (
	"encoding/csv"
	"strings"

	"github.com/matrixorigin/matrixone/pkg/common/moerr"
)

const ServiceNameSeparator = ":"

type Path struct {
	Service          string
	ServiceArguments []string
	File             string
}

func ParsePath(s string) (path Path, err error) {
	// split
	i := strings.LastIndex(s, ServiceNameSeparator)
	if i == -1 {
		// no service part
		path.File = s
	} else {
		// with service part
		path.Service, path.ServiceArguments, err = parseService(s[:i])
		if err != nil {
			return
		}
		path.File = s[i+1:]
	}

	// validate
	for _, r := range path.File {
		// most common patterns first
		if r >= '0' && r <= '9' ||
			r >= 'a' && r <= 'z' ||
			r >= 'A' && r <= 'Z' ||
			r == '@' ||
			r == '/' {
			continue
		}
		switch r {
		case '!', '-', '_', '.', '*', '\'', '(', ')':
			continue
		}
		err = moerr.NewInvalidPathNoCtx(path.File)
		return
	}

	path.File = strings.TrimLeft(path.File, "/") // trim leading /

	return
}

func parseService(str string) (service string, arguments []string, err error) {
	r := csv.NewReader(strings.NewReader(str))
	records, err := r.ReadAll()
	if err != nil {
		return
	}
	if len(records) != 1 {
		err = moerr.NewInvalidInputNoCtx("bad service: %v", str)
		return
	}
	service = records[0][0]
	arguments = records[0][1:]
	return
}

func ParsePathAtService(s string, serviceName string) (path Path, err error) {
	path, err = ParsePath(s)
	if err != nil {
		return
	}
	if serviceName != "" &&
		path.Service != "" &&
		!strings.EqualFold(path.Service, serviceName) {
		err = moerr.NewWrongServiceNoCtx(serviceName, path.Service)
		return
	}
	return
}

func JoinPath(serviceName string, path string) string {
	buf := new(strings.Builder)
	buf.WriteString(serviceName)
	buf.WriteString(ServiceNameSeparator)
	buf.WriteString(path)
	return buf.String()
}
