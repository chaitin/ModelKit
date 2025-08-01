package promptparser

import (
	"fmt"
	"testing"
)

func TestTaskParse(t *testing.T) {
	prompt := "The user denied this operation and provided the following feedback:\n\u003cfeedback\u003e\n直接写到 'src/notion/mod.rs' (see below for file content) 这个文件\n\u003c/feedback\u003e\n\u003cfiles\u003e\n\u003cfile\u003e\u003cpath\u003eCargo.toml\u003c/path\u003e\u003cstatus\u003eDenied by user\u003c/status\u003e\u003c/file\u003e\n\u003c/files\u003e\n\n\u003cfile_content path=\"src/notion/mod.rs\"\u003e\n\n\u003c/file_content\u003e"

	tp := TaskParse{}
	info, err := tp.Parse(prompt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", info)
}
