// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"net/url"
	"strconv"
	"time"
)

var MIN_TIME time.Time = time.Unix(0, 0)
var MAX_TIME time.Time = time.Now().AddDate(1000, 0, 0)

//ToJSON convert v to a string containing the equivalent json rappresentation
func ToJSON(v interface{}) string {
	raw, _ := json.Marshal(v)
	return string(raw)
}

//FromJSON convert a json str to interface containing the equivalent json rappresentation
func FromJSON(str []byte) interface{} {
	var out map[string]interface{}
	json.Unmarshal(str, &out)
	return out
}

//Intptr return a point to the int passed in the argument
func Intptr(v int64) *int64 {
	return &v
}

// Contains return true if a contains x, otherwise false.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// Str2bool parse a string to a boolean
func Str2bool(in string, defaultValue bool) (bool, AdvancedErrorInterface) {
	if in == "" {
		return defaultValue, nil
	} else if val, err := strconv.ParseBool(in); err != nil {
		return false, NewAdvancedErrorPtr(err, "Unable to parse string to bool")
	} else {
		return val, nil
	}
}

// Str2int parse a string to a int
func Str2int(in string, defaultValue int) (int, AdvancedErrorInterface) {
	if in == "" {
		return defaultValue, nil
	} else if val, err := strconv.ParseInt(in, 10, 32); err != nil {
		return -1, NewAdvancedErrorPtr(err, "Unable to parse string to int")
	} else {
		return int(val), nil
	}
}

// Str2time parse a string to a time
func Str2time(in string, defaultValue time.Time) (time.Time, AdvancedErrorInterface) {
	if in == "" {
		return defaultValue, nil
	} else if val, err := time.Parse(time.RFC3339, in); err != nil {
		return time.Time{}, NewAdvancedErrorPtr(err, "Unable to parse string to time.Time")
	} else {
		return val, nil
	}
}

// NewAPIUrl return a new url crafted using the parameters
func NewAPIUrl(baseURL string, username string, password string, path string, params url.Values) *url.URL {
	u := NewAPIUrlNoParams(baseURL, username, password, path)
	u.RawQuery = params.Encode()

	return u
}

// NewAPIUrlNoParams return a new url crafted using the parameters
func NewAPIUrlNoParams(baseURL string, username string, password string, path string) *url.URL {
	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	u.User = url.UserPassword(username, password)
	u.Path = path

	return u
}
