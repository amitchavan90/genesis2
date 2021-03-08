package genesis

import (
	"bytes"
	"fmt"
	"genesis/config"
	"genesis/graphql"
	"net/http"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gofrs/uuid"
)

// DownloadSpreadSheetOrder generates and returns a spreadsheet of products for an order
func DownloadSpreadSheetOrder(OrderStore OrderStorer, Auther AuthProvider, API *config.API) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Auth check
		user, err := Auther.UserFromContext(r.Context())
		if err != nil {
			restWriteError(r.Context(), w, http.StatusForbidden, err, "There was a problem reading your credentials. Please sign in and try again.")
			return
		}

		hasPerm := false
		for _, perm := range user.R.Role.Permissions {
			if perm == graphql.PermOrderRead.String() {
				hasPerm = true
				break
			}
		}
		if !hasPerm {
			restWriteError(r.Context(), w, http.StatusForbidden, err, "You are not authorized to do this action.")
			return
		}

		// Get order
		id := r.URL.Query().Get("id")
		if id == "" {
			restWriteError(r.Context(), w, http.StatusBadRequest, err, "no id provided")
			return
		}

		orderUUID, err := uuid.FromString(id)
		if err != nil {
			restWriteError(r.Context(), w, http.StatusBadRequest, err, "invalid id provided")
			return
		}

		order, err := OrderStore.Get(orderUUID)
		if err != nil {
			restWriteError(r.Context(), w, http.StatusInternalServerError, err, "failed to get order")
			return
		}

		// Get products
		products, err := OrderStore.Products(order)
		if err != nil {
			restWriteError(r.Context(), w, http.StatusInternalServerError, err, "failed to get products")
			return
		}

		// Create spreadsheet
		fileName := fmt.Sprintf("genesis_order_%s.xlsx", order.Code)
		sheet := "Sheet1"
		f := excelize.NewFile()

		// - Headers
		f.SetCellValue(sheet, "A1", "Code")
		f.SetCellValue(sheet, "B1", "SKU Code")
		f.SetCellValue(sheet, "C1", "SKU Name")
		f.SetCellValue(sheet, "D1", "View URL (QR on back of package)")
		f.SetCellValue(sheet, "E1", "Register URL (QR hidden under product)")

		// - Styles
		headerStyle, err := f.NewStyle(`{"fill":{"type":"pattern","color":["#808080"],"pattern":1}, "font": {"color": "#f7f7f7", "bold": true, "size": 14 } }`)
		if err != nil {
			restWriteError(r.Context(), w, http.StatusInternalServerError, err, "failed create spreadsheet style")
			return
		}

		f.SetColWidth(sheet, "B", "C", 15)
		f.SetColWidth(sheet, "D", "E", 100)
		f.SetRowHeight(sheet, 1, 30)

		f.SetCellStyle(sheet, "A1", "E1", headerStyle)

		// - Data
		for index, product := range products {
			i := index + 2
			if product == nil {
				f.SetCellValue(sheet, fmt.Sprintf("A%d", i), "null")
				continue
			}
			f.SetCellValue(sheet, fmt.Sprintf("A%d", i), product.Code)
			if product.R.Sku != nil {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", i), product.R.Sku.Code)
				f.SetCellValue(sheet, fmt.Sprintf("C%d", i), product.R.Sku.Name)
			}
			f.SetCellValue(sheet, fmt.Sprintf("D%d", i), fmt.Sprintf("%s/api/steak/view?productID=%s", API.ConsumerHost, product.ID))
			f.SetCellValue(sheet, fmt.Sprintf("E%d", i), fmt.Sprintf("%s/api/steak/detail?steakID=%s", API.ConsumerHost, product.RegisterID))
		}

		// Download spreadsheet
		w.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
		buf, err := f.WriteToBuffer()
		if err != nil {
			restWriteError(r.Context(), w, http.StatusInternalServerError, err, "Failed to write spreadsheet to buffer")
			return
		}
		rdr := bytes.NewReader(buf.Bytes())
		http.ServeContent(w, r, fileName, time.Now(), rdr)

	}
	return fn
}

const (
	spreadSheetTypeCarton    = "carton"
	spreadSheetTypePallet    = "pallet"
	spreadSheetTypeContainer = "container"
)

// GeneralItem for spreadsheets
type GeneralItem struct {
	ID   string
	Code string
}

// DownloadSpreadSheet generates and returns a spreadsheet of cartons, pallets or containers
func DownloadSpreadSheet(
	CartonStore CartonStorer,
	PalletStore PalletStorer,
	ContainerStore ContainerStorer,
	Auther AuthProvider,
	API *config.API,
) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Get arguments
		itemType := r.URL.Query().Get("type")
		if itemType == "" || !(itemType == spreadSheetTypeCarton || itemType == spreadSheetTypePallet || itemType == spreadSheetTypeContainer) {
			restWriteError(r.Context(), w, http.StatusBadRequest, nil, "invalid arguments")
			return
		}
		from := r.URL.Query().Get("from")
		if from == "" {
			restWriteError(r.Context(), w, http.StatusBadRequest, nil, "invalid arguments")
			return
		}
		to := r.URL.Query().Get("to")
		if to == "" {
			restWriteError(r.Context(), w, http.StatusBadRequest, nil, "invalid arguments")
			return
		}

		// Auth check
		user, err := Auther.UserFromContext(r.Context())
		if err != nil {
			restWriteError(r.Context(), w, http.StatusForbidden, err, "There was a problem reading your credentials. Please sign in and try again.")
			return
		}

		hasPerm := false
		for _, perm := range user.R.Role.Permissions {
			if perm == graphql.PermCartonCreate.String() {
				hasPerm = true
				break
			}
		}
		if !hasPerm {
			restWriteError(r.Context(), w, http.StatusForbidden, err, "You are not authorized to do this action.")
			return
		}

		// Get Cartons
		items := []*GeneralItem{}

		switch itemType {
		case spreadSheetTypeCarton:
			cartons, err := CartonStore.GetRange(from, to)
			if err != nil {
				restWriteError(r.Context(), w, http.StatusBadRequest, err, "failed to get cartons")
			}
			for _, carton := range cartons {
				items = append(items, &GeneralItem{ID: carton.ID, Code: carton.Code})
			}
		case spreadSheetTypePallet:
			pallets, err := PalletStore.GetRange(from, to)
			if err != nil {
				restWriteError(r.Context(), w, http.StatusBadRequest, err, "failed to get pallets")
			}
			for _, pallet := range pallets {
				items = append(items, &GeneralItem{ID: pallet.ID, Code: pallet.Code})
			}
		case spreadSheetTypeContainer:
			containers, err := ContainerStore.GetRange(from, to)
			if err != nil {
				restWriteError(r.Context(), w, http.StatusBadRequest, err, "failed to get containers")
			}
			for _, container := range containers {
				items = append(items, &GeneralItem{ID: container.ID, Code: container.Code})
			}
		}

		// Create spreadsheet
		fileName := fmt.Sprintf("genesis_%s_%s_to_%s.xlsx", itemType, from, to)
		sheet := "Sheet1"
		f := excelize.NewFile()

		// - Headers
		f.SetCellValue(sheet, "A1", "Code")
		f.SetCellValue(sheet, "B1", "UUID")
		f.SetCellValue(sheet, "C1", "URL (QR)")

		// - Styles
		headerStyle, err := f.NewStyle(`{"fill":{"type":"pattern","color":["#808080"],"pattern":1}, "font": {"color": "#f7f7f7", "bold": true, "size": 14 } }`)
		if err != nil {
			restWriteError(r.Context(), w, http.StatusInternalServerError, err, "failed create spreadsheet style")
			return
		}

		f.SetColWidth(sheet, "A", "A", 15)
		f.SetColWidth(sheet, "B", "B", 40)
		f.SetColWidth(sheet, "C", "C", 100)
		f.SetRowHeight(sheet, 1, 30)

		f.SetCellStyle(sheet, "A1", "C1", headerStyle)

		// - Data
		for index, item := range items {
			i := index + 2
			if item == nil {
				f.SetCellValue(sheet, fmt.Sprintf("A%d", i), "null")
				continue
			}
			f.SetCellValue(sheet, fmt.Sprintf("A%d", i), item.Code)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", i), item.ID)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", i), fmt.Sprintf("%s/api/q/%s", API.AdminHost, item.ID))
		}

		// Download spreadsheet
		w.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
		buf, err := f.WriteToBuffer()
		if err != nil {
			restWriteError(r.Context(), w, http.StatusInternalServerError, err, "Failed to write spreadsheet to buffer")
			return
		}
		rdr := bytes.NewReader(buf.Bytes())
		http.ServeContent(w, r, fileName, time.Now(), rdr)

	}
	return fn
}
