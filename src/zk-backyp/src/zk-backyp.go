package main

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"

	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.Println("ZOOKEEPER BACKYP DAEMON started")
	// GETTING ENV VARS
	backup_interval := os.Getenv("ZK_BACKYP_INTERVAL")
	s3_bucket := os.Getenv("ZK_BACKYP_S3_BUCKET")
	s3_endpoint := os.Getenv("ZK_BACKYP_S3_ENDPOINT")
	zk_txlog_path := os.Getenv("ZK_TXLOG_PATH")
	zk_snapshot_path := os.Getenv("ZK_SNAPSHOT_PATH")

	// DEFINING AWS REGION
	var region = aws.Region{}
	region.S3Endpoint = s3_endpoint
	region.S3BucketEndpoint = ""
	region.S3LocationConstraint = true
	region.S3LowercaseBucket = true

	// LET THE BACKUP BEGIN!
	interval, errIntv := strconv.ParseInt(backup_interval, 10, 64)
	if errIntv != nil {
		log.Fatal("Could not parse environment variable ZK_BACKYP_INTERVAL")
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	quit := make(chan struct{})

	// execute one backup now
	archivefilename := os.Getenv("ZK_BACKYP_PREFIX") + "-" + time.Now().Format("20060102T150424") + ".tar.gz"
	go executeBackup(zk_txlog_path, zk_snapshot_path, archivefilename, s3_bucket, region)

	// execute the others later by interval
	go func() {
		for {
			select {
			case <-ticker.C:
				archivefilename := os.Getenv("ZK_BACKYP_PREFIX") + "-" + time.Now().Format("20060102T150424") + ".tar.gz"
				executeBackup(zk_txlog_path, zk_snapshot_path, archivefilename, s3_bucket, region)

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	select {} // block forever
}

func executeBackup(zk_txlog_path string, zk_snapshot_path string, archivefilename string, s3_bucket string, region aws.Region) {
	log.Println("Executing Backup")
	// STOPPING ZOOKEEPER
	err := stopZookeeper()
	if err != nil {
		log.Fatal(err)
	}

	// COMPRESS SNAPSHOTS AND TX LOGS
	if zk_txlog_path == zk_snapshot_path {
		err = tar(archivefilename, []string{zk_txlog_path})
	} else {
		err = tar(archivefilename, []string{zk_txlog_path, zk_snapshot_path})
	}
	if err != nil {
		log.Fatal(err)
	}

	// STARTING ZOOKEEPER
	err = startZookeeper()
	if err != nil {
		log.Fatal(err)
	}

	// UPLOADING BACKUP TO S3
	err = uploadToS3(archivefilename, s3_bucket, region)
	if err != nil {
		log.Fatal(err)
	}

	err = deleteFile(archivefilename)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Backup successfully done")
}

func stopZookeeper() error {
	log.Println("Stopping Zookeeper...")
	err := exec.Command("monit", "stop", "zookeeper").Run()

	if err != nil {
		return err
	}

	timeout := 60 // Seconds
	zk_running := true
	for zk_running && timeout > 0 && err == nil {
		var out []byte
		out, err = exec.Command("monit","summary").Output()

		if err != nil {
			return err
		}

		lines := strings.Split(string(out[:]),"\n")
		for _, line := range lines {
			if strings.Contains(line,"zookeeper") && strings.Contains(line, "ot monitored") {
				zk_running = false
			}
		}

		time.Sleep(1 * time.Second)
		timeout = timeout - 1
	}

	if zk_running {
		if err != nil {
			return err
		}
		return errors.New("Stopping zookeeper timed out!")
	}

	log.Println("Zookeeper is stopped")
	return nil
}

func startZookeeper() error {
	log.Println("Starting Zookeeper...")
	err := exec.Command("monit","start","zookeeper").Run()

	if err != nil {
		return err
	}

	timeout := 60 // seconds
	zk_running := false
	for zk_running == false && timeout > 0 && err == nil {
		var out []byte
		out, err = exec.Command("monit","summary").Output()

		if err != nil {
			return err
		}

		lines := strings.Split(string(out[:]),"\n")
		for _, line := range lines {
			if strings.Contains(line,"zookeeper") && strings.Contains(line, "unning") {
				zk_running = true
			}
		}

		time.Sleep(1 * time.Second)
		timeout = timeout - 1
	}

	if !zk_running {
		return errors.New("Starting zookeeper timed out!")
	}

	log.Println("Zookeeper is running")
	return nil
}

func tar(tarFile string, files []string) error {
	log.Println("creating TAR archive...")
	cmd := "tar czf " + tarFile
	for _, f := range files {
		cmd += " " + f
	}

	err := exec.Command("tar","czf",tarFile,files[0]).Run()
	if err != nil {
		return err
	}

	log.Println("tar archiving done!")
	return nil
}

func uploadToS3(file string, s3_bucket string, awsRegion aws.Region) error {
	log.Println("Uploading to S3...")
	auth, err := aws.EnvAuth()
	if err != nil {
		return err
	}
	client := s3.New(auth, awsRegion)
	bucket := client.Bucket(s3_bucket)

	r, errF := os.Open(file)
	fi, _ := r.Stat()
	length := fi.Size()

	if errF != nil {
		return errF
	}

	path := file
	perm := s3.ACL("private")
	contType := "application/x-compressed"

	err = bucket.PutReader(path, r, length, contType, perm)
	if err != nil {
		return err
	}

	log.Println("Upload successful")
	return nil
}

func deleteFile(filename string) error {
	log.Println("Deleting file "+filename)

	err := os.Remove(filename)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("deletion successful")
	return nil
}
