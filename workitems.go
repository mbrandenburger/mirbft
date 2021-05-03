/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mirbft

import (
	"github.com/hyperledger-labs/mirbft/pkg/pb/msgs"
	"github.com/hyperledger-labs/mirbft/pkg/pb/state"
	"github.com/hyperledger-labs/mirbft/pkg/statemachine"
)

type WorkItems struct {
	walActions     *statemachine.ActionList
	netActions     *statemachine.ActionList
	hashActions    *statemachine.ActionList
	clientActions  *statemachine.ActionList
	appActions     *statemachine.ActionList
	reqStoreEvents *statemachine.EventList
	resultEvents   *statemachine.EventList
}

func NewWorkItems() *WorkItems {
	return &WorkItems{
		walActions:     &statemachine.ActionList{},
		netActions:     &statemachine.ActionList{},
		hashActions:    &statemachine.ActionList{},
		clientActions:  &statemachine.ActionList{},
		appActions:     &statemachine.ActionList{},
		reqStoreEvents: &statemachine.EventList{},
		resultEvents:   &statemachine.EventList{},
	}
}

func (wi *WorkItems) ClearWALActions() {
	wi.walActions = &statemachine.ActionList{}
}

func (wi *WorkItems) ClearNetActions() {
	wi.netActions = &statemachine.ActionList{}
}

func (wi *WorkItems) ClearHashActions() {
	wi.hashActions = &statemachine.ActionList{}
}

func (wi *WorkItems) ClearClientActions() {
	wi.clientActions = &statemachine.ActionList{}
}

func (wi *WorkItems) ClearAppActions() {
	wi.appActions = &statemachine.ActionList{}
}

func (wi *WorkItems) ClearReqStoreEvents() {
	wi.reqStoreEvents = &statemachine.EventList{}
}

func (wi *WorkItems) ClearResultEvents() {
	wi.resultEvents = &statemachine.EventList{}
}

func (wi *WorkItems) WALActions() *statemachine.ActionList {
	if wi.walActions == nil {
		wi.walActions = &statemachine.ActionList{}
	}
	return wi.walActions
}

func (wi *WorkItems) NetActions() *statemachine.ActionList {
	if wi.netActions == nil {
		wi.netActions = &statemachine.ActionList{}
	}
	return wi.netActions
}

func (wi *WorkItems) HashActions() *statemachine.ActionList {
	if wi.hashActions == nil {
		wi.hashActions = &statemachine.ActionList{}
	}
	return wi.hashActions
}

func (wi *WorkItems) ClientActions() *statemachine.ActionList {
	if wi.clientActions == nil {
		wi.clientActions = &statemachine.ActionList{}
	}
	return wi.clientActions
}

func (wi *WorkItems) AppActions() *statemachine.ActionList {
	if wi.appActions == nil {
		wi.appActions = &statemachine.ActionList{}
	}
	return wi.appActions
}

func (wi *WorkItems) ReqStoreEvents() *statemachine.EventList {
	if wi.reqStoreEvents == nil {
		wi.reqStoreEvents = &statemachine.EventList{}
	}
	return wi.reqStoreEvents
}

func (wi *WorkItems) ResultEvents() *statemachine.EventList {
	if wi.resultEvents == nil {
		wi.resultEvents = &statemachine.EventList{}
	}
	return wi.resultEvents
}

func (wi *WorkItems) AddHashResults(events *statemachine.EventList) {
	wi.ResultEvents().PushBackList(events)
}

func (wi *WorkItems) AddNetResults(events *statemachine.EventList) {
	wi.ResultEvents().PushBackList(events)
}

func (wi *WorkItems) AddAppResults(events *statemachine.EventList) {
	wi.ResultEvents().PushBackList(events)
}

func (wi *WorkItems) AddClientResults(events *statemachine.EventList) {
	wi.ReqStoreEvents().PushBackList(events)
}

func (wi *WorkItems) AddWALResults(actions *statemachine.ActionList) {
	wi.NetActions().PushBackList(actions)
}

func (wi *WorkItems) AddReqStoreResults(events *statemachine.EventList) {
	wi.ResultEvents().PushBackList(events)
}

func (wi *WorkItems) AddStateMachineResults(actions *statemachine.ActionList) {
	// First we'll handle everything that's not a network send
	iter := actions.Iterator()
	for action := iter.Next(); action != nil; action = iter.Next() {
		switch t := action.Type.(type) {
		case *state.Action_Send:
			walDependent := false
			// TODO, make sure this switch captures all the safe ones
			switch t.Send.Msg.Type.(type) {
			case *msgs.Msg_RequestAck:
			case *msgs.Msg_Checkpoint:
			case *msgs.Msg_FetchBatch:
			case *msgs.Msg_ForwardBatch:
			default:
				walDependent = true
			}
			if walDependent {
				wi.WALActions().PushBack(action)
			} else {
				wi.NetActions().PushBack(action)
			}
		case *state.Action_Hash:
			wi.HashActions().PushBack(action)
		case *state.Action_AppendWriteAhead:
			wi.WALActions().PushBack(action)
		case *state.Action_TruncateWriteAhead:
			wi.WALActions().PushBack(action)
		case *state.Action_Commit:
			wi.AppActions().PushBack(action)
		case *state.Action_Checkpoint:
			wi.AppActions().PushBack(action)
		case *state.Action_AllocatedRequest:
			wi.ClientActions().PushBack(action)
		case *state.Action_CorrectRequest:
			wi.ClientActions().PushBack(action)
		case *state.Action_StateApplied:
			wi.ClientActions().PushBack(action)
			// TODO, create replicas
		case *state.Action_ForwardRequest:
			// XXX address
		case *state.Action_StateTransfer:
			wi.AppActions().PushBack(action)
		}
	}
}
