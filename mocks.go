package golog

import (
	"os"

	"github.com/Cappta/debugo"
)

var (
	debugGetCaller = debugo.GetCaller
	osHostname     = os.Hostname
)
