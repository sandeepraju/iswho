// Copyright (c) 2016, Sandeep Raju Prabhakar <me@sandeepraju.in>
// All rights reserved.
package iswho

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/sandeepraju/iswho/src/util"
)

const (
	VERSION = "0.0.1-dev"
)

// IsWhoConfig defines the configuration
// to be used for searching.
type IsWhoConfig struct {
	Host    string
	Port    int
	Verbose bool
}

// isWho exposes methods to interact
// with the  whois directory service.
type isWho struct {
	Host        string
	Port        int
	Verbose     bool
	QueryServer string
	Query       string
	conn        net.Conn
}

// initializeConnection initializes
// a tcp connection to the whois server.
func (i *isWho) initializeConnection() {
	// determine the query server
	var conn net.Conn
	conn, err := net.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", i.Host, i.Port),
	)
	if err != nil {
		// pass on the error.
		panic(err)
	}
	i.conn = conn
}

// determineQueryServer determines the whois service
// server to query based on the query string.
func (i *isWho) determineQueryServer(query string) {
	if util.IsValidDomainName(query) {
		tld, err := util.GetTLD(query)
		if err != nil {
			// pass on the error.
			panic(err)
		}
		// set the query server.
		i.QueryServer = util.GetQueryServer(tld)
	} else if util.IsValidIPv4Address(query) {
		// set the query server.
		// TODO: how do we get the query server here?
	} else {
		// something else. maybe an IPv6 address?
		// or an NIC query? RIPE query?
		// all these are not implemented for now.
		panic(errors.New(
			"Error: iswho only supports " +
				"Fully Qualified Domain Names (FQDN) and " +
				"IPv4 addresses for now!",
		))
	}
}

func (i *isWho) Search(query string) (string, error) {
	// set the query string for all internal functions to access
	i.Query = query
	// determine the server to query based on the query string.
	if i.Host == "" {
		// host is not provided explicitly.
		i.determineQueryServer(query)
	}

	if i.Port == 0 {
		// port is not provided explicitly.
		// fallback to the default 43.
		i.Port = 43
	}

	// initialize the tcp socket connection to the server.
	i.initializeConnection()

	// initiate a query
	fmt.Fprintf(i.conn, "%s\r\n", i.Query)

	// read the results
	reader := bufio.NewReader(i.conn)

	// process the results
	// TODO: what should the ideal value of size be?
	results := make([]byte, 1)
	isEOF := false // to detect the end of file
	// read the data line by line
	for !isEOF {
		// 10 is the ascii code for \n
		// TODO: see if we can make it
		// more readable by replacing with '\n'
		line, err := reader.ReadBytes(10)
		if err != nil {
			if err == io.EOF {
				// we have reached the end of file.
				// toggle our flag to exit loop.
				isEOF = true
			} else {
				// pass on the error.
				return "", err
			}
		}
		// append the line to the results
		results = append(results, line...)
	}
	return string(results), nil
}

// NewIsWho returns a new iswho instance configured
// as per the configuration options passed.
func NewIsWho(config *IsWhoConfig) *isWho {
	// create a new iswho instance and pass it back.
	return &isWho{
		Host:    config.Host,
		Port:    config.Port,
		Verbose: config.Verbose,
	}
}
