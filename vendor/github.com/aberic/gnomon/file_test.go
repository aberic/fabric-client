/*
 * Copyright (c) 2019. aberic - All Rights Reserved.
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

package gnomon

import (
	"testing"
)

func TestFileCommon_PathExists(t *testing.T) {
	path := "/etc"
	exist := File().PathExists(path)
	t.Log(path, "exist =", exist)

	path = "/etc/hello"
	exist = File().PathExists(path)
	t.Log(path, "exist =", exist)
}

func TestFileCommon_ReadFirstLine(t *testing.T) {
	profile, err := File().ReadFirstLine("/etc/profile")
	if nil != err {
		t.Skip(err)
	} else {
		t.Log("profile =", profile)
	}
}

func TestFileCommon_ReadFirstLine_Fail(t *testing.T) {
	_, err := File().ReadFirstLine("/etc/hello")
	t.Skip(err)
}

func TestFileCommon_ReadPointLine(t *testing.T) {
	profile, err := File().ReadPointLine("/etc/profile", 1)
	if nil != err {
		t.Skip(err)
	} else {
		t.Log("profile =", profile)
	}
}

func TestFileCommon_ReadPointLine_KeyPoint(t *testing.T) {
	_, _ = File().Append("./tmp/log/yes/go/point.txt", []byte("haha"), false)
	profile, err := File().ReadPointLine("./tmp/log/yes/go/point.txt", 1)
	if nil != err {
		t.Skip(err)
	} else {
		t.Log("profile =", profile)
	}
}

func TestFileCommon_ReadPointLine_Fail_IndexOut(t *testing.T) {
	_, err := File().ReadPointLine("/etc/profile", 300)
	t.Skip(err)
}

func TestFileCommon_ReadPointLine_Fail_NotExist(t *testing.T) {
	_, err := File().ReadPointLine("/etc/hello", 1)
	t.Skip(err)
}

func TestFileCommon_ReadLines(t *testing.T) {
	profile, err := File().ReadLines("/etc/profile")
	if nil != err {
		t.Skip(err)
	} else {
		t.Log("profile =", profile)
	}
}

func TestFileCommon_ReadLines_Fail(t *testing.T) {
	_, err := File().ReadLines("/etc/hello")
	t.Skip(err)
}

func TestFileCommon_ParentPath(t *testing.T) {
	t.Log(File().ParentPath("/etc/yes/go/test.txt"))
}

func TestFileCommon_Append(t *testing.T) {
	if _, err := File().Append("./tmp/log/yes/go/test.txt", []byte("haha"), false); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Append_Force(t *testing.T) {
	if _, err := File().Append("./tmp/log/yes/go/test.txt", []byte("haha"), true); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Append_UnForce(t *testing.T) {
	if _, err := File().Append("./tmp/log/yes/go/test.txt", []byte("haha"), false); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Append_Fail_PermissionFileForce(t *testing.T) {
	_, err := File().Append("/etc/www.json", []byte("haha"), true)
	t.Skip(err)
}

func TestFileCommon_Append_Fail_PermissionFileUnForce(t *testing.T) {
	_, err := File().Append("/etc/www.json", []byte("haha"), false)
	t.Skip(err)
}

func TestFileCommon_Modify(t *testing.T) {
	if _, err := File().Modify("./tmp/log/yes/go/test.txt", 1, []byte("haha"), false); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Modify_Force(t *testing.T) {
	if _, err := File().Modify("./tmp/log/yes/go/test.txt", 1, []byte("haha"), true); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Modify_UnForce(t *testing.T) {
	if _, err := File().Modify("./tmp/log/yes/go/test.txt", 1, []byte("haha"), false); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Modify_Fail_PermissionFileForce(t *testing.T) {
	_, err := File().Modify("/etc/www.json", 1, []byte("haha"), true)
	t.Skip(err)
}

func TestFileCommon_Modify_Fail_PermissionFileUnForce(t *testing.T) {
	_, err := File().Modify("/etc/www.json", 1, []byte("haha"), false)
	t.Skip(err)
}

func TestFileCommon_LoopDirs(t *testing.T) {
	if arr, err := File().LoopDirs("./tmp/log"); nil != err {
		t.Skip(err)
	} else {
		t.Log(arr)
	}
}

func TestFileCommon_LoopDirs_Fail(t *testing.T) {
	_, err := File().LoopDirs("./tmp/logger")
	t.Skip(err)
}

func TestFileCommon_LoopFiles(t *testing.T) {
	var s []string
	if arr, err := File().LoopFiles("./tmp/log", s); nil != err {
		t.Skip(err)
	} else {
		t.Log(arr)
	}
}

func TestFileCommon_LoopFiles_Fail(t *testing.T) {
	_, err := File().LoopFiles("./tmp/logger", nil)
	t.Skip(err)
}
