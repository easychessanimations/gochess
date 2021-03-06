package uciengine

import "fmt"

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// about

func AboutStr() string {
	return fmt.Sprintf(
		"\n--------------\n%s %s by %s\n--------------\n",
		ENGINE_NAME,
		ENGINE_DESCRIPTION,
		ENGINE_AUTHOR,
	)
}

func About() {
	fmt.Println(AboutStr())
}

/////////////////////////////////////////////////////////////////////
