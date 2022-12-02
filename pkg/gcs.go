package cloudlock

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

var (
	GcloudProjectId      string
	GcloudRegion         string
	GcloudServiceAccount string
	CloudlockBucket      string
)

func init() {
	GcloudProjectId = os.Getenv("GCP_PROJECT_ID")
	GcloudRegion = os.Getenv("GCP_REGION")
	GcloudServiceAccount = os.Getenv("GCP_SERVICE_ACCOUNT")
	CloudlockBucket = os.Getenv("CLOUDLOCK_BUCKET")
	if CloudlockBucket == "" {
		CloudlockBucket = "cloudlock"
	}
}

// GCS is an implementation of the Locker interface backed by Google Cloud Storage.
type GCS struct {
	ctx    context.Context
	client *storage.Client

	bucketName  string
	bucket      *storage.BucketHandle
	bucketAttrs *storage.BucketAttrs
}

// NewGCS returns a new GCS instance.
func NewGCS(bucketName string) *GCS {
	ctx := context.Background()
	// todo(jsam): Implement a timeout context.
	// ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	// defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Using project %s", GcloudProjectId)
	log.Printf("Using bucket %s", bucketName)

	bucket := client.Bucket(bucketName)
	bucketAttrs, err := bucket.Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		log.Println("Bucket not found - creating bucket")
		if err := bucket.Create(ctx, GcloudProjectId, nil); err != nil {
			_err := fmt.Errorf("Bucket(%q).Create: %v", bucketName, err)
			log.Fatal(_err)
		}
	}

	return &GCS{
		ctx:         ctx,
		client:      client,
		bucketName:  bucketName,
		bucket:      bucket,
		bucketAttrs: bucketAttrs,
	}
}

// Teardown cleans up the GCS connections.
func (gcs *GCS) Teardown() error {
	return gcs.client.Close()
}

// Lock creates a lock file in the GCS.
func (gcs *GCS) Lock(lockname string, payload []byte, wait bool) error {
	o := gcs.bucket.Object(lockname)

	objectAttrs, err := o.Attrs(gcs.ctx)
	if err == storage.ErrObjectNotExist {
		log.Println("Lockfile not found")
		o = o.If(storage.Conditions{DoesNotExist: true})
		wc := o.NewWriter(gcs.ctx)
		log.Println("Acquiring lockfile")
		if _, err := wc.Write(payload); err != nil {
			return err
		}
		if err := wc.Close(); err != nil {
			return err
		}
		return nil
	}

	if err != nil && err != storage.ErrObjectNotExist {
		log.Fatal(err)
	}

	if wait && objectAttrs != nil {
		log.Println("Lockfile already exists")
		isReleased, returnValue := gcs.waitForRelease(lockname, payload)
		// todo(jsam): Implement a timeout context and lock expiry.
		if isReleased {
			return returnValue
		}
	} else if !wait && objectAttrs != nil {
		os.Exit(2)
	}

	return nil
}

func (gcs *GCS) waitForRelease(lockname string, payload []byte) (bool, error) {
	for {
		log.Printf("Waiting for lockfile `%s` to be released\n", lockname)
		time.Sleep(5 * time.Second)
		rc, err := gcs.bucket.Object(lockname).NewReader(gcs.ctx)
		if err == storage.ErrObjectNotExist {
			log.Println("Lockfile released")
			return true, gcs.Lock(lockname, payload, false)
		}
		defer rc.Close()
		body, err := io.ReadAll(rc)
		if err != nil {
			log.Println("Error reading lockfile")
		}
		log.Printf("Lockfile still exists: [%s]", body)
	}
}

// Unlock removes a lock file from the GCS.
func (gcs *GCS) Unlock(lockname string) ([]byte, error) {
	obj := gcs.bucket.Object(lockname)

	attrs, err := obj.Attrs(gcs.ctx)
	if err == storage.ErrObjectNotExist {
		log.Println("Lockfile not found - nothing to unlock")
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	obj = obj.If(storage.Conditions{GenerationMatch: attrs.Generation})

	rc, err := obj.NewReader(gcs.ctx)
	if err == storage.ErrObjectNotExist {
		log.Println("Lockfile not found - nothing to unlock")
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	defer rc.Close()

	body, err := io.ReadAll(rc)
	if err != nil {
		log.Println("Error reading lockfile")
		return nil, err
	}

	if err := obj.Delete(gcs.ctx); err != nil {
		log.Println("Error deleting lockfile")
		return nil, err
	}

	log.Printf("Lockfile contents deleted: [%s]", body)
	return body, nil
}

// IsLocked returns true if the lock file exists in the GCS.
func (gcs *GCS) IsLocked(lockname string) (bool, error) {
	_, err := gcs.bucket.Object(lockname).NewReader(gcs.ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}
