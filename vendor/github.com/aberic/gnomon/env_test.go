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
	"os"
	"testing"
)

func TestEnvCommon_Get(t *testing.T) {
	_ = os.Setenv("HELLO", "hello")
	t.Log("HELLO =", Env().Get("HELLO"))
}

func TestEnvCommon_GetD(t *testing.T) {
	_ = os.Setenv("HELLO", "hello")
	t.Log("HELLO =", Env().GetD("HELLO", "WORLD"))
	t.Log("WORLD =", Env().GetD("WORLD", "HELLO"))
}

func TestEnvCommon_GetInt(t *testing.T) {
	_ = os.Setenv("HELLO", "100")
	i, _ := Env().GetInt("HELLO")
	t.Log("HELLO =", i)
	_ = os.Setenv("HELLO", "WORLD")
	_, err := Env().GetInt("HELLO")
	t.Skip(err)
}

func TestEnvCommon_GetIntD(t *testing.T) {
	_ = os.Setenv("HELLO", "100")
	t.Log("HELLO =", Env().GetIntD("HELLO", 10))
	_ = os.Setenv("HELLO", "WORLD")
	t.Log("HELLO =", Env().GetIntD("HELLO", 10))
}

func TestEnvCommon_GetInt64(t *testing.T) {
	_ = os.Setenv("HELLO", "100")
	i, _ := Env().GetInt64("HELLO")
	t.Log("HELLO =", i)
	_ = os.Setenv("HELLO", "WORLD")
	_, err := Env().GetInt64("HELLO")
	t.Skip(err)
}

func TestEnvCommon_GetInt64D(t *testing.T) {
	_ = os.Setenv("HELLO", "100")
	t.Log("HELLO =", Env().GetInt64D("HELLO", 10))
	_ = os.Setenv("HELLO", "WORLD")
	t.Log("HELLO =", Env().GetInt64D("HELLO", 10))
}

func TestEnvCommon_GetUint64(t *testing.T) {
	_ = os.Setenv("HELLO", "100")
	i, _ := Env().GetUint64("HELLO")
	t.Log("HELLO =", i)
	_ = os.Setenv("HELLO", "WORLD")
	_, err := Env().GetUint64("HELLO")
	t.Skip(err)
}

func TestEnvCommon_GetUint64D(t *testing.T) {
	_ = os.Setenv("HELLO", "100")
	t.Log("HELLO =", Env().GetUint64D("HELLO", 10))
	_ = os.Setenv("HELLO", "WORLD")
	t.Log("HELLO =", Env().GetUint64D("HELLO", 10))
}

func TestEnvCommon_GetFloat64(t *testing.T) {
	_ = os.Setenv("HELLO", "100.3254")
	i, _ := Env().GetFloat64("HELLO")
	t.Log("HELLO =", i)
	_ = os.Setenv("HELLO", "WORLD")
	_, err := Env().GetFloat64("HELLO")
	t.Skip(err)
}

func TestEnvCommon_GetFloat64D(t *testing.T) {
	_ = os.Setenv("HELLO", "100.3254")
	t.Log("HELLO =", Env().GetFloat64D("HELLO", 100.3254))
	_ = os.Setenv("HELLO", "WORLD")
	t.Log("HELLO =", Env().GetFloat64D("HELLO", 100.32541))
}

func TestEnvCommon_GetBool(t *testing.T) {
	_ = os.Setenv("HELLO", "true")
	t.Log("HELLO =", Env().GetBool("HELLO"))
	_ = os.Setenv("HELLO", "false")
	t.Log("HELLO =", Env().GetBool("HELLO"))
	_ = os.Setenv("HELLO", "WORLD")
	t.Log("HELLO =", Env().GetBool("HELLO"))
}
