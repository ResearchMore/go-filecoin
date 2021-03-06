package gengen

import (
	"fmt"

	"github.com/filecoin-project/go-filecoin/internal/pkg/constants"
)

// MakeCommitCfgs creates n gengen commit configs, casting strings to cids.
func MakeCommitCfgs(n int) ([]*CommitConfig, error) {
	cfgs := make([]*CommitConfig, n)
	for i := 0; i < n; i++ {
		commP, err := constants.DefaultCidBuilder.Sum([]byte(fmt.Sprintf("commP: %d", i)))
		if err != nil {
			return nil, err
		}
		commR, err := constants.DefaultCidBuilder.Sum([]byte(fmt.Sprintf("commR: %d", i)))
		if err != nil {
			return nil, err
		}
		commD, err := constants.DefaultCidBuilder.Sum([]byte(fmt.Sprintf("commD: %d", i)))
		if err != nil {
			return nil, err
		}

		dealCfg := &DealConfig{
			CommP:     commP,
			PieceSize: uint64(1),
			EndEpoch:  int64(1024),
		}

		cfgs[i] = &CommitConfig{
			CommR:     commR,
			CommD:     commD,
			SectorNum: uint64(i),
			DealCfg:   dealCfg,
		}
	}
	return cfgs, nil
}
