package main

import (
	"os"

	"github.com/integrii/flaggy"
)

type Word struct {
	ID      string `json:"id"`
	English string `json:"english"`
	Polish  string `json:"polish"`
}

var cfg struct {
	addressFlag string
}

var (
	appName    string = "lidi"
	appVersion string = "0.0.1"
	appDesc    string = "A little dictionary client app"
	appWord    string
	cmdEnglish *flaggy.Subcommand
	cmdPolish  *flaggy.Subcommand
	//cmdAdd     *flaggy.Subcommand
)

func init() {
	cfg.addressFlag = os.Getenv("DICTIONARY_SERVER")
	if cfg.addressFlag == "" {
		cfg.addressFlag = "http://localhost:8080"
	}

	flaggy.SetName(appName)
	flaggy.SetDescription(appDesc)
	flaggy.SetVersion(appVersion)
	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	flaggy.String(&cfg.addressFlag, "-s", "server", "dictionary server address")

	cmdEnglish = flaggy.NewSubcommand("en")
	cmdEnglish.Description = "translate from English to Polish"
	cmdEnglish.AddPositionalValue(&appWord, "word", 1, true, "word to translate")
	flaggy.AttachSubcommand(cmdEnglish, 1)

	cmdPolish = flaggy.NewSubcommand("pl")
	cmdPolish.Description = "translate from Polish to English"
	cmdPolish.AddPositionalValue(&appWord, "word", 1, true, "word to translate")
	flaggy.AttachSubcommand(cmdPolish, 1)

	flaggy.Parse()
}

func main() {
	if cmdEnglish.Used {
		translateEnglish(appWord)
	} else if cmdPolish.Used {
		translatePolish(appWord)
	} else {
		/*
			if no correct subcommand is given, a general help is displayed
			and the program will terminate
		*/
		flaggy.ShowHelp("")
		os.Exit(1)
	}
}
