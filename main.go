package main

import (
    "context"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"

    "cloud.google.com/go/storage"
    "github.com/joho/godotenv"
    "google.golang.org/api/option"
)

func main() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    keyFilePath := os.Getenv("KEY_FILE_PATH")
    if keyFilePath == "" {
        log.Fatalf("KEY_FILE_PATH environment variable is not set")
    }

    bucketName := os.Getenv("BUCKET_NAME")
    if bucketName == "" {
        log.Fatalf("BUCKET_NAME environment variable is not set")
    }

    // Replace with the directory you want to upload
    dirPath := os.Args[1]

    ctx := context.Background()

    client, err := storage.NewClient(ctx, option.WithCredentialsFile(keyFilePath))
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer client.Close()

    bucket := client.Bucket(bucketName)

    // Read all files in the directory
    files, err := ioutil.ReadDir(dirPath)
    if err != nil {
        log.Fatalf("Failed to read directory: %v", err)
    }

    // Iterate over each file
    for _, file := range files {
        filePath := filepath.Join(dirPath, file.Name())
        destFileName := file.Name()

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
}