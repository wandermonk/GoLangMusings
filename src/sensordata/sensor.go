package main

import (
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"
	"encoding/gob"
	"dto"
	"bytes"
	"qutils"
	"github.com/streadway/amqp"
)
var url = "amqp://guest:guest@localhost:15672"
var name = flag.String("name", "sensor", "name of sensor")
var freq = flag.Uint("freq", 5, "update frequency in cycles/sec")
var max = flag.Float64("max", 5, "maximum value for generated frequency")
var min = flag.Float64("min", 1, "minimum value for generated frequency")
var stepSize = flag.Float64("step", 0.1, "maximum allowable change per measurement")

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var value = r.Float64()*(*max-*min) + *min
var nom = (*max-*min)/2 + *min

buf:= new(bytes.Buffer)
enc:= gob.NewDecoder(buf)

func main() {
	flag.Parse()

	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	dataQueue := qutils.GetQueue(*name, ch)

	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")
	signal := time.Tick(dur)
	for range signal {
		calcValue()
		reading := dto.SensorMessage{
			Name: *name,
			Value: value,
			Time: time.Now(),
		}
		buf.Reset()
		enc.Encode(reading)

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		ch.Publish(
			"",
			dataueue.Name,
			false,
			false,
			msg)

		log.Printf("Reading the sent value : %v\n", value)
	}
}
func calcValue() {
	var maxStep, minStep float64
	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}
	value += r.Float64()*(maxStep-minStep) + minStep
}
