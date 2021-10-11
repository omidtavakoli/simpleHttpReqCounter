package main

import (
	"encoding/csv"
	"os"
	"testing"
	"time"
)

func TestWriteLogs(t *testing.T) {
	file := "test.csv"
	logs := [][]string{{"180"}}
	WriteLogs(&logs, file)

	data, err := os.Open(file)
	if err != nil {
		t.Error(err)
	}
	defer data.Close()
	lines, err := csv.NewReader(data).ReadAll()
	if err != nil {
		t.Error(err)
	}
	for _, line := range lines {
		if line[0] != "180" {
			t.Error("not found")
		}
	}
	os.Remove(file)
}

func TestLoadLogs(t *testing.T) {
	file := "test.csv"
	logs := [][]string{{"180"}}
	WriteLogs(&logs, file)

	var logs2 [][]string
	err := loadLogs(&logs2, file)
	if err != nil {
		t.Error(err)
	}
	if logs2[0][0] != "180" {
		t.Error("not found")
	}
	os.Remove(file)
}

func TestCalculateLastNSeconds(t *testing.T) {
	file := "test.csv"
	logs := [][]string{
		{time.Now().Add(-time.Second * 3).Format(Layout)},
		{time.Now().Add(-time.Second * 5).Format(Layout)},
		{time.Now().Add(-time.Second * 15).Format(Layout)},
		{time.Now().Add(-time.Second * 55).Format(Layout)},
		{time.Now().Add(-time.Second * 125).Format(Layout)},
		{time.Now().Add(-time.Second * 65).Format(Layout)},
	}
	WriteLogs(&logs, file)
	count := CalculateLastNSeconds(&logs)
	if count != 4 {
		t.Error("count is wrong")
	}
	os.Remove(file)
}
