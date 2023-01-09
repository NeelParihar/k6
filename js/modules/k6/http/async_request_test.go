package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsyncRequest(t *testing.T) {
	t.Parallel()
	tb, _, _, rt, _ := newRuntime(t)

	sr := func(input string) string {
		return tb.Replacer.Replace(
			`(async () => {
    ` + input + `
})()`)
	}
	t.Run("HTTPRequest", func(t *testing.T) {
		t.Run("EmptyBody", func(t *testing.T) {
			_, err := rt.RunString(sr(`
				var reqUrl = "HTTPBIN_URL/cookies"
				var res = http.get(reqUrl);
				var jar = new http.CookieJar();

				jar.set("HTTPBIN_URL/cookies", "key", "value");
				res = await http.request("GET", "HTTPBIN_URL/cookies", null, { cookies: { key2: "value2" }, jar: jar });

				if (res.json().key != "value") { throw new Error("wrong cookie value: " + res.json().key); }

				if (res.status != 200) { throw new Error("wrong status: " + res.status); }
				if (res.request["method"] !== "GET") { throw new Error("http request method was not \"GET\": " + JSON.stringify(res.request)) }
				if (res.request["body"].length != 0) { throw new Error("http request body was not null: " + JSON.stringify(res.request["body"])) }
				if (res.request["url"] != reqUrl) {
					throw new Error("wrong http request url: " + JSON.stringify(res.request))
				}
				if (res.request["cookies"]["key2"][0].name != "key2") { throw new Error("wrong http request cookies: " + JSON.stringify(JSON.stringify(res.request["cookies"]["key2"]))) }
				if (res.request["headers"]["User-Agent"][0] != "TestUserAgent") { throw new Error("wrong http request headers: " + JSON.stringify(res.request)) }
				`))
			assert.NoError(t, err)
		})
		t.Run("NonEmptyBody", func(t *testing.T) {
			_, err := rt.RunString(sr(`
				var res = await http.request("HTTPBIN_URL/post", {a: "a", b: 2}, {headers: {"Content-Type": "application/x-www-form-urlencoded; charset=utf-8"}});
				if (res.status != 200) { throw new Error("wrong status: " + res.status); }
				if (res.request["body"] != "a=a&b=2") { throw new Error("http request body was not set properly: " + JSON.stringify(res.request))}
				`))
			assert.NoError(t, err)
		})
	})
}
