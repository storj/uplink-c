// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"storj.io/common/pb"
	"storj.io/common/testcontext"
	"storj.io/drpc/drpcmux"
	"storj.io/drpc/drpcserver"
)

type dRPCServerMock struct {
	pb.DRPCEdgeAuthServer
}

func (g *dRPCServerMock) RegisterAccess(context.Context, *pb.EdgeRegisterAccessRequest) (*pb.EdgeRegisterAccessResponse, error) {
	return &pb.EdgeRegisterAccessResponse{
		AccessKeyId: "l5pucy3dmvzxgs3fpfewix27l5pq",
		SecretKey:   "l5pvgzldojsxis3fpfpv6x27l5pv6x27l5pv6x27l5pv6",
		Endpoint:    "https://gateway.example",
	}, nil
}

func TestEdgeBindings(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	certificatePEM, privateKeyPEM := createSelfSignedCertificate(t, "localhost")

	certificate, err := tls.X509KeyPair(certificatePEM, privateKeyPEM)
	require.NoError(t, err)

	cancelCtx, authCancel := context.WithCancel(ctx)
	defer authCancel()

	portTLS := startMockAuthServiceTLS(t, ctx, cancelCtx, certificate)
	portUnencrypted := startMockAuthServiceUnencrypted(t, ctx, cancelCtx)

	authServiceTLSAddr := "localhost:" + strconv.Itoa(portTLS)
	authServiceUnencryptedAddr := "localhost:" + strconv.Itoa(portUnencrypted)

	libuplinkInclude := CompileSharedAt(ctx, t, "../", "uplink", "storj.io/uplink-c")

	currentdir, err := os.Getwd()
	require.NoError(t, err)

	definition := Include{
		Header: filepath.Join(currentdir, "..", "uplink_definitions.h"),
	}

	ctests, err := filepath.Glob(filepath.Join("edge", "*_test.c"))
	require.NoError(t, err)

	t.Run("Edge", func(t *testing.T) {
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

				cmd := exec.Command(testexe)
				cmd.Dir = filepath.Dir(testexe)
				cmd.Env = append(os.Environ(),
					"AUTH_SERVICE_TLS_ADDR="+authServiceTLSAddr,
					"AUTH_SERVICE_UNENCRYPTED_ADDR="+authServiceUnencryptedAddr,
					"AUTH_SERVICE_CERT="+string(certificatePEM),
					"INSECURE_UNENCRYPTED_CONNECTION="+strconv.FormatBool(true),
				)

				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Error(string(out))
					t.Fatal(err)
				} else {
					t.Log(string(out))
				}
			})
		}
	})
}

func startMockAuthServiceTLS(t *testing.T, testCtx *testcontext.Context, cancelCtx context.Context, certificate tls.Certificate) (port int) {
	serverTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}

	// start a server with the certificate
	tcpListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Logf("listening on %s", tcpListener.Addr())
	port = tcpListener.Addr().(*net.TCPAddr).Port

	drpcListener := tls.NewListener(tcpListener, serverTLSConfig)

	mux := drpcmux.New()
	err = pb.DRPCRegisterEdgeAuth(mux, &dRPCServerMock{})
	require.NoError(t, err)

	server := drpcserver.New(mux)
	testCtx.Go(func() error {
		return server.Serve(cancelCtx, drpcListener)
	})

	return port
}

func startMockAuthServiceUnencrypted(t *testing.T, testCtx *testcontext.Context, cancelCtx context.Context) (port int) {
	tcpListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Logf("listening on %s", tcpListener.Addr())
	port = tcpListener.Addr().(*net.TCPAddr).Port

	mux := drpcmux.New()
	err = pb.DRPCRegisterEdgeAuth(mux, &dRPCServerMock{})
	require.NoError(t, err)

	server := drpcserver.New(mux)
	testCtx.Go(func() error {
		return server.Serve(cancelCtx, tcpListener)
	})

	return port
}

func createSelfSignedCertificate(t *testing.T, hostname string) (certificatePEM []byte, privateKeyPEM []byte) {
	notAfter := time.Now().Add(3 * time.Minute)

	template := x509.Certificate{
		Subject: pkix.Name{
			CommonName: hostname,
		},
		DNSNames:              []string{hostname},
		SerialNumber:          big.NewInt(1337),
		BasicConstraintsValid: false,
		IsCA:                  true,
		NotAfter:              notAfter,
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	certificateDERBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		&privateKey.PublicKey,
		privateKey,
	)
	require.NoError(t, err)

	certificatePEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificateDERBytes})

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	require.NoError(t, err)
	privateKeyPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyBytes})

	return certificatePEM, privateKeyPEM
}
