package database

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/imgasm/server/config"
	"github.com/spf13/viper"
	"os"
	"sync"
)

var (
	datastoreClient *datastore.Client
	doOnce          sync.Once
)

func init() {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(config.ViperConfigPath)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %s \n", err))
	}
	datastoreClient, err = configureDatastoreClient(v)
	if err != nil {
		panic(fmt.Errorf("fatal error configuring datastore client: %s \n", err))
	}
}

func newDatastoreClient() {
	runtimeViper := config.RuntimeViper
	var err error
	datastoreClient, err = configureDatastoreClient(runtimeViper)
	if err != nil {
		// todo: error handling
	}

}

func configureDatastoreClient(v *viper.Viper) (*datastore.Client, error) {
	var projectID string

	env := v.GetString("environment")
	switch env {
	case "dev":
		os.Setenv("DATASTORE_DATASET", v.GetString("gcloud.datastore.local.dataset"))
		os.Setenv("DATASTORE_EMULATOR_HOST_PATH", v.GetString("gcloud.datastore.local.emulator_host_path"))
		os.Setenv("DATASTORE_HOST", v.GetString("gcloud.datastore.local.host"))
		os.Setenv("DATASTORE_EMULATOR_HOST", v.GetString("gcloud.datastore.local.emulator_host"))
		os.Setenv("DATASTORE_PROJECT_ID", v.GetString("gcloud.datastore.local.project_id"))
		projectID = v.GetString("gcloud.datastore.local.project_id")
	case "prod":
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", v.GetString("gcloud.datastore.production.credentials_file"))
		projectID = v.GetString("gcloud.datastore.production.project_id")
	default:
		// todo: error handling
	}

	ctx := context.Background()
	return datastore.NewClient(ctx, projectID)
}

func DatastoreClient() *datastore.Client {
	if datastoreClient == nil {
		doOnce.Do(newDatastoreClient)
	}
	return datastoreClient
}
