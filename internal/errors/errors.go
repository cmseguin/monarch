package errors

import "github.com/cmseguin/khata"

var FatalError = khata.NewTemplate().
	SetCode(1).
	SetMessage("Fatal error").
	SetExitCode(1)

var WarningError = khata.NewTemplate().
	SetCode(2).
	SetMessage("Warning").
	SetExitCode(0)
