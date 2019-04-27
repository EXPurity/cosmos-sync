package task

import (
	conf "cosmos-sync/conf/server"
	"cosmos-sync/logger"
	"cosmos-sync/store"
	"cosmos-sync/store/document"
	"cosmos-sync/util/constant"
	"cosmos-sync/util/helper"
)

func syncProposalStatus() {
	var status = []string{constant.StatusDepositPeriod, constant.StatusVotingPeriod}
	if proposals, err := document.QueryByStatus(status); err == nil {
		for _, proposal := range proposals {
			propo, err := helper.GetProposal(proposal.ProposalId)
			if err != nil {
				store.Delete(proposal)
				return
			}
			if propo.Status != proposal.Status {
				propo.SubmitTime = proposal.SubmitTime
				propo.Votes = proposal.Votes
				store.SaveOrUpdate(propo)
			}
		}
	}
}

func MakeSyncProposalStatusTask() Task {
	return NewLockTaskFromEnv(conf.SyncProposalStatus, "sync_proposal_status_lock", func() {
		logger.Debug("========================task's trigger [SyncProposalStatus] begin===================")
		syncProposalStatus()
		logger.Debug("========================task's trigger [SyncProposalStatus] end===================")
	})
}
