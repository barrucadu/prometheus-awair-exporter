package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	score = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_score",
			Help: "Awair score (%)",
		},
		[]string{"sensor"},
	)

	dewPoint = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_dew_point",
			Help: "Dew point (C)",
		},
		[]string{"sensor"},
	)

	temp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_temp",
			Help: "Temperature (C)",
		},
		[]string{"sensor"},
	)

	humid = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_humid",
			Help: "Relative humidity (%)",
		},
		[]string{"sensor"},
	)

	absHumid = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_abs_humid",
			Help: "Absolute humidity (g/m^3)",
		},
		[]string{"sensor"},
	)

	co2 = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_co2",
			Help: "CO2 level (ppm)",
		},
		[]string{"sensor"},
	)

	co2Est = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_co2_est",
			Help: "?",
		},
		[]string{"sensor"},
	)

	co2EstBaseline = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_co2_est_baseline",
			Help: "?",
		},
		[]string{"sensor"},
	)

	voc = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_voc",
			Help: "Volatile organic compounds (ppb)",
		},
		[]string{"sensor"},
	)

	vocBaseline = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_voc_baseline",
			Help: "?",
		},
		[]string{"sensor"},
	)

	vocH2Raw = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_voc_h2_raw",
			Help: "?",
		},
		[]string{"sensor"},
	)

	vocEthanolRaw = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_voc_ethanol_raw",
			Help: "?",
		},
		[]string{"sensor"},
	)

	pm25 = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_pm25",
			Help: "Particulate matter (ug/m^3)",
		},
		[]string{"sensor"},
	)

	pm10Est = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awair_pm10_est",
			Help: "?",
		},
		[]string{"sensor"},
	)
)

type SensorData struct {
	Timestamp      time.Time `json:"timestamp"`
	Score          int       `json:"score"`
	DewPoint       float64   `json:"dew_point"`
	Temp           float64   `json:"temp"`
	Humid          float64   `json:"humid"`
	AbsHumid       float64   `json:"abs_humid"`
	Co2            int       `json:"co2"`
	Co2Est         int       `json:"co2_est"`
	Co2EstBaseline int       `json:"co2_est_baseline"`
	Voc            int       `json:"voc"`
	VocBaseline    int       `json:"voc_baseline"`
	VocH2Raw       int       `json:"voc_h2_raw"`
	VocEthanolRaw  int       `json:"voc_ethanol_raw"`
	Pm25           int       `json:"pm25"`
	Pm10Est        int       `json:"pm10_est"`
}

func recordMetricsForSensor(client http.Client, name string, address string) {
	resp, err := client.Get("http://" + address + "/air-data/latest")
	if err != nil {
		log.Printf("[%v:%v] request failed: %v", name, address, err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[%v:%v] request failed: %v", name, address, err)
		return
	}

	var sensorData SensorData
	if err := json.Unmarshal(body, &sensorData); err != nil {
		log.Printf("[%v:%v] could not unmarshal json: %v", name, address, err)
		return
	}

	score.With(prometheus.Labels{"sensor": name}).Set(float64(sensorData.Score))
	dewPoint.With(prometheus.Labels{"sensor": name}).Set(sensorData.DewPoint)
	temp.With(prometheus.Labels{"sensor": name}).Set(sensorData.Temp)
	humid.With(prometheus.Labels{"sensor": name}).Set(sensorData.Humid)
	absHumid.With(prometheus.Labels{"sensor": name}).Set(sensorData.AbsHumid)
	co2.With(prometheus.Labels{"sensor": name}).Set(float64(sensorData.Co2))
	co2Est.With(prometheus.Labels{"sensor": name}).Set(float64(sensorData.Co2Est))
	co2EstBaseline.With(prometheus.Labels{"sensor": name}).Set(float64(sensorData.Co2EstBaseline))
	voc.With(prometheus.Labels{"sensor": name}).Set(float64(sensorData.Voc))
	vocBaseline.With(prometheus.Labels{"sensor": name}).Set(float64(sensorData.VocBaseline))
	vocH2Raw.With(prometheus.Labels{"sensor": name}).Set(float64(sensorData.VocH2Raw))
	vocEthanolRaw.With(prometheus.Labels{"sensor": name}).Set(float64(sensorData.VocEthanolRaw))
	pm25.With(prometheus.Labels{"sensor": name}).Set(float64(sensorData.Pm25))
	pm10Est.With(prometheus.Labels{"sensor": name}).Set(float64(sensorData.Pm10Est))
}

func recordMetricsLoop(sensors map[string]string, delay time.Duration) {
	client := http.Client{Timeout: delay}

	for {
		for name, address := range sensors {
			go recordMetricsForSensor(client, name, address)
		}

		time.Sleep(delay)
	}
}

func main() {
	var opts struct {
		Address string            `short:"a" long:"address" default:"127.0.0.1:8888" description:"Address to listen on"`
		Sensors map[string]string `short:"s" long:"sensor" required:"true" description:"Sensor names and IP addresses"`
		Delay   time.Duration     `short:"d" long:"delay" default:"5s" description:"Delay between attempts to refresh metrics"`
	}

	if _, err := flags.Parse(&opts); err != nil {
		// it only seems to return an error when `-h` /
		// `--help` is passed, and it already prints the help
		// text in that case, so there's no need to print the
		// message again.
		os.Exit(1)
	}

	go recordMetricsLoop(opts.Sensors, opts.Delay)

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("listening on %v", opts.Address)
	log.Fatal(http.ListenAndServe(opts.Address, nil))
}
