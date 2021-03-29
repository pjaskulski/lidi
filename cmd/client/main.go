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
	speakFlag   bool
	showIdFlag  bool
}

var (
	appName          string = "lidi"
	appVersion       string = "0.0.1"
	appDesc          string = "A little dictionary client app"
	appWord          string
	appTranslationID string
	cmdEnglish       *flaggy.Subcommand
	cmdPolish        *flaggy.Subcommand
	cmdSpeak         *flaggy.Subcommand
	cmdAdd           *flaggy.Subcommand
	cmdUpdate        *flaggy.Subcommand
	cmdDelete        *flaggy.Subcommand
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

	flaggy.String(&cfg.addressFlag, "s", "server", "Dictionary server address")
	flaggy.Bool(&cfg.speakFlag, "p", "speak", "Speak English after translate")
	flaggy.Bool(&cfg.showIdFlag, "i", "id", "show record id")

	cmdEnglish = flaggy.NewSubcommand("en")
	cmdEnglish.Description = "Translate from English to Polish"
	cmdEnglish.AddPositionalValue(&appWord, "word", 1, true, "word to translate")
	flaggy.AttachSubcommand(cmdEnglish, 1)

	cmdPolish = flaggy.NewSubcommand("pl")
	cmdPolish.Description = "Translate from Polish to English"
	cmdPolish.AddPositionalValue(&appWord, "word", 1, true, "word to translate")
	flaggy.AttachSubcommand(cmdPolish, 1)

	cmdSpeak = flaggy.NewSubcommand("speak")
	cmdSpeak.Description = "Say in English (Google API is used)"
	cmdSpeak.AddPositionalValue(&appWord, "word", 1, true, "word to speak")
	flaggy.AttachSubcommand(cmdSpeak, 1)

	cmdAdd = flaggy.NewSubcommand("add")
	cmdAdd.Description = "Add new item to dictionary (English=Polish)"
	cmdAdd.AddPositionalValue(&appWord, "word", 1, true, "translation in form: English=Polish")
	flaggy.AttachSubcommand(cmdAdd, 1)

	cmdUpdate = flaggy.NewSubcommand("update")
	cmdUpdate.Description = "Update item in dictionary (ID English=Polish)"
	cmdUpdate.AddPositionalValue(&appTranslationID, "id", 1, true, "translation ID from database")
	cmdUpdate.AddPositionalValue(&appWord, "word", 2, true, "new translation in form: English=Polish")
	flaggy.AttachSubcommand(cmdUpdate, 1)

	cmdDelete = flaggy.NewSubcommand("delete")
	cmdDelete.Description = "Delete item in dictionary (English=Polish)"
	cmdDelete.AddPositionalValue(&appTranslationID, "id", 1, true, "translation ID from database")
	flaggy.AttachSubcommand(cmdDelete, 1)

	flaggy.Parse()
}

func main() {
	if (cmdEnglish.Used || cmdPolish.Used || cmdSpeak.Used || cmdAdd.Used || cmdUpdate.Used) && appWord == "" {
		flaggy.ShowHelp("")
		os.Exit(1)
	}

	if cmdEnglish.Used {
		translateEnglish(appWord, cfg.speakFlag, cfg.showIdFlag)
	} else if cmdPolish.Used {
		translatePolish(appWord, cfg.speakFlag, cfg.showIdFlag)
	} else if cmdSpeak.Used {
		speak(appWord)
	} else if cmdAdd.Used {
		addTranslation(appWord)
	} else if cmdUpdate.Used {
		updateTranslation(appTranslationID, appWord)
	} else if cmdDelete.Used {
		deleteTranslation(appTranslationID)
	} else {
		/*
			if no correct subcommand is given, a general help is displayed
			and the program will terminate
		*/
		flaggy.ShowHelp("")
		os.Exit(1)
	}
}
