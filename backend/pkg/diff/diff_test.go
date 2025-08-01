package diff

import (
	"fmt"
	"testing"
)

func TestParseConflictsAndCountLines(t *testing.T) {
	conflictText := `<<<<<<< HEAD
old line 1
=======
new line 1
new line 2
>>>>>>> branch1
normal line
<<<<<<< HEAD
old line 2
old line 3
=======
new line 3
>>>>>>> branch2`

	addedLines := ParseConflictsAndCountLines(conflictText)
	fmt.Println(addedLines)
}
