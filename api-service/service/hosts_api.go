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
	"strings"
	"time"

	"github.com/ercole-io/ercole/api-service/database"
	"github.com/ercole-io/ercole/utils"
)

// SearchHosts search hosts
func (as *APIService) SearchHosts(mode string, search string, otherFilters database.SearchHostsFilters, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.SearchHosts(mode, strings.Split(search, " "), otherFilters, sortBy, sortDesc, page, pageSize, location, environment, olderThan)
}

// GetHost return the host specified in the hostname param
func (as *APIService) GetHost(hostname string, olderThan time.Time, raw bool) (interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetHost(hostname, olderThan, raw)
}

// ListLocations list locations
func (as *APIService) ListLocations(location string, environment string, olderThan time.Time) ([]string, utils.AdvancedErrorInterface) {
	return as.Database.ListLocations(location, environment, olderThan)
}

// ListEnvironments list environments
func (as *APIService) ListEnvironments(location string, environment string, olderThan time.Time) ([]string, utils.AdvancedErrorInterface) {
	return as.Database.ListEnvironments(location, environment, olderThan)
}

// ArchiveHost archive the specified host
func (as *APIService) ArchiveHost(hostname string) utils.AdvancedErrorInterface {
	return as.Database.ArchiveHost(hostname)
}
