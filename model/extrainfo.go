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

package model

import "go.mongodb.org/mongo-driver/bson"

// ExtraInfo holds various informations.
type ExtraInfo struct {
	Databases   []Database    `bson:"Databases"`
	Filesystems []Filesystem  `bson:"Filesystems"`
	Clusters    []ClusterInfo `bson:"Clusters"`
	Exadata     *Exadata      `bson:"Exadata"`
}

// ExtraInfoBsonValidatorRules contains mongodb validation rules for extraInfo
var ExtraInfoBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"Filesystems",
	}},
	{"properties", bson.D{
		{"Databases", bson.D{
			{"anyOf", bson.A{
				bson.D{
					{"bsonType", "array"},
					{"items", DatabaseBsonValidatorRules},
				},
				bson.D{{"type", "null"}},
			}},
		}},
		{"Filesystems", bson.D{
			{"bsonType", "array"},
			{"items", FilesystemBsonValidatorRules},
		}},
		{"Clusters", bson.D{
			{"anyOf", bson.A{
				bson.D{
					{"bsonType", "array"},
					{"items", ClusterInfoBsonValidatorRules},
				},
				bson.D{{"type", "null"}},
			}},
		}},
		{"Exadata", bson.D{
			{"anyOf", bson.A{
				ExadataBsonValidatorRules,
				bson.D{{"type", "null"}},
			}},
		}},
	}},
}

// ExadataBsonValidatorRules
