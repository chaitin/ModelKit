package machine

import (
	"testing"
)

func TestGenerateShortMachineID(t *testing.T) {
	machineID, err := GenerateMachineID()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(machineID)
}
