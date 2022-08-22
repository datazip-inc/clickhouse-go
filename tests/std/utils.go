// Licensed to ClickHouse, Inc. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. ClickHouse, Inc. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package std

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	clickhouse_tests "github.com/ClickHouse/clickhouse-go/v2/tests"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	seed := time.Now().UnixNano()
	fmt.Printf("using random seed %d for std tests\n", seed)
	rand.Seed(seed)
}

func CheckMinServerVersion(conn *sql.DB, major, minor, patch uint64) error {
	var version struct {
		Major uint64
		Minor uint64
		Patch uint64
	}
	var res string
	if err := conn.QueryRow("SELECT version()").Scan(&res); err != nil {
		panic(err)
	}
	for i, v := range strings.Split(res, ".") {
		switch i {
		case 0:
			version.Major, _ = strconv.ParseUint(v, 10, 64)
		case 1:
			version.Minor, _ = strconv.ParseUint(v, 10, 64)
		case 2:
			version.Patch, _ = strconv.ParseUint(v, 10, 64)
		}
	}
	if version.Major < major || (version.Major == major && version.Minor < minor) || (version.Major == major && version.Minor == minor && version.Patch < patch) {
		return fmt.Errorf("unsupported server version %d.%d.%d < %d.%d.%d", version.Major, version.Minor, version.Patch, major, minor, patch)
	}
	return nil
}

func GetDSNConnection(protocol clickhouse.Protocol, secure bool) (*sql.DB, error) {
	env, err := clickhouse_tests.GetTestEnvironment("std")
	if err != nil {
		return nil, err
	}
	switch protocol {
	case clickhouse.HTTP:
		switch secure {
		case true:
			return sql.Open("clickhouse", fmt.Sprintf(fmt.Sprintf("https://%s:%s@%s:%d?secure=true", env.Username, env.Password, env.Host, env.HttpsPort)))
		case false:
			return sql.Open("clickhouse", fmt.Sprintf(fmt.Sprintf("http://%s:%s@%s:%d", env.Username, env.Password, env.Host, env.HttpPort)))
		}
	case clickhouse.Native:
		switch secure {
		case true:
			return sql.Open("clickhouse", fmt.Sprintf(fmt.Sprintf("clickhouse://%s:%s@%s:%d?secure=true", env.Username, env.Password, env.Host, env.SslPort)))
		case false:
			return sql.Open("clickhouse", fmt.Sprintf(fmt.Sprintf("clickhouse://%s:%s@%s:%d", env.Username, env.Password, env.Host, env.Port)))
		}
	}
	return nil, fmt.Errorf("unsupport protocol - %s", protocol.String())
}
