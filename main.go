package main

import (
	"context"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"github.com/quhar/bme280"
	"golang.org/x/exp/io/i2c"
	"google.golang.org/api/option"
)

//センサーからデータを取得
func GetTPH() (temperature float64, pressure float64, humidity float64, err error) {

	d, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, bme280.I2CAddr)
	if err != nil {
		log.Fatal(err)
	}

	b := bme280.New(d)
	err = b.Init()

	temperature, pressure, humidity, err = b.EnvData()

	return
}

type RoomData struct {
	temperature float64
	pressure    float64
	humidity    float64
}

func recordData(roomData *RoomData) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("path of your json file")
	app, err := firebase.NewApp(ctx, nil, opt)
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = client.Collection("conditions").Add(ctx, map[string]interface{}{
		"createdAt:":  time.Now(),
		"temperature": roomData.temperature,
		"pressure":    roomData.pressure,
		"humidity":    roomData.humidity,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
}

func main() {
	log.Println("Start")
	t, p, h, _ := GetTPH()
	log.Printf("Temperature: %f C\n", t)
	log.Printf("Pressure: %f hPa\n", p)
	log.Printf("Humidity: %f %%rh\n", h)

	data := RoomData{
		temperature: t,
		pressure:    p,
		humidity:    h,
	}

	recordData(&data)
	log.Println("Finished")
}
