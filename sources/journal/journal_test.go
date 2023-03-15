package journal

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"log-aggregator/types"
)

func TestEntryToRecord(t *testing.T) {
	entry := &JournalEntry{
		RealtimeTimestamp: 18446744073709551615,
		Fields: map[string]string{
			"_SOURCE_REALTIME_TIMESTAMP": "abcdefghejgjslfk",
			"MY_FIELD":                   "foobar",
		},
		Cursor: "mycursor",
	}
	record := entryToRecord(entry)

	if record.Cursor != types.Cursor("mycursor") {
		t.Errorf("Expected cursor to be mycursor, but got %s", record.Cursor)
	}

	if record.Time.Unix() != 18446744073709 || record.Time.Nanosecond() != 551615000 {
		t.Errorf("Expected time to equal 18446744073709.551615000, but got %d.%d",
			record.Time.Unix(), record.Time.Nanosecond())
	}

	if _, ok := record.Fields["_SOURCE_REALTIME_TIMESTAMP"]; ok {
		t.Errorf("Expected field _SOURCE_REALTIME_TIMESTAMP to be removed, but it is stil present")
	}

	if val, ok := record.Fields["MY_FIELD"]; !ok || val != "foobar" {
		t.Errorf("Expected field MY_FIELD to be foobar, but got '%s'", val)
	}
}

func TestEntryToTime(t *testing.T) {
	entry := &JournalEntry{
		RealtimeTimestamp: 18446744073709551615,
	}

	entryTime := entryToTime(entry)
	if entryTime.Unix() != 18446744073709 || entryTime.Nanosecond() != 551615000 {
		t.Errorf("Expected time to equal 18446744073709.551615000, but got %d.%d",
			entryTime.Unix(), entryTime.Nanosecond())
	}

	// Prioritize _SOURCE_REALTIME_TIMESTAMP field if its available.
	entry = &JournalEntry{
		Fields: map[string]string{
			"_SOURCE_REALTIME_TIMESTAMP": "18446744073709551615",
		},
		RealtimeTimestamp: 1234567890,
	}

	entryTime = entryToTime(entry)
	if entryTime.Unix() != 18446744073709 || entryTime.Nanosecond() != 551615000 {
		t.Errorf("Expected time to equal 18446744073709.551615000, but got %d.%d",
			entryTime.Unix(), entryTime.Nanosecond())
	}

	// Fall back to RealtimeTimestamp if parsing _SOURCE_REALTIME_TIMESTAMP fails
	entry = &JournalEntry{
		Fields: map[string]string{
			"_SOURCE_REALTIME_TIMESTAMP": "18446abcd744073709551615",
		},
		RealtimeTimestamp: 1234567890,
	}

	entryTime = entryToTime(entry)
	if entryTime.Unix() != 1234 || entryTime.Nanosecond() != 567890000 {
		t.Errorf("Expected time to equal 1234.567890000, but got %d.%d",
			entryTime.Unix(), entryTime.Nanosecond())
	}
}

func TestIgnoreSystemUnitConfiguration(t *testing.T) {
	// create maps of ignored units
	sysUnits := map[string]int{}

	services := strings.Split("httpd.service,httpd", ",")
	for _, v := range services {
		sysUnits[v] = 1
	}

	if len(sysUnits) != 2 {
		t.Errorf("Expected 2 units to be ignored, but got %d", len(sysUnits))
	}

	sysUnits = make(map[string]int)
	os.Setenv("SYSTEMD_UNITS_IGNORE", "")
	emptyEnvVar := os.Getenv("SYSTEMD_UNITS_IGNORE")

	services = strings.Split(emptyEnvVar, ",")
	for _, v := range services {
		sysUnits[v] = 1
	}

	if len(sysUnits) != 1 {
		t.Errorf("Expected 1 units because split of empty string produce the same string %d", len(sysUnits))
	}

}

func TestIgnoreSystemUnitWork(t *testing.T) {
	sysUnits := map[string]int{}

	services := strings.Split("httpd.service", ",")
	for _, v := range services {
		sysUnits[v] = 1
	}

	unit := "httpd.service"

	_, ok := sysUnits[unit]
	if ok {
		fmt.Printf("Unit %s is ignored\n", unit)
		return
	}
	t.Errorf("Unit %s is not ignored\n", unit)
	return

}
