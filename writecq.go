package main

import (
  "compress/zlib"
  "flag"
  "fmt"
  "io"
  "os"
  "runtime"
  "time"

  "github.com/padster/go-sound/cq"
  f "github.com/padster/go-sound/file"
  s "github.com/padster/go-sound/sounds"
)

// Runs CQ to generate the spectrogram (without interpolation) and writes to file.
func main() {
  runtime.GOMAXPROCS(4)

  // Parse flags...
  sampleRate := s.CyclesPerSecond
  octaves := 7
  minFreq := flag.Float64("minFreq", 55.0, "minimum frequency")
  maxFreq := flag.Float64("maxFreq", 55.0*float64(cq.UnsafeShift(octaves)), "maximum frequency")
  bpo := flag.Int("bpo", 24, "Buckets per octave")
  zip := flag.Bool("zip", false, "Whether to zip the output")
  flag.Parse()

  remainingArgs := flag.Args()
  if len(remainingArgs) < 1 || len(remainingArgs) > 2 {
    panic("Required: <input> [<input>] filename arguments")
  }
  inputFile := remainingArgs[0]
  outputFile := "out.raw"
  if len(remainingArgs) == 2 {
    outputFile = remainingArgs[1]
  }

  inputSound := f.Read(inputFile)
  inputSound.Start()
  defer inputSound.Stop()

  // minFreq, maxFreq, bpo := 110.0, 14080.0, 24
  params := cq.NewCQParams(sampleRate, *minFreq, *maxFreq, *bpo)
  constantQ := cq.NewConstantQ(params)

  startTime := time.Now()
  columns := constantQ.ProcessChannel(inputSound.GetSamples())
  writeSamples(outputFile, *zip, constantQ.OutputLatency, columns)
  elapsedSeconds := time.Since(startTime).Seconds()

  fmt.Printf("elapsed time (not counting init): %f sec\n", elapsedSeconds)
}

func writeSamples(outputFile string, compress bool, latency int, samples <-chan []complex128) {
  // BIG HACK
  latency = 0

  file, err := os.Create(outputFile)
  if err != nil {
    panic(err)
  }
  defer file.Close()

  framesWritten := 0
  maxHeight := 0

  var writer io.Writer
  if compress {
    zip := zlib.NewWriter(file) 
    defer zip.Close()
    writer = zip
  } else {
    writer = file
  }

  fmt.Printf("Latency = %d\n", latency)

  for sample := range samples {
    if latency > 0 {
      latency--
    } else {
      if len(sample) > maxHeight {
        maxHeight = len(sample)
      }
      cq.WriteComplexArray(writer, sample)
      framesWritten++
      if framesWritten % 10000 == 0 {
        fmt.Printf("Written frame %d\n", framesWritten)
      }
    }
  }
  fmt.Printf("Result: %d by %d\n", framesWritten, maxHeight)
}