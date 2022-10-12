package logger

import (
	"os"
	"testing"
)

func TestInitLog(t *testing.T) {
	flog, err := os.OpenFile(`client.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	_ = InitLog(flog)
}
