package bucket

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/imgasm/server/config"
	"github.com/spf13/viper"
	"sync"
)

var (
	storageBucket *storage.BucketHandle
	doOnce        sync.Once
)

func init() {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(config.ViperConfigPath)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %s \n", err))
	}
	storageBucket, err = configureStorageBucket(v)
	if err != nil {
		panic(fmt.Errorf("fatal error configuring storage bucket: %s \n", err))
	}
}

func newStorageBucket() {
	runtimeViper := config.RuntimeViper
	var err error
	storageBucket, err = configureStorageBucket(runtimeViper)
	if err != nil {
		// todo: error handling
	}
}

func configureStorageBucket(v *viper.Viper) (*storage.BucketHandle, error) {
	var bucketID string

	env := v.GetString("environment")
	switch env {
	case "dev":
		// todo (currently no emulator available)
		return nil, nil
	case "prod":
		bucketID = v.GetString("gcloud.storage.bucket.name")
	default:
		// todo: error handling
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.Bucket(bucketID), nil
}

func StorageBucket() *storage.BucketHandle {
	if storageBucket == nil {
		doOnce.Do(newStorageBucket)
	}
	return storageBucket
}
