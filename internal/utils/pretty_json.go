package utils

import (
    "bytes"
    "encoding/json"
)

// Utility function nicely prints string data of JSON formatted contents.
// SOURCE: https://stackoverflow.com/a/36544455
func JsonPrettyPrint(in string) string {
    var out bytes.Buffer
    err := json.Indent(&out, []byte(in), "", "  ")
    if err != nil {
        return in
    }
    return out.String()
}
