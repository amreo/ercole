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

// Package service is a package that provides methods for manipulating host informations

package service

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHostInfo_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				PublisherUsername: "publ1sh3r",
				PublisherPassword: "M0stS3cretP4ssw0rd",
				RemoteEndpoint:    "http://ercole.example.org",
			},
			DataService: config.DataService{
				EnablePatching:       true,
				LogInsertingHostdata: true,
			},
		},
		Version: "1.6.6",
		Log:     utils.NewLogger("TEST"),
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	db.EXPECT().ArchiveHost("rac1_x").Return(nil, nil).Times(1)
	db.EXPECT().ArchiveHost(gomock.Any()).Times(0)
	db.EXPECT().FindPatchingFunction("rac1_x").Return(model.PatchingFunction{}, nil)
	db.EXPECT().InsertHostData(gomock.Any()).Return(&mongo.InsertOneResult{InsertedID: utils.Str2oid("5dd3a8db184dbf295f0376f2")}, nil).Do(func(newHD model.HostDataBE) {
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.ID.Timestamp())
		assert.False(t, newHD.Archived)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.CreatedAt)
		assert.Equal(t, model.SchemaVersion, newHD.ServerSchemaVersion)
		assert.Equal(t, "1.6.6", newHD.ServerVersion)
		assert.Equal(t, hd.Hostname, newHD.Hostname)
		assert.Equal(t, hd.Environment, newHD.Environment)
		//I assume that other fields are correct
	}).Times(1)
	db.EXPECT().InsertHostData(gomock.Any()).Times(0)
	http.DefaultClient = NewHTTPTestClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "http://publ1sh3r:M0stS3cretP4ssw0rd@ercole.example.org/queue/host-data-insertion/5dd3a8db184dbf295f0376f2", req.URL.String())

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
			Header:     make(http.Header),
		}, nil
	})

	res, err := hds.UpdateHostInfo(hd)
	require.NoError(t, err)
	assert.Equal(t, utils.Str2oid("5dd3a8db184dbf295f0376f2"), res)
}

func TestUpdateHostInfo_DatabaseError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				PublisherUsername: "publ1sh3r",
				PublisherPassword: "M0stS3cretP4ssw0rd",
				RemoteEndpoint:    "http://ercole.example.org",
			},
		},
		Version: "1.6.6",
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")
	db.EXPECT().ArchiveHost("rac1_x").Return(nil, aerrMock).Times(1)
	db.EXPECT().ArchiveHost(gomock.Any()).Times(0)
	db.EXPECT().FindPatchingFunction(gomock.Any()).Return(model.PatchingFunction{}, nil).Times(0)

	_, err := hds.UpdateHostInfo(hd)
	require.Equal(t, aerrMock, err)
}

func TestUpdateHostInfo_DatabaseError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				PublisherUsername: "publ1sh3r",
				PublisherPassword: "M0stS3cretP4ssw0rd",
				RemoteEndpoint:    "http://ercole.example.org",
			},
		},
		Version: "1.6.6",
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	db.EXPECT().ArchiveHost("rac1_x").Return(nil, nil).Times(1)
	db.EXPECT().ArchiveHost(gomock.Any()).Times(0)
	db.EXPECT().FindPatchingFunction(gomock.Any()).Return(model.PatchingFunction{}, nil).Times(0)
	db.EXPECT().InsertHostData(gomock.Any()).Return(nil, aerrMock).Do(func(newHD model.HostDataBE) {
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.ID.Timestamp())
		assert.False(t, newHD.Archived)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.CreatedAt)
		assert.Equal(t, model.SchemaVersion, newHD.ServerSchemaVersion)
		assert.Equal(t, "1.6.6", newHD.ServerVersion)
		assert.Equal(t, hd.Hostname, newHD.Hostname)
		assert.Equal(t, hd.Environment, newHD.Environment)
		//I assume that other fields are correct
	}).Times(1)
	db.EXPECT().InsertHostData(gomock.Any()).Times(0)

	_, err := hds.UpdateHostInfo(hd)
	require.Equal(t, aerrMock, err)
}

func TestUpdateHostInfo_DatabaseError3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				PublisherUsername: "publ1sh3r",
				PublisherPassword: "M0stS3cretP4ssw0rd",
				RemoteEndpoint:    "http://ercole.example.org",
			},
			DataService: config.DataService{
				EnablePatching:       true,
				LogInsertingHostdata: true,
			},
		},
		Version: "1.6.6",
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	db.EXPECT().FindPatchingFunction("rac1_x").Return(model.PatchingFunction{}, aerrMock)

	_, err := hds.UpdateHostInfo(hd)
	require.Equal(t, aerrMock, err)
}

func TestUpdateHostInfo_HttpError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				PublisherUsername: "publ1sh3r",
				PublisherPassword: "M0stS3cretP4ssw0rd",
				RemoteEndpoint:    "http://ercole.example.org",
			},
		},
		Version: "1.6.6",
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	db.EXPECT().ArchiveHost("rac1_x").Return(nil, nil).Times(1)
	db.EXPECT().ArchiveHost(gomock.Any()).Times(0)
	db.EXPECT().FindPatchingFunction(gomock.Any()).Return(model.PatchingFunction{}, nil).Times(0)
	db.EXPECT().InsertHostData(gomock.Any()).Return(&mongo.InsertOneResult{InsertedID: utils.Str2oid("5dd3a8db184dbf295f0376f2")}, nil).Do(func(newHD model.HostDataBE) {
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.ID.Timestamp())
		assert.False(t, newHD.Archived)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.CreatedAt)
		assert.Equal(t, model.SchemaVersion, newHD.ServerSchemaVersion)
		assert.Equal(t, "1.6.6", newHD.ServerVersion)
		assert.Equal(t, hd.Hostname, newHD.Hostname)
		assert.Equal(t, hd.Environment, newHD.Environment)
		//I assume that other fields are correct
	}).Times(1)
	db.EXPECT().InsertHostData(gomock.Any()).Times(0)
	http.DefaultClient = NewHTTPTestClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "http://publ1sh3r:M0stS3cretP4ssw0rd@ercole.example.org/queue/host-data-insertion/5dd3a8db184dbf295f0376f2", req.URL.String())
		return nil, errMock
	})

	_, err := hds.UpdateHostInfo(hd)
	require.Equal(t, "EVENT ENQUEUE", err.ErrorClass())
	require.Contains(t, err.Error(), "http://publ1sh3r:***@ercole.example.org/queue/host-data-insertion/5dd3a8db184dbf295f0376f2")
	require.Contains(t, err.Error(), "MockError")
}

func TestUpdateHostInfo_HttpError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				PublisherUsername: "publ1sh3r",
				PublisherPassword: "M0stS3cretP4ssw0rd",
				RemoteEndpoint:    "http://ercole.example.org",
			},
		},
		Version: "1.6.6",
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	db.EXPECT().ArchiveHost("rac1_x").Return(nil, nil).Times(1)
	db.EXPECT().ArchiveHost(gomock.Any()).Times(0)
	db.EXPECT().FindPatchingFunction(gomock.Any()).Return(model.PatchingFunction{}, nil).Times(0)
	db.EXPECT().InsertHostData(gomock.Any()).Return(&mongo.InsertOneResult{InsertedID: utils.Str2oid("5dd3a8db184dbf295f0376f2")}, nil).Do(func(newHD model.HostDataBE) {
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.ID.Timestamp())
		assert.False(t, newHD.Archived)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.CreatedAt)
		assert.Equal(t, model.SchemaVersion, newHD.ServerSchemaVersion)
		assert.Equal(t, "1.6.6", newHD.ServerVersion)
		assert.Equal(t, hd.Hostname, newHD.Hostname)
		assert.Equal(t, hd.Environment, newHD.Environment)
		//I assume that other fields are correct
	}).Times(1)
	db.EXPECT().InsertHostData(gomock.Any()).Times(0)
	http.DefaultClient = NewHTTPTestClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "http://publ1sh3r:M0stS3cretP4ssw0rd@ercole.example.org/queue/host-data-insertion/5dd3a8db184dbf295f0376f2", req.URL.String())
		return &http.Response{
			StatusCode: 500,
			Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
			Header:     make(http.Header),
		}, nil
	})

	_, err := hds.UpdateHostInfo(hd)
	require.Equal(t, "EVENT ENQUEUE", err.ErrorClass())
	require.EqualError(t, err, "Failed to enqueue event")
}

func TestPatchHostData_SuccessNoPatchingFunction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	db.EXPECT().FindPatchingFunction("rac1_x").Return(model.PatchingFunction{}, nil)

	res, err := hds.PatchHostData(hd)
	require.NoError(t, err)
	assert.Equal(t, hd, res)
}

func TestPatchHostData_SuccessPatchingFunction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			DataService: config.DataService{
				LogDataPatching: true,
			},
		},
		Log: utils.NewLogger("TEST"),
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")
	patchedHd := hd
	patchedHd.Tags = []string{"topolino", "pluto"}

	objID := utils.Str2oid("5ef9b4bcda4e04c0c1a94e9e")
	db.EXPECT().FindPatchingFunction("rac1_x").Return(model.PatchingFunction{
		ID:        &objID,
		CreatedAt: utils.P("2020-06-29T09:30:55+00:00"),
		Hostname:  "rac1_x",
		Vars: map[string]interface{}{
			"tags": []string{"topolino", "pluto"},
		},
		Code: `
			hostdata.tags = vars.tags;
		`,
	}, nil)

	res, err := hds.PatchHostData(hd)
	require.NoError(t, err)
	assert.Equal(t, patchedHd, res)
}

func TestPatchHostData_FailPatchingFunction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	db.EXPECT().FindPatchingFunction("rac1_x").Return(model.PatchingFunction{}, aerrMock)

	_, err := hds.PatchHostData(hd)
	require.Equal(t, aerrMock, err)
}

func TestPatchHostData_FailPatchingFunction2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			DataService: config.DataService{
				LogDataPatching: true,
			},
		},
		Log: utils.NewLogger("TEST"),
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	objID := utils.Str2oid("5ef9b4bcda4e04c0c1a94e9e")
	db.EXPECT().FindPatchingFunction("rac1_x").Return(model.PatchingFunction{
		ID:        &objID,
		CreatedAt: utils.P("2020-06-29T09:30:55+00:00"),
		Hostname:  "rac1_x",
		Vars: map[string]interface{}{
			"tags": []string{"topolino", "pluto"},
		},
		Code: `
			sdfsdasdfsdf
		`,
	}, nil)

	_, err := hds.PatchHostData(hd)
	assert.Error(t, err)
}
