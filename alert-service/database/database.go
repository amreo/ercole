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

// Database contains methods used to perform CRUD operations to the MongoDB database
package database

import (
	"context"
	"log"
	"time"

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabaseInterface is a interface that wrap methods used to perform CRUD operations in the mongodb database
type MongoDatabaseInterface interface {
	// Init initializes the connection to the database
	Init()
	// FindHostData find a host data
	FindHostData(id primitive.ObjectID) (model.HostData, utils.AdvancedErrorInterface)
	// FindMostRecentHostDataOlderThan return the most recest hostdata that is older than t
	FindMostRecentHostDataOlderThan(hostname string, t time.Time) (model.HostData, utils.AdvancedErrorInterface)
	// InsertAlert inserr the alert in the database
	InsertAlert(alert model.Alert) (*mongo.InsertOneResult, utils.AdvancedErrorInterface)
	// FindOldCurrentHost return the list of current hosts that haven't sent hostdata after time t
	FindOldCurrentHosts(t time.Time) ([]string, utils.AdvancedErrorInterface)
	// ExistNoDataAlertByHost return true if the host has associated a new NO_DATA alert
	ExistNoDataAlertByHost(hostname string) (bool, utils.AdvancedErrorInterface)
}

// MongoDatabase is a implementation
type MongoDatabase struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Client contain the mongodb client
	Client *mongo.Client
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
}

// Init initializes the connection to the database
func (md *MongoDatabase) Init() {
	//Connect to mongodb
	md.ConnectToMongodb()
	log.Println("MongoDatabase is connected to MongoDB!", md.Config.Mongodb.URI)
}

// ConnectToMongodb connects to the MongoDB and return the connection
func (md *MongoDatabase) ConnectToMongodb() {
	var err error

	//Set client options
	clientOptions := options.Client().ApplyURI(md.Config.Mongodb.URI)

	//Connect to MongoDB
	md.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	//Check the connection
	err = md.Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// FindHostData find a host data
func (md *MongoDatabase) FindHostData(id primitive.ObjectID) (model.HostData, utils.AdvancedErrorInterface) {
	//Find the hostdata
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").FindOne(context.TODO(), bson.M{
		"_id": id,
	})
	if res.Err() != nil {
		return model.HostData{}, utils.NewAdvancedErrorPtr(res.Err(), "DB ERROR")
	}

	//Decode the data
	var out model.HostData
	if err := res.Decode(&out); err != nil {
		return model.HostData{}, utils.NewAdvancedErrorPtr(res.Err(), "DB ERROR")
	}

	//Return it!
	return out, nil
}

// FindMostRecentHostDataOlderThan return the most recest hostdata that is older than t
func (md *MongoDatabase) FindMostRecentHostDataOlderThan(hostname string, t time.Time) (model.HostData, utils.AdvancedErrorInterface) {
	var out model.HostData

	//Find the most recent HostData older than t
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		bson.A{
			bson.M{"$match": bson.M{
				"hostname": hostname,
				"created_at": bson.M{
					"$lt": t,
				},
			}},
			bson.M{"$sort": bson.M{
				"created_at": -1,
			}},
			bson.M{"$limit": 1},
		},
	)
	if err != nil {
		return model.HostData{}, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return model.HostData{}, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return model.HostData{}, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out, nil
}

// InsertAlert inser the alert in the database
func (md *MongoDatabase) InsertAlert(alert model.Alert) (*mongo.InsertOneResult, utils.AdvancedErrorInterface) {
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("alerts").InsertOne(context.TODO(), alert)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}
	return res, nil
}

// FindOldCurrentHosts return the list of current hosts that haven't sent hostdata after time t
func (md *MongoDatabase) FindOldCurrentHosts(t time.Time) ([]string, utils.AdvancedErrorInterface) {
	//Get the list of old current hosts
	values, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Distinct(
		context.TODO(),
		"hostname",
		bson.M{
			"archived": false,
			"created_at": bson.M{
				"$lt": t,
			},
		})
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Convert the slice of interface{} to []string
	var hosts []string
	for _, val := range values {
		hosts = append(hosts, val.(string))
	}

	//Return it
	return hosts, nil
}

// ExistNoDataAlertByHost return true if the host has associated a new NO_DATA alert
func (md *MongoDatabase) ExistNoDataAlertByHost(hostname string) (bool, utils.AdvancedErrorInterface) {
	//Count the number of new NO_DATA alerts associated to the host
	val, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("alerts").CountDocuments(context.TODO(), bson.D{
		{"alert_code", model.AlertCodeNoData},
		{"alert_status", model.AlertStatusNew},
		{"other_info.hostname", hostname},
	}, &options.CountOptions{
		Limit: utils.Intptr(1),
	})
	if err != nil {
		return false, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Return true if the count > 0
	return val > 0, nil
}