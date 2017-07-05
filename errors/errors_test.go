package errorkits

import "testing"

func TestInit(t *testing.T) {
	ErrorMessageFile = "error_msg.yaml"
	t.Logf("%v", Init()["ErrorLoginExpired"])
}
