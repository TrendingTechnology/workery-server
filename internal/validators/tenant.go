package validators

import (
	"encoding/json"
	"unicode/utf8"

	"github.com/over55/workery-server/internal/idos"
)

func ValidateTenantSaveFromRequest(dirtyData *idos.TenantIDO) (bool, string) {
	e := make(map[string]string)

	if dirtyData.SchemaName == "" {
		e["schema_name"] = "missing value"
	} else {
		if utf8.RuneCountInString(dirtyData.SchemaName) > 63 {
			e["schema_name"] = "character count over 63"
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
