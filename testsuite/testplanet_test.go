// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	libuplinkInclude := CompileSharedAt(ctx, t, "../", "uplink", "storj.io/uplink-c")

	currentdir, err := os.Getwd()
	require.NoError(t, err)

	definition := Include{
		Header: filepath.Join(currentdir, "..", "uplink_definitions.h"),
	}

	ctests, err := filepath.Glob(filepath.Join("testplanet", "*_test.c"))
	require.NoError(t, err)

	t.Run("Testplanet", func(t *testing.T) {
		for _, ctest := range ctests {
			ctest := ctest
			testName := filepath.Base(ctest)
			t.Run(testName, func(t *testing.T) {
				t.Parallel()

				testexe := CompileC(ctx, t, CompileCOptions{
					Dest:    testName,
					Sources: []string{ctest},
					Includes: []Include{
						libuplinkInclude,
						definition,
						CLibMath,
						{
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

// CompileShared compiles pkg as c-shared.
//
// Note: cgo header paths are currently relative to package root.
// TODO: support inclusion from other directories.
func CompileSharedAt(ctx *testcontext.Context, t *testing.T, workDir, name string, pkg string) Include {
	t.Helper()

	if absDir, err := filepath.Abs(workDir); err == nil {
		workDir = absDir
	} else {
		t.Fatal(err)
	}

	base := ctx.File("build", name)

	args := []string{"build", "-buildmode", "c-shared"}
	args = append(args, "-o", base+".so", pkg)

	// not using race detector for c-shared
	cmd := exec.Command("go", args...)
	cmd.Dir = workDir

	t.Log("exec:", cmd.Args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(string(out))
		t.Fatal(err)
	}
	t.Log(string(out))

	return Include{Header: base + ".h", Library: base + ".so"}
}

// CLibMath is the standard C math library (see `man math.h`).
var CLibMath = Include{Standard: true, Library: "m"}

// CompileCOptions stores options for compiling C source to an executable.
type CompileCOptions struct {
	Dest     string
	Sources  []string
	Includes []Include
	NoWarn   bool
}

// CompileC compiles file as with gcc and adds the includes.
func CompileC(ctx *testcontext.Context, t *testing.T, opts CompileCOptions) string {
	t.Helper()

	exe := ctx.File("build", opts.Dest+".exe")

	args := []string{}
	if !opts.NoWarn {
		args = append(args, "-Wall")
		args = append(args, "-Wextra")
		args = append(args, "-Wpedantic")
		args = append(args, "-Werror")
	}
	args = append(args, "-ggdb")
	args = append(args, "-o", exe)

	// include headers
	for _, inc := range opts.Includes {
		if inc.Header != "" {
			args = append(args, "-I", filepath.Dir(inc.Header))
		}
	}

	// include sources
	args = append(args, opts.Sources...)

	// include libraries
	for _, inc := range opts.Includes {
		if inc.Library != "" {
			if inc.Standard {
				args = append(args,
					"-l"+inc.Library,
				)
				continue
			}
			if runtime.GOOS == "windows" {
				args = append(args,
					"-L"+filepath.Dir(inc.Library),
					"-l:"+filepath.Base(inc.Library),
				)
			} else {
				args = append(args, inc.Library)
			}
		}
	}

	/* #nosec G204 */ // This package is used for testing and the parameter's value are controlled by the above logic
	cmd := exec.Command("gcc", args...)
	t.Log("exec:", cmd.Args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(string(out))
		t.Fatal(err)
	}
	t.Log(string(out))

	return exe
}

// Include defines an includable library for gcc.
type Include struct {
	Header   string
	Library  string
	Standard bool
}
