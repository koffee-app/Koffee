package store

import "koffee/internal/models"

import "koffee/internal/auth"

import "github.com/jmoiron/sqlx"

import "sync"

type driverConnection struct {
	Driver       *models.Driver
	connectionID uint32
	// todo More will come like ubication and things like that...
}

type driverID uint32

type driverStorageConnected struct {
	mutex                 sync.RWMutex
	driversModelConnected map[driverID]models.Driver
	driversConnected      map[uint32]driverConnection
}

var driversModelConnectedV map[driverID]models.Driver = make(map[driverID]models.Driver, 10000)
var driversConnectedV map[uint32]driverConnection = make(map[uint32]driverConnection, 10000)
var driverStorage driverStorageConnected = driverStorageConnected{driversModelConnected: driversModelConnectedV, driversConnected: driversConnectedV}

// DriverModelToConnection .
func driverModelToConnection(d *models.Driver) driverConnection {
	return driverConnection{Driver: d}
}

// Connection stores a new connection to drivers and returns if it was succesful...
func Connection(db *sqlx.DB, token string, connectionID uint32) bool {
	parsed, errParsing := auth.VerifyToken(token)
	if errParsing != nil {
		return false
	}
	user := models.UserJWTToUser(parsed)
	driver := models.GetDriverByID(db, user.UserID)
	if driver == nil {
		return false
	}
	driverStorage.mutex.Lock()
	defer driverStorage.mutex.Unlock()
	driverStorage.driversModelConnected[driverID(driver.UserID)] = *driver
	d, _ := driverStorage.driversModelConnected[driverID(driver.UserID)]
	driverConnect := driverModelToConnection(&d)
	driverConnect.connectionID = connectionID
	driverStorage.driversConnected[connectionID] = driverConnect
	return true
}

// DeleteConnection deletes a connection from the store
func DeleteConnection(connectionID uint32) {
	driverStorage.mutex.Lock()
	defer driverStorage.mutex.Unlock()
	delete(driverStorage.driversConnected, connectionID)
	// todo maybe store the number of connections of each driver so we can check when to delete a driver.
	// drivers,
}
