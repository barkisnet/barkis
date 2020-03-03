package types

import "fmt"

const (
	RewardUpgrade                 = "RewardUpgrade"
	TokenIssueUpgrade             = "TokenIssueUpgrade"
	UpdateVotingPeriodHeight      = "UpdateVotingPeriodHeight"
	UpdateTokenSymbolRulesHeight  = "UpdateTokenSymbolRulesHeight"
	TokenDesLenLimitUpgradeHeight = "TokenDesLenLimitUpgradeHeight"
)

var GlobalUpgradeMgr = NewUpgradeManager()

type UpgradeConfig struct {
	UpgradeHeight  map[string]int64
	NewStoreHeight map[string]int64
	NewMsgHeight   map[string]int64

	BeginBlockersFirst map[int64][]func(ctx Context)
	BeginBlockersLast  map[int64][]func(ctx Context)

	EndBlockersFirst map[int64][]func(ctx Context)
	EndBlockersLast  map[int64][]func(ctx Context)
}

type UpgradeManager struct {
	Config      UpgradeConfig
	BlockHeight int64
}

func NewUpgradeManager() *UpgradeManager {
	return &UpgradeManager{
		Config: UpgradeConfig{
			UpgradeHeight:      make(map[string]int64),
			NewStoreHeight:     make(map[string]int64),
			NewMsgHeight:       make(map[string]int64),
			BeginBlockersFirst: make(map[int64][]func(ctx Context)),
			BeginBlockersLast:  make(map[int64][]func(ctx Context)),
			EndBlockersFirst:   make(map[int64][]func(ctx Context)),
			EndBlockersLast:    make(map[int64][]func(ctx Context)),
		},
		BlockHeight: 0,
	}
}

func (mgr *UpgradeManager) SetBlockHeight(height int64) {
	mgr.BlockHeight = height
}

func (mgr *UpgradeManager) GetBlockHeight() int64 {
	return mgr.BlockHeight
}

// BeginBlockers for upgrade
func (mgr *UpgradeManager) BeginBlockersFirst(ctx Context) {
	if beginBlockers, ok := mgr.Config.BeginBlockersFirst[mgr.GetBlockHeight()]; ok {
		for _, beginBlocker := range beginBlockers {
			if beginBlocker == nil {
				continue
			}
			beginBlocker(ctx)
		}
	}
}

func (mgr *UpgradeManager) BeginBlockersLast(ctx Context) {
	if beginBlockers, ok := mgr.Config.BeginBlockersLast[mgr.GetBlockHeight()]; ok {
		for _, beginBlocker := range beginBlockers {
			if beginBlocker == nil {
				continue
			}
			beginBlocker(ctx)
		}
	}
}

func (mgr *UpgradeManager) RegisterBeginBlockerFirst(name string, beginBlocker func(Context)) {
	if beginBlocker == nil {
		return
	}
	height := mgr.GetUpgradeHeight(name)
	if height == 0 {
		panic(fmt.Sprintf("no upgrade for %s at %d", name, height))
	}

	beginBlockers := mgr.Config.BeginBlockersFirst[height]
	mgr.Config.BeginBlockersFirst[height] = append(beginBlockers, beginBlocker)
}

func (mgr *UpgradeManager) RegisterBeginBlockerLast(name string, beginBlocker func(Context)) {
	if beginBlocker == nil {
		return
	}
	height := mgr.GetUpgradeHeight(name)
	if height == 0 {
		panic(fmt.Sprintf("no upgrade for %s at %d", name, height))
	}

	beginBlockers := mgr.Config.BeginBlockersLast[height]
	mgr.Config.BeginBlockersLast[height] = append(beginBlockers, beginBlocker)
}

// EndBlockers for upgrade
func (mgr *UpgradeManager) EndBlockersFirst(ctx Context) {
	if endBlockers, ok := mgr.Config.EndBlockersFirst[mgr.GetBlockHeight()]; ok {
		for _, endBlocker := range endBlockers {
			if endBlocker == nil {
				continue
			}
			endBlocker(ctx)
		}
	}
}

func (mgr *UpgradeManager) EndBlockersLast(ctx Context) {
	if endBlockers, ok := mgr.Config.EndBlockersLast[mgr.GetBlockHeight()]; ok {
		for _, endBlocker := range endBlockers {
			if endBlocker == nil {
				continue
			}
			endBlocker(ctx)
		}
	}
}

func (mgr *UpgradeManager) RegisterEndBlockerFirst(name string, endBlocker func(Context)) {
	if endBlocker == nil {
		return
	}
	height := mgr.GetUpgradeHeight(name)
	if height == 0 {
		panic(fmt.Sprintf("no upgrade for %s at %d", name, height))
	}

	if mgr.Config.EndBlockersFirst == nil {
		panic("EndBlockersFirst is null")
	}

	endBlockers := mgr.Config.EndBlockersFirst[height]
	mgr.Config.EndBlockersFirst[height] = append(endBlockers, endBlocker)
}

func (mgr *UpgradeManager) RegisterEndBlockerLast(name string, endBlocker func(Context)) {
	if endBlocker == nil {
		return
	}
	height := mgr.GetUpgradeHeight(name)
	if height == 0 {
		panic(fmt.Sprintf("no upgrade for %s at %d", name, height))
	}

	if mgr.Config.EndBlockersLast == nil {
		panic("EndBlockersLast is null")
	}

	endBlockers := mgr.Config.EndBlockersLast[height]
	mgr.Config.EndBlockersLast[height] = append(endBlockers, endBlocker)
}

// Add new upgrade
func (mgr *UpgradeManager) RegisterUpgradeHeight(name string, height int64) {
	//if mgr.Config.UpgradeHeight[name] != 0 {
	//	panic("duplicated upgrade name")
	//}
	mgr.Config.UpgradeHeight[name] = height
}

func (mgr *UpgradeManager) GetUpgradeHeight(name string) int64 {
	return mgr.Config.UpgradeHeight[name]
}

func (mgr *UpgradeManager) RegisterNewStore(upgradeName string, newStores ...string) {
	height := mgr.GetUpgradeHeight(upgradeName)
	if height == 0 {
		panic(fmt.Sprintf("no upgrade for %s", upgradeName))
	}

	for _, store := range newStores {
		mgr.Config.NewStoreHeight[store] = height
	}
}

func (mgr *UpgradeManager) GetStoreHeight(storeName string) int64 {
	return mgr.Config.UpgradeHeight[storeName]
}

func (mgr *UpgradeManager) RegisterNewMsg(upgradeName string, msgTypes ...string) {
	height := mgr.GetUpgradeHeight(upgradeName)
	if height == 0 {
		panic(fmt.Sprintf("no upgrade for %s", upgradeName))
	}

	for _, msgType := range msgTypes {
		mgr.Config.NewMsgHeight[msgType] = height
	}
}

func (mgr *UpgradeManager) GetMsgHeight(msgType string) int64 {
	return mgr.Config.NewMsgHeight[msgType]
}

func (mgr *UpgradeManager) IsUpgradeApplied(upgradeName string) bool {
	height, ok := mgr.Config.UpgradeHeight[upgradeName]
	if !ok {
		return false
	}
	return mgr.BlockHeight >= height
}

func (mgr *UpgradeManager) IsOnUpgradeHeight(upgradeName string) bool {
	height, ok := mgr.Config.UpgradeHeight[upgradeName]
	if !ok {
		return false
	}
	return mgr.BlockHeight == height
}

func (mgr *UpgradeManager) MsgCheck(msgType string) bool {
	height, ok := mgr.Config.NewMsgHeight[msgType]
	if !ok {
		return true
	}
	return mgr.BlockHeight >= height
}

func (mgr *UpgradeManager) StoreCheck(storeName string) bool {
	height, ok := mgr.Config.NewStoreHeight[storeName]
	if !ok {
		return true
	}
	return mgr.BlockHeight >= height
}

func (mgr *UpgradeManager) IsOnStoreStartHeight(storeName string) bool {
	height, ok := mgr.Config.NewStoreHeight[storeName]
	if !ok {
		return false
	}
	return mgr.BlockHeight == height
}
