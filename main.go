package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Layout   = "2006-01-02 15:04:05"
	fileName = "logs.csv"
)

func main() {
	Worker()
}

func Worker() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	var logs [][]string

	err := loadLogs(&logs, fileName)
	if err != nil {
		close(done)
	}
	go handler(&logs)

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		WriteLogs(&logs, fileName)
		cancel()
	}()
	ctx.Done()

	log.Print("Server Exited Gracefully")
}

func handler(logs *[][]string) {
	log.Print("Server Started")
	http.HandleFunc("/fetch", func(w http.ResponseWriter, r *http.Request) {
		*logs = append(*logs, []string{time.Now().Format(Layout)})
		count := CalculateLastNSeconds(logs)
		msg := fmt.Sprintf("requests in past %dsec:%d", 60, count)
		fmt.Println(msg)
		json.NewEncoder(w).Encode(count)
	})
	http.ListenAndServe(":8080", nil)
}

func loadLogs(logs *[][]string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return err
	}

	*logs = lines
	return nil
}

func WriteLogs(logs *[][]string, file string) {
	fmt.Printf("logs len for writing:%d\n", len(*logs))
	os.Remove(file)
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(f)
	err = w.WriteAll(*logs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("persisting data done")
}

func CalculateLastNSeconds(logs *[][]string) int {
	count := 0
	fmt.Printf("calculating request count since: %s\n", time.Now().Add(-time.Second*60))
	for _, logLine := range *logs {
		now, err := time.Parse(Layout, time.Now().Format(Layout))
		if err != nil {
			fmt.Printf("error parsing now")
		}
		logTime, err := time.Parse(Layout, logLine[0])
		if err != nil {
			fmt.Printf("error parsing logTime")
		}
		if now.Sub(logTime).Seconds() < 60 {
			count += 1
		}
		//else{
		//rotate logs if required
		//	*logs = remove(*logs, i)
		//}
	}
	return count
}

func remove(slice [][]string, s int) [][]string {
	return append(slice[:s], slice[s+1:]...)
}
