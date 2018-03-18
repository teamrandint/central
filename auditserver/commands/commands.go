package commands

import (
	"encoding/xml"
	"fmt"
)

// Command contains types of commands user can run
type Command interface {
	String() string
	Byte() []byte
}

type UserCommand struct {
	XMLName        xml.Name `xml:"userCommand"`
	Timestamp      int64    `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum string   `xml:"transactionNum"`
	Command        string   `xml:"command"`
	Username       string   `xml:"username,omitempty"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Filename       string   `xml:"filename,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
}

type QuoteServer struct {
	XMLName         xml.Name `xml:"quoteServer"`
	Timestamp       int64    `xml:"timestamp"`
	Server          string   `xml:"server"`
	TransactionNum  string   `xml:"transactionNum"`
	Price           string   `xml:"price"`
	StockSymbol     string   `xml:"stockSymbol"`
	Username        string   `xml:"username"`
	QuoteServerTime string   `xml:"quoteServerTime"`
	Cryptokey       string   `xml:"cryptokey"`
}

type AccountTransaction struct {
	XMLName        xml.Name `xml:"accountTransaction"`
	Timestamp      int64    `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum string   `xml:"transactionNum"`
	Action         string   `xml:"action"`
	Username       string   `xml:"username,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
}

type SystemEvent struct {
	XMLName        xml.Name `xml:"systemEvent"`
	Timestamp      int64    `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum string   `xml:"transactionNum"`
	Command        string   `xml:"command"`
	Username       string   `xml:"username,omitempty"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Filename       string   `xml:"filename,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
}

type ErrorEvent struct {
	XMLName        xml.Name `xml:"errorEvent"`
	Timestamp      int64    `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum string   `xml:"transactionNum"`
	Command        string   `xml:"command"`
	Username       string   `xml:"username,omitempty"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Filename       string   `xml:"filename,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
	ErrorMessage   string   `xml:"errorMessage,omitempty"`
}

// String returns a string representation of userCommand
func (u *UserCommand) String() string {
	return string(u.Byte())
}

// Byte returns byte array of usercommand
func (u *UserCommand) Byte() []byte {
	output, err := xml.MarshalIndent(u, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return output
}

// String returns a string representation of QuoteServer
func (u *QuoteServer) String() string {
	return string(u.Byte())
}

// Byte returns byte array of QuoteServer
func (u *QuoteServer) Byte() []byte {
	output, err := xml.MarshalIndent(u, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return output
}

// String returns a string representation of AccountTransaction
func (u *AccountTransaction) String() string {
	return string(u.Byte())
}

// Byte returns byte array of AccountTransaction
func (u *AccountTransaction) Byte() []byte {
	output, err := xml.MarshalIndent(u, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return output
}

// String returns a string representation of SystemEvent
func (u *SystemEvent) String() string {
	return string(u.Byte())
}

// Byte returns byte array of SystemEvent
func (u *SystemEvent) Byte() []byte {
	output, err := xml.MarshalIndent(u, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return output
}

// String returns a string representation of ErrorEvent
func (u *ErrorEvent) String() string {
	return string(u.Byte())
}

// Byte returns byte array of ErrorEvent
func (u *ErrorEvent) Byte() []byte {
	output, err := xml.MarshalIndent(u, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return output
}
