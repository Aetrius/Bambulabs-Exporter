package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var data BambuLabsX1C
var datav2 BambuLabsX1C

var username string
var password string
var broker string
var mqtt_topic string

// var humidity float64
// var ams_temp float64
// var ams_bed_temp float64
var layer_number float64
var print_error float64

var wifi_signal float64

var big_fan1_speed float64
var big_fan2_speed float64
var chamber_temper float64
var cooling_fan_speed float64
var fail_reason float64
var fan_gear float64

// var gcode_state string
var mc_percent float64
var mc_print_error_code float64
var mc_print_stage float64
var mc_print_sub_stage float64
var mc_remaining_time float64
var nozzle_target_temper float64
var nozzle_temper float64

var unmarshal bool

type bambulabsCollector struct {
	amsHumidityMetric     *prometheus.Desc
	amsTempMetric         *prometheus.Desc
	amsBedTempMetric      *prometheus.Desc
	amsColorMetric        *prometheus.Desc //Custom color metric with multiple labels
	layerNumberMetric     *prometheus.Desc
	printErrorMetric      *prometheus.Desc
	wifiSignalMetric      *prometheus.Desc
	bigFan1SpeedMetric    *prometheus.Desc
	bigFan2SpeedMetric    *prometheus.Desc
	chamberTemperMetric   *prometheus.Desc
	coolingFanSpeedMetric *prometheus.Desc
	failReasonMetric      *prometheus.Desc
	fanGearMetric         *prometheus.Desc
	//gCodeStateMetric       *prometheus.Desc
	mcPercentMetric          *prometheus.Desc
	mcPrintErrorCodeMetric   *prometheus.Desc
	mcPrintStageMetric       *prometheus.Desc
	mcPrintSubStageMetric    *prometheus.Desc
	mcRemainingTimeMetric    *prometheus.Desc
	nozzleTargetTemperMetric *prometheus.Desc
	nozzleTemperMetric       *prometheus.Desc
}

func env(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// You must create a constructor for you collector that
// initializes every descriptor and returns a pointer to the collector
func newBambulabsCollector() *bambulabsCollector {
	return &bambulabsCollector{
		amsHumidityMetric: prometheus.NewDesc("ams_humidity_metric",
			"humidity of the ams",
			[]string{"ams_number"}, nil,
		),
		amsTempMetric: prometheus.NewDesc("ams_temp_metric",
			"temperature of the ams",
			[]string{"ams_number"}, nil,
		),
		amsColorMetric: prometheus.NewDesc("ams_tray_color_metric",
			"ID of the ams with color hex values",
			[]string{"ams_number", "tray_number", "tray_color", "tray_type"}, nil,
		),
		amsBedTempMetric: prometheus.NewDesc("ams_bed_temp_metric",
			"temperature of the ams bed",
			[]string{"ams_number", "tray_number"}, nil,
		),
		layerNumberMetric: prometheus.NewDesc("layer_number_metric",
			"layer number of the print head in gcode",
			nil, nil,
		),
		printErrorMetric: prometheus.NewDesc("print_error_metric",
			"Print error int",
			nil, nil,
		),
		wifiSignalMetric: prometheus.NewDesc("wifi_signal_metric",
			"Wifi signal in dBm",
			nil, nil,
		),
		bigFan1SpeedMetric: prometheus.NewDesc("big_fan1_speed_metric",
			"Big Fan 1 Speed",
			nil, nil,
		),
		bigFan2SpeedMetric: prometheus.NewDesc("big_fan2_speed_metric",
			"Big Fan 2 Speed",
			nil, nil,
		),
		chamberTemperMetric: prometheus.NewDesc("chamber_temper_metric",
			"Chamber Temperature of Printer",
			nil, nil,
		),
		coolingFanSpeedMetric: prometheus.NewDesc("cooling_fan_speed_metric",
			"Cooling Fan Speed",
			nil, nil,
		),
		failReasonMetric: prometheus.NewDesc("fail_reason_metric",
			"Print Failure Reason",
			nil, nil,
		),
		fanGearMetric: prometheus.NewDesc("fan_gear_metric",
			"Fan Gear",
			nil, nil,
		),
		mcPercentMetric: prometheus.NewDesc("mc_percent_metric",
			"Percentage of Progress of print",
			nil, nil,
		),
		mcPrintErrorCodeMetric: prometheus.NewDesc("mc_print_error_code_metric",
			"Print Progress Error Code",
			nil, nil,
		),
		mcPrintStageMetric: prometheus.NewDesc("mc_print_stage_metric",
			"Print Progress Stage",
			nil, nil,
		),
		mcPrintSubStageMetric: prometheus.NewDesc("mc_print_sub_stage_metric",
			"Print Progress Sub Stage",
			nil, nil,
		),
		mcRemainingTimeMetric: prometheus.NewDesc("mc_remaining_time_metric",
			"Print Progress Remaining Time in minutes",
			nil, nil,
		),
		nozzleTargetTemperMetric: prometheus.NewDesc("nozzle_target_temper_metric",
			"Nozzle Target Temperature Metric",
			nil, nil,
		),
		nozzleTemperMetric: prometheus.NewDesc("nozzle_temper_metric",
			"Nozzle Temperature Metric",
			nil, nil,
		),
	}
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *bambulabsCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.amsHumidityMetric
	ch <- collector.amsTempMetric
	ch <- collector.amsColorMetric
	ch <- collector.amsBedTempMetric
	ch <- collector.layerNumberMetric
	ch <- collector.printErrorMetric
	ch <- collector.wifiSignalMetric
	ch <- collector.bigFan1SpeedMetric
	ch <- collector.bigFan2SpeedMetric
	ch <- collector.chamberTemperMetric
	ch <- collector.coolingFanSpeedMetric
	ch <- collector.failReasonMetric
	ch <- collector.fanGearMetric
	ch <- collector.mcPercentMetric
	ch <- collector.mcPrintErrorCodeMetric
	ch <- collector.mcPrintStageMetric
	ch <- collector.mcPrintSubStageMetric
	ch <- collector.mcRemainingTimeMetric
	ch <- collector.nozzleTargetTemperMetric
	ch <- collector.nozzleTemperMetric
}

// Collect implements required collect function for all prometheus collectors
func (collector *bambulabsCollector) Collect(ch chan<- prometheus.Metric) {

	//var broker = broker
	var port = 8883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	opts.SetTLSConfig(newTLSConfig())
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client)
	//defer client.Disconnect(250)
	//defer token.Done()
	token.Wait()
	time.Sleep(time.Second)
	defer client.Disconnect(250)
	defer token.Done()
	//fmt.Printf("\nHumidity: %s", data.Print.Ams.Ams[0].Humidity)

	if reflect.ValueOf(data).IsZero() == true {
		//Loop through the AMS
		for x := 0; x < len(datav2.Print.Ams.Ams); x++ {

			ams_temp, _ := strconv.ParseFloat(datav2.Print.Ams.Ams[x].Temp, 64)
			ams_temp_1 := prometheus.MustNewConstMetric(collector.amsTempMetric, prometheus.GaugeValue, ams_temp, strconv.Itoa(x))
			ch <- ams_temp_1

			humidity, _ := strconv.ParseFloat(datav2.Print.Ams.Ams[x].Humidity, 64)
			humidity_1 := prometheus.MustNewConstMetric(collector.amsHumidityMetric, prometheus.GaugeValue, humidity, strconv.Itoa(x))
			ch <- humidity_1

			// loop through the Trays
			for i := 0; i < len(datav2.Print.Ams.Ams[x].Tray); i++ {

				ams_bed_temp, _ := strconv.ParseFloat(datav2.Print.Ams.Ams[x].Tray[i].BedTemp, 64)
				ams_bed_temp_1 := prometheus.MustNewConstMetric(collector.amsBedTempMetric, prometheus.GaugeValue, ams_bed_temp, strconv.Itoa(x), strconv.Itoa(i))
				ch <- ams_bed_temp_1

				ams_tray_color := datav2.Print.Ams.Ams[x].Tray[i].TrayColor
				ams_tray_type := datav2.Print.Ams.Ams[x].Tray[i].TrayType
				//ams_tray_id, _ := strconv.ParseFloat(datav2.Print.Ams.Ams[x].Tray[i].ID, 64)
				ams_color_1 := prometheus.MustNewConstMetric(collector.amsColorMetric, prometheus.GaugeValue, 1, strconv.Itoa(x), strconv.Itoa(i), ams_tray_color, ams_tray_type)
				ch <- ams_color_1

			}
		}

		//fmt.Println("\nHumidity: ", ams_temp)
		// humidity_1 := prometheus.MustNewConstMetric(collector.amsHumidityMetric, prometheus.GaugeValue, humidity)
		// ch <- humidity_1

		// ams_temp_1 := prometheus.MustNewConstMetric(collector.amsTempMetric, prometheus.GaugeValue, ams_temp)
		// ch <- ams_temp_1

		// ams_bed_temp_1 := prometheus.MustNewConstMetric(collector.amsBedTempMetric, prometheus.GaugeValue, ams_bed_temp)
		// ch <- ams_bed_temp_1

		layer_number_1 := prometheus.MustNewConstMetric(collector.layerNumberMetric, prometheus.GaugeValue, layer_number)
		ch <- layer_number_1

		print_error_1 := prometheus.MustNewConstMetric(collector.printErrorMetric, prometheus.GaugeValue, print_error)
		ch <- print_error_1

		wifi_signal_1 := prometheus.MustNewConstMetric(collector.wifiSignalMetric, prometheus.GaugeValue, wifi_signal)
		ch <- wifi_signal_1

		big_fan1_speed_1 := prometheus.MustNewConstMetric(collector.bigFan1SpeedMetric, prometheus.GaugeValue, big_fan1_speed)
		ch <- big_fan1_speed_1

		big_fan2_speed_1 := prometheus.MustNewConstMetric(collector.bigFan2SpeedMetric, prometheus.GaugeValue, big_fan2_speed)
		ch <- big_fan2_speed_1

		chamber_temper_1 := prometheus.MustNewConstMetric(collector.chamberTemperMetric, prometheus.GaugeValue, chamber_temper)
		ch <- chamber_temper_1

		cooling_fan_speed_1 := prometheus.MustNewConstMetric(collector.coolingFanSpeedMetric, prometheus.GaugeValue, cooling_fan_speed)
		ch <- cooling_fan_speed_1

		fail_reason_metric_1 := prometheus.MustNewConstMetric(collector.failReasonMetric, prometheus.GaugeValue, fail_reason)
		ch <- fail_reason_metric_1

		fan_gear_metric_1 := prometheus.MustNewConstMetric(collector.fanGearMetric, prometheus.GaugeValue, fan_gear)
		ch <- fan_gear_metric_1

		mc_percent_1 := prometheus.MustNewConstMetric(collector.mcPercentMetric, prometheus.GaugeValue, mc_percent)
		ch <- mc_percent_1

		mc_print_error_code_1 := prometheus.MustNewConstMetric(collector.mcPrintErrorCodeMetric, prometheus.GaugeValue, mc_print_error_code)
		ch <- mc_print_error_code_1

		mc_print_stage_metric_1 := prometheus.MustNewConstMetric(collector.mcPrintStageMetric, prometheus.GaugeValue, mc_print_stage)
		ch <- mc_print_stage_metric_1

		mc_print_sub_stage_metric_1 := prometheus.MustNewConstMetric(collector.mcPrintSubStageMetric, prometheus.GaugeValue, mc_print_sub_stage)
		ch <- mc_print_sub_stage_metric_1

		mc_remaining_time_metric_1 := prometheus.MustNewConstMetric(collector.mcRemainingTimeMetric, prometheus.GaugeValue, mc_remaining_time)
		ch <- mc_remaining_time_metric_1

		nozzle_target_temper_metric_1 := prometheus.MustNewConstMetric(collector.nozzleTargetTemperMetric, prometheus.GaugeValue, nozzle_target_temper)
		ch <- nozzle_target_temper_metric_1

		nozzle_temper_metric_1 := prometheus.MustNewConstMetric(collector.nozzleTemperMetric, prometheus.GaugeValue, nozzle_temper)
		ch <- nozzle_temper_metric_1

		client.Disconnect(1)
		token.Done()
	} else {
		fmt.Printf("\ndata might be empty")
	}

}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//fmt.Printf("Payload %s\n", msg.Payload())
	s := msg.Payload()
	data := BambuLabsX1C{}
	json.Unmarshal([]byte(s), &data)

	//if reflect.ValueOf(data).IsZero() == false {
	//fmt.Printf("\nHumidity: %s", data.Print.Ams.Ams[0].Humidity)
	if data.Print.WifiSignal == "" {
		//fmt.Println("\nWifi Signal was empty")
	} else {
		datav2 = data

		//humidity, _ = strconv.ParseFloat(data.Print.Ams.Ams[0].Humidity, 64)
		//ams_temp, _ = strconv.ParseFloat(data.Print.Ams.Ams[0].Temp, 64)
		//ams_bed_temp, _ = strconv.ParseFloat(data.Print.Ams.Ams[0].Tray[0].BedTemp, 64)
		layer_number = float64(data.Print.LayerNum)
		print_error = float64(data.Print.PrintError)
		wifi_signal, _ = strconv.ParseFloat(strings.ReplaceAll(data.Print.WifiSignal, "dBm", ""), 64)
		big_fan1_speed, _ = strconv.ParseFloat(data.Print.BigFan1Speed, 64)
		big_fan2_speed, _ = strconv.ParseFloat(data.Print.BigFan2Speed, 64)
		chamber_temper = data.Print.ChamberTemper
		cooling_fan_speed, _ = strconv.ParseFloat(data.Print.CoolingFanSpeed, 64)
		fail_reason, _ = strconv.ParseFloat(data.Print.FailReason, 64)
		fan_gear = float64(data.Print.FanGear)
		mc_percent = float64(data.Print.McPercent)
		mc_print_error_code, _ = strconv.ParseFloat(data.Print.McPrintErrorCode, 64)
		mc_print_stage, _ = strconv.ParseFloat(data.Print.McPrintStage, 64)
		mc_print_sub_stage = float64(data.Print.McPrintSubStage)
		mc_remaining_time = float64(data.Print.McRemainingTime)
		nozzle_temper = float64(data.Print.NozzleTemper)

	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	dt := time.Now()
	fmt.Println("\nConnected: ", dt.String())
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("\nConnect lost: %+v", err)
}

func main() {
	dt := time.Now()
	fmt.Printf("\nStarting Exporter: ", dt.String())
	godotenv.Load()

	broker = env("BAMBU_PRINTER_IP")
	username = env("USERNAME")
	password = env("PASSWORD")
	mqtt_topic = env("MQTT_TOPIC")

	if broker == "" {
		broker = os.Getenv("BAMBU_PRINTER_IP")
	}

	if password == "" {
		password = os.Getenv("PASSWORD")
	}

	if mqtt_topic == "device/<>/report" {
		mqtt_topic = os.Getenv("MQTT_TOPIC")
	}

	fmt.Printf("\nEnv Vars Loaded")

	fmt.Printf("\nRegistering collector")
	bambulabs := newBambulabsCollector()
	prometheus.MustRegister(bambulabs)
	http.HandleFunc("/", home)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))

}

const body = `<html><head><title>BambuLabs Exporter Metrics</title></head><body>
			<h1>BambuLabs Exporter</h1><p><a href='` + "/metrics" + `'>Metrics</a></p></body>
			</html>`

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, body)
}

func newTLSConfig() *tls.Config {
	return &tls.Config{InsecureSkipVerify: true}
}

func sub(client mqtt.Client) {
	// Subscribe to the LWT connection status
	topic := mqtt_topic
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to LWT %s", topic)
}

type BambuLabsX1C struct {
	Print struct {
		Ams struct {
			Ams []struct {
				Humidity string `json:"humidity"`
				ID       string `json:"id"`
				Temp     string `json:"temp"`
				Tray     []struct {
					BedTemp       string `json:"bed_temp"`
					BedTempType   string `json:"bed_temp_type"`
					DryingTemp    string `json:"drying_temp"`
					DryingTime    string `json:"drying_time"`
					ID            string `json:"id"`
					NozzleTempMax string `json:"nozzle_temp_max"`
					NozzleTempMin string `json:"nozzle_temp_min"`
					Remain        int    `json:"remain"`
					TagUID        string `json:"tag_uid"`
					TrayColor     string `json:"tray_color"`
					TrayDiameter  string `json:"tray_diameter"`
					TrayIDName    string `json:"tray_id_name"`
					TrayInfoIdx   string `json:"tray_info_idx"`
					TraySubBrands string `json:"tray_sub_brands"`
					TrayType      string `json:"tray_type"`
					TrayUUID      string `json:"tray_uuid"`
					TrayWeight    string `json:"tray_weight"`
					XcamInfo      string `json:"xcam_info"`
				} `json:"tray"`
			} `json:"ams"`
			AmsExistBits     string `json:"ams_exist_bits"`
			InsertFlag       bool   `json:"insert_flag"`
			PowerOnFlag      bool   `json:"power_on_flag"`
			TrayExistBits    string `json:"tray_exist_bits"`
			TrayIsBblBits    string `json:"tray_is_bbl_bits"`
			TrayNow          string `json:"tray_now"`
			TrayReadDoneBits string `json:"tray_read_done_bits"`
			TrayReadingBits  string `json:"tray_reading_bits"`
			TrayTar          string `json:"tray_tar"`
			Version          int    `json:"version"`
		} `json:"ams"`
		AmsRfidStatus           int     `json:"ams_rfid_status"`
		AmsStatus               int     `json:"ams_status"`
		BedTargetTemper         float64 `json:"bed_target_temper"`
		BedTemper               float64 `json:"bed_temper"`
		BigFan1Speed            string  `json:"big_fan1_speed"`
		BigFan2Speed            string  `json:"big_fan2_speed"`
		ChamberTemper           float64 `json:"chamber_temper"`
		Command                 string  `json:"command"`
		CoolingFanSpeed         string  `json:"cooling_fan_speed"`
		FailReason              string  `json:"fail_reason"`
		FanGear                 int     `json:"fan_gear"`
		ForceUpgrade            bool    `json:"force_upgrade"`
		GcodeFile               string  `json:"gcode_file"`
		GcodeFilePreparePercent string  `json:"gcode_file_prepare_percent"`
		GcodeStartTime          string  `json:"gcode_start_time"`
		GcodeState              string  `json:"gcode_state"`
		HeatbreakFanSpeed       string  `json:"heatbreak_fan_speed"`
		Hms                     []any   `json:"hms"`
		HomeFlag                int     `json:"home_flag"`
		HwSwitchState           int     `json:"hw_switch_state"`
		Ipcam                   struct {
			IpcamDev    string `json:"ipcam_dev"`
			IpcamRecord string `json:"ipcam_record"`
			Resolution  string `json:"resolution"`
			Timelapse   string `json:"timelapse"`
		} `json:"ipcam"`
		LayerNum     int    `json:"layer_num"`
		Lifecycle    string `json:"lifecycle"`
		LightsReport []struct {
			Mode string `json:"mode"`
			Node string `json:"node"`
		} `json:"lights_report"`
		Maintain            int     `json:"maintain"`
		McPercent           int     `json:"mc_percent"`
		McPrintErrorCode    string  `json:"mc_print_error_code"`
		McPrintStage        string  `json:"mc_print_stage"`
		McPrintSubStage     int     `json:"mc_print_sub_stage"`
		McRemainingTime     int     `json:"mc_remaining_time"`
		MessProductionState string  `json:"mess_production_state"`
		NozzleTargetTemper  float64 `json:"nozzle_target_temper"`
		NozzleTemper        float64 `json:"nozzle_temper"`
		Online              struct {
			Ahb  bool `json:"ahb"`
			Rfid bool `json:"rfid"`
		} `json:"online"`
		PrintError       int    `json:"print_error"`
		PrintGcodeAction int    `json:"print_gcode_action"`
		PrintRealAction  int    `json:"print_real_action"`
		PrintType        string `json:"print_type"`
		ProfileID        string `json:"profile_id"`
		ProjectID        string `json:"project_id"`
		Sdcard           bool   `json:"sdcard"`
		SequenceID       string `json:"sequence_id"`
		SpdLvl           int    `json:"spd_lvl"`
		SpdMag           int    `json:"spd_mag"`
		Stg              []int  `json:"stg"`
		StgCur           int    `json:"stg_cur"`
		SubtaskID        string `json:"subtask_id"`
		SubtaskName      string `json:"subtask_name"`
		TaskID           string `json:"task_id"`
		TotalLayerNum    int    `json:"total_layer_num"`
		UpgradeState     struct {
			AhbNewVersionNumber string `json:"ahb_new_version_number"`
			AmsNewVersionNumber string `json:"ams_new_version_number"`
			ConsistencyRequest  bool   `json:"consistency_request"`
			DisState            int    `json:"dis_state"`
			ErrCode             int    `json:"err_code"`
			ForceUpgrade        bool   `json:"force_upgrade"`
			Message             string `json:"message"`
			Module              string `json:"module"`
			NewVersionState     int    `json:"new_version_state"`
			OtaNewVersionNumber string `json:"ota_new_version_number"`
			Progress            string `json:"progress"`
			SequenceID          int    `json:"sequence_id"`
			Status              string `json:"status"`
		} `json:"upgrade_state"`
		Upload struct {
			FileSize      int    `json:"file_size"`
			FinishSize    int    `json:"finish_size"`
			Message       string `json:"message"`
			OssURL        string `json:"oss_url"`
			Progress      int    `json:"progress"`
			SequenceID    string `json:"sequence_id"`
			Speed         int    `json:"speed"`
			Status        string `json:"status"`
			TaskID        string `json:"task_id"`
			TimeRemaining int    `json:"time_remaining"`
			TroubleID     string `json:"trouble_id"`
		} `json:"upload"`
		WifiSignal string `json:"wifi_signal"`
		Xcam       struct {
			AllowSkipParts           bool   `json:"allow_skip_parts"`
			BuildplateMarkerDetector bool   `json:"buildplate_marker_detector"`
			FirstLayerInspector      bool   `json:"first_layer_inspector"`
			HaltPrintSensitivity     string `json:"halt_print_sensitivity"`
			PrintHalt                bool   `json:"print_halt"`
			PrintingMonitor          bool   `json:"printing_monitor"`
			SpaghettiDetector        bool   `json:"spaghetti_detector"`
		} `json:"xcam"`
		XcamStatus string `json:"xcam_status"`
	} `json:"print"`
}
