package ip_repo

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const SHEET_RANGE = "A1:G"
const (
	IP_ADDRESS_TYPE_LOCAL  = "local"
	IP_ADDRESS_TYPE_GLOBAL = "global"
)

type IPRepo struct {
	spreadSheetID   string
	manageSheetName string
	sheetService    *sheets.Service
}

type IPAddress struct {
	Address           string
	Subnet            string
	GatewayAddress    string
	AddressType       string
	Using             bool
	AutoAssignEnabled bool
	Description       string
}

func NewRepo(ctx context.Context, serviceAccountCredential string, spreadSheetID string, manageSheetName string) (*IPRepo, error) {
	credential, err := google.CredentialsFromJSON(
		ctx,
		[]byte(serviceAccountCredential),
		sheets.DriveScope,
		sheets.SpreadsheetsScope,
	)
	if err != nil {
		return nil, err
	}

	service, err := sheets.NewService(ctx, option.WithCredentials(credential))
	if err != nil {
		return nil, err
	}

	return &IPRepo{
		manageSheetName: manageSheetName,
		spreadSheetID:   spreadSheetID,
		sheetService:    service,
	}, nil
}

func (r IPRepo) GetSheetData() ([]IPAddress, error) {
	values, err := r.sheetService.Spreadsheets.Values.Get(r.spreadSheetID, fmt.Sprintf("'%s'!%s", r.manageSheetName, SHEET_RANGE)).Do()
	if err != nil {
		return nil, err
	}

	colOrder := make(map[string]int)
	headerRow := values.Values[0]
	for i, i2 := range headerRow {
		cellData, valid := i2.(string)
		if !valid {
			continue
		}

		colOrder[cellData] = i
	}

	addressList := make([]IPAddress, 0)
	for i := 1; i < len(values.Values); i++ {
		row := values.Values[i]

		ipAddressData := IPAddress{
			Address:           row[colOrder["IP"]].(string),
			Subnet:            row[colOrder["Subnet"]].(string),
			GatewayAddress:    row[colOrder["GW"]].(string),
			AddressType:       row[colOrder["種別"]].(string),
			Using:             row[colOrder["利用中"]].(string) == "TRUE",
			AutoAssignEnabled: row[colOrder["自動割当対象"]].(string) == "TRUE",
		}
		if len(row) == 7 {
			ipAddressData.Description = row[colOrder["用途"]].(string)
		}

		addressList = append(addressList, ipAddressData)
	}

	return addressList, nil
}
