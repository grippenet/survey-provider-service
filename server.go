package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/grippenet/survey-provider-service/surveys"
)

func main() {

	DIR := os.Getenv("SURVEY_DIR")

	if _, err := os.Stat(DIR); os.IsNotExist(err) {
		log.Fatalf("Directory '%s' doesnt exists", DIR)
	}

	surveyList := surveys.NewSurveyList(DIR)

	err := surveyList.Update()
	if err != nil {
		log.Fatalf("Unable to read dir '%s' : %s", DIR, err)
	}

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/list", func(c *gin.Context) {
		list := surveyList.GetList()
		c.JSON(http.StatusOK, list)
	})

	r.GET("/update", func(c *gin.Context) {
		surveyList.Update()
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	r.GET("/survey", func(c *gin.Context) {

		code := c.Query("id")

		file, err := surveyList.Get(code)

		if file == "" {
			c.JSON(404, gin.H{
				"msg":   "Unknown file",
				"error": err,
			})
			return
		}

		f, err := os.Open(file)

		if err != nil {
			c.JSON(404, gin.H{
				"msg":   "Unable to read file",
				"error": err,
			})
			return
		}

		fi, err := f.Stat()
		if err != nil {
			c.JSON(404, gin.H{
				"msg":   "Unable to read file",
				"error": err,
			})
			return
		}

		c.DataFromReader(http.StatusOK, fi.Size(), gin.MIMEJSON, f, nil)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
