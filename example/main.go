package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/skoef/gop1"
)

var (
	powerConsumed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "instantaneous_power_consumed",
		Help: "Instantaneous power consumed per phase in W",
	}, []string{"phase"})
	powerGenerated = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "instantaneous_power_generated",
		Help: "Instantaneous power generated per phase in W",
	}, []string{"phase"})
	currentConsumed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "instantaneous_current",
		Help: "Instantaneous current per phase in A",
	}, []string{"phase"})
	voltageConsumed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "instantaneous_voltage",
		Help: "Instantaneous voltage per phase in V",
	}, []string{"phase"})
	tariffIndicator = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tariff_indicator",
		Help: "Tariff indicator electricity",
	})
	electricityConsumed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "electricity_consumed",
		Help: "Electricity consumed per tariff in kWh",
	}, []string{"tariff"})
	electricityGenerated = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "electricity_generated",
		Help: "Electricity generated per tariff in kWh",
	}, []string{"tariff"})
	gasConsumed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gas_consumed",
		Help: "Gas consumed in m3",
	})
)

func floatValue(input string) (fval float64) {
	fval, _ = strconv.ParseFloat(input, 64)
	return
}

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the info severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {

	var deviceName, metricsPath string
	var metricsPort int

	flag.StringVar(&deviceName, "device", "/dev/ttyUSB0", "Serial device towards P1 port")
	flag.IntVar(&metricsPort, "metrics-port", 2112, "Prometheus metrics port")
	flag.StringVar(&metricsPath, "metrics-path", "/metrics", "Prometheus metrics path")

	flag.Parse()

	if deviceName == "" {
		flag.Usage()
		os.Exit(2)
	}

	// set up prometheus metrics
	http.Handle(metricsPath, promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", metricsPort), nil)

	// open connection to serial port
	p1, err := gop1.New(gop1.P1Config{
		USBDevice: deviceName,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err.Error(),
			"device": deviceName,
		}).Error("failed to open serial device")
		return
	}

	// start reading from P1 port
	p1.Start()

	for {

		select {
		case tgram := <-p1.Incoming:
			for _, obj := range tgram.Objects {
				switch obj.Type {

				case gop1.OBISTypeInstantaneousPowerDeliveredL1:
					powerConsumed.With(prometheus.Labels{"phase": "l1"}).Set((floatValue(obj.Values[0].Value) * 1000))
				case gop1.OBISTypeInstantaneousPowerDeliveredL2:
					powerConsumed.With(prometheus.Labels{"phase": "l2"}).Set((floatValue(obj.Values[0].Value) * 1000))
				case gop1.OBISTypeInstantaneousPowerDeliveredL3:
					powerConsumed.With(prometheus.Labels{"phase": "l3"}).Set((floatValue(obj.Values[0].Value) * 1000))

				case gop1.OBISTypeInstantaneousPowerGeneratedL1:
					powerGenerated.With(prometheus.Labels{"phase": "l1"}).Set((floatValue(obj.Values[0].Value) * 1000))
				case gop1.OBISTypeInstantaneousPowerGeneratedL2:
					powerGenerated.With(prometheus.Labels{"phase": "l2"}).Set((floatValue(obj.Values[0].Value) * 1000))
				case gop1.OBISTypeInstantaneousPowerGeneratedL3:
					powerGenerated.With(prometheus.Labels{"phase": "l3"}).Set((floatValue(obj.Values[0].Value) * 1000))

				case gop1.OBISTypeInstantaneousCurrentL1:
					currentConsumed.With(prometheus.Labels{"phase": "l1"}).Set(floatValue(obj.Values[0].Value))
				case gop1.OBISTypeInstantaneousCurrentL2:
					currentConsumed.With(prometheus.Labels{"phase": "l2"}).Set(floatValue(obj.Values[0].Value))
				case gop1.OBISTypeInstantaneousCurrentL3:
					currentConsumed.With(prometheus.Labels{"phase": "l3"}).Set(floatValue(obj.Values[0].Value))

				case gop1.OBISTypeInstantaneousVoltageL1:
					voltageConsumed.With(prometheus.Labels{"phase": "l1"}).Set(floatValue(obj.Values[0].Value))
				case gop1.OBISTypeInstantaneousVoltageL2:
					voltageConsumed.With(prometheus.Labels{"phase": "l2"}).Set(floatValue(obj.Values[0].Value))
				case gop1.OBISTypeInstantaneousVoltageL3:
					voltageConsumed.With(prometheus.Labels{"phase": "l3"}).Set(floatValue(obj.Values[0].Value))

				case gop1.OBISTypeElectricityTariffIndicator:
					tariffIndicator.Set(floatValue(obj.Values[0].Value))

				case gop1.OBISTypeElectricityDeliveredTariff1:
					electricityConsumed.With(prometheus.Labels{"tariff": "1"}).Set(floatValue(obj.Values[0].Value))
				case gop1.OBISTypeElectricityDeliveredTariff2:
					electricityConsumed.With(prometheus.Labels{"tariff": "2"}).Set(floatValue(obj.Values[0].Value))

				case gop1.OBISTypeElectricityGeneratedTariff1:
					electricityGenerated.With(prometheus.Labels{"tariff": "1"}).Set(floatValue(obj.Values[0].Value))
				case gop1.OBISTypeElectricityGeneratedTariff2:
					electricityGenerated.With(prometheus.Labels{"tariff": "2"}).Set(floatValue(obj.Values[0].Value))

				case gop1.OBISTypeGasDelivered:
					gasConsumed.Set(floatValue(obj.Values[0].Value))
				}
			}
		}
	}
}
