package nanogo

import (
	"github.com/dan-and-dna/singleinstmodule"

	// modules
	_ "singleinst/modules/grpc"
	_ "singleinst/modules/http"
	_ "singleinst/modules/mvc"
	_ "singleinst/modules/network"
)

func Poll() {
	singleinstmodule.Run(false)
}
