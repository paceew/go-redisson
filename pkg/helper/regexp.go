package helper

import "regexp"

var (
	regexpStore = map[string]*regexp.Regexp{
		"qq": regexp.MustCompile(`[1-9]\d{4,10}`), "wechat": regexp.MustCompile(`[a-zA-Z][a-zA-Z\d_-]{5,19}`),
		//^(13[0-9]|14[5|7]|15[0|1|2|3|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\d{8}$
		//国际^((\d{3,4}-)|\d{3.4}-)?\d{7,8}$
		//国内带区号\d{3}-\d{8}|\d{4}-\d{7}
		"phone": regexp.MustCompile(`1[34578]\d{9}`),
		"email": regexp.MustCompile(`\w+@[a-z0-9A-Z]+\.[a-z]+`),

		"id":   regexp.MustCompile(`([1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx])|([1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3})`),
		"name": regexp.MustCompile("([\u4e00-\u9fa5]{2,20}|[a-zA-Z.\\s]{2,20})"),
	}
)

// RegExpCollect 正则收集，qq、wx、phone、email、id、name
func RegExpCollect(str, style string) (string, bool) {
	var result string
	var eff bool
	if Regexp, ok := regexpStore[style]; ok {
		result = Regexp.FindString(str)
	}

	if result != "" {
		eff = true
	}
	return result, eff
}
