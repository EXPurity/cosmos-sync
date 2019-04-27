package task

import (
	"cosmos-sync/conf/server"
	"cosmos-sync/logger"
	"cosmos-sync/service/handler"
	"cosmos-sync/store"
	"cosmos-sync/store/document"
)

func MakeUpdateDelegatorTask() Task {
	return NewLockTaskFromEnv(server.CronUpdateDelegator, "save_update_delegator_lock", func() {
		logger.Debug("========================task's trigger [MakeUpdateDelegatorTask] begin===================")
		updateDelegator()
		logger.Debug("========================task's trigger [MakeUpdateDelegatorTask] end===================")
	})
}

func updateDelegator() {
	var delegatorStore document.Delegator
	delegators := delegatorStore.QueryUnbonding()
	if len(delegators) == 0 {
		logger.Info("no delegator is unbonding")
		return
	}

	for _, d := range delegators {
		ubd := handler.BuildUnbondingDelegation(d.Address, d.ValidatorAddr)
		d.UnbondingDelegation = ubd
		if d.BondedHeight < 0 &&
			d.UnbondingDelegation.CreationHeight < 0 {
			store.Delete(d)
			logger.Info("delete delegator", logger.String("delAddress", d.Address), logger.String("valAddress", d.ValidatorAddr))
		} else {
			store.Update(d)
			logger.Info("Update delegator", logger.String("delAddress", d.Address), logger.String("valAddress", d.ValidatorAddr))
		}
	}
}
