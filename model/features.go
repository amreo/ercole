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

// Features holds various informations about the features of the host.
type Features struct {
	Oracle    *OracleFeature         `json:"oracle,omitempty" bson:"oracle,omitempty"`
	Microsoft *MicrosoftFeature      `json:"microsoft,omitempty" bson:"microsoft,omitempty"`
	Unknown   *UnknownFeature        `json:"unknown,omitempty" bson:"unknown,omitempty"`
	OtherInfo map[string]interface{} `json:"-" bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v Features) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *Features) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v Features) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *Features) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// FeaturesBsonValidatorRules contains mongodb validation rules for Features
var FeaturesBsonValidatorRules = bson.M{
	"bsonType": "object",
	"properties": bson.M{
		"oracle": bson.M{
			"anyOf": bson.A{
				bson.M{
					"bsonType": "null",
				},
				OracleFeatureBsonValidatorRules,
			},
		},
		"microsoft": bson.M{
			"anyOf": bson.A{
				bson.M{
					"bsonType": "null",
				},
				MicrosoftFeatureBsonValidatorRules,
			},
		},
		"unknown": bson.M{
			"anyOf": bson.A{
				bson.M{
					"bsonType": "null",
				},
				UnknownFeatureBsonValidatorRules,
			},
		},
	},
}
