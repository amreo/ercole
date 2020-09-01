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

package apimodel

import "go.mongodb.org/mongo-driver/bson/primitive"

// OracleDatabaseAgreementsAddRequest contains the informations needed to add new agreements
type OracleDatabaseAgreementsAddRequest struct {
	AgreementID     string   `json:"agreementID" bson:"agreementID"`
	PartsID         []string `json:"partsID" bson:"partsID"`
	CSI             string   `json:"csi" bson:"csi"`
	ReferenceNumber string   `json:"referenceNumber" bson:"referenceNumber"`
	Unlimited       bool     `json:"unlimited" bson:"unlimited"`
	Count           int      `json:"count" bson:"count"`
	CatchAll        bool     `json:"catchAll" bson:"catchAll"`
	Hosts           []string `json:"hosts" bson:"hosts"`
}

// OracleDatabaseAgreementsFE contains the informations about a agreement
type OracleDatabaseAgreementsFE struct {
	ID              primitive.ObjectID                         `json:"id" bson:"_id"`
	AgreementID     string                                     `json:"agreementID" bson:"agreementID"`
	PartID          string                                     `json:"partID" bson:"partID"`
	ItemDescription string                                     `json:"itemDescription" bson:"itemDescription"`
	Metrics         string                                     `json:"metrics" bson:"metrics"`
	CSI             string                                     `json:"csi" bson:"csi"`
	ReferenceNumber string                                     `json:"referenceNumber" bson:"referenceNumber"`
	Unlimited       bool                                       `json:"unlimited" bson:"unlimited"`
	Count           int                                        `json:"-" bson:"count"`
	LicensesCount   int                                        `json:"licensesCount" bson:"licensesCount"`
	UsersCount      int                                        `json:"usersCount" bson:"usersCount"`
	AvailableCount  int                                        `json:"availableCount" bson:"availableCount"`
	CatchAll        bool                                       `json:"catchAll" bson:"catchAll"`
	Hosts           []OracleDatabaseAgreementsAssociatedHostFE `json:"hosts" bson:"hosts"`
}

// OracleDatabaseAgreementsAssociatedHostFE contains the informations about a associated host in agreement
type OracleDatabaseAgreementsAssociatedHostFE struct {
	Hostname                  string `json:"hostname" bson:"hostname"`
	CoveredLicensesCount      int    `json:"coveredLicensesCount" bson:"coveredLicensesCount"`
	TotalCoveredLicensesCount int    `json:"totalCoveredLicensesCount" bson:"totalCoveredLicensesCount"`
}

// SearchOracleDatabaseAgreementsFilters contains the filters used to get the list of Oracle/Database agreements
type SearchOracleDatabaseAgreementsFilters struct {
	AgreementID       string
	PartID            string
	ItemDescription   string
	CSI               string
	Metrics           string
	ReferenceNumber   string
	Unlimited         string //"" -> Ignore, "true" -> true, "false" -> false
	CatchAll          string //"" -> Ignore, "true" -> true, "false" -> false
	LicensesCountLTE  int
	LicensesCountGTE  int
	UsersCountLTE     int
	UsersCountGTE     int
	AvailableCountLTE int
	AvailableCountGTE int
}