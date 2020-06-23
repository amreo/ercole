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

package model

import (
	"reflect"

	godynstruct "github.com/amreo/go-dyn-struct"
	"go.mongodb.org/mongo-driver/bson"
)

// VMInfo holds info about the vm
type VMInfo struct {
	Name               string                 `bson:"Name"`
	Hostname           string                 `bson:"Hostname"` //Hostname or IP address
	CappedCPU          bool                   `bson:"CappedCPU"`
	VirtualizationNode string                 `bson:"VirtualizationNode"`
	OtherInfo          map[string]interface{} `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v VMInfo) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *VMInfo) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v VMInfo) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *VMInfo) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// VMInfoBsonValidatorRules contains mongodb validation rules for VMInfo
var VMInfoBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"Name",
		"Hostname",
		"CappedCPU",
		"VirtualizationNode",
	},
	"properties": bson.M{
		"Name": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 128,
		},
		"Hostname": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 253,
			"pattern":   "^(([a-zA-Z]|[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]).)*([A-Za-z]|[A-Za-z][A-Za-z0-9-]*[A-Za-z0-9])$",
		},
		"CappedCPU": bson.M{
			"bsonType": "bool",
		},
		"VirtualizationNode": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 253,
			"pattern":   "^(([a-zA-Z]|[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]).)*([A-Za-z]|[A-Za-z][A-Za-z0-9-]*[A-Za-z0-9])$",
		},
	},
}
