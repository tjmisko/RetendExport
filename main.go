package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func getAllRetendFiles(directoryString string) {
	directoryContents, err := os.ReadDir(directoryString)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range directoryContents {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), ".retend") || strings.Contains(file.Name(), ".schedule") {
			processRetendFile(directoryString + "/" + file.Name())
		}
	}

	if err != nil {
		fmt.Println(err)
	}
}

func processRetendFile(filepath string) {
	fmt.Println("Processing file:", filepath)
	fileContents, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer fileContents.Close()

	scanner := bufio.NewScanner(fileContents)

    var lastCategory, lastTitle, runningNotes string = "", "", ""
    var lastStart time.Time
    var duration int = 15
    var timeBlockList []TimeBlock

	for scanner.Scan() {
		lineText := scanner.Text()
		splitList := strings.SplitN(lineText, " | ", 3)
		if len(splitList) != 3 {
			log.Fatal("Invalid line in file:", filepath)
		}
		category, title, notes := processRetendInfo(splitList[2])
        if lastStart.IsZero() {
            lastStart = processDateTimeString(splitList[0], splitList[1])
            lastCategory = category
            lastTitle = title
            runningNotes = notes
            continue
        } else if lastCategory == category && lastTitle == title {
            runningNotes += notes
            duration += 15
        } else { 
            startTime := processDateTimeString(splitList[0], splitList[1])
            newBlock := TimeBlock{
                StartTime:  lastStart,
                EndTime:    lastStart.Add(time.Duration(duration) * time.Minute),
                Duration:   time.Duration(duration) * time.Minute,
                Category:   lastCategory,
                Title:      lastTitle,
                Notes:      runningNotes,
            }
            timeBlockList = append(timeBlockList, newBlock)
            duration = 15
            lastCategory = category
            lastTitle = title
            runningNotes = notes
            lastStart = startTime
        }
	}
    finalBlock := TimeBlock{
        StartTime:  lastStart,
        EndTime:    lastStart.Add(time.Duration(duration) * time.Minute),
        Duration:   time.Duration(duration) * time.Minute,
        Category:   lastCategory,
        Title:      lastTitle,
        Notes:      runningNotes,
    }
    timeBlockList = append(timeBlockList, finalBlock)
    for _, block := range timeBlockList {
        fmt.Println("====================================")
        fmt.Println("StartTime:", block.StartTime)
        fmt.Println("EndTime:", block.EndTime)
        fmt.Println("Duration:", block.Duration)
        fmt.Println("Category:", block.Category)
        fmt.Println("Title:", block.Title)
        fmt.Println("Notes:", block.Notes)

    }
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

type TimeBlock struct {
    StartTime time.Time
    EndTime time.Time
    Duration time.Duration
    Category string
    SubCategory string
    Title string
    Notes string
}

func processDateTimeString(dateString string, timeString string) time.Time {
	startTime, location, err := parseTime(timeString)
	if err != nil {
        fmt.Println("Error parsing time:", timeString)
	}
	startDate, err := time.ParseInLocation("2006-01-02", dateString, location)
	if err != nil {
		fmt.Println("Error parsing date:", dateString)
	}
	return time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startTime.Hour(), startTime.Minute(), startTime.Second(), 0, location)
}

func parseTime(timeString string) ( time.Time, *time.Location, error ) {
    timeList := strings.Split(timeString, " ")
    startTime, err := time.Parse("15:04", timeList[0])
    if err != nil {
        log.Fatal("Invalid time format:", timeString)
    }
	if len(timeList) == 1 {
        location, err := extractLocation("")
        return startTime, location, err
	} else if len(timeList) > 2 {
        location, err := extractLocation("")
		return startTime, location, err
	}
    location, err := extractLocation(timeList[1])
    return startTime, location, err
}

func extractLocation(timeZoneString string) (*time.Location, error) {
	var locationName string = ""
	switch timeZoneString {
	case "EDT":
		locationName = "America/New_York"
	case "EST":
		locationName = "America/New_York"
	case "PDT":
		locationName = "America/Los_Angeles"
	case "PST":
		locationName = "America/Los_Angeles"
	case "SST":
		locationName = "Singapore"
	default:
		locationName = "America/Los_Angeles"
	}
	location, err := time.LoadLocation(locationName)
	return location, err
}

func processRetendInfo(timeBlockString string) (string, string, string) {
    categoryStart := strings.Index(timeBlockString, "<{")
    categoryEnd := strings.Index(timeBlockString, "}")
    titleStart := strings.Index(timeBlockString, "(")
    titleEnd := strings.Index(timeBlockString, ")>")
    category := timeBlockString[categoryStart + 2:categoryEnd] 
    title := timeBlockString[titleStart + 1:titleEnd]
    notes := timeBlockString[titleEnd + 2:]
    if len(notes) < 2 {
        notes = ""
    } else if notes[0:2] == " \"" {
        notes = notes[2:]
    } else {
        log.Fatal("Invalid notes format:", notes)
    }
    return category, title, notes
}
        
func main() {
	retendDirectoryPaths := os.Args[1:]
	for _, retendDirectoryPath := range retendDirectoryPaths {
		getAllRetendFiles(retendDirectoryPath)
	}
}
