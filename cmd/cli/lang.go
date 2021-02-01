package main

const (
	ERROR = iota
	HEADER
	FLAG_CONFIG
	FLAG_CREDENTIALS
	FLAG_VERBOSITY
	MARKER
	READING_CONFIG
	FOUND_CONFIGURATION
	USER
	TITLE
	MISSING
	FOUND
	UPDATING_BG_COLOR
	UPDATING_TEXT_COLOR
	MISMATCH
	DELETING
	UPDATING
)

type LangPack map[int]string
type LangPacks map[string]LangPack

var language = LangPacks{
	"colors": LangPack{
		ERROR:               "\033[0;31m[ERROR] %s\033[0m\n",
		HEADER:              "\033[2m======= \033[0m\033[1mStickers\033[0m\033[2m =======\033[0m\n",
		FLAG_CONFIG:         "\033[2mflagConfig: \033[0m%s\n",
		FLAG_CREDENTIALS:    "\033[2mflagCredentials: \033[0m%s\n",
		FLAG_VERBOSITY:      "\033[2mflagVerbosity: \033[0m%v\n",
		MARKER:              "\033[2m------------------------\033[0m\n\n",
		READING_CONFIG:      "\033[2mReading configuration file: \033[0m\033[1m%s\033[0m\n",
		FOUND_CONFIGURATION: "\033[1mFound configuration.\033[0m\n\033[2m# => %#v\033[0m\n",
		USER:                "\033[2m[ User: \033[0m\033[1m %s \033[0m\033[2m]\033[0m\n",
		TITLE:               "\t\033[0;34m%s\033[0m\n",
		MISSING:             "\t\033[0;31m- Missing:\033[0m \033[1m%s\033[0m\n",
		FOUND:               "\t\033[2m- Found:\033[0m \033[1m%s\033[0m\n",
		UPDATING_BG_COLOR:   "\t\033[2m Found:\033[0m \033[1m%s\033[0m\n",
		UPDATING_TEXT_COLOR: "\t\033[2m- Found:\033[0m \033[1m%s\033[0m\n",
		MISMATCH:            "\t\033[2m Update required\n",
		DELETING:            "\t\t\033[2m Deleting\033[0m \033[1m%s\033[0m\n",
		UPDATING:            "\t\t\033[2m Updating\033[0m \033[1m%s\033[0m\n",
	},
	"en-US": LangPack{
		ERROR:               "[ERROR] %s\n",
		HEADER:              "======= ]Stickers[ =======\n",
		FLAG_CONFIG:         "flagConfig: %s\n",
		FLAG_CREDENTIALS:    "flagCredentials: %s\n",
		FLAG_VERBOSITY:      "flagVerbosity: %v\n",
		MARKER:              "--------------------------\n\b",
		READING_CONFIG:      "Reading configuration file: %s\n",
		FOUND_CONFIGURATION: "Found configuration.\n# => %#v\n",
		USER:                "[ User: %s ]\n",
		TITLE:               "\t%s\n",
		MISSING:             "\t- Missing: %s\n",
		FOUND:               "\t- Found: %s\n",
		UPDATING_BG_COLOR:   "\t Found: %s\n",
		UPDATING_TEXT_COLOR: "\t- Found:%s\n",
		MISMATCH:            "\t\tUpdate Required\n",
		DELETING:            "\t\tDeleting %s\n",
		UPDATING:            "\t\tUpdating %s\n",
	},
}
