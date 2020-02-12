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
	"time"

	"github.com/amreo/ercole-services/utils"
	"github.com/sirupsen/logrus"

	"github.com/amreo/ercole-services/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabaseInterface is a interface that wrap methods used to perform CRUD operations in the mongodb database
type MongoDatabaseInterface interface {
	// Init initializes the connection to the database
	Init()
	// SearchHosts search hosts
	SearchHosts(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetHost fetch all informations about a host in the database
	GetHost(hostname string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface)
	// SearchAlerts search alerts
	SearchAlerts(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, severity string, status string, from time.Time, to time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchClusters search clusters
	SearchClusters(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchAddms search addms
	SearchAddms(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchSegmentAdvisors search segment advisors
	SearchSegmentAdvisors(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchPatchAdvisors search patch advisors
	SearchPatchAdvisors(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchDatabases search databases
	SearchDatabases(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchExadata search exadata
	SearchExadata(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// ListLicenses list licenses
	ListLicenses(full bool, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)

	// GetEnvironmentStats return a array containing the number of hosts per environment
	GetEnvironmentStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTypeStats return a array containing the number of hosts per operating system
	GetOperatingSystemStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTypeStats return a array containing the number of hosts per type
	GetTypeStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDatabaseEnvironmentStats return a array containing the number of databases per environment
	GetDatabaseEnvironmentStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDatabaseVersionStats return a array containing the number of databases per version
	GetDatabaseVersionStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTopReclaimableDatabaseStats return a array containing the total sum of reclaimable of segments advisors of the top reclaimable databases
	GetTopReclaimableDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDatabasePatchStatusStats return a array containing the number of databases per patch status
	GetDatabasePatchStatusStats(location string, windowTime time.Time, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTopWorkloadDatabaseStats return a array containing top databases by workload
	GetTopWorkloadDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDatabaseDataguardStatusStats return a array containing the number of databases per dataguard status
	GetDatabaseDataguardStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDatabaseRACStatusStats return a array containing the number of databases per RAC status
	GetDatabaseRACStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDatabaseLicenseComplianceStatusStats return the status of the compliance of licenses of databases
	GetDatabaseLicenseComplianceStatusStats(location string, environment string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface)
	// GetDatabaseArchivelogStatusStats return a array containing the number of databases per archivelog status
	GetDatabaseArchivelogStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTotalDatabaseWorkStats return the total work of databases
	GetTotalDatabaseWorkStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface)
	// GetTotalDatabaseMemorySizeStats return the total of memory size of databases
	GetTotalDatabaseMemorySizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface)
	// GetTotalDatabaseDatafileSizeStats return the total size of datafiles of databases
	GetTotalDatabaseDatafileSizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface)
	// GetTotalDatabaseSegmentSizeStats return the total size of segments of databases
	GetTotalDatabaseSegmentSizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface)
	// GetTotalExadataMemorySizeStats return the total size of memory of exadata
	GetTotalExadataMemorySizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface)
	// GetTotalExadataCPUStats return the total cpu of exadata
	GetTotalExadataCPUStats(location string, environment string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface)
	// GetAvegageExadataStorageUsageStats return the average usage of cell disks of exadata
	GetAvegageExadataStorageUsageStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface)
	// GetExadataStorageErrorCountStatusStats return a array containing the number of cell disks of exadata per error count status
	GetExadataStorageErrorCountStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetExadataPatchStatusStats return a array containing the number of exadata per patch status
	GetExadataPatchStatusStats(location string, environment string, windowTime time.Time, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)

	// SetLicenseCount set the count of a certain license
	SetLicenseCount(name string, count int) utils.AdvancedErrorInterface
}

// MongoDatabase is a implementation
type MongoDatabase struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Client contain the mongodb client
	Client *mongo.Client
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// OperatingSystemAggregationRules contains rules used to aggregate various operating systems
	OperatingSystemAggregationRules []config.AggregationRule
	// Log contains logger formatted
	Log *logrus.Logger
}

// Init initializes the connection to the database
func (md *MongoDatabase) Init() {
	//Connect to mongodb
	md.ConnectToMongodb()
	md.Log.Info("MongoDatabase is connected to MongoDB! ", md.Config.Mongodb.URI)
}

// ConnectToMongodb connects to the MongoDB and return the connection
func (md *MongoDatabase) ConnectToMongodb() {
	var err error

	//Set client options
	clientOptions := options.Client().ApplyURI(md.Config.Mongodb.URI)

	//Connect to MongoDB
	md.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		md.Log.Fatal(err)
	}

	//Check the connection
	err = md.Client.Ping(context.TODO(), nil)
	if err != nil {
		md.Log.Fatal(err)
	}
}
