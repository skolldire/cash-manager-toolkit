package file_util

import (
	"github.com/skolldire/cash-manager-toolkit/pkg/kit"
	"net/http"
)

var (
	WRE = kit.StatusCode{Code: "WRE-500", Msg: "write %s file error", HttpCode: http.StatusInternalServerError}
	RRE = kit.StatusCode{Code: "RRE-500", Msg: "read %s file error", HttpCode: http.StatusInternalServerError}
	CFE = kit.StatusCode{Code: "CFE-500", Msg: "create %s file error", HttpCode: http.StatusInternalServerError}
)
