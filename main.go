package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	v8 "rogchap.com/v8go"
)

func main() {
	iso, err := v8.NewIsolate()
	realmScript, err := ioutil.ReadFile("./lib/reactivesearch-realm.min.js") // the file is inside the local directory
	if err != nil {
		fmt.Println("Err", err)
	}

	realmCtx, err := v8.NewContext(iso)
	if err != nil {
		fmt.Println("Err initializing realm ctx", err)
	}

	tessaractCtx, err := v8.NewContext(iso)
	if err != nil {
		fmt.Println("Err initializing tessaract ctx", err)
	}

	r := gin.Default()

	ocrScript, err := ioutil.ReadFile("./js/tesseract.js")
	if err != nil {
		fmt.Println("Err loading tessaract js", err)
	}

	r.POST("/_ocr", func(c *gin.Context) {
		t, err := tessaractCtx.RunScript(string(ocrScript), "tesseract.js")
		fmt.Println(t, err)
		script := `
			console.log("hello world")
			console.log(Tessarct)
			console.log(Tessarct.recognize)
			var res = Tesseract.recognize(
				'https://tesseract.projectnaptha.com/img/eng_bw.png',
				'eng'
			).then(function(data) {
				return data;
			})

			console.log(res);
		`

		r, err := tessaractCtx.RunScript(script, "script.js")
		fmt.Println("=> r:", r, err)

		promiseValue, err := tessaractCtx.RunScript("res", "value.js")
		fmt.Println(promiseValue, err)
		// p, err := promiseValue.AsPromise()
		// if err != nil {
		// 	fmt.Println("=> error executing as promise:", err)
		// }
		// value := p.Result()
		// fmt.Println("=> promise", value)

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/_reactivesearch/validate", func(c *gin.Context) {
		realmCtx.RunScript(string(realmScript), "reactivesearch.js")

		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}

		var reqBody map[string]interface{}
		err = json.Unmarshal(jsonData, &reqBody)
		if err != nil {
			panic(err)
		}

		query := reqBody["query"]
		queryBytes, _ := json.Marshal(query)

		mongodb := reqBody["mongodb"].(map[string]interface{})
		fmt.Println(mongodb)

		script := fmt.Sprintf(`
			var ref = new reactivesearch.ReactiveSearch({
				client: {},
				database: "%v",
				collection: "%v"
			});
			var data = ref.translate(%s);
			var res = JSON.stringify(data);`, mongodb["db"], mongodb["collection"], string(queryBytes))

		realmCtx.RunScript(script, "main.js")
		val, _ := realmCtx.RunScript("res", "value.js")
		obj := val.String()

		var resBody map[string]interface{}
		json.Unmarshal([]byte(obj), &resBody)
		c.JSON(200, resBody)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
