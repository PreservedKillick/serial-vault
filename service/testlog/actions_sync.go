// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2018 Canonical Ltd
 * License granted by Canonical Limited
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package testlog

import (
	"encoding/base64"
	"net/http"

	"github.com/CanonicalLtd/serial-vault/datastore"
	"github.com/CanonicalLtd/serial-vault/service/auth"
	"github.com/CanonicalLtd/serial-vault/service/response"
)

func syncLogHandler(w http.ResponseWriter, user datastore.User, apiCall bool, testLog datastore.TestLog) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := auth.CheckUserPermissions(user, datastore.SyncUser, apiCall)
	if err != nil {
		response.FormatStandardResponse(false, "error-auth", "", "", w)
		return
	}

	if len(testLog.Data) == 0 {
		response.FormatStandardResponse(false, "error-testlog-data", "", "No file data provided", w)
		return
	}

	// Check we have something that's decodeable
	_, err = base64.StdEncoding.DecodeString(testLog.Data)
	if err != nil {
		response.FormatStandardResponse(false, "error-testlog-data", "", err.Error(), w)
		return
	}

	// Create the test log record
	err = datastore.Environ.DB.CreateTestLog(testLog)
	if err != nil {
		response.FormatStandardResponse(false, "error-testlog-create", "", err.Error(), w)
		return
	}

	// Return successful JSON response
	w.WriteHeader(http.StatusOK)
	response.FormatStandardResponse(true, "", "", "", w)
}

func syncUpdateLogHandler(w http.ResponseWriter, user datastore.User, apiCall bool, logID int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := auth.CheckUserPermissions(user, datastore.SyncUser, apiCall)
	if err != nil {
		response.FormatStandardResponse(false, "error-auth", "", "", w)
		return
	}

	// Update the test log record to indicate that it's been synced
	err = datastore.Environ.DB.UpdateAllowedTestLog(logID, user)
	if err != nil {
		response.FormatStandardResponse(false, "error-testlog-update", "", err.Error(), w)
		return
	}

	// Return successful JSON response
	w.WriteHeader(http.StatusOK)
	response.FormatStandardResponse(true, "", "", "", w)
}
