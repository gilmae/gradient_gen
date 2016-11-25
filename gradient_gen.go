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
)

func main() {

  var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

  var points = uint8(rng.Float64() * 4.0) + 3
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

  sort.Float64s(xSequence)

  var redInterpolant interpolation.MonotonicCubic = interpolation.CreateMonotonicCubic(xSequence, redpoints)
  var greenInterpolant interpolation.MonotonicCubic = interpolation.CreateMonotonicCubic(xSequence, greenpoints)
  var blueInterpolant interpolation.MonotonicCubic = interpolation.CreateMonotonicCubic(xSequence, bluepoints)

  bounds := image.Rect(0,0,300,50)
  b := image.NewNRGBA(bounds)
  draw.Draw(b, bounds, image.NewUniform(color.Black), image.ZP, draw.Src)

  for x:=0 ; x < 300 ; x++ {
    for y:= 0; y < 50 ; y++ {
      var point = float64(x) / 300.0
      var redpoint = redInterpolant(point)
      var greenpoint = greenInterpolant(point)
      var bluepoint = blueInterpolant(point)

      c := color.NRGBA{uint8(redpoint), uint8(greenpoint), uint8(bluepoint), 255}

      b.Set(x,y,c)
    }
  }

  file, err := os.Create(os.Args[1])
  if err != nil {
    fmt.Println(err)
  }

  if err = jpeg.Encode(file,b, &jpeg.Options{jpeg.DefaultQuality}); err != nil {
    fmt.Println(err)
  }

  if err = file.Close();err != nil {
    fmt.Println(err)
  }


  fmt.Println(xSequence)
  fmt.Println(redpoints)
  fmt.Println(greenpoints)
  fmt.Println(bluepoints)
}
