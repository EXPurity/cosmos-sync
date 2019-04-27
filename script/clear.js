// drop collection
 db.account.drop();
 db.block.drop();
 db.power_change.drop();
 db.proposal.drop();
 db.stake_role_candidate.drop();
 db.stake_role_delegator.drop();
 db.sync_task.drop();
 db.tx_common.drop();
 db.validator_up_time.drop();
 db.tx_gas.drop();
 db.tx_msg.drop();
 db.uptime_change.drop();
 db.mgo_txn.drop();
 db.mgo_txn.stash.drop();

 //remove collection data
 db.account.remove({});
 db.block.remove({});
 db.power_change.remove({});
 db.proposal.remove({});
 db.stake_role_candidate.remove({});
 db.stake_role_delegator.remove({});
 db.sync_task.remove({});
 db.tx_common.remove({});
 db.validator_up_time.remove({});
 db.tx_gas.remove({});
 db.tx_msg.remove({});
 db.uptime_change.remove({});
 db.mgo_txn.remove({});
 db.mgo_txn.stash.remove({});










