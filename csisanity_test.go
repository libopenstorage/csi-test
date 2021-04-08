/*
CSI Interface for OSD
Copyright 2017 Portworx

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
package csi

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/api/server/sdk"
	"github.com/libopenstorage/openstorage/cluster"
	clustermanager "github.com/libopenstorage/openstorage/cluster/manager"
	"github.com/libopenstorage/openstorage/config"
	"github.com/libopenstorage/openstorage/csi"
	"github.com/libopenstorage/openstorage/pkg/auth"
	"github.com/libopenstorage/openstorage/pkg/grpcserver"
	"github.com/libopenstorage/openstorage/pkg/role"
	"github.com/libopenstorage/openstorage/pkg/storagepolicy"
	volumedrivers "github.com/libopenstorage/openstorage/volume/drivers"
	"github.com/portworx/kvdb"
	"github.com/portworx/kvdb/mem"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	mockcluster "github.com/libopenstorage/openstorage/cluster/mock"
	mockdriver "github.com/libopenstorage/openstorage/volume/drivers/mock"

	"github.com/kubernetes-csi/csi-test/pkg/sanity"
	"github.com/kubernetes-csi/csi-test/utils"
)

const (
	mockDriverName   = "mock"
	testSharedSecret = "mysecret"
	fakeWithSched    = "fake-sched"
)

var (
	cm              cluster.Cluster
	systemUserToken string
)

func init() {
	setupFakeDriver()
}

// testServer is a simple struct used abstract
// the creation and setup of the gRPC CSI service
type testServer struct {
	conn   *grpc.ClientConn
	server grpcserver.Server
	m      *mockdriver.MockVolumeDriver
	c      *mockcluster.MockCluster
	mc     *gomock.Controller
	sdk    *sdk.Server
	port   string
	gwport string
	uds    string
}

func setupFakeDriver() {
	kv, err := kvdb.New(mem.Name, "fake_test", []string{}, nil, kvdb.LogFatalErrorCB)
	if err != nil {
		logrus.Panicf("Failed to initialize KVDB")
	}
	if err := kvdb.SetInstance(kv); err != nil {
		logrus.Panicf("Failed to set KVDB instance")
	}
	// Need to setup a fake cluster. No need to start it.
	clustermanager.Init(config.ClusterConfig{
		ClusterId: "fakecluster",
		NodeId:    "fakeNode",
	})
	cm, err = clustermanager.Inst()
	if err != nil {
		logrus.Panicf("Unable to initialize cluster manager: %v", err)
	}

	// Requires a non-nil cluster
	if err := volumedrivers.Register("fake", map[string]string{}); err != nil {
		logrus.Panicf("Unable to start volume driver fake: %v", err)
	}
}

func (s *testServer) setPorts() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	port := r.Intn(20000) + 10000

	s.port = fmt.Sprintf("%d", port)
	s.gwport = fmt.Sprintf("%d", port+1)
	s.uds = fmt.Sprintf("/tmp/osd-csi-ut-%d.sock", port)
}

func (s *testServer) MockDriver() *mockdriver.MockVolumeDriver {
	return s.m
}

func (s *testServer) MockCluster() *mockcluster.MockCluster {
	return s.c
}

func (s *testServer) Stop() {
	// Remove from registry
	volumedrivers.Remove("mock")

	// Shutdown servers
	s.conn.Close()
	s.server.Stop()
	s.sdk.Stop()

	// Check mocks
	s.mc.Finish()
}

func (s *testServer) Conn() *grpc.ClientConn {
	return s.conn
}

func (s *testServer) Server() grpcserver.Server {
	return s.server
}

func TestCSISanity(t *testing.T) {
	tester := &testServer{}
	tester.setPorts()
	tester.mc = gomock.NewController(&utils.SafeGoroutineTester{})

	clustermanager.Init(config.ClusterConfig{
		ClusterId: "fakecluster",
		NodeId:    "fakeNode",
	})
	cm, err := clustermanager.Inst()
	go func() {
		cm.Start(false, "9002", "")
	}()
	defer cm.Shutdown()

	// Setup sdk server
	kv, err := kvdb.New(mem.Name, "test", []string{}, nil, kvdb.LogFatalErrorCB)
	assert.NoError(t, err)
	kvdb.SetInstance(kv)
	stp, err := storagepolicy.Init()
	if err != nil {
		stp, _ = storagepolicy.Inst()
	}
	assert.NotNil(t, stp)
	rm, err := role.NewSdkRoleManager(kv)
	assert.NoError(t, err)

	selfsignedJwt, err := auth.NewJwtAuth(&auth.JwtAuthConfig{
		SharedSecret:  []byte(testSharedSecret),
		UsernameClaim: auth.UsernameClaimTypeName,
	})

	_ = rm
	_ = selfsignedJwt

	// setup sdk server
	sdk, err := sdk.New(&sdk.ServerConfig{
		DriverName:    "fake",
		Net:           "tcp",
		Address:       ":" + tester.port,
		RestPort:      tester.gwport,
		Cluster:       cm,
		Socket:        tester.uds,
		StoragePolicy: stp,
		AccessOutput:  ioutil.Discard,
		AuditOutput:   ioutil.Discard,
		// Auth disabled for now.
		// We're only sanity testing Client -> CSI -> SDK (No Auth)
		/*Security: &sdk.SecurityConfig{
			Role: rm,
			Authenticators: map[string]auth.Authenticator{
				"openstorage.io": selfsignedJwt,
			},
		},*/
	})
	assert.Nil(t, err)

	err = sdk.Start()
	assert.Nil(t, err)
	defer sdk.Stop()

	// Start CSI Server
	server, err := csi.NewOsdCsiServer(&csi.OsdCsiServerConfig{
		DriverName: "fake",
		Net:        "tcp",
		Address:    "127.0.0.1:0",
		Cluster:    cm,
		SdkUds:     tester.uds,
	})
	if err != nil {
		t.Fatalf("Unable to start csi server: %v", err)
	}
	server.Start()
	defer server.Stop()

	timeout := time.After(30 * time.Second)
	for {
		select {
		case <-timeout:
			t.Fatal("Timemout waiting for cluster to be ready")
		default:
		}
		cl, err := cm.Enumerate()
		if err != nil {
			t.Fatal("Unable to get cluster status")
		}
		if cl.Status == api.Status_STATUS_OK {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Start CSI Sanity test
	targetPath := "/tmp/mnt/csi"
	sanity.Test(t, &sanity.Config{
		Address:    server.Address(),
		TargetPath: targetPath,
		CreateTargetDir: func(p string) (string, error) {
			os.MkdirAll(p+"/target", os.FileMode(0755))
			return p, nil
		},
	})
}
