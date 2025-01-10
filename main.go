package main

import (
	"context"
	"fmt"
	"github.com/bear-san/ipman/internal/server"
	ipmanpb "github.com/bear-san/ipman/pkg/grpc"
	"github.com/bear-san/ipman/pkg/ip_repo"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	ctx := context.Background()
	googleCredentialPath := os.Getenv("GOOGLE_CREDENTIAL_PATH")
	if len(googleCredentialPath) == 0 {
		// use default path
		googleCredentialPath = "/google-credential.json"
	}
	credentialString, err := os.ReadFile(googleCredentialPath)
	if err != nil {
		panic(err)
	}

	manageSheetName := os.Getenv("MANAGE_SHEET_NAME")
	if len(manageSheetName) == 0 {
		panic(fmt.Errorf("invalid manage sheet name: %s", manageSheetName))
	}
	spreadSheetID := os.Getenv("SPREADSHEET_ID")
	if len(spreadSheetID) == 0 {
		panic(fmt.Errorf("invalid spreadsheet id: %s", spreadSheetID))
	}

	ipRepo, err := ip_repo.NewRepo(ctx, string(credentialString), spreadSheetID, manageSheetName)
	if err != nil {
		panic(err)
	}

	ipmanServer := server.IPManServer{
		IPRepo: ipRepo,
	}
	grpcSrv := grpc.NewServer()

	ipmanpb.RegisterIPAddressServiceServer(grpcSrv, &ipmanServer)

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	if err := grpcSrv.Serve(listener); err != nil {
		panic(err)
	}
}
