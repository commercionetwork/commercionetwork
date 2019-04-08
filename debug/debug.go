package debug

import (
	standLog "log"
	"os"
)

func Log(value string) {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		standLog.Fatal(err)
	}
	defer file.Close()
	standLog.SetOutput(file)
	standLog.Print(value)
}
