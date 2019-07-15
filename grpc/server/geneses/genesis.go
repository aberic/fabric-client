/*
 * Copyright (c) 2019. ENNOO - All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package geneses

import (
	"github.com/ennoo/fabric-client/geneses"
	pb "github.com/ennoo/fabric-client/grpc/proto/geneses"
	"golang.org/x/net/context"
)

type GenesisServer struct {
}

func (g *GenesisServer) GenerateYml(ctx context.Context, in *pb.Generate) (*pb.String, error) {
	return &pb.String{}, geneses.GenerateYml(in)
}

func (g *GenesisServer) GenerateCryptoFiles(ctx context.Context, in *pb.Crypto) (*pb.CMDOut, error) {
	line, cmds, err := geneses.GenerateCryptoFiles(in)
	return &pb.CMDOut{Line: int32(line), Prints: cmds}, err
}

func (g *GenesisServer) GenerateGenesisBlock(ctx context.Context, in *pb.Crypto) (*pb.CMDOut, error) {
	line, cmds, err := geneses.GenerateGenesisBlock(in)
	return &pb.CMDOut{Line: int32(line), Prints: cmds}, err
}

func (g *GenesisServer) GenerateChannelTX(ctx context.Context, in *pb.ChannelTX) (*pb.CMDOut, error) {
	line, cmds, err := geneses.GenerateChannelTX(in)
	return &pb.CMDOut{Line: int32(line), Prints: cmds}, err
}
