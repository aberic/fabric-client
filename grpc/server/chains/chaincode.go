/*
 * Copyright (c) 2019.. Aberic - All Rights Reserved.
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

package chains

import (
	"errors"
	"github.com/aberic/fabric-client/config"
	"github.com/aberic/fabric-client/core"
	"github.com/aberic/fabric-client/geneses"
	pb "github.com/aberic/fabric-client/grpc/proto/chain"
	"github.com/aberic/fabric-client/service"
	"github.com/aberic/gnomon"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ChainCodeServer struct {
}

func (c *ChainCodeServer) UploadCC(stream pb.LedgerChainCode_UploadCCServer) error {
	var (
		upload *pb.Upload
		data   []byte
	)
	upload = &pb.Upload{}
	data = make([]byte, 0)
	for {
		uploadRecv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		gnomon.Log().Debug("code data", gnomon.Log().Field("bytes", uploadRecv.Data))
		data = append(data, uploadRecv.Data...)
		upload.LedgerName = uploadRecv.LedgerName
		upload.Name = uploadRecv.Name
		upload.Version = uploadRecv.Version
	}
	if gnomon.String().IsEmpty(upload.LedgerName) || gnomon.String().IsEmpty(upload.Name) || gnomon.String().IsEmpty(upload.Version) {
		return errors.New("upload nil")
	}
	source, path, zipPath := geneses.ChainCodePath(upload.LedgerName, upload.Name, upload.Version)
	if _, err := gnomon.File().Append(zipPath, data, true); nil != err {
		return err
	}
	if err := gnomon.File().DeCompressZip(zipPath, strings.Join([]string{source, "src", path}, "/")); nil != err {
		return err
	}
	if err := os.Remove(zipPath); nil != err {
		return err
	}
	lastIndex := strings.LastIndex(zipPath, "/")
	zipParentPath := zipPath[0:lastIndex]
	codePath := getSinglePath(zipParentPath)
	paths := strings.Split(codePath, "/src/")
	path = paths[len(paths)-1]
	return stream.SendAndClose(&pb.ResultUpload{Code: pb.Code_Success, Source: source, Path: path})
}

func getSinglePath(path string) string {
	fileInfo, _ := ioutil.ReadDir(path)
	fileInfoLen := len(fileInfo)
	if fileInfoLen != 1 {
		fileNum := 0
		var fileInfoPath string
		for _, info := range fileInfo {
			if !strings.HasPrefix(info.Name(), ".") {
				fileInfoPath = info.Name()
				fileNum++
			}
		}
		if fileNum != 1 {
			return path
		}
		return getSinglePath(filepath.Join(path, fileInfoPath))
	} else if fileInfoLen == 1 && !fileInfo[0].IsDir() {
		return path
	}
	return getSinglePath(filepath.Join(path, fileInfo[0].Name()))
}

func (c *ChainCodeServer) InstallCC(ctx context.Context, in *pb.Install) (*pb.Result, error) {
	var (
		res  *sdk.Result
		conf *config.Config
	)
	if conf = service.Configs[in.ConfigID]; nil == conf {
		return nil, errors.New("config client is not exist")
	}
	if res = sdk.Install(in.OrgName, in.OrgUser, in.PeerName, in.Name, in.Source, in.Path, in.Version, service.GetBytes(in.ConfigID)); res.ResultCode == sdk.Success {
		gnomon.Log().Info("InstallCC", gnomon.Log().Field("success", res))
		return &pb.Result{Code: pb.Code_Success, Data: res.Data.(string)}, nil
	}
	return &pb.Result{Code: pb.Code_Fail, ErrMsg: res.Msg}, errors.New(res.Msg)
}

func (c *ChainCodeServer) InstalledCC(ctx context.Context, in *pb.Installed) (*pb.ResultCCList, error) {
	var (
		res              *sdk.Result
		chainCodeInfoArr *sdk.ChainCodeInfoArr
		conf             *config.Config
	)
	if conf = service.Configs[in.ConfigID]; nil == conf {
		return nil, errors.New("config client is not exist")
	}
	if res = sdk.Installed(in.OrgName, in.OrgUser, in.PeerName, service.GetBytes(in.ConfigID)); res.ResultCode == sdk.Success {
		chainCodeInfoArr = res.Data.(*sdk.ChainCodeInfoArr)
		gnomon.Log().Info("InstalledCC", gnomon.Log().Field("chainCodeInfoArr", chainCodeInfoArr))
		chainCodes := chainCodeInfoArr.ChainCodes
		gnomon.Log().Info("InstalledCC", gnomon.Log().Field("chainCodes", chainCodes))
		data := make([]*pb.ChainCodeInfo, len(chainCodes))
		for index, code := range chainCodes {
			data[index] = &pb.ChainCodeInfo{}
			data[index].Name = code.Name
			data[index].Version = code.Version
			data[index].Path = code.Path
			data[index].Input = code.Input
			data[index].Escc = code.Escc
			data[index].Vscc = code.Vscc
			data[index].Id = code.Id
		}
		return &pb.ResultCCList{Code: pb.Code_Success, List: &pb.CCList{Data: data}}, nil
	}
	return &pb.ResultCCList{Code: pb.Code_Fail, ErrMsg: res.Msg}, errors.New(res.Msg)
}

func (c *ChainCodeServer) InstantiateCC(ctx context.Context, in *pb.Instantiate) (*pb.Result, error) {
	var (
		res  *sdk.Result
		conf *config.Config
	)
	if conf = service.Configs[in.ConfigID]; nil == conf {
		return nil, errors.New("config client is not exist")
	}
	if res = sdk.Instantiate(in.OrgName, in.OrgUser, in.PeerName, in.ChannelID, in.Name, in.Path, in.Version, in.OrgPolicies,
		in.Args, service.GetBytes(in.ConfigID)); res.ResultCode == sdk.Success {
		gnomon.Log().Info("InstantiateCC", gnomon.Log().Field("success", res))
		return &pb.Result{Code: pb.Code_Success, Data: res.Data.(string)}, nil
	}
	return &pb.Result{Code: pb.Code_Fail, ErrMsg: res.Msg}, errors.New(res.Msg)
}

func (c *ChainCodeServer) InstantiatedCC(ctx context.Context, in *pb.Instantiated) (*pb.ResultCCList, error) {
	var (
		res              *sdk.Result
		chainCodeInfoArr *sdk.ChainCodeInfoArr
		conf             *config.Config
	)
	if conf = service.Configs[in.ConfigID]; nil == conf {
		return nil, errors.New("config client is not exist")
	}
	if res = sdk.Instantiated(in.OrgName, in.OrgUser, in.ChannelID, in.PeerName, service.GetBytes(in.ConfigID)); res.ResultCode == sdk.Success {
		chainCodeInfoArr = res.Data.(*sdk.ChainCodeInfoArr)
		gnomon.Log().Info("InstantiatedCC", gnomon.Log().Field("chainCodeInfoArr", chainCodeInfoArr))
		chainCodes := chainCodeInfoArr.ChainCodes
		gnomon.Log().Info("InstantiatedCC", gnomon.Log().Field("chainCodes", chainCodes))
		data := make([]*pb.ChainCodeInfo, len(chainCodes))
		for index, code := range chainCodes {
			data[index] = &pb.ChainCodeInfo{}
			data[index].Name = code.Name
			data[index].Version = code.Version
			data[index].Path = code.Path
			data[index].Input = code.Input
			data[index].Escc = code.Escc
			data[index].Vscc = code.Vscc
			data[index].Id = code.Id
		}
		return &pb.ResultCCList{Code: pb.Code_Success, List: &pb.CCList{Data: data}}, nil
	}
	return &pb.ResultCCList{Code: pb.Code_Fail, ErrMsg: res.Msg}, errors.New(res.Msg)
}

func (c *ChainCodeServer) UpgradeCC(ctx context.Context, in *pb.Upgrade) (*pb.Result, error) {
	var (
		res  *sdk.Result
		conf *config.Config
	)
	if conf = service.Configs[in.ConfigID]; nil == conf {
		return nil, errors.New("config client is not exist")
	}
	if res = sdk.Upgrade(in.OrgName, in.OrgUser, in.PeerName, in.ChannelID, in.Name, in.Path, in.Version, in.OrgPolicies,
		in.Args, service.GetBytes(in.ConfigID)); res.ResultCode == sdk.Success {
		gnomon.Log().Info("UpgradeCC", gnomon.Log().Field("success", res))
		return &pb.Result{Code: pb.Code_Success, Data: res.Data.(string)}, nil
	}
	return &pb.Result{Code: pb.Code_Fail, ErrMsg: res.Msg}, errors.New(res.Msg)
}

func (c *ChainCodeServer) InvokeCC(ctx context.Context, in *pb.Invoke) (*pb.Result, error) {
	var (
		res  *sdk.Result
		conf *config.Config
	)
	if conf = service.Configs[in.ConfigID]; nil == conf {
		return nil, errors.New("config client is not exist")
	}
	if res = sdk.Invoke(in.ChainCodeID, in.OrgName, in.OrgUser, in.ChannelID, in.Fcn, in.Args, in.TargetEndpoints,
		service.GetBytes(in.ConfigID)); res.ResultCode == sdk.Success {
		gnomon.Log().Debug("InvokeCC", gnomon.Log().Field("success", res))
		return &pb.Result{Code: pb.Code_Success, Data: res.Data.(string)}, nil
	}
	return &pb.Result{Code: pb.Code_Fail, ErrMsg: res.Msg}, errors.New(res.Msg)
}

func (c *ChainCodeServer) InvokeCCAsync(ctx context.Context, in *pb.InvokeAsync) (*pb.Result, error) {
	var (
		res  *sdk.Result
		conf *config.Config
	)
	if conf = service.Configs[in.ConfigID]; nil == conf {
		return nil, errors.New("config client is not exist")
	}
	if res = sdk.InvokeAsync(in.ChainCodeID, in.OrgName, in.OrgUser, in.ChannelID, in.Callback, in.Fcn, in.Args, in.TargetEndpoints,
		service.GetBytes(in.ConfigID)); res.ResultCode == sdk.Success {
		gnomon.Log().Debug("InvokeCCAsync", gnomon.Log().Field("success", res))
		return &pb.Result{Code: pb.Code_Success, Data: res.Data.(string)}, nil
	}
	return &pb.Result{Code: pb.Code_Fail, ErrMsg: res.Msg}, errors.New(res.Msg)
}

func (c *ChainCodeServer) QueryCC(ctx context.Context, in *pb.Query) (*pb.Result, error) {
	var (
		res  *sdk.Result
		conf *config.Config
	)
	if conf = service.Configs[in.ConfigID]; nil == conf {
		return nil, errors.New("config client is not exist")
	}
	if res = sdk.Query(in.ChainCodeID, in.OrgName, in.OrgUser, in.ChannelID, in.Fcn, in.Args, in.TargetEndpoints,
		service.GetBytes(in.ConfigID)); res.ResultCode == sdk.Success {
		gnomon.Log().Debug("QueryCC", gnomon.Log().Field("success", res))
		return &pb.Result{Code: pb.Code_Success, Data: res.Data.(string)}, nil
	}
	return &pb.Result{Code: pb.Code_Fail, ErrMsg: res.Msg}, errors.New(res.Msg)
}
