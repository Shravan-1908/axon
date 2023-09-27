package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/thatisuday/commando"
)

const (
	NAME    = "axon"
	VERSION = "1.2.0"
)

func main() {
	fmt.Println(NAME, VERSION)
	deletePreviousInstallation()

	// setting up executable details
	commando.
		SetExecutableName(NAME).
		SetVersion(VERSION).
		SetDescription("axon is a command line utility to organise and pretty your file system quickly and reliably.")

	// root command
	commando.
		Register(nil).
		SetShortDescription("Run axon with default options.").
		SetDescription("Run axon with default options, organising all files and folders.").
		AddArgument("dirs...", "The directory to be organized.", "./").
		AddFlag("prettify,p", "Prettify all files with a desired casing.", commando.String, "none").
		AddFlag("organise,o", "Organise the directory.", commando.Bool, true).
		AddFlag("rename,r", "Rename the files numerically with a certain alias.", commando.String, "none").
		AddFlag("regex,x", "Filter files using regular expressions.", commando.String, ":_none_:").
		AddFlag("move,m", "Move selected files to a directiry.", commando.String, ":none:").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {

			// getting all arg and flag values
			dirs := strings.Split(args["dirs"].Value, ",")

			prettify, e := flags["prettify"].GetString()
			if e != nil {
				prettify = "none"
			}
			organise, e := flags["organise"].GetBool()
			if e != nil {
				organise = false
			}

			rename, e := flags["rename"].GetString()
			if e != nil {
				rename = "none"
			}

			moveToDir, e := flags["move"].GetString()
			if e != nil {
				moveToDir = ":_none_:"
			}

			regexPattern, e := flags["regex"].GetString()
			if e != nil {
				regexPattern = ""
			}
			regex, err := regexp.Compile(regexPattern)
			if err != nil {
				fmt.Printf("unable to parse the given regex. please check it again.\n%v", err)
				return
			}
			// making a buffered channel
			ch := make(chan string, len(dirs))
			// todo better actions report
			// todo move before organising
			// organising the files
			for _, dir := range dirs {
				go func(dir string) {
					if validPath(dir) {
						fo := FileOrganizer{
							path:  dir,
							regex: regex,
						}

						if prettify != "none" {
							fo.prettify(prettify)
						}

						if moveToDir != ":_none_:" {
							fo.move(moveToDir)
						}

						// todo better renaming
						if rename != "none" {
							fo.renameDir(rename)
						}

						if organise {
							fo.organize()
						}

						ch <- fo.showActions()

					} else {
						ch <- fmt.Sprintf("Skipping `%s` since it's not a valid directory.", dir)
					}

				}(dir)
			}

			// waiting for all the goroutines to finish
			for i := 0; i < len(dirs); i++ {
				fmt.Println(<-ch)
			}

		})

	// up command
	commando.
		Register("up").
		SetShortDescription("Update axon.").
		SetDescription("Update axon to the latest version.").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			update()
		})

	commando.Parse(nil)
}
