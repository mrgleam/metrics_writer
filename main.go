package main

import (
	"crypto/sha512"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	size = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "my_company",
		Subsystem: "storage",
		Name:      "documents_total_size_bytes",
		Help:      "The total size of all documents in the storage.",
	})
)

var counter = make(map[string]float64)
var indexed = make(map[string]prometheus.Gauge)

func main() {
	router := gin.Default()
	router.GET("/metrics", gin.WrapH(prometheus.Handler()))
	router.GET("/", Index)
	router.POST("/counter/metrics", Counter)

	router.Run(":8080")
}

func register(indexed prometheus.Gauge) {
	prometheus.Unregister(indexed)
}

func Index(c *gin.Context) {
	c.String(http.StatusOK, "Welcome!")
}

func Counter(c *gin.Context) {
	name := c.PostForm("name")
	inputlabels := c.PostForm("labels")
	keepuniqindexed := ""
	h512 := sha512.New()
	io.WriteString(h512, name+inputlabels)
	uniq := string(h512.Sum(nil))
	h2512 := sha512.New()
	var labels []string
	byt := []byte(inputlabels)
	var dat map[string]string
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	dat["instance"] = c.ClientIP()
	for key, value := range dat {
		log.Println("Key:", key, "Value:", value)
		labels = append(labels, key)
		keepuniqindexed = keepuniqindexed + key
	}

	io.WriteString(h2512, name+keepuniqindexed)
	indexeduniq := string(h2512.Sum(nil))

	indexed[indexeduniq] = prometheus.NewGauge(prometheus.GaugeOpts{
		ConstLabels: dat,
		Namespace:   "my_company",
		Subsystem:   "indexer",
		Name:        name,
		Help:        "The number of documents indexed.",
	})

	prometheus.Unregister(indexed[indexeduniq])
	prometheus.Register(indexed[indexeduniq])
	counter[uniq] = counter[uniq] + 1.00
	indexed[indexeduniq].Set(counter[uniq])

	c.JSON(http.StatusOK, gin.H{
		"status": "counter",
		"name":   name,
		"value":  counter[uniq],
	})
}
