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

// Package service is a package that provides methods for querying data
package service

import (
	"time"

	"github.com/ercole-io/ercole/api-service/apimodel"
	"github.com/ercole-io/ercole/api-service/database"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/config"
)

// APIServiceInterface is a interface that wrap methods used to querying data
type APIServiceInterface interface {
	// Init initialize the service
	Init()
	// SearchHosts search hosts
	SearchHosts(mode string, search string, otherFilters database.SearchHostsFilters, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface)
	// GetHost return the host specified in the hostname param
	GetHost(hostname string, olderThan time.Time, raw bool) (interface{}, utils.AdvancedErrorInterface)
	// ListManagedTechnologies returns the list of technologies with some stats
	ListManagedTechnologies(sortBy string, sortDesc bool, location string, environment string, olderThan time.Time) ([]model.TechnologyStatus, utils.AdvancedErrorInterface)
	// SearchAlerts search alerts
	SearchAlerts(mode string, search string, sortBy string, sortDesc bool, page int, pageSize int, severity string, status string, from time.Time, to time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface)
	// SearchClusters search clusters
	SearchClusters(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface)
	// GetCluster return the cluster specified in the clusterName param
	GetCluster(clusterName string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface)
	// SearchOracleDatabaseAddms search addm
	SearchOracleDatabaseAddms(search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface)
	// SearchOracleDatabaseSegmentAdvisors search segment advisors
	SearchOracleDatabaseSegmentAdvisors(search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface)
	// SearchOracleDatabasePatchAdvisors search patch advisors
	SearchOracleDatabasePatchAdvisors(search string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string, olderThan time.Time, status string) ([]map[string]interface{}, utils.AdvancedErrorInterface)
	// SearchOracleDatabases search databases
	SearchOracleDatabases(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface)
	// SearchOracleExadata search exadata
	SearchOracleExadata(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchLicenses search licenses
	SearchLicenses(location string, environment string, olderThan time.Time) ([]apimodel.OracleDatabaseLicenseInfo, utils.AdvancedErrorInterface)
	// SearchOracleDatabaseConsumedLicenses return the list of consumed licenses
	SearchOracleDatabaseConsumedLicenses(sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)

	// GetLicense return the license specified in the name param
	GetLicense(name string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface)
	// ListLocations list locations
	ListLocations(location string, environment string, olderThan time.Time) ([]string, utils.AdvancedErrorInterface)
	// ListEnvironments list environments
	ListEnvironments(location string, environment string, olderThan time.Time) ([]string, utils.AdvancedErrorInterface)
	// SearchOracleDatabaseLicenseModifiers search license modifiers
	SearchOracleDatabaseLicenseModifiers(search string, sortBy string, sortDesc bool, page int, pageSize int) ([]map[string]interface{}, utils.AdvancedErrorInterface)

	// GetPatchingFunction return the patching function specified in the hostname param
	GetPatchingFunction(hostname string) (interface{}, utils.AdvancedErrorInterface)

	// GetHostsCountStats return the number of the non-archived hosts
	GetHostsCountStats(location string, environment string, olderThan time.Time) (int, utils.AdvancedErrorInterface)
	// GetEnvironmentStats return a array containing the number of hosts per environment
	GetEnvironmentStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOperatingSystemStats return a array containing the number of hosts per operating system
	GetOperatingSystemStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTypeStats return a array containing the number of hosts per type
	GetTypeStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTopUnusedOracleDatabaseInstanceResourceStats return a array containing top unused instance resource by workload
	GetTopUnusedOracleDatabaseInstanceResourceStats(location string, environment string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabaseEnvironmentStats return a array containing the number of databases per environment
	GetOracleDatabaseEnvironmentStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabaseHighReliabilityStats return a array containing the number of databases per high-reliability status
	GetOracleDatabaseHighReliabilityStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabaseVersionStats return a array containing the number of databases per version
	GetOracleDatabaseVersionStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTopReclaimableOracleDatabaseStats return a array containing the total sum of reclaimable of segments advisors of the top reclaimable databases
	GetTopReclaimableOracleDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTotalTechnologiesComplianceStats return the total compliance of all technologie
	GetTotalTechnologiesComplianceStats(location string, environment string, olderThan time.Time) (map[string]interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabasePatchStatusStats return a array containing the number of databases per patch status
	GetOracleDatabasePatchStatusStats(location string, windowTime time.Time, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTopWorkloadOracleDatabaseStats return a array containing top databases by workload
	GetTopWorkloadOracleDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabaseDataguardStatusStats return a array containing the number of databases per dataguard status
	GetOracleDatabaseDataguardStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabaseRACStatusStats return a array containing the number of databases per RAC status
	GetOracleDatabaseRACStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabaseArchivelogStatusStats return a array containing the number of databases per archivelog status
	GetOracleDatabaseArchivelogStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTotalOracleDatabaseWorkStats return the total work of databases
	GetTotalOracleDatabaseWorkStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface)
	// GetTotalOracleDatabaseMemorySizeStats return the total of memory size of databases
	GetTotalOracleDatabaseMemorySizeStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface)
	// GetTotalOracleDatabaseDatafileSizeStats return the total size of datafiles of databases
	GetTotalOracleDatabaseDatafileSizeStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface)
	// GetTotalOracleDatabaseSegmentSizeStats return the total size of segments of databases
	GetTotalOracleDatabaseSegmentSizeStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface)
	// GetOracleDatabaseLicenseComplianceStatusStats return the status of the compliance of licenses of databases
	GetOracleDatabaseLicenseComplianceStatusStats(location string, environment string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface)
	// GetTotalOracleExadataMemorySizeStats return the total size of memory of exadata
	GetTotalOracleExadataMemorySizeStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface)
	// GetTotalOracleExadataCPUStats return the total cpu of exadata
	GetTotalOracleExadataCPUStats(location string, environment string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface)
	// GetAverageOracleExadataStorageUsageStats return the average usage of cell disks of exadata
	GetAverageOracleExadataStorageUsageStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface)
	// GetOracleExadataStorageErrorCountStatusStats return a array containing the number of cell disks of exadata per error count status
	GetOracleExadataStorageErrorCountStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleExadataPatchStatusStats return a array containing the number of exadata per patch status
	GetOracleExadataPatchStatusStats(location string, environment string, windowTime time.Time, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDefaultDatabaseTags return the default list of database tags from configuration
	GetDefaultDatabaseTags() ([]string, utils.AdvancedErrorInterface)
	// GetErcoleFeatures return a map of active/inactive features
	GetErcoleFeatures() (map[string]bool, utils.AdvancedErrorInterface)
	// GetErcoleFeatures return the list of technologies
	GetTechnologyList() ([]model.TechnologyInfo, utils.AdvancedErrorInterface)
	// SetLicenseCostPerProcessor set the cost per processor of a certain license
	SetLicenseCostPerProcessor(name string, costPerProcessor float64) utils.AdvancedErrorInterface

	// GetOracleDatabaseAgreementPartsList return the list of Oracle/Database agreement parts
	GetOracleDatabaseAgreementPartsList() ([]model.OracleDatabaseAgreementPart, utils.AdvancedErrorInterface)
	// AddOracleDatabaseAgreements add an Oracle Database Agreement
	AddOracleDatabaseAgreements(req apimodel.OracleDatabaseAgreementsAddRequest) (interface{}, utils.AdvancedErrorInterface)
	// UpdateOracleDatabaseAgreement update an Oracle Database Agreement
	UpdateOracleDatabaseAgreement(agreement model.OracleDatabaseAgreement) utils.AdvancedErrorInterface
	// DeleteOracleDatabaseAgreement remove an Oracle/Database agreement
	DeleteOracleDatabaseAgreement(id primitive.ObjectID) utils.AdvancedErrorInterface

	// SearchOracleDatabaseAgreements search Oracle/Database agreements
	SearchOracleDatabaseAgreements(search string, filters apimodel.SearchOracleDatabaseAgreementsFilters) ([]apimodel.OracleDatabaseAgreementsFE, utils.AdvancedErrorInterface)
	// AddAssociatedHostToOracleDatabaseAgreement add a new host to the list of associated hosts of the agreement
	AddAssociatedHostToOracleDatabaseAgreement(id primitive.ObjectID, hostname string) utils.AdvancedErrorInterface
	// RemoveAssociatedHostToOracleDatabaseAgreement remove the host from the list of associated hosts of the agreement
	RemoveAssociatedHostToOracleDatabaseAgreement(id primitive.ObjectID, hostname string) utils.AdvancedErrorInterface

	// SetPatchingFunction set the patching function of a host
	SetPatchingFunction(hostname string, pf model.PatchingFunction) (interface{}, utils.AdvancedErrorInterface)
	// DeletePatchingFunction delete the patching function of a host
	DeletePatchingFunction(hostname string) utils.AdvancedErrorInterface

	// AddTagToOracleDatabase add the tag to the database if it hasn't the tag
	AddTagToOracleDatabase(hostname string, dbname string, tagname string) utils.AdvancedErrorInterface
	// DeleteTagOfOracleDatabase delete the tag from the database if it hasn't the tag
	DeleteTagOfOracleDatabase(hostname string, dbname string, tagname string) utils.AdvancedErrorInterface
	// SetOracleDatabaseLicenseModifier set the value of certain license to newValue
	SetOracleDatabaseLicenseModifier(hostname string, dbname string, licenseName string, newValue int) utils.AdvancedErrorInterface
	// DeleteOracleDatabaseLicenseModifier delete the modifier of a certain license
	DeleteOracleDatabaseLicenseModifier(hostname string, dbname string, licenseName string) utils.AdvancedErrorInterface
	// AckAlerts ack the specified alerts
	AckAlerts(ids []primitive.ObjectID) utils.AdvancedErrorInterface
	// ArchiveHost archive the specified host
	ArchiveHost(hostname string) utils.AdvancedErrorInterface

	// GetInfoForFrontendDashboard return all informations needed for the frontend dashboard page
	GetInfoForFrontendDashboard(location string, environment string, olderThan time.Time) (map[string]interface{}, utils.AdvancedErrorInterface)
}

// APIService is the concrete implementation of APIServiceInterface.
type APIService struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Version of the saved data
	Version string
	// Database contains the database layer
	Database database.MongoDatabaseInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Log contains logger formatted
	Log *logrus.Logger
	// TechnologyInfos contains the list of technologies with their informations
	TechnologyInfos []model.TechnologyInfo
	// OracleDatabaseAgreementParts contains the list of Oracle/Database agreeement parts
	OracleDatabaseAgreementParts []model.OracleDatabaseAgreementPart
}

// Init initializes the service and database
func (as *APIService) Init() {
	as.LoadManagedTechnologiesList()
	as.LoadOracleDatabaseAgreementPartsList()
}
