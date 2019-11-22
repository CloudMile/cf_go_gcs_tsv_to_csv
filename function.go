package function

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/functions/metadata"
	"cloud.google.com/go/storage"
)

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Bucket         string    `json:"bucket"`
	Name           string    `json:"name"`
	Metageneration string    `json:"metageneration"`
	ResourceState  string    `json:"resourceState"`
	TimeCreated    time.Time `json:"timeCreated"`
	Updated        time.Time `json:"updated"`
}

// GcsTSVtoCSV is change tsv to csv
func GcsTSVtoCSV(ctx context.Context, e GCSEvent) error {
	_, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %s", err)
	}

	if !strings.HasSuffix(e.Name, ".tsv") {
		log.Printf("file %s is not a tsv", e.Name)
		return nil
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("GCS client err: %s", err)
	}

	bucket := client.Bucket(e.Bucket)
	reader, err := getFile(ctx, bucket, e.Name)
	if err != nil {
		return fmt.Errorf("reader err: %s", err)
	}
	defer reader.Close()

	fileName := strings.TrimSuffix(e.Name, ".tsv") + ".csv"
	writer := writeFile(ctx, bucket, fileName)
	defer writer.Close()

	toCSV(reader, writer)
	return nil
}

func getFile(ctx context.Context, bucket *storage.BucketHandle, objName string) (reader *storage.Reader, err error) {
	obj := bucket.Object(objName)
	reader, err = obj.NewReader(ctx)
	return
}

func writeFile(ctx context.Context, bucket *storage.BucketHandle, objName string) (writer *storage.Writer) {
	obj := bucket.Object(objName)
	writer = obj.NewWriter(ctx)
	return
}

func toCSV(reader *storage.Reader, writer *storage.Writer) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = '\t'

	csvWriter := csv.NewWriter(writer)
	csvWriter.Comma = ','

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			log.Printf("done!")
			break
		}
		csvWriter.Write(record)
	}
	csvWriter.Flush()
}
