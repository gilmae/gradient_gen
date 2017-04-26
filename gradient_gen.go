package main

import (
  "fmt"
  "image"
  "image/color"
  "image/draw"
  "image/jpeg"
  "github.com/gilmae/interpolation"
  "os"
  "math/rand"
  "time"
  "sort"
  "encoding/hex"
  "encoding/json"
  "flag"
)

func generate_swatch(filename string, drawLine bool, redInterpolant interpolation.MonotonicCubic, greenInterpolant interpolation.MonotonicCubic, blueInterpolant interpolation.MonotonicCubic) {
  bounds := image.Rect(0,0,300,50)
  b := image.NewNRGBA(bounds)
  draw.Draw(b, bounds, image.NewUniform(color.Black), image.ZP, draw.Src)

  line_colour := color.NRGBA{255,255,255, 255}

  for x:=0 ; x < 640 ; x++ {
    var point = float64(x) / 300.0
    var redpoint = redInterpolant(point)
    var greenpoint = greenInterpolant(point)
    var bluepoint = blueInterpolant(point)

    for y:= 0; y < 255 ; y++ {
      c := color.NRGBA{uint8(redpoint), uint8(greenpoint), uint8(bluepoint), 255}
      b.Set(x,y,c)
    }

    if (drawLine) {
      b.Set(x, 255-int(redpoint), line_colour)
      b.Set(x, 255-int(greenpoint), line_colour)
      b.Set(x, 255-int(bluepoint), line_colour)
    }
  }

  file, err := os.Create(filename)
  if err != nil {
    fmt.Println(err)
  }

  if err = jpeg.Encode(file,b, &jpeg.Options{jpeg.DefaultQuality}); err != nil {
    fmt.Println(err)
  }

  if err = file.Close();err != nil {
    fmt.Println(err)
  }
}

func main() {

  var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
  var filename string
  var showRgbLines bool

  flag.StringVar(&filename, "f", "", "Generate a swatch to the input filepath.")
  flag.BoolVar(&showRgbLines, "l", false, "Show the RGB curve in the swatch.")

  flag.Parse()

  var points = uint8(rng.Float64() * 3.0) + 4
  var xSequence = make([]float64, points)
  var redpoints = make([]float64, points)
  var greenpoints = make([]float64, points)
  var bluepoints = make([]float64, points)

  for i,_:= range xSequence {
    xSequence[i] = rng.Float64()
    redpoints[i] = rng.Float64()*255.0
    greenpoints[i] = rng.Float64()*255.0
    bluepoints[i] = rng.Float64()*255.0
  }

xSequence[0] = 0.0
xSequence[1] = 1.0
  sort.Float64s(xSequence)

  var redInterpolant interpolation.MonotonicCubic = interpolation.CreateMonotonicCubic(xSequence, redpoints)
  var greenInterpolant interpolation.MonotonicCubic = interpolation.CreateMonotonicCubic(xSequence, greenpoints)
  var blueInterpolant interpolation.MonotonicCubic = interpolation.CreateMonotonicCubic(xSequence, bluepoints)

  if (filename != "") {
    generate_swatch(filename, showRgbLines, redInterpolant, greenInterpolant, blueInterpolant)
  }

  var gradient  [][]string

  gradient = make([][]string, points)

  for i,_ := range gradient {
    var hexRGB = hex.EncodeToString([]byte {uint8(redpoints[i]), uint8(greenpoints[i],), uint8(bluepoints[i])})
    gradient[i] = []string {fmt.Sprintf("%.6f", xSequence[i]), hexRGB}
  }

  j, _ := json.Marshal(gradient)
  var jsonString = string(j[:])

  fmt.Println(jsonString)
  }
