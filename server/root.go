package server

import (
	"net/http"
	"strconv"

	"github.com/chazapp/prometheus-whistleblower/collector"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(port int) error {
	router := gin.Default()
	router.LoadHTMLGlob("ui/templates/*")
	router.Static("/assets", "ui/assets")
	router.StaticFile("/favicon.ico", "ui/assets/logo.png")

	registry := prometheus.NewRegistry()
	coll := collector.NewWhistleblowerCollector()

	registry.MustRegister(&coll)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
		})
	})

	router.GET("/metrics", gin.WrapH(handler))
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.POST("/metric", func(c *gin.Context) {
		var metric collector.Metric

		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		coll.AddMetric(metric)
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.DELETE("/metric/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		err = coll.DeleteMetric(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	router.GET("/metrics/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"metrics": coll.Metrics,
		})
	})

	return router.Run("0.0.0.0:" + strconv.Itoa(port))
}
