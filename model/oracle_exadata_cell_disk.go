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

// OracleExadataCellDisk holds info about a exadata cell disk
type OracleExadataCellDisk struct {
	ErrCount  int                    `bson:"ErrCount"`
	Name      string                 `bson:"Name"`
	Status    string                 `bson:"Status"`
	UsedPerc  int                    `bson:"UsedPerc"`
	OtherInfo map[string]interface{} `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v OracleExadataCellDisk) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *OracleExadataCellDisk) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v OracleExadataCellDisk) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *OracleExadataCellDisk) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// OracleExadataCellDiskBsonValidatorRules contains mongodb validation rules for OracleExadataCellDisk
var OracleExadataCellDiskBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"ErrCount",
		"Name",
		"Status",
		"UsedPerc",
	},
	"properties": bson.M{
		"ErrCount": bson.M{
			"bsonType": "number",
			"minimum":  0,
		},
		"Name": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 64,
		},
		"Status": bson.M{
			"bsonType": "string",
		},
		"UsedPerc": bson.M{
			"bsonType": "number",
			"minimum":  0,
			"maximum":  100,
		},
	},
}