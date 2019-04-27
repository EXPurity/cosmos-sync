package helper

import (
	"errors"
	"cosmos-sync/store/document"
	"cosmos-sync/types"
	"cosmos-sync/util/constant"
)

func GetProposal(proposalID uint64) (proposal document.Proposal, err error) {
	cdc := types.GetCodec()

	res, err := Query(types.KeyProposal(proposalID), "gov", constant.StoreDefaultEndPath)
	if len(res) == 0 || err != nil {
		return proposal, errors.New("no data")
	}
	var propo types.Proposal
	cdc.UnmarshalBinaryLengthPrefixed(res, &propo) //TODO
	proposal.ProposalId = proposalID
	proposal.Title = propo.ProposalContent.GetTitle()
	proposal.Type = propo.ProposalContent.ProposalType().String()
	proposal.Description = propo.ProposalContent.GetDescription()
	proposal.Status = propo.Status.String()

	proposal.SubmitTime = propo.SubmitTime
	proposal.VotingStartTime = propo.VotingStartTime
	proposal.VotingEndTime = propo.VotingEndTime
	proposal.DepositEndTime = propo.DepositEndTime
	proposal.TotalDeposit = types.ParseCoins(propo.TotalDeposit.String())
	proposal.Votes = []document.PVote{}
	return
}

func GetVotes(proposalID uint64) (pVotes []document.PVote, err error) {
	cdc := types.GetCodec()

	res, err := QuerySubspace(types.KeyVotesSubspace(proposalID), "gov")
	if len(res) == 0 || err != nil {
		return pVotes, err
	}
	for i := 0; i < len(res); i++ {
		var vote types.SdkVote
		cdc.UnmarshalBinaryLengthPrefixed(res[i].Value, &vote)
		v := document.PVote{
			Voter:  vote.Voter.String(),
			Option: vote.Option.String(),
		}
		pVotes = append(pVotes, v)
	}
	return
}
