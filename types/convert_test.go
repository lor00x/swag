// Copyright © 2024 zc2638 <zc2638@qq.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	Register("json.Number", Integer)
	assert.Equal(t, Integer, Get("json.Number"))
}

func TestParse(t *testing.T) {
	pointer := &[]int{1, 2, 3}
	pt := Parse(reflect.TypeOf(pointer))
	assert.Equal(t, Array, pt)

	jn := json.Number("123")
	pt = Parse(reflect.TypeOf(jn))
	assert.Equal(t, Get("json.Number"), pt)

	other := make(chan struct{})
	pt = Parse(reflect.TypeOf(other))
	assert.Equal(t, Unknown, pt)
}
