package mgi18n

import (
	"fmt"
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/logs"
	"github.com/maczh/mgcache"
	"github.com/maczh/mgconfig"
	"github.com/maczh/mgerr"
	"github.com/maczh/mgi18n/xlang"
	"time"
)

var appName, defaultLanguage string

func Init() {
	appName = mgconfig.GetConfigString("go.xlang.appName")
	if appName == "" {
		appName = mgconfig.GetConfigString("go.application.name")
	}
	defaultLanguage = mgconfig.GetConfigString("go.xlang.default")
	if defaultLanguage == "" {
		defaultLanguage = "zh-cn"
	}
	version, err := xlang.GetAppXLangVersion(appName)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	initCache(version)
	//设置定时任务自动检查
	ticker := time.NewTicker(time.Minute * 5)
	go func() {
		for _ = range ticker.C {
			reflushXLangCache()
		}
	}()
}

func initCache(version string) {
	mgcache.OnGetCache("x-lang").Add("version", version, 0)
	langs, err := xlang.GetAppXLangStringsAll(appName)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	for k, v := range langs {
		mgcache.OnGetCache("x-lang").Add(k, v, 0)
	}
	mgcache.OnGetCache("x-lang").Add("success:zh-cn", "成功", 0)
	mgcache.OnGetCache("x-lang").Add("success:en-us", "Success", 0)
	mgcache.OnGetCache("x-lang").Add("success:zh-tw", "成功", 0)
	mgcache.OnGetCache("x-lang").Add("success:ja", "成功", 0)
	mgcache.OnGetCache("x-lang").Add("success:fr", "Succès", 0)
	mgcache.OnGetCache("x-lang").Add("success:it", "Successo", 0)
	mgcache.OnGetCache("x-lang").Add("success:de", "der Erfolg", 0)
	mgcache.OnGetCache("x-lang").Add("success:ko", "성공", 0)
}

func GetXLangString(stringId, lang string) string {
	key := fmt.Sprintf("%s:%s", stringId, lang)
	str, ok := mgcache.OnGetCache("x-lang").Value(key)
	if ok {
		return str.(string)
	}
	key = fmt.Sprintf("%s:%s", stringId, defaultLanguage)
	str, ok = mgcache.OnGetCache("x-lang").Value(key)
	if ok {
		return str.(string)
	}
	return ""
}

func reflushXLangCache() {
	version, err := xlang.GetAppXLangVersion(appName)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	oldVersion, ok := mgcache.OnGetCache("x-lang").Value("version")
	if ok {
		if oldVersion != version {
			mgcache.OnGetCache("x-lang").Clear()
			initCache(version)
		}
	}
}

func Error(code int, messageId string) mgresult.Result {
	return mgresult.Error(code, String(messageId))
}

func ErrorWithMsg(code int, messageId, msg string) mgresult.Result {
	return mgresult.Error(code, fmt.Sprintf("%s:%s", String(messageId), msg))
}

func Success(data interface{}) mgresult.Result {
	return mgresult.SuccessWithMsg(String("success"), data)
}

func SuccessWithPage(data interface{}, count, index, size, total int) mgresult.Result {
	return mgresult.Result{
		Status: 1,
		Msg:    String("success"),
		Data:   data,
		Page: &mgresult.ResultPage{
			Count: count,
			Index: index,
			Size:  size,
			Total: total,
		},
	}
}

func String(messageId string) string {
	lang := mgerr.GetCurrentLanguage()
	return GetXLangString(messageId, lang)
}
