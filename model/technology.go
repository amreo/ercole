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

// Technology names
const (
	TechnologyOracleDatabase           string = "Oracle/Database"
	TechnologyOracleExadata            string = "Oracle/Exadata"
	TechnologyMicrosoftSQLServer       string = "Microsoft/SQLServer"
	TechnologyMariaDBFoundationMariaDB string = "MariaDBFoundation/MariaDB"
	TechnologyPostgreSQLPostgreSQL     string = "PostgreSQL/PostgreSQL"
	TechnologyOracleMySQL              string = "Oracle/MySQL"
	TechnologyOracleVM                 string = "Oracle/VM"
	TechnologyVMWare                   string = "VMWare/VMWare"
	TechnologyUnknownOperatingSystem   string = "Unknown/Unknown"
)

// Pointers to technology names
var (
	TechnologyOracleDatabasePtr           *string = str2CopyPtr(TechnologyOracleDatabase)
	TechnologyOracleExadataPtr            *string = str2CopyPtr(TechnologyOracleExadata)
	TechnologyMicrosoftSQLServerPrt       *string = str2CopyPtr(TechnologyMicrosoftSQLServer)
	TechnologyMariaDBFoundationMariaDBPrt *string = str2CopyPtr(TechnologyMariaDBFoundationMariaDB)
	TechnologyPostgreSQLPostgreSQLPrt     *string = str2CopyPtr(TechnologyPostgreSQLPostgreSQL)
	TechnologyOracleMySQLPrt              *string = str2CopyPtr(TechnologyOracleMySQL)
	TechnologyOracleVMPrt                 *string = str2CopyPtr(TechnologyOracleVM)
	TechnologyVMWarePrt                   *string = str2CopyPtr(TechnologyVMWare)
	TechnologyUnknownOperatingSystemPrt   *string = str2CopyPtr(TechnologyUnknownOperatingSystem)
)

// TechnologyInfo contains the informations about a technology
type TechnologyInfo struct {
	Product    string `json:"product"`
	PrettyName string `json:"prettyName"`
	Color      string `json:"color"`
	Logo       string `json:"logo"`
}

// TechnologySupportedMetrics contains the informations about the supported metrics of a technology
type TechnologySupportedMetrics struct {
	Product string   `json:"product"`
	Metrics []string `json:"metrics"`
}

// TechnologiesSupportedMetricsMap contains all metrics of all technology
var TechnologiesSupportedMetricsMap map[string]TechnologySupportedMetrics = map[string]TechnologySupportedMetrics{
	TechnologyOracleDatabase: {
		Product: TechnologyOracleDatabase,
		Metrics: []string{"work", "version"},
	},
}
