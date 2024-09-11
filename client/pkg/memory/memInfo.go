package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"linux-client/models"
	"os"
	"strings"
	// "strings"
)

//	func memoryInfoInitialization(key string,value string, unit string) (models.MemoryInfoModels){
//		return models.MemoryInfoModels{
//			key:key,value:value,unit
//		}
//	}
var MemInfoData []models.MemoryInfoModels

func MemInfo() ([]models.MemoryInfoModels, error) {

	file, err := os.Open("/proc/meminfo")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return []models.MemoryInfoModels{}, err
	}
	defer file.Close()

	// Create a map to hold memory info entries
	memInfoMap := make(map[string]models.MemoryInfoModels)

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into key and value+unit
		parts := strings.Fields(line)
		// fmt.Println(len(parts),"parts::")
		if len(parts) < 2 {
			continue
		}

		key := strings.TrimSuffix(parts[0], ":") // Remove the colon from the key
		value := parts[1]                        // This is the numeric value
		unit := ""
		if len(parts) == 3 { // If there is a unit, it's the third part
			unit = parts[2]
		}

		// Convert value from string to uint64
		var val uint64
		fmt.Sscanf(value, "%d", &val)

		memInfoMap[key] = models.MemoryInfoModels{
			Key:   key,
			Value: val,
			Unit:  unit,
		}
	}


	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return []models.MemoryInfoModels{}, err
	}

	
	for _, info := range memInfoMap {
		// fmt.Printf("%s: %d %s\n", info.Key, info.Value, info.Unit)
		MemInfoData = append(MemInfoData, models.MemoryInfoModels{info.Key, info.Value, info.Unit})

	}
	// fmt.Println(MemInfoData)
	
	b,err := json.Marshal(MemInfoData)
	if err!=nil{
		return []models.MemoryInfoModels{},err
	}
	fmt.Println(string(b))

	return MemInfoData, nil
}
