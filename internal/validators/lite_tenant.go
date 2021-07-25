package validators

import (
	"encoding/json"
	"unicode/utf8"

	"github.com/over55/workery-server/internal/models"
)

func ValidateTenantSaveFromRequest(dirtyData *models.Tenant) (bool, string) {
	e := make(map[string]string)

	if dirtyData.Schema == "" {
		e["schema"] = "missing value"
	} else {
		if utf8.RuneCountInString(dirtyData.Schema) > 63 {
			e["schema"] = "character count over 63"
		}
	}
	if dirtyData.Name == "" {
		e["name"] = "missing value"
	} else {
		if utf8.RuneCountInString(dirtyData.Name) > 127 {
			e["name"] = "character count over 127"
		}
	}
	if dirtyData.State == 0 {
		e["state"] = "missing value"
	}

	if len(e) != 0 {
		b, err := json.Marshal(e)
		if err != nil { // Defensive code
			return false, err.Error()
		}
		return false, string(b)
	}
	return true, ""
}
