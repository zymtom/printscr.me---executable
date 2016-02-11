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
    //"net/url"
    //"io/ioutil"
    "encoding/json"
    "errors"
    "os/exec"
    "runtime"
    //"github.com/getlantern/systray"
    //"github.com/getlantern/systray/example/icon"
)

func main() {
    //systray.Run(onReady)
    f, e := makeScreenshot()
    if e != nil {
        panic(e)
    }
    res, err := Upload("http://printscr.me/api.php", "./"+f)
    if err != nil {
        panic(err)
    }
    if res["upload"] == "failed" {
        panic(res["reason"])
    }
    if str, ok := res["location"].(string); !ok {
        panic("something went wrong with getting the location")
    }else{
        var err error
        switch runtime.GOOS {
        case "linux":
            err = exec.Command("xdg-open", str).Start()
        case "windows", "darwin":
            err = exec.Command("cmd", "/c", "start", str).Start()
        default:
            err = fmt.Errorf("unsupported platform")
        }
        if err != nil {
            panic(err)
        }
    }
    
    
}
/*func onReady() {
    systray.SetIcon(icon.Data)
    systray.SetTitle("Awesome App")
    systray.SetTooltip("Pretty awesome")
    mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
    go func() {
		<-mQuit.ClickedCh
		systray.Quit()
		fmt.Println("Quit now...")
	}()
}*/
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
    var b bytes.Buffer
    w := multipart.NewWriter(&b)
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
    if fw, err = w.CreateFormField("key"); err != nil {
        return nil, err
    }
    if _, err = fw.Write([]byte("KEY")); err != nil {
        return nil, err
    }
    w.Close()
    req, err := http.NewRequest("POST", reqUrl, &b)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", w.FormDataContentType())
    //proxyUrl, err := url.Parse("http://127.0.0.1:8888")
    //client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
    client := &http.Client{}
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