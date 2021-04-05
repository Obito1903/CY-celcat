package celcat

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/Obito1903/CY-celcat/celcat/common"
	"github.com/Obito1903/CY-celcat/celcat/ics"
	"github.com/Obito1903/CY-celcat/celcat/svg"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func GetConfig(path string) common.Config {
	var result map[string]interface{}
	var config common.Config

	configFile, err := os.Open(path)
	common.CheckErr(err)
	configFileStream, err := ioutil.ReadAll(configFile)
	common.CheckErr(err)

	json.Unmarshal(configFileStream, &result)
	config.UserId, config.UserPassword = fmt.Sprintf("%s", result["userId"]), fmt.Sprintf("%s", result["userPassword"])
	return config
}

func ToICS(calendar common.Calendar) {
	icsCal := ics.CreatICS(calendar)
	ics.SaveICS(icsCal)
}

func ToSVG(calendar common.Calendar, path string) {
	f, err := os.Create(path)
	common.CheckErr(err)
	defer f.Close()
	fmt.Fprint(f, svg.CreateSVG(calendar))

}

func ToPNG(calendar common.Calendar, width int, height int) {
	calSVG, _ := oksvg.ReadIconStream(svg.CreateSVG(calendar))
	calSVG.SetTarget(0, 0, float64(width), float64(height))
	calPNG := image.NewRGBA(image.Rect(0, 0, width, height))
	calSVG.Draw(rasterx.NewDasher(width, height, rasterx.NewScannerGV(width, height, calPNG, calPNG.Bounds())), 1)

	out, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = png.Encode(out, calPNG)
	if err != nil {
		panic(err)
	}

}
