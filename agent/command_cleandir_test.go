/*
 * Copyright 2016 ThoughtWorks, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package agent_test

import (
	"bytes"
	"github.com/bmatcuk/doublestar"
	. "github.com/gocd-contrib/gocd-golang-agent/agent"
	"github.com/xli/assert"
	"io/ioutil"
	"path/filepath"
	"sort"
	"testing"
)

func TestCleandir(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "cleandir-test")
	assert.Nil(t, err)
	createTestProject(tmpDir)

	var log bytes.Buffer
	err = Cleandir(&log, tmpDir, "src/hello", "test/world2")
	assert.Nil(t, err)

	matches, err := doublestar.Glob(filepath.Join(tmpDir, "**/*.txt"))
	assert.Nil(t, err)
	sort.Strings(matches)
	expected := []string{
		"src/hello/3.txt",
		"src/hello/4.txt",
		"test/world2/10.txt",
		"test/world2/11.txt",
	}

	for i, f := range matches {
		actual := f[len(tmpDir)+1:]
		assert.Equal(t, expected[i], actual)
	}
	expectedLog := `Deleting file 0.txt
Deleting file src/1.txt
Deleting file src/2.txt
Keeping folder src/hello
Deleting file test/5.txt
Deleting file test/6.txt
Deleting file test/7.txt
Deleting file test/world/10.txt
Deleting file test/world/11.txt
Deleting file test/world/8.txt
Deleting file test/world/9.txt
Keeping folder test/world2
`
	assert.Equal(t, expectedLog, log.String())
}

func TestShouldFailWhenCleandirAllowsContainsPathThatIsOutsideOfBaseDir(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "cleandir-test2")
	assert.Nil(t, err)
	createTestProject(tmpDir)

	var log bytes.Buffer
	err = Cleandir(&log, tmpDir, "test/world2", "./../")
	assert.NotNil(t, err)
	assert.Equal(t, "", log.String())
}
