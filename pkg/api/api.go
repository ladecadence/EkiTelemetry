package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/ladecadence/EkiTelemetry/pkg/telemetry"
)

const (
	uploadPathData = "api/newdata"
	uploadPathImg  = "api/imgupload"
)

type API struct {
	Server   string
	User     string
	Password string
}

func (a *API) DataUpload(t telemetry.Telemetry) error {
	// prepare data
	json, err := json.Marshal(&t)
	if err != nil {
		return err
	}
	body := new(bytes.Buffer)
	body.Write(json)

	// post
	post, err := http.NewRequest("POST", a.Server+"/"+uploadPathData, body)
	if err != nil {
		return err
	}
	post.SetBasicAuth(a.User, a.Password)

	// make request
	client := &http.Client{}
	resp, err := client.Do(post)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Problem with request: %d", resp.StatusCode))
	}

	return nil
}

func (a *API) ImageUpload(imgPath string, mission string) error {
	// get image

	fmt.Printf("Upload image: %s\n", imgPath)

	// prepare data
	file, err := os.Open(imgPath)
	if err != nil {
		return err
	}
	fileContents, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	file.Close()

	// file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return err
	}
	part.Write(fileContents)
	part, err = writer.CreateFormField("mission")
	part.Write([]byte(mission))

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	writer.Close()

	// post
	post, err := http.NewRequest("POST", a.Server+"/"+uploadPathImg, body)
	if err != nil {
		return err
	}
	post.Header.Add("Content-Type", writer.FormDataContentType())
	post.SetBasicAuth(a.User, a.Password)

	// make request
	client := &http.Client{}
	resp, err := client.Do(post)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Problem with request: %d", resp.StatusCode))
	}

	return nil
}
