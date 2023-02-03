// Copyright 2022 The Matrix.org Foundation C.I.C.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlite3

import (
	"database/sql"

	"github.com/matrix-org/dendrite/internal/caching"
	"github.com/matrix-org/dendrite/internal/sqlutil"
	"github.com/matrix-org/dendrite/relayapi/storage/shared"
	"github.com/matrix-org/dendrite/setup/base"
	"github.com/matrix-org/dendrite/setup/config"
	"github.com/matrix-org/gomatrixserverlib"
)

// Database stores information needed by the federation sender
type Database struct {
	shared.Database
	db     *sql.DB
	writer sqlutil.Writer
}

// NewDatabase opens a new database
func NewDatabase(
	base *base.BaseDendrite,
	dbProperties *config.DatabaseOptions,
	cache caching.FederationCache,
	isLocalServerName func(gomatrixserverlib.ServerName) bool,
) (*Database, error) {
	var d Database
	var err error
	if d.db, d.writer, err = base.DatabaseConnection(dbProperties, sqlutil.NewExclusiveWriter()); err != nil {
		return nil, err
	}
	queue, err := NewSQLiteRelayQueueTable(d.db)
	if err != nil {
		return nil, err
	}
	queueJSON, err := NewSQLiteRelayQueueJSONTable(d.db)
	if err != nil {
		return nil, err
	}
	d.Database = shared.Database{
		DB:                d.db,
		IsLocalServerName: isLocalServerName,
		Cache:             cache,
		Writer:            d.writer,
		RelayQueue:        queue,
		RelayQueueJSON:    queueJSON,
	}
	return &d, nil
}