package backend

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db       *gorm.DB
	nativeDB *sql.DB
)

func openDatabase() (err error) {
	// TODO: validate this dns... if using timestamps add: "?parseTime=true"
	dsn := viper.GetString("database.dsn")

	// disable gorm logger for now...
	gormConfig := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}

	// init the correct db driver
	switch driver := viper.GetString("database.driver"); driver {
	case "mysql":
		log.Infof("opening database connection using driver: %s", driver)
		log.Debugf("dsn: %s", dsn)
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
	default:
		log.Fatalf("database driver is unsupported: %s", driver)
	}

	// used for DB health check
	if err == nil {
		nativeDB, err = db.DB()
	}
	return
}

func migrateDatabaseSchema() error {
	return db.AutoMigrate(&Client{})
}

func dbHealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return nativeDB.PingContext(ctx)
}

func dbCreateClient(client *Client) (err error) {
	result := db.Create(&client)
	if result.Error != nil {
		err = fmt.Errorf("failed to create new client: %v", result.Error)
		log.Warningf(err.Error())
		return
	}
	log.Debugf("created client with id: %s", client.ID)
	return
}

func dbGetClient(id string) (client *Client, err error) {
	result := db.Find(&client, "id = ?", id)
	if result.Error != nil {
		err = fmt.Errorf("failed to find client (%s): %v", id, result.Error)
		log.Warningf(err.Error())
		return
	}
	return
}

func dbUpdateClient(client *Client) (err error) {
	result := db.Save(&client)
	if result.Error != nil {
		err = fmt.Errorf("failed to update client (%s): %v", client.ID, result.Error)
		log.Warningf(err.Error())
		return
	}
	return
}

func dbDeleteClient(id string) (err error) {
	result := db.Delete(Client{}, id)
	if result.Error != nil {
		err = fmt.Errorf("failed to delete client (%s): %v", id, result.Error)
		log.Warningf(err.Error())
		return
	}
	log.Debugf("delete client with id: %s", id)
	return
}

func dbGetAllClients() (clients []Client, err error) {
	result := db.Find(&clients)
	if result.Error != nil {
		err = fmt.Errorf("failed to get all clients: %v", result.Error)
		log.Warningf(err.Error())
		return
	}
	return
}
