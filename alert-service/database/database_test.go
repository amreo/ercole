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
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/stretchr/testify/suite"
)

func TestMongodbSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test for mongodb database(alert-service)")
	}

	mongodbHandlerSuiteTest := &MongodbSuite{}

	suite.Run(t, mongodbHandlerSuiteTest)
}

func TestConnectToMongodb_FailToConnect(t *testing.T) {
	logger := utils.NewLogger("TEST")
	logger.ExitFunc = func(int) {
		panic("log.Fatal called by test")
	}

	db := MongoDatabase{
		Config: config.Configuration{
			Mongodb: config.Mongodb{
				URI:    "wronguri:1234/test",
				DBName: fmt.Sprintf("ercole_test_%d", rand.Int()),
			},
		},
		Log: logger,
	}

	assert.PanicsWithValue(t, "log.Fatal called by test", db.ConnectToMongodb)
}

func (m *MongodbSuite) TestFindHostData_SuccessExist() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_dataservice_mongohostdata_01.json")

	m.InsertHostData(hd)

	hd2, err := m.db.FindHostData(utils.Str2oid("5dd41aa2b2fa40163c878538"))
	require.NoError(m.T(), err)

	assert.Equal(m.T(), "itl-csllab-112.sorint.localpippo", hd2.Hostname)
	assert.False(m.T(), hd2.Archived)
	assert.Equal(m.T(), utils.P("2019-11-19T16:38:58Z"), hd2.CreatedAt)
}

func (m *MongodbSuite) TestFindHostData_FailWrongID() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_dataservice_mongohostdata_01.json")
	m.InsertHostData(hd)

	notExistingID := utils.Str2oid("8a46027b2ddab34ed01a8c56")

	hd2, err := m.db.FindHostData(notExistingID)
	require.Error(m.T(), err)

	assert.Equal(m.T(), model.HostDataBE{}, hd2)
}

func (m *MongodbSuite) TestFindMostRecentHostDataOlderThan_OnlyOne() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_dataservice_mongohostdata_01.json")
	m.InsertHostData(hd)

	foundHd, err := m.db.FindMostRecentHostDataOlderThan("itl-csllab-112.sorint.localpippo", utils.P("2020-04-25T11:45:23Z"))
	require.NoError(m.T(), err)
	assert.Equal(m.T(), "itl-csllab-112.sorint.localpippo", foundHd.Hostname)
	assert.False(m.T(), foundHd.Archived)
	assert.Equal(m.T(), utils.P("2019-11-19T16:38:58Z"), foundHd.CreatedAt)
}

func (m *MongodbSuite) TestFindMostRecentHostDataOlderThan_MoreThanOne() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_dataservice_mongohostdata_04.json")
	m.InsertHostData(hd)

	m.T().Run("Should_find_hd_even_if_archived", func(t *testing.T) {
		foundHd, err := m.db.FindMostRecentHostDataOlderThan("itl-csllab-112.sorint.localpippo", utils.P("2019-11-20T15:38:58Z"))
		require.NoError(t, err)
		assert.True(t, foundHd.Archived)
		assert.Equal(t, "itl-csllab-112.sorint.localpippo", foundHd.Hostname)
		assert.Equal(t, utils.P("2019-11-19T15:38:58Z"), foundHd.CreatedAt)
	})

	hd2 := utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_dataservice_mongohostdata_02.json")
	m.InsertHostData(hd2)
	assert.NotEqual(m.T(), hd, hd2)

	m.T().Run("Should_find_hd2_more_inserts", func(t *testing.T) {
		foundHd, err := m.db.FindMostRecentHostDataOlderThan("itl-csllab-112.sorint.localpippo", utils.P("2019-12-20T15:38:58Z"))
		require.NoError(t, err)
		assert.False(t, foundHd.Archived)
		assert.Equal(t, "itl-csllab-112.sorint.localpippo", foundHd.Hostname)
		assert.Equal(t, utils.P("2019-12-19T15:38:58Z"), foundHd.CreatedAt)
	})
}

func (m *MongodbSuite) TestInsertAlert_Success() {
	_, err := m.db.InsertAlert(alert1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("alerts").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("alerts").FindOne(context.TODO(), bson.M{
		"_id": alert1.ID,
	})
	require.NoError(m.T(), val.Err())

	var out model.Alert
	val.Decode(&out)

	assert.Equal(m.T(), alert1, out)
}

func (m *MongodbSuite) TestFindOldCurrentHosts() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_dataservice_mongohostdata_01.json")
	m.InsertHostData(hd)

	hd2 := utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_dataservice_mongohostdata_03.json")
	m.InsertHostData(hd2)

	m.T().Run("Should not find any", func(t *testing.T) {
		hosts, err := m.db.FindOldCurrentHosts(utils.P("2019-11-18T16:38:58Z"))
		require.NoError(t, err)
		assert.Empty(t, hosts)
	})

	m.T().Run("Should find one", func(t *testing.T) {
		hosts, err := m.db.FindOldCurrentHosts(utils.P("2019-12-04T16:38:58Z"))
		require.NoError(t, err)

		assert.Len(t, hosts, 1)

		assert.Equal(m.T(), "itl-csllab-112.sorint.localpippo", hosts[0].Hostname)
		assert.False(m.T(), hosts[0].Archived)
		assert.Equal(m.T(), utils.P("2019-11-19T16:38:58Z"), hosts[0].CreatedAt)
	})

	m.T().Run("Should find two", func(t *testing.T) {
		hosts, err := m.db.FindOldCurrentHosts(utils.P("2020-01-14T15:38:58Z"))
		require.NoError(t, err)

		assert.Len(t, hosts, 2)

		sort.Slice(hosts, func(i, j int) bool {
			return hosts[i].ID.String() < hosts[j].ID.String()
		})

		assert.Equal(m.T(), "itl-csllab-112.sorint.localpippo", hosts[0].Hostname)
		assert.False(m.T(), hosts[0].Archived)
		assert.Equal(m.T(), utils.P("2019-11-19T16:38:58Z"), hosts[0].CreatedAt)

		assert.Equal(m.T(), "itl-csllab-223.sorint.localpippo", hosts[1].Hostname)
		assert.False(m.T(), hosts[1].Archived)
		assert.Equal(m.T(), utils.P("2020-01-13T15:38:58Z"), hosts[1].CreatedAt)
	})
}

func (m *MongodbSuite) TestExistNoDataAlert_SuccessNotExist() {
	_, err := m.db.InsertAlert(alert1)
	defer m.db.Client.Database(m.dbname).Collection("alerts").DeleteMany(context.TODO(), bson.M{})
	require.NoError(m.T(), err)

	exist, err := m.db.ExistNoDataAlertByHost("myhost")
	require.NoError(m.T(), err)

	assert.False(m.T(), exist)
}

func (m *MongodbSuite) TestExistNoDataAlert_SuccessExist() {
	_, err := m.db.InsertAlert(alert2)
	defer m.db.Client.Database(m.dbname).Collection("alerts").DeleteMany(context.TODO(), bson.M{})
	require.NoError(m.T(), err)

	exist, err := m.db.ExistNoDataAlertByHost("myhost")
	require.NoError(m.T(), err)

	assert.True(m.T(), exist)
}

func (m *MongodbSuite) TestDeleteNoDataAlertByHost_Success() {
	defer m.db.Client.Database(m.dbname).Collection("alerts").DeleteMany(context.TODO(), bson.M{})

	_, err := m.db.InsertAlert(alert1)
	require.NoError(m.T(), err)

	_, err = m.db.InsertAlert(alert3)
	require.NoError(m.T(), err)

	m.T().Run("There is still alert3", func(t *testing.T) {
		alert1Hostname := alert1.OtherInfo["hostname"].(string)
		err = m.db.DeleteNoDataAlertByHost(alert1Hostname)
		require.NoError(m.T(), err)

		val, err2 := m.db.Client.Database(m.dbname).Collection("alerts").
			Find(context.TODO(), bson.M{"alertCode": model.AlertCodeNoData})
		require.NoError(m.T(), err2)

		alerts := make([]model.Alert, 0)
		err3 := val.All(context.TODO(), &alerts)
		require.NoError(m.T(), err3)

		require.Equal(m.T(), 1, len(alerts))
		require.Equal(m.T(), alerts[0], alert3)
	})

	m.T().Run("There are no more alerts", func(t *testing.T) {
		alert3Hostname := alert3.OtherInfo["hostname"].(string)
		err = m.db.DeleteNoDataAlertByHost(alert3Hostname)

		require.NoError(m.T(), err)
		val, err2 := m.db.Client.Database(m.dbname).Collection("alerts").
			Find(context.TODO(), bson.M{"alertCode": model.AlertCodeNoData})
		require.NoError(m.T(), err2)

		alerts := make([]model.Alert, 0)
		err3 := val.All(context.TODO(), &alerts)
		require.NoError(m.T(), err3)
		require.Equal(m.T(), 0, len(alerts))
	})
}

func (m *MongodbSuite) TestDeleteAllNoDataAlerts_Success() {
	defer m.db.Client.Database(m.dbname).Collection("alerts").DeleteMany(context.TODO(), bson.M{})

	_, err := m.db.InsertAlert(alert1)
	require.NoError(m.T(), err)

	_, err = m.db.InsertAlert(alert3)
	require.NoError(m.T(), err)
	_, err = m.db.InsertAlert(alert4)
	require.NoError(m.T(), err)

	err = m.db.DeleteAllNoDataAlerts()
	require.NoError(m.T(), err)

	// Check that there are no more AlertCodeNoData alerts
	val, erro := m.db.Client.Database(m.dbname).Collection("alerts").
		Find(context.TODO(), bson.M{"alertCode": model.AlertCodeNoData})
	require.NoError(m.T(), erro)

	res := make([]model.Alert, 0)
	erro = val.All(context.TODO(), &res)
	require.NoError(m.T(), erro)
	require.Equal(m.T(), 0, len(res))

	// Check that there's still alert1
	val, erro = m.db.Client.Database(m.dbname).Collection("alerts").
		Find(context.TODO(), bson.M{})
	require.NoError(m.T(), erro)

	res = make([]model.Alert, 0)
	erro = val.All(context.TODO(), &res)
	require.NoError(m.T(), erro)
	require.Equal(m.T(), 1, len(res))
	require.Equal(m.T(), alert1, (res)[0])
}
