package singleinst

import (
	"github.com/dan-and-dna/singleinstmodule"

	// modules
	_ "singleinst/modules/config"
	_ "singleinst/modules/mvc"
)

func Poll() {
	singleinstmodule.Run(false)
}
