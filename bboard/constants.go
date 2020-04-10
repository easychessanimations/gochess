package minboard

/////////////////////////////////////////////////////////////////////
// imports

import (
	"github.com/easychessanimations/gochess/utils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// constants

const ADD_SAN = true

const DEFAULT_SEARCH_DEPTH = 10

const DEFAULT_UCI_VARIANT_STRING = "standard"

var UCI_OPTIONS = []utils.UciOption{
	{
		Kind:      "combo",
		Name:      "UCI_Variant",
		ValueKind: "string",
		Default:   DEFAULT_UCI_VARIANT_STRING,
		Value:     DEFAULT_UCI_VARIANT_STRING,
		DefaultStringArray: []string{
			"standard",
			"atomic",
			"seirawan",
			"eightpiece",
		},
	},
}

var UCI_COMMAND_ALIASES = map[string]string{
	"v":  "setoption name UCI_Variant value standard\ni",
	"va": "setoption name UCI_Variant value atomic\ni",
	"vs": "setoption name UCI_Variant value seirawan\ni",
	"ve": "setoption name UCI_Variant value eightpiece\ni",
}

/////////////////////////////////////////////////////////////////////
