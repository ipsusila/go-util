package util

import (
	"encoding/json"

	log "github.com/ipsusila/slog"
	"github.com/tidwall/pretty"
)

// Pretty convert value to pretty json
func Pretty(v interface{}) []byte {
	js, err := json.Marshal(v)
	if err != nil {
		log.Warnw("marshal json error", "error", err.Error())
	}
	return pretty.Pretty(js)
}

// PrettyStr convert value to pretty json string
func PrettyStr(v interface{}) string {
	return string(Pretty(v))
}

// PrettyColor convert value to colored json
func PrettyColor(v interface{}) []byte {
	jsFmt := Pretty(v)
	return pretty.Color(jsFmt, nil)
}

// PrettyColorStr convert type to JSON with color
func PrettyColorStr(v interface{}) string {
	return string(PrettyColor(v))
}
