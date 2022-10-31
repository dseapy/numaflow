package cmd

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/numaproj/numaflow/pkg/shared/logging"
	sharedtls "github.com/numaproj/numaflow/pkg/shared/tls"
	"github.com/numaproj/numaflow/server/routes"
)

var (
	rewritePathPrefixes = []string{
		"/namespaces",
	}
)

func Start(insecure bool, port int, namespaced bool, managedNamespace string, baseHRef string) {
	logger := logging.NewLogger().Named("server")

	if err := setBaseHRef("ui/build/index.html", baseHRef); err != nil {
		panic(err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.RedirectTrailingSlash = true
	router.Use(static.Serve("/", static.LocalFile("./ui/build", true)))
	if namespaced {
		router.Use(Namespace(managedNamespace))
	}
	routes.Routes(router)
	router.Use(UrlRewrite(router, baseHRef))
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	if insecure {
		logger.Infof("Starting server (TLS disabled) on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	} else {
		cert, err := sharedtls.GenerateX509KeyPair()
		if err != nil {
			panic(err)
		}
		server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{*cert}, MinVersion: tls.VersionTLS12}

		logger.Infof("Starting server on %s", server.Addr)
		if err := server.ListenAndServeTLS("", ""); err != nil {
			panic(err)
		}
	}
}

func needToRewrite(path string) bool {
	for _, p := range rewritePathPrefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func UrlRewrite(r *gin.Engine, baseHRef string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if needToRewrite(c.Request.URL.Path) { //TODO fix
			c.Request.URL.Path = baseHRef
			r.HandleContext(c)
		} else if baseHRef != "/" &&
			!strings.HasPrefix(c.Request.URL.Path, baseHRef) &&
			!strings.HasPrefix("/"+c.Request.URL.Path, baseHRef) {
			if strings.HasPrefix(c.Request.URL.Path, "/") {
				c.Request.URL.Path = baseHRef + strings.TrimPrefix(c.Request.URL.Path, "/")
				r.HandleContext(c)
			} else {
				c.Request.URL.Path = strings.TrimPrefix(baseHRef, "/") + c.Request.URL.Path
				r.HandleContext(c)
			}
		}
		c.Next()
	}
}

func Namespace(ns string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("namespace", ns)
		c.Next()
	}
}

func setBaseHRef(filename string, baseHRef string) error {
	if baseHRef == "/" {
		return nil
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	prevHRef := "<base href=\"/\">"
	newHRef := fmt.Sprintf("<base href=\"%s\">", baseHRef)
	file = bytes.Replace(file, []byte(prevHRef), []byte(newHRef), -1)

	err = os.WriteFile(filename, file, 0666)
	return err
}
