/*
Copyright 2023 cuisongliu@qq.com.

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

package operators

import (
	"github.com/nebstudio/sealos/test/e2e/testhelper/cmd"
)

var _ FakeClusterInterface = &fakeClusterClient{}

func newClusterClient(sealosCmd *cmd.SealosCmd, clusterName string) FakeClusterInterface {
	return &fakeClusterClient{
		SealosCmd:   sealosCmd,
		clusterName: clusterName,
	}
}

type fakeClusterClient struct {
	*cmd.SealosCmd
	clusterName string
}

func (c *fakeClusterClient) Run(images ...string) error {
	return c.SealosCmd.Run(&cmd.RunOptions{
		Cluster: c.clusterName,
		Images:  images,
	})
}

func (c *fakeClusterClient) Apply(file string) error {
	return c.SealosCmd.Apply(&cmd.ApplyOptions{
		Clusterfile: file,
	})
}

func (c *fakeClusterClient) Reset() error {
	return c.SealosCmd.Reset(&cmd.ResetOptions{
		Cluster: c.clusterName,
		Force:   true,
	})
}
