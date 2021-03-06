// Copyright (c) 2020 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package utils

import (
	"encoding/json"
	"net/http"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

// ErrorResponseFE is a struct that contains informations about a error
type ErrorResponseFE struct {
	// Error contains the (generic) class of the error
	Error string `json:"error"`
	// ErrorDescription contains detailed informations about the error
	ErrorDescription string `json:"errorDescription"`
	// File contains the filename of the source code where the error was detected
	SourceFilename string `json:"sourceFilename"`
	// LineNumber contains the number of the line where the error was detected
	LineNumber int `json:"lineNumber"`
}

// WriteAndLogError write the error to the w with the statusCode as statusCode and log the error to the stdout
func WriteAndLogError(log *logrus.Logger, w http.ResponseWriter, statusCode int, err AdvancedErrorInterface) {
	resp := ErrorResponseFE{
		Error:            err.ErrorClass(),
		ErrorDescription: err.Error(),
		LineNumber:       err.LineNumber(),
		SourceFilename:   err.SourceFilename(),
	}

	LogErr(log, err)

	WriteJSONResponse(w, statusCode, resp)
}

// WriteJSONResponse write the statuscode and the response to w
func WriteJSONResponse(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

// WriteExtJSONResponse write the statuscode and the response to w
func WriteExtJSONResponse(log *logrus.Logger, w http.ResponseWriter, statusCode int, resp interface{}) {
	raw, err := bson.MarshalExtJSON(resp, true, false)
	if err != nil {
		WriteAndLogError(log, w, http.StatusInternalServerError, NewAdvancedErrorPtr(err, "MARSHAL_EXT_JSON"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(raw)
}

// WriteXLSXResponse for .xlsx fils
func WriteXLSXResponse(w http.ResponseWriter, resp *excelize.File) {
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.WriteHeader(http.StatusOK)

	resp.Write(w)
}

// WriteXLSMResponse for .xlsm files
func WriteXLSMResponse(w http.ResponseWriter, resp *excelize.File) {
	w.Header().Set("Content-Type", "application/vnd.ms-excel.sheet.macroEnabled.12")
	w.WriteHeader(http.StatusOK)

	resp.Write(w)
}
