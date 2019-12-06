/*

Copyright 2016 All in Bits, Inc
Copyright 2018 public-chain
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

package record

import (
	"github.com/xar-network/xar-network/x/record/client"
	"github.com/xar-network/xar-network/x/record/internal/keeper"
	"github.com/xar-network/xar-network/x/record/internal/types"
)

type (
	Keeper     = keeper.Keeper
	RecordInfo = types.RecordInfo
)

var (
	NewKeeper       = keeper.NewKeeper
	NewModuleClient = client.NewModuleClient
	RegisterCodec   = types.RegisterCodec
	ModuleCdc       = types.ModuleCdc
	ModuleName      = types.ModuleName
)

const (
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	DefaultCodespace  = types.DefaultCodespace
)
