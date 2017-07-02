package database

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/imgasm/server/config"
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"
)

var (
	datastoreClient *datastore.Client
	doOnce          sync.Once
)

func init() {
	newDatastoreClient()
}

func newDatastoreClient() {
	ctx := context.Background()
	viper.SetConfigName("config")
	viper.AddConfigPath(config.ViperConfigPath)
	err := viper.ReadInConfig()
	if err != nil {
		// todo: better error handling
		log.Println("Fatal error config file: %s \n", err)
	}
	env := viper.GetString("environment")
	switch env {
	case "dev":
		os.Setenv("DATASTORE_DATASET", viper.GetString("gcloud.datastore.local.dataset"))
		os.Setenv("DATASTORE_EMULATOR_HOST_PATH", viper.GetString("gcloud.datastore.local.emulator_host_path"))
		os.Setenv("DATASTORE_HOST", viper.GetString("gcloud.datastore.local.host"))
		os.Setenv("DATASTORE_EMULATOR_HOST", viper.GetString("gcloud.datastore.local.emulator_host"))
		os.Setenv("DATASTORE_PROJECT_ID", viper.GetString("gcloud.datastore.local.project_id"))
		datastoreClient, err = datastore.NewClient(ctx, viper.GetString("gcloud.datastore.local.project_id"))
		if err != nil {
			// todo: better error handling
			log.Println("err creating datastore client:", err)
		}
	case "prod":
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", viper.GetString("gcloud.datastore.production.credentials_file"))
		datastoreClient, err = datastore.NewClient(ctx, viper.GetString("gcloud.datastore.production.project_id"))
		if err != nil {
			// todo: better error handling
			log.Println("err creating datastore client:", err)
		}
	default:
		// todo: better error handling
		log.Println("environment not set")
	}
}

func DatastoreClient() *datastore.Client {
	if datastoreClient == nil {
		doOnce.Do(newDatastoreClient)
	}
	return datastoreClient
}
