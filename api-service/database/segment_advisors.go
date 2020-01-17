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

package database

import (
	"context"
	"time"

	"github.com/amreo/ercole-services/utils"
	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"
)

// SearchSegmentAdvisors search segment advisors
func (md *MongoDatabase) SearchSegmentAdvisors(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APSearchFilterStage([]string{"hostname", "database.name"}, keywords),
			mu.APProject(bson.M{
				"hostname":                  true,
				"location":                  true,
				"environment":               true,
				"created_at":                true,
				"database.name":             true,
				"database.segment_advisors": true,
			}),
			mu.APUnwind("$database.segment_advisors"),
			mu.APProject(bson.M{
				"hostname":       true,
				"location":       true,
				"environment":    true,
				"created_at":     true,
				"dbname":         "$database.name",
				"reclaimable":    "$database.segment_advisors.reclaimable",
				"segment_owner":  "$database.segment_advisors.segment_owner",
				"segment_name":   "$database.segment_advisors.segment_name",
				"segment_type":   "$database.segment_advisors.segment_type",
				"partition_name": "$database.segment_advisors.partition_name",
				"recommendation": "$database.segment_advisors.recommendation",
			}),
			mu.APOptionalSortingStage(sortBy, sortDesc),
			mu.APOptionalPagingStage(page, pageSize),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}