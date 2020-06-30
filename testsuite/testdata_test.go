// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"storj.io/common/storj"
	"storj.io/common/testcontext"
	"storj.io/storj/private/testplanet"
	"storj.io/uplink"
)

func TestC(t *testing.T) {
	ctx := testcontext.NewWithTimeout(t, 5*time.Minute)
	defer ctx.Cleanup()

	libuplinkInclude := ctx.CompileShared(t, "uplink", "storj.io/uplink-c")

	currentdir, err := os.Getwd()
	require.NoError(t, err)

	definition := testcontext.Include{
		Header: filepath.Join(currentdir, "..", "uplink_definitions.h"),
	}

	ctests, err := filepath.Glob(filepath.Join("testdata", "*_test.c"))
	require.NoError(t, err)

	t.Run("ALL", func(t *testing.T) {
		for _, ctest := range ctests {
			ctest := ctest
			testName := filepath.Base(ctest)
			t.Run(testName, func(t *testing.T) {
				t.Parallel()

				testexe := ctx.CompileC(t, testcontext.CompileCOptions{
					Dest:    testName,
					Sources: []string{ctest},
					Includes: []testcontext.Include{
						libuplinkInclude,
						definition,
						testcontext.CLibMath,
						testcontext.Include{
							Standard: true,
							Library:  "pthread",
						},
					},
				})

				testplanet.Run(t, testplanet.Config{
					SatelliteCount: 1, StorageNodeCount: 5, UplinkCount: 1,
					Reconfigure: testplanet.DisablePeerCAWhitelist,
				}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
					satellite := planet.Satellites[0]
					satelliteNodeURL := storj.NodeURL{
						ID:      satellite.ID(),
						Address: satellite.Addr(),
					}.String()

					apikey := planet.Uplinks[0].APIKey[satellite.ID()]
					uplinkConfig := uplink.Config{}

					access, err := uplinkConfig.RequestAccessWithPassphrase(ctx, satelliteNodeURL, apikey.Serialize(), "mypassphrase")
					require.NoError(t, err)
					accessString, err := access.Serialize()
					require.NoError(t, err)

					cmd := exec.Command(testexe)
					cmd.Dir = filepath.Dir(testexe)
					cmd.Env = append(os.Environ(),
						"SATELLITE_0_ADDR="+satelliteNodeURL,
						"UPLINK_0_APIKEY="+apikey.Serialize(),
						"UPLINK_0_ACCESS="+accessString,
						"TMP_DIR="+ctx.Dir("c_temp"),
					)

					out, err := cmd.CombinedOutput()
					if err != nil {
						t.Error(string(out))
						t.Fatal(err)
					} else {
						t.Log(string(out))
					}
				})
			})
		}
	})
}
