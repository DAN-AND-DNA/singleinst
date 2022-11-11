package grpc

import (
	"github.com/dan-and-dna/singleinstmodule"
	"singleinst/modules/grpc/internal"
)

type Grpc = internal.Grpc

func ModuleLock() singleinstmodule.ModuleCore {
	return internal.GetSingleInst().ModuleLock()
}

func ModuleUnlock() {
	internal.GetSingleInst().ModuleUnlock()
}
