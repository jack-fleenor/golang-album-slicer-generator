package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// func sliceTrack() error {

// 	return nil
// }

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

//	func progressBar(items []string) {
//		if len(items) == 0 {
//			return
//		}
//		bar := progressbar.Default(int64(len(items)))
//		for i := 0; i < len(items); i++ {
//			bar.Add(1)
//			// fmt.Println(items[i])
//			time.Sleep(40 * time.Millisecond)
//		}
//	}

func getFileDuration(path *string) error {
	// Replace "your_audio_file_path" with the path to your audio file

	// Run ffprobe command to get audio file duration
	cmd := exec.Command("ffprobe", "-i", *path, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running ffprobe command:", err)
		return err
	}

	// Extract duration from output
	durationStr := strings.TrimSpace(string(output))
	var duration float64
	_, err = fmt.Sscanf(durationStr, "%f", &duration)
	if err != nil {
		fmt.Println("Error parsing duration:", err)
		return err
	}

	fmt.Printf("Duration of %s: %.2f seconds\n", *path, duration)
	return nil
}

func timestampToSeconds(timestamp string) (int, error) {
	parts := strings.Split(timestamp, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid timestamp format: %s", timestamp)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}

	totalSeconds := hours*3600 + minutes*60 + seconds
	return totalSeconds, nil
}

func durationBetweenTimestamps(startTimestamp, endTimestamp string) (int, error) {
	startSeconds, err := timestampToSeconds(startTimestamp)
	if err != nil {
		return 0, err
	}

	endSeconds, err := timestampToSeconds(endTimestamp)
	if err != nil {
		return 0, err
	}

	duration := endSeconds - startSeconds
	return duration, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occurred while reading input. Please try again", err)
		return
	}
	fmt.Print("You entered: ", text)
	records := readCsvFile("./queue/test.csv")
	path := "./queue/test.mp3"
	if err := getFileDuration(&path); err != nil {
		fmt.Println("Error getting file duration:", err)
	}
	for _, record := range records {
		duration, err := durationBetweenTimestamps(record[2], record[3])
		if err != nil {
			fmt.Println("Error calculating duration between timestamps:", err)
		}
		fmt.Println("This song: ", record[0], " by ", record[1], " will start at: ", record[2], " and end at: ", record[3], " with a duration of: ", strconv.Itoa(duration), " seconds")
	}
}
