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

package database

import (
	"context"
	"time"

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchAlerts search alerts
func (md *MongoDatabase) SearchAlerts(mode string, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, severity string, status string, from time.Time, to time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface) {
	var out []map[string]interface{} = make([]map[string]interface{}, 0)

	//Find the matching alerts
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("alerts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APOptionalStage(status != "", mu.APMatch(bson.M{
				"alertStatus": status,
			})),
			mu.APOptionalStage(severity != "", mu.APMatch(bson.M{
				"alertSeverity": severity,
			})),
			mu.APMatch(bson.M{
				"date": bson.M{
					"$gte": from,
					"$lt":  to,
				},
			}),
			mu.APSearchFilterStage([]interface{}{
				"$description",
				"$alertCode",
				"$alertSeverity",
				"$otherInfo.Hostname",
				"$otherInfo.Dbname",
				"$otherInfo.Features",
			}, keywords),
			mu.APSet(bson.M{
				"hostname": "$otherInfo.hostname",
			}),
			mu.APOptionalStage(mode == "aggregated-code-severity", mu.MAPipeline(
				mu.APGroup(bson.M{
					"_id": bson.M{
						"code":     "$alertCode",
						"severity": "$alertSeverity",
						"category": "$alertCategory",
					},
					"count": mu.APOSum(1),
					"oldestAlert": bson.M{
						"$min": "$date",
					},
					"affectedHosts": bson.M{
						"$addToSet": "$hostname",
					},
				}),
				mu.APProject(bson.M{
					"_id":           false,
					"category":      "$_id.category",
					"code":          "$_id.code",
					"severity":      "$_id.severity",
					"count":         true,
					"affectedHosts": mu.APOSize("$affectedHosts"),
					"oldestAlert":   true,
				}),
			)),
			mu.APOptionalStage(mode == "aggregated-category-severity", mu.MAPipeline(
				mu.APGroup(bson.M{
					"_id": bson.M{
						"severity": "$alertSeverity",
						"category": "$alertCategory",
					},
					"count": mu.APOSum(1),
					"oldestAlert": bson.M{
						"$min": "$date",
					},
					"affectedHosts": bson.M{
						"$addToSet": "$hostname",
					},
				}),
				mu.APProject(bson.M{
					"_id":           false,
					"category":      "$_id.category",
					"severity":      "$_id.severity",
					"count":         true,
					"affectedHosts": mu.APOSize("$affectedHosts"),
					"oldestAlert":   true,
				}),
			)),
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
		out = append(out, item)
	}
	return out, nil
}

// UpdateAlertsStatus change the status of the specified alerts
func (md *MongoDatabase) UpdateAlertsStatus(ids []primitive.ObjectID, newStatus string) utils.AdvancedErrorInterface {
	bsonIds := bson.A{}
	for _, id := range ids {
		bsonIds = append(bsonIds, bson.M{"_id": id})
	}
	filter := bson.M{"$or": bsonIds}

	count, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("alerts").
		CountDocuments(context.TODO(), filter)
	if count != int64(len(ids)) {
		return utils.AerrAlertNotFound
	}

	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("alerts").
		UpdateMany(context.TODO(),
			filter,
			mu.UOSet(bson.M{
				"alertStatus": newStatus,
			}))
	if err != nil || res.MatchedCount != int64(len(ids)) {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return nil
}
