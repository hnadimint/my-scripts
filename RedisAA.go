package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	redisURL     = "http://download.redis.io/releases/redis-6.2.6.tar.gz"
	downloadPath = "/tmp/redis.tar.gz"
	extractPath  = "/tmp/redis"
)

func main() {
	// Step 1: Download Redis
	fmt.Println("Downloading Redis...")
	err := downloadFile(redisURL, downloadPath)
	if err != nil {
		fmt.Printf("Error downloading Redis: %v\n", err)
		return
	}
	fmt.Println("Download completed.")

	// Step 2: Extract the tar.gz file
	fmt.Println("Extracting Redis...")
	err = extractTarGz(downloadPath, extractPath)
	if err != nil {
		fmt.Printf("Error extracting Redis: %v\n", err)
		return
	}
	fmt.Println("Extraction completed.")

	// Step 3: Build Redis
	fmt.Println("Building Redis...")
	err = buildRedis(extractPath)
	if err != nil {
		fmt.Printf("Error building Redis: %v\n", err)
		return
	}
	fmt.Println("Build completed.")

	// Step 4: Run Redis server
	fmt.Println("Starting Redis server...")
	err = runRedis(extractPath)
	if err != nil {
		fmt.Printf("Error running Redis: %v\n", err)
		return
	}
	fmt.Println("Redis server started.")
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractTarGz(gzipPath, destPath string) error {
	file, err := os.Open(gzipPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tarReader := tar.NewReader(gzr)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(destPath, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(path)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
	return nil
}

func buildRedis(path string) error {
	cmd := exec.Command("make")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runRedis(path string) error {
	cmd := exec.Command(filepath.Join(path, "src/redis-server"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}