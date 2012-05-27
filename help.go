package main

import (
	"fmt"
	"os"
	"strings"
)

var cmdFetchUpdate = &Command{
	Run:   runFetchUpdate,
	Usage: "fetch-update",
	Long:  `Downloads the next version of hk for later installation.`,
}

func runFetchUpdate(cmd *Command, args []string) {
	if len(os.Args) != 2 {
		unrecArg(os.Args[2], "fetch-update")
	}

	updater.fetchAndApply()
}

var cmdVersion = &Command{
	Run:   runVersion,
	Usage: "version",
	Short: "show hk version",
	Long:  `Version shows the hk client version string.`,
}

func runVersion(cmd *Command, args []string) {
	fmt.Println(Version)
}

var cmdHelp = &Command{
	Usage: "help [command]",
	Short: "show help",
	Long:  `Help shows usage for a command.`,
}

func init() {
	cmdHelp.Run = runHelp // break init loop
}

func runHelp(cmd *Command, args []string) {
	if len(args) == 0 {
		printUsage()
		return // not os.Exit(2); success
	}
	if len(args) != 1 {
		errorf("too many arguments")
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			fmt.Printf("Usage: hk %s\n\n", cmd.Usage)
			fmt.Println(strings.TrimSpace(cmd.Long))
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic: %q. Run 'hk help'.\n", args[0])
	os.Exit(2)
}

func printUsage() {
	fmt.Printf("Usage: hk <command> [options] [arguments]\n\n")

	fmt.Printf("Supported options are:\n\n")
	fmt.Printf("  -a APP     name of the app to operate on\n")
	fmt.Println()

	fmt.Printf("Supported commands are:\n\n")
	for _, cmd := range commands {
		if cmd.Short != "" {
			fmt.Printf("  %-8s   %s\n", cmd.Name(), cmd.Short)
		}
	}
	fmt.Println()

	fmt.Printf("See 'hk help [command]' for more information about a command.\n")
}

func usage() {
	printUsage()
	os.Exit(2)
}
