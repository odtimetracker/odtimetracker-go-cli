// Copyright 2015 Ondřej Doněk. All rights reserved.
// See LICENSE file for more details about licensing.

// odTimeTracker is simple time-tracking tool.
package main

import (
	"database/sql"
	"fmt"
	"github.com/odTimeTracker/odtimetracker-go-lib"
	"github.com/odTimeTracker/odtimetracker-go-lib/database"
	"os"
	"os/user"
	"path"
)

var (
	// Application's name
	AppName = "odTimeTracker"
	// Application's short name (system name)
	AppShortName = "odtimetracker"
	AppVersion = odtimetracker.Version{ Major: 0, Minor: 1, Maintenance: 0, }
	// Application's info line
	AppInfo = AppName + " " + AppVersion.String()
	// Application's description
	AppDesc = "Simple tool for time-tracking."
)

// Simple struct representing command
type Command struct {
	Name      string // Name of a command
	Desc      string // Description of a command
	UsageDesc string // Usage description (arguments)

	// Runs the command self.
	Run func(cmd *Command, db *sql.DB, args []string)

	// Prints help on the command.
	Help func(cmd *Command)
}

// Prints usage information for the command.
func (cmd *Command) Usage(prefix string, suffix string) {
	fmt.Printf("%s%s %s\t%s\n%s", prefix,
		cmd.Name, cmd.UsageDesc, cmd.Desc, suffix)
}

// All commands supported by this tool
var commands = []*Command{
	cmdInfo,
	cmdList,
	cmdReport,
	cmdStart,
	cmdStop,
}

// Main (entry) function.
func main() {
	fmt.Println(AppInfo)

	if len(os.Args) <= 1 {
		usage()
		return
	}

	if os.Args[1] == "help" {
		help(os.Args[2:])
		return
	}

	path, _ := databasePath()
	db, err := database.InitStorage(path)
	if err != nil {
		fmt.Printf("Error occured during initializing database connection:\n\n%s\n\n", err.Error())
		return
	}
	defer db.Close()

	for _, cmd := range commands {
		if os.Args[1] == cmd.Name {
			cmd.Run(cmd, db, os.Args[2:])
			return
		}
	}

	fmt.Printf("Unknown command '%s'.\n\nRun '%s help' for usage.\n", os.Args[1], AppShortName)
}

// Returns path to the SQLite database file.
func databasePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(usr.HomeDir, ".odtimetracker.sqlite"), nil
}

// Prints usage informations.
func usage() {
	fmt.Printf("\n%s\n\n", AppDesc)
	fmt.Printf("Usage:\n\n")
	fmt.Printf("\t%s command [arguments]\n\n", AppShortName)
	fmt.Printf("Available commands:\n\n")
	for _, cmd := range commands {
		fmt.Printf("\t%s\t%s\n", cmd.Name, cmd.Desc)
	}
	fmt.Printf("\nUse \"%s help [command]\" for more information about a command.\n\n", AppShortName)
}

// Implements the 'help' command.
func help(args []string) {
	if len(args) == 0 {
		usage()
		os.Exit(0)
	}

	if len(args) > 1 {
		fmt.Printf("\nUsage:\n\n\t%s help [command]\n\nToo many arguments given!\n", AppShortName)
		os.Exit(1)
	}

	for _, cmd := range commands {
		if args[0] == cmd.Name {
			cmd.Help(cmd)
			os.Exit(0)
		}
	}

	fmt.Printf("Unknown command name %s given\nRun '%s help' for usage.\n", args[0], AppShortName)
	os.Exit(1)
}
