package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/dragonflyoss/Dragonfly/dfdaemon/config"
	logger "github.com/sirupsen/logrus"
)

func newReverseProxy(mirror *config.RegistryMirror) *httputil.ReverseProxy {
	reverseProxy := httputil.NewSingleHostReverseProxy(mirror.Remote.URL)
	reverseProxy.Director = newDynamicDirector(mirror.Remote.URL)

	return reverseProxy
}

func newDynamicDirector(remote *url.URL) func(*http.Request) {
	director := func(req *http.Request) {
		var target = remote
		targetQuery := target.RawQuery
		reg := req.Header.Get("X-Dragonfly-Registry")
		if len(reg) > 0 {
			logger.Debugf("Replacing URL host with X-Dragonfly-Registry: %s", reg)
			target.Host = reg
		}
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return director
}

// singleJoiningSlash is from net/http/httputil/reverseproxy.go
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// joinURLPath is from net/http/httputil/reverseproxy.go
func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}
