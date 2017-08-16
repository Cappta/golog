package golog

import (
	"os"

	"github.com/Cappta/Cappta.Common.Go/Debug"
)

var (
	debugGetCaller = Debug.GetCaller
	osHostname     = os.Hostname
)
