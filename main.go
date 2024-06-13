package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "os"

    "cloud.google.com/go/storage"
    "github.com/joho/godotenv"
    "google.golang.org/api/option"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    if len(os.Args) != 2 {
        log.Fatalf("Usage: go run main.go <filePath>")
    }
    // Replace with the path to your service account key file
    keyFilePath := os.Getenv("KEY_FILE_PATH")
    if keyFilePath == "" {
        log.Fatalf("KEY_FILE_PATH environment variable is not set")
    }

    // Replace with your GCS bucket name
    bucketName := os.Getenv("BUCKET_NAME")
    if bucketName == "" {
        log.Fatalf("BUCKET_NAME environment variable is not set")
    }

    // Replace with the file you want to upload
    filePath := os.Args[1]
    destFileName := filePath

    ctx := context.Background()

    // Creates a client using the service account key file
    client, err := storage.NewClient(ctx, option.WithCredentialsFile(keyFilePath))
    if err != nil {
            log.Fatalf("Failed to create client: %v", err)
    }
    defer client.Close()

    bucket := client.Bucket(bucketName)
    object := bucket.Object(destFileName)

    f, err := os.Open(filePath)
    if err != nil {
            log.Fatalf("Failed to open file: %v", err)
    }
    defer f.Close()

    wc := object.NewWriter(ctx)
    if _, err = io.Copy(wc, f); err != nil {
            log.Fatalf("Failed to copy file to bucket: %v", err)
    }
    if err := wc.Close(); err != nil {
            log.Fatalf("Failed to close bucket writer: %v", err)
    }

    fmt.Printf("%s uploaded to %s\n", filePath, bucketName)
}