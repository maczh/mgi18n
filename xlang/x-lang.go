package xlang

import (
	"errors"
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/logs"
	"github.com/maczh/mgcall"
	"github.com/maczh/utils"
)

const (
	SERVICE_X_LANGUAGE          = "x-lang"
	URI_LIST_STRINGS_BY_APPNAME = "/lang/string/list"
	URI_GET_APP_STRINGS_VERSION = "/lang/string/app/version"
)

func mgCall(service, uri string, params map[string]string) mgresult.Result {
	res, err := mgcall.Call(service, uri, params)
	if err != nil {
		logs.Error("微服务{}{}调用异常:{}", service, uri, err.Error())
		return mgresult.Error(-1, "xlang service unavailable")
	}
	var result mgresult.Result
	utils.FromJSON(res, &result)
	return result
}

func GetAppXLangVersion(appName string) (string, error) {
	params := map[string]string{
		"appName": appName,
	}
	rs := mgCall(SERVICE_X_LANGUAGE, URI_GET_APP_STRINGS_VERSION, params)
	if rs.Status != 1 {
		return "", errors.New(rs.Msg)
	}
	data := make(map[string]string)
	utils.FromJSON(utils.ToJSON(rs.Data), &data)
	return data["version"], nil
}

func GetAppXLangStringsAll(appName string) (map[string]string, error) {
	params := map[string]string{
		"appName": appName,
	}
	rs := mgCall(SERVICE_X_LANGUAGE, URI_LIST_STRINGS_BY_APPNAME, params)
	if rs.Status != 1 {
		return nil, errors.New(rs.Msg)
	}
	result := make(map[string]string)
	utils.FromJSON(utils.ToJSON(rs.Data), &result)
	return result, nil
}
