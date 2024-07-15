/*
Copyright 2022 cuisongliu@qq.com.

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

package apply

import (
	"reflect"
	"testing"

	v2 "github.com/nebstudio/sealos/pkg/types/v1beta1"
)

func TestPreProcessIPList(t *testing.T) {
	type args struct {
		joinArgs *Cluster
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "node",
			args: args{
				joinArgs: &Cluster{
					Masters:     "",
					Nodes:       "192.168.1.1",
					ClusterName: "",
				},
			},
			wantErr: false,
		},
		{
			name: "master",
			args: args{
				joinArgs: &Cluster{
					Masters:     "192.168.1.1",
					Nodes:       "",
					ClusterName: "",
				},
			},
			wantErr: false,
		},
		{
			name: "node list",
			args: args{
				joinArgs: &Cluster{
					Masters:     "",
					Nodes:       "192.168.1.1,192.168.1.2,192.168.1.5",
					ClusterName: "",
				},
			},
			wantErr: false,
		},
		{
			name: "master list",
			args: args{
				joinArgs: &Cluster{
					Masters:     "192.168.1.1,192.168.1.2,192.168.1.5",
					Nodes:       "",
					ClusterName: "",
				},
			},
			wantErr: false,
		},
		{
			name: "node range",
			args: args{
				joinArgs: &Cluster{
					Masters:     "",
					Nodes:       "192.168.1.1-192.168.1.5",
					ClusterName: "",
				},
			},
			wantErr: false,
		},
		{
			name: "master range",
			args: args{
				joinArgs: &Cluster{
					Masters:     "192.168.1.1-192.168.1.5",
					Nodes:       "",
					ClusterName: "",
				},
			},
			wantErr: false,
		},
		{
			name: "node cidr",
			args: args{
				joinArgs: &Cluster{
					Masters:     "",
					Nodes:       "192.168.1.1/28",
					ClusterName: "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PreProcessIPList(tt.args.joinArgs); (err != nil) != tt.wantErr {
				t.Errorf("PreProcessIPList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsIPList(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			name: "single",
			args: "192.168.1.1",
			want: true,
		},
		{
			name: "multi",
			args: "192.168.1.2",
			want: true,
		},
		{
			name: "single with port",
			args: "192.168.1.1:22",
			want: true,
		},
		{
			name: "multi with port",
			args: "192.168.1.1:22,192.168.1.2:22",
			want: true,
		},
		{
			name: "invalid",
			args: "xxxx",
			want: false,
		},
		{
			name: "invalid with port",
			args: "xxxx:xx",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ok := IsIPList(tt.args); ok != tt.want {
				t.Errorf("IsIPList() = %v, want %v", ok, tt.want)
			}
		})
	}
}

func TestGetImagesDiff(t *testing.T) {
	current := []string{
		"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
		"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
		"hub.sealos.cn/nebstudio/helm:v3.11.0",
		"hub.sealos.cn/nebstudio/calico:v3.24.5",
	}
	desired := []string{
		"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
		"hub.sealos.cn/nebstudio/helm:v3.11.0",
		"hub.sealos.cn/nebstudio/helm:v3.11.0",
		"hub.sealos.cn/nebstudio/calico:v3.24.5",
		"hub.sealos.cn/nebstudio/nginx:v1.23.3",
	}

	diff := GetImagesDiff(current, desired)

	expected := []string{"hub.sealos.cn/nebstudio/nginx:v1.23.3"}

	if !reflect.DeepEqual(diff, expected) {
		t.Errorf("Unexpected diff. Expected %v, but got %v", expected, diff)
	}
}

func TestCompareImageSpecHash(t *testing.T) {
	currentImages := []string{
		"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
		"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
		"hub.sealos.cn/nebstudio/helm:v3.11.0",
		"hub.sealos.cn/nebstudio/calico:v3.24.5",
	}
	newImages := []string{
		"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
		"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
		"hub.sealos.cn/nebstudio/helm:v3.11.0",
		"hub.sealos.cn/nebstudio/calico:v3.24.5",
	}

	if !CompareImageSpecHash(currentImages, newImages) {
		t.Errorf("CompareImageSpecHash(%v, %v) = false; want true", currentImages, newImages)
	}

	newImages = []string{
		"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
		"hub.sealos.cn/nebstudio/helm:v3.11.0",
		"hub.sealos.cn/nebstudio/helm:v3.11.0",
		"hub.sealos.cn/nebstudio/calico:v3.24.5",
		"hub.sealos.cn/nebstudio/nginx:v1.23.3",
		"hub.sealos.cn/nebstudio/nginx:v1.23.3",
		"hub.sealos.cn/nebstudio/nginx:v1.23.3",
		"hub.sealos.cn/nebstudio/nginx:v1.23.5",
		"hub.sealos.cn/nebstudio/nginx:v1.23.4",
		"hub.sealos.cn/nebstudio/nginx:v1.23.3",
		"hub.sealos.cn/nebstudio/nginx:v1.23.2",
		"hub.sealos.cn/nebstudio/nginx:v1.23.1",
	}

	if CompareImageSpecHash(currentImages, newImages) {
		t.Errorf("CompareImageSpecHash(%v, %v) = true; want false", currentImages, newImages)
	}
}

func TestGetNewImages(t *testing.T) {
	// Test case 1
	actual := GetNewImages(nil, nil)
	if actual != nil {
		t.Errorf("GetNewImages(nil, nil) = %v, expected nil", actual)
	}

	// Test case 2
	currentCluster := &v2.Cluster{
		Spec: v2.ClusterSpec{
			Image: []string{
				"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
				"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
				"hub.sealos.cn/nebstudio/helm:v3.11.0",
				"hub.sealos.cn/nebstudio/calico:v3.24.5",
			},
		},
	}
	actual = GetNewImages(currentCluster, nil)
	if actual != nil {
		t.Errorf("GetNewImages(currentCluster, nil) = %v, expected nil", actual)
	}

	// Test case 3
	desiredCluster := &v2.Cluster{
		Spec: v2.ClusterSpec{
			Image: []string{
				"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
				"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
				"hub.sealos.cn/nebstudio/helm:v3.11.0",
				"hub.sealos.cn/nebstudio/calico:v3.24.5",
			},
		},
	}
	actual = GetNewImages(currentCluster, desiredCluster)
	if actual != nil {
		t.Errorf("GetNewImages(currentCluster, desiredCluster) = %v, expected nil", actual)
	}

	// Test case 4
	desiredCluster.Spec.Image = []string{
		"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
		"hub.sealos.cn/nebstudio/helm:v3.11.0",
		"hub.sealos.cn/nebstudio/helm:v3.11.0",
		"hub.sealos.cn/nebstudio/calico:v3.24.5",
		"hub.sealos.cn/nebstudio/nginx:v1.23.3",
	}

	expected := []string{"hub.sealos.cn/nebstudio/nginx:v1.23.3"}
	actual = GetNewImages(currentCluster, desiredCluster)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("GetNewImages(currentCluster, desiredCluster) = %v, expected %v", actual, expected)
	}

	// Test case 5
	desiredCluster.Spec.Image = []string{
		"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
		"hub.sealos.cn/nebstudio/helm:v3.11.0",
		"hub.sealos.cn/nebstudio/helm:v3.11.0",
		"hub.sealos.cn/nebstudio/calico:v3.24.5",
		"hub.sealos.cn/nebstudio/nginx:v1.23.3",
	}

	expected = []string{
		"hub.sealos.cn/nebstudio/kubernetes:v1.25.6",
		"hub.sealos.cn/nebstudio/helm:v3.11.0",
		"hub.sealos.cn/nebstudio/helm:v3.11.0",
		"hub.sealos.cn/nebstudio/calico:v3.24.5",
		"hub.sealos.cn/nebstudio/nginx:v1.23.3",
	}
	actual = GetNewImages(nil, desiredCluster)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("GetNewImages(nil, desiredCluster) = %v, expected %v", actual, expected)
	}
}
