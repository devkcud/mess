package therapy

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	affirmations = []string{
		"you are more than your file structure.",
		"you are not defined by your dotfiles.",
		"a broken symlink is still trying.",
		"mkdir your dreams.",
		"you matter, even if your folders don't.",
		"you've made worse decisions than this. and survived.",
		"git status: enough.",
		"you are not a failed shell script.",
		"breathe. hydrate. maybe go outside.",
		"your commit history does not define you.",
		"even if the build fails, you are a success.",
		"remember: rm -rf can only delete files, not your worth.",
		"no linter can judge the code in your heart.",
		"today you are the main branch.",
		"even vim users get lost sometimes.",
		"you are not a merge conflict.",
		"your value is not measured in LOC (lines of code).",
		"take a breakpoint. you deserve it.",
	}
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func Header() {
	fmt.Println("🤍 mess cares about you.")
}

func RandomQuote() {
	fmt.Println("\"" + affirmations[rng.Intn(len(affirmations))] + "\"")
}

func Footer() {
	fmt.Println("go take good care of yourself now!")
}

func HelpMe() {
	Header()
	fmt.Println()
	RandomQuote()
	fmt.Println()
	Footer()
}
