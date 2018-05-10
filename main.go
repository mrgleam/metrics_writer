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

var counter = make(map[string]float64)
var indexed = make(map[string]prometheus.Gauge)

func main() {
	router := gin.Default()
	router.GET("/metrics", gin.WrapH(prometheus.Handler()))
	router.GET("/", Index)
	router.POST("/counter/metrics", HandlerCounter)

	router.Run(":8080")
}

func restartRegister(indexed prometheus.Gauge) {
	prometheus.Unregister(indexed)
	prometheus.Register(indexed)
}

func GetHash(content string) string {
	h512 := sha512.New()
	io.WriteString(h512, content)
	hash := string(h512.Sum(nil))
	return hash
}
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Welcome!")
}

func Counter(name string, inputlabels string, clientIP string) (string, string) {
	keepuniqindexed := ""
	uniq := GetHash(name + inputlabels)
	var labels []string
	byt := []byte(inputlabels)
	var dat map[string]string
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	dat["instance"] = clientIP
	for key, value := range dat {
		log.Println("Key:", key, "Value:", value)
		labels = append(labels, key)
		keepuniqindexed = keepuniqindexed + key
	}
	indexeduniq := GetHash(name + keepuniqindexed)

	indexed[indexeduniq] = prometheus.NewGauge(prometheus.GaugeOpts{
		ConstLabels: dat,
		Namespace:   "my_company",
		Subsystem:   "indexer",
		Name:        name,
		Help:        "The number of documents indexed.",
	})
	restartRegister(indexed[indexeduniq])
	return uniq, indexeduniq
}

func HandlerCounter(c *gin.Context) {
	name := c.PostForm("name")
	inputlabels := c.PostForm("labels")
	uniq, indexeduniq := Counter(name, inputlabels, c.ClientIP())
	counter[uniq] = counter[uniq] + 1.00
	indexed[indexeduniq].Set(counter[uniq])

	c.JSON(http.StatusOK, gin.H{
		"status": "counter",
		"name":   name,
		"value":  counter[uniq],
	})
}
