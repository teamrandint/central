package log

import (
	"encoding/xml"
	"fmt"
	"io"
	"seng468/auditserver/commands"
)

// Log contains a list of user commands
type Log struct {
	XMLName xml.Name `xml:"log"`
	Entries []commands.Command
}

// Write takes in a writer object and writes the log to a file
func (l *Log) Write(w io.Writer) {
	w.Write(l.Byte())
}

// String returns an XML representation of the log
func (l *Log) String() string {
	return string(l.Byte())
}

// Byte returns an XML representation of the log
func (l *Log) Byte() []byte {
	output, err := xml.MarshalIndent(l, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return output
}

// Insert takes a command object and inserts it into the log
func (l *Log) Insert(c commands.Command) {
	l.Entries = append(l.Entries, c)
}
