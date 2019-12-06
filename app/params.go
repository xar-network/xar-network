/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package app

// Simulation parameter constants
const (
	StakePerAccount                        = "stake_per_account"
	InitiallyBondedValidators              = "initially_bonded_validators"
	OpWeightDeductFee                      = "op_weight_deduct_fee"
	OpWeightMsgSend                        = "op_weight_msg_send"
	OpWeightMsgMultiSend                   = "op_weight_msg_multisend"
	OpWeightMsgSetWithdrawAddress          = "op_weight_msg_set_withdraw_address"
	OpWeightMsgWithdrawDelegationReward    = "op_weight_msg_withdraw_delegation_reward"
	OpWeightMsgWithdrawValidatorCommission = "op_weight_msg_withdraw_validator_commission"
	OpWeightSubmitTextProposal             = "op_weight_submit_text_proposal"
	OpWeightSubmitCommunitySpendProposal   = "op_weight_submit_community_spend_proposal"
	OpWeightSubmitParamChangeProposal      = "op_weight_submit_param_change_proposal"
	OpWeightMsgDeposit                     = "op_weight_msg_deposit"
	OpWeightMsgVote                        = "op_weight_msg_vote"
	OpWeightMsgCreateValidator             = "op_weight_msg_create_validator"
	OpWeightMsgEditValidator               = "op_weight_msg_edit_validator"
	OpWeightMsgDelegate                    = "op_weight_msg_delegate"
	OpWeightMsgUndelegate                  = "op_weight_msg_undelegate"
	OpWeightMsgBeginRedelegate             = "op_weight_msg_begin_redelegate"
	OpWeightMsgUnjail                      = "op_weight_msg_unjail"
)
