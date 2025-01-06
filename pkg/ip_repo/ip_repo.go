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
	IP_ADDRESS_TYPE_LOCAL  = "Local"
	IP_ADDRESS_TYPE_GLOBAL = "Global"
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

func (r IPRepo) ReadFromSheet() (map[string]int, [][]interface{}, error) {
	values, err := r.sheetService.Spreadsheets.Values.Get(r.spreadSheetID, fmt.Sprintf("'%s'!%s", r.manageSheetName, SHEET_RANGE)).Do()
	if err != nil {
		return nil, nil, err
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

	dataList := make([][]interface{}, 0)
	for i := 1; i < len(values.Values); i++ {
		dataList = append(dataList, values.Values[i])
	}

	return colOrder, dataList, nil
}

func (r IPRepo) WriteToSheet(ipAddress IPAddress) error {
	headers, rows, err := r.ReadFromSheet()
	if err != nil {
		return err
	}

	ipAddresses, err := r.ParseRows(headers, rows)
	if err != nil {
		return err
	}

	rowIndex := -1
	for i, state := range ipAddresses {
		if state.Address != ipAddress.Address || state.Subnet != ipAddress.Subnet || state.GatewayAddress != ipAddress.GatewayAddress {
			continue
		}

		rowIndex = i
		break
	}

	if rowIndex < 0 {
		return fmt.Errorf("invalid IP Address %s(%s)", ipAddress.Address, ipAddress.AddressType)
	}

	newRow := make([]interface{}, 0)
	for s, i := range headers {
		if len(newRow) <= i {
			newRow = append(newRow, make([]interface{}, i-len(newRow)+1)...)
		}

		switch s {
		case "IP":
			newRow[i] = ipAddress.Address
		case "Subnet":
			newRow[i] = ipAddress.Subnet
		case "GW":
			newRow[i] = ipAddress.GatewayAddress
		case "種別":
			newRow[i] = ipAddress.AddressType
		case "利用中":
			if ipAddress.Using {
				newRow[i] = "TRUE"
			} else {
				newRow[i] = "FALSE"
			}
		case "自動割当対象":
			if ipAddress.AutoAssignEnabled {
				newRow[i] = "TRUE"
			} else {
				newRow[i] = "FALSE"
			}
		case "用途":
			newRow[i] = ipAddress.Description
		default:
			break
		}
	}

	v := &sheets.ValueRange{
		Values: [][]interface{}{
			newRow,
		},
	}
	if _, err := r.sheetService.Spreadsheets.Values.Update(
		r.spreadSheetID,
		fmt.Sprintf("'%s'!A%[2]d:%[2]d", r.manageSheetName, rowIndex+3),
		v,
	).ValueInputOption("USER_ENTERED").Do(); err != nil {
		return err
	}

	return nil
}

func (r IPRepo) ParseRows(headers map[string]int, rows [][]interface{}) ([]IPAddress, error) {
	addressList := make([]IPAddress, 0)
	for i := 1; i < len(rows); i++ {
		row := rows[i]

		ipAddressData := IPAddress{
			Address:           row[headers["IP"]].(string),
			Subnet:            row[headers["Subnet"]].(string),
			GatewayAddress:    row[headers["GW"]].(string),
			AddressType:       row[headers["種別"]].(string),
			Using:             row[headers["利用中"]].(string) == "TRUE",
			AutoAssignEnabled: row[headers["自動割当対象"]].(string) == "TRUE",
		}
		if len(row) == 7 {
			ipAddressData.Description = row[headers["用途"]].(string)
		}

		addressList = append(addressList, ipAddressData)
	}

	return addressList, nil
}

func (r IPRepo) GetAddresses() ([]IPAddress, error) {
	headers, dataRows, err := r.ReadFromSheet()
	if err != nil {
		return nil, err
	}

	return r.ParseRows(headers, dataRows)
}
