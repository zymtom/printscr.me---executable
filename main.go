package main

import (
	"github.com/vova616/screenshot"
	"image/png"
	"os"
    "bytes"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "net/url"
    //"io/ioutil"
    "encoding/json"
    "errors"
    "github.com/getlantern/systray"
    "github.com/getlantern/systray/example/icon"
)

func main() {
    systray.Run(onReady)
    f, e := makeScreenshot()
    if e != nil {
        panic(e)
    }
    res, err := Upload("http://127.0.0.1:6969/projects/printscreen.me---website/api.php", "./"+f)
    if err != nil {
        panic(err)
    }
    if res["upload"] == "failed" {
        panic(res["reason"])
    }
    fmt.Print(res["location"])
}
func onReady() {
    systray.SetIcon(icon.Data)
    systray.SetTitle("Awesome App")
    systray.SetTooltip("Pretty awesome超级棒")
    mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
    go func() {
		<-mQuit.ClickedCh
		systray.Quit()
		fmt.Println("Quit now...")
	}()
}
func makeScreenshot()(file string, err error){
    
    file = "ss123.png"
    img, err := screenshot.CaptureScreen()
	if err != nil {
		return file, err
	}
    
	f, err := os.Create("./"+file)
	if err != nil {
		return file, err
	}
	err = png.Encode(f, img)
	if err != nil {
		return file, err
	}
	f.Close()
    return file, nil
}

func Upload(reqUrl, file string) (response map[string]interface{}, err error) {
    // Prepare a form that you will submit to that URL.
    var b bytes.Buffer
    w := multipart.NewWriter(&b)
    // Add your image file
    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    fw, err := w.CreateFormFile("image", file)
    if err != nil {
        return nil, err
    }
    if _, err = io.Copy(fw, f); err != nil {
        return nil, err
    }
    // Add the other fields
    if fw, err = w.CreateFormField("key"); err != nil {
        return nil, err
    }
    if _, err = fw.Write([]byte("KEY")); err != nil {
        return nil, err
    }
    // Don't forget to close the multipart writer.
    // If you don't close it, your request will be missing the terminating boundary.
    w.Close()

    // Now that you have a form, you can submit it to your handler.
    req, err := http.NewRequest("POST", reqUrl, &b)
    if err != nil {
        return nil, err
    }
    // Don't forget to set the content type, this will contain the boundary.
    req.Header.Set("Content-Type", w.FormDataContentType())

    proxyUrl, err := url.Parse("http://127.0.0.1:8888")
    client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
    res, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    if res.StatusCode != http.StatusOK {
        err = fmt.Errorf("bad status: %s", res.Status)
        return nil, err
    }
    /*contents, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }*/
    dec := json.NewDecoder(res.Body)
    if dec == nil {
        return nil, errors.New("error decoding json")
    }

    json_map := make(map[string]interface{})
    err = dec.Decode(&json_map)
    if err != nil {
        return nil, errors.New("error decoding json")
    }

    return json_map, nil
}