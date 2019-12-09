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

package service

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/amreo/ercole-services/model"

	"github.com/amreo/ercole-services/utils"
)

//go:generate mockgen -source ../database/database.go -destination=fake_database.go -package=service
//go:generate mockgen -source service.go -destination=fake_service.go -package=service

//Common data
var errMock error = errors.New("MockError")
var aerrMock utils.AdvancedErrorInterface = utils.NewAdvancedErrorPtr(errMock, "mock")

var hostData1 model.HostData = model.HostData{
	ID:        utils.Str2oid("5dc3f534db7e81a98b726a52"),
	Hostname:  "superhost1",
	Archived:  false,
	CreatedAt: utils.P("2019-11-05T14:02:03Z"),
}

var hostData2 model.HostData = model.HostData{
	ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	Hostname:  "superhost1",
	Archived:  true,
	CreatedAt: utils.P("2019-11-05T12:02:03Z"),
}

var hostData3 model.HostData = model.HostData{
	ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	Hostname:  "superhost1",
	Archived:  true,
	CreatedAt: utils.P("2019-11-05T16:02:03Z"),
	Extra: model.ExtraInfo{
		Databases: []model.Database{
			model.Database{
				Name: "acd",
			},
		},
	},
}

var hostData4 model.HostData = model.HostData{
	ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	Hostname:  "superhost1",
	Archived:  true,
	CreatedAt: utils.P("2019-11-05T18:02:03Z"),
	Extra: model.ExtraInfo{
		Databases: []model.Database{
			model.Database{
				Name: "acd",
				Licenses: []model.License{
					model.License{Name: "Oracle ENT", Count: 10},
					model.License{Name: "Driving", Count: 100},
				},
				Features: []model.Feature{
					model.Feature{Name: "Driving", Status: true},
				},
			},
		},
	},
}

type alertSimilarTo struct{ al model.Alert }

func (sa *alertSimilarTo) Matches(x interface{}) bool {
	if val, ok := x.(model.Alert); !ok {
		return false
	} else if val.AlertCode != sa.al.AlertCode {
		return false
	} else {
		return reflect.DeepEqual(sa.al.OtherInfo, val.OtherInfo)
	}
}

func (sa *alertSimilarTo) String() string {
	return fmt.Sprintf("is similar to %v", sa.al)
}