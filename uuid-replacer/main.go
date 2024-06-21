// This utility replaces all UUIDs found in a file with a new deterministic UUID that
// is derived from the previous UUID plus a salt that gets randomly generated 
// once per execution.
// The purpose of this was to quickly duplicate a database dump so that
// we can restore x amount of desired objects directly to the database
// without needing to use the API, which is **much** quicker when dealing with
// millions of lines.
// Although its single threaded, I measured the performance
// of this util to about 260k lines per second
package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha1"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// generate a deterministic UUID based on an input UUID and salt
func generateDeterministicUUID(inputUUID uuid.UUID, salt []byte) uuid.UUID {
	hasher := sha1.New()
	hasher.Write(salt)         // Write the static salt bytes
	hasher.Write(inputUUID[:]) // Write the input UUID
	hash := hasher.Sum(nil)

	var uuidBytes [16]byte
	copy(uuidBytes[:], hash[:16])
	uuidBytes[6] = (uuidBytes[6] & 0x0f) | 0x50 // Version 5
	uuidBytes[8] = (uuidBytes[8] & 0x3f) | 0x80 // Variant is 10

	return uuidBytes
}

func main() {
	inputFileName := flag.String("input", "data/input-sql-dumps/dump_instance_small.sql", "sql dump of instance")
	outputFileName := flag.String("output", "data/output-sql-dumps/instance.sql", "output sql dump file")
	flag.Parse()
	fmt.Printf("Input File: %s\nOutput File: %s\n", *inputFileName, *outputFileName)

	file, err := os.Open(*inputFileName)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	outputFile, err := os.Create(*outputFileName)
	if err != nil {
		log.Fatalf("Failed to create output file: %s", err)
	}
	defer outputFile.Close()

	// Regular expression to match UUIDs
	uuidRegex := regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)

	// UUIDs to skip (taken from the platform dump file)
	skipUUIDs := map[string]struct{}{
		// platform id
		"c398ef31-15d9-4208-af18-e56a67e3e131": {},
		// coa template ids
		"aa589154-8214-49b3-9629-fcba297240ce": {},
		"71a0201a-1bd2-437b-8312-6a4ed5005691": {},
		"69034792-3241-49b4-b4e7-189d2405f4af": {},
		"81b3232f-ada2-4db2-8e9f-74068ac0cfef": {},
		// tag group ids
		"9c045f2d-0b69-463f-93f4-e28558ac09d1": {},
		"cf3ad1b8-e442-4ff6-88f0-74bc6661187d": {},
	}

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outputFile)

	// generate a random salt so that this instance UUIDs won't conflict with another run
	// this salt is unique for each execution (not each line!)
	salt := make([]byte, 8)
	rand.Read(salt)

	for scanner.Scan() {
		line := scanner.Text()
		matches := uuidRegex.FindAllString(line, -1)
		for _, match := range matches {
			if _, shouldSkip := skipUUIDs[match]; shouldSkip {
				continue
			}
			originalUUID, err := uuid.Parse(match)
			if err != nil {
				log.Fatalf("Failed to parse UUID: %s", err)
			}
			deterministicUUID := generateDeterministicUUID(originalUUID, salt)
			line = strings.ReplaceAll(line, match, deterministicUUID.String())
		}
		writer.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input file: %s", err)
	}
	writer.Flush()

	fmt.Printf("UUID replacement completed. Check '%s' for results.\n", *outputFileName)
}
