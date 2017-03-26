package main

import (
	"github.com/spf13/viper"
	"os"
	"net/http"
	"fmt"
	"bytes"
	"io/ioutil"
	"time"
	"encoding/json"
)

func readJson() ([]json.RawMessage, error) {
	data := make([]json.RawMessage, 0)

	file, err := os.Open(viper.GetString("data.userSrc"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func main() {
	viper.Set("scim.baseUrl", "http://localhost:8080/v2")
	viper.Set("data.userSrc", "/Users/davidiamyou/Downloads/scim_user.json")

	data, err := readJson()
	if err != nil {
		panic(err)
	}

	for i, body := range data {
		client := &http.Client{}

		req, err :=  http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/Users", viper.GetString("scim.baseUrl")),
			bytes.NewBuffer(body),
		)
		if err != nil {
			fmt.Printf("\nSkipped row %d due to %s, data: %s\n\n", i+1, err.Error(), string(body))
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("\nResponse error for row %d: %s\n\n", i+1, err.Error())
			continue
		}

		if resp.StatusCode < 299 {
			fmt.Printf("code:%d\n", resp.StatusCode)
		} else {
			b, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("[x]code:%d, row:%d, body:%s\n", resp.StatusCode, i+1, string(b))
		}

		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Done!")
}