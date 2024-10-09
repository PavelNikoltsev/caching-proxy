package main

import (
	"bytes"
	"caching-proxy/cache"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cache.Init()
	port := flag.Int("port", 8080, "Port to listen on")
	origin := flag.String("origin", "*", "Origin for redirects")
	autoClear := flag.Int("auto-clear-cache", 30, "Cache clear interval in minutes (default 30 minutes)")
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println("Error: No flags provided")
		flag.Usage()
		os.Exit(1)
	}
	cache.Cache.StartAutoClear(time.Duration(*autoClear) * time.Minute)
	serve(origin, port)
}

func serve(origin *string, port *int) {
	r := gin.Default()

	proxy, targetURL, err := createProxy(origin)
	if err != nil {
		panic(err)
	}
	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/clear-cache", func(c *gin.Context) {
			cache.Cache.Clear()
			c.JSON(http.StatusOK, gin.H{"message": "Cache cleared"})
		})
	}
	proxyGroup := r.Group("/proxy")
	{
		proxyGroup.Any("/*path", proxyHandler(proxy, targetURL))
	}

	if err := r.Run(fmt.Sprintf(":%d", *port)); err != nil {
		panic(err)
	}
}

func createProxy(origin *string) (*httputil.ReverseProxy, *url.URL, error) {
	targetURL, err := url.Parse(*origin)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid origin URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.ModifyResponse = func(response *http.Response) error {
		return handleCache(response)
	}

	return proxy, targetURL, nil
}

func handleCache(response *http.Response) error {
	key := fmt.Sprintf("%s-%s?%s", response.Request.Method, response.Request.URL.Path, response.Request.URL.RawQuery)

	if value, ok := cache.Cache.Get(key); ok {
		response.Body = io.NopCloser(bytes.NewBuffer(value))
		response.Header.Add("X-Cache", "HIT")
		return nil
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	cache.Cache.Set(key, bodyBytes)
	response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	response.Header.Add("X-Cache", "MISS")

	return nil
}

func proxyHandler(proxy *httputil.ReverseProxy, targetURL *url.URL) gin.HandlerFunc {
	return func(c *gin.Context) {
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = targetURL.Host
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.URL.Path = c.Param("path")
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
