package utils

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"
	"strconv"
	"strings"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// types

func (uo *UciOption) ToUciString() string {
	if uo.Kind == "button" {
		return fmt.Sprintf("option name %s type button", uo.Name)
	} else if uo.Kind == "combo" {
		return fmt.Sprintf("option name %s type combo default %s var %s", uo.Name, uo.Default, strings.Join(uo.DefaultStringArray, " var "))
	} else if uo.Kind == "spin" {
		return fmt.Sprintf("option name %s type spin default %d min %d max %d", uo.Name, uo.DefaultInt, uo.MinInt, uo.MaxInt)
	} else if uo.ValueKind == "int" {
		return fmt.Sprintf("option name %s type %s default %d", uo.Name, uo.Kind, uo.DefaultInt)
	} else if uo.ValueKind == "bool" {
		return fmt.Sprintf("option name %s type %s default %v", uo.Name, uo.Kind, uo.DefaultBool)
	} else if uo.ValueKind == "stringarray" {
		return fmt.Sprintf("option name %s type %s default %s", uo.Name, uo.Kind, strings.Join(uo.DefaultStringArray, " "))
	} else {
		return fmt.Sprintf("option name %s type %s default %s", uo.Name, uo.Kind, uo.Default)
	}
}

func (uo *UciOption) PrintUci() {
	fmt.Println(uo.ToUciString())
}

func (uo *UciOption) SetFromString(value string) {
	if uo.ValueKind == "string" {
		uo.Value = value
	} else if uo.ValueKind == "int" {
		i, err := strconv.ParseInt(value, 10, 32)

		if err == nil {
			uo.ValueInt = int(i)
		} else {
			fmt.Println("option value should be an integer")
		}
	} else if uo.ValueKind == "bool" {
		uo.ValueBool = false

		if value == "true" {
			uo.ValueBool = true
		}
	}
}

/////////////////////////////////////////////////////////////////////
