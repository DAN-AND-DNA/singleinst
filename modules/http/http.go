package http

import (
	"github.com/dan-and-dna/singleinstmodule"
	"singleinst/modules/http/internal"
)

func ModuleLock() singleinstmodule.ModuleCore {
	return internal.GetSingleInst().ModuleLock()
}

func ModuleUnlock() {
	internal.GetSingleInst().ModuleUnlock()
}
