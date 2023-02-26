package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var data BambuLabsX1C

type bambulabsCollector struct {
	humidityMetric *prometheus.Desc
}

// You must create a constructor for you collector that
// initializes every descriptor and returns a pointer to the collector
func newBambulabsCollector() *bambulabsCollector {
	return &bambulabsCollector{
		humidityMetric: prometheus.NewDesc("humidity_metric",
			"humidity of the ams",
			nil, nil,
		),
	}
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *bambulabsCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.humidityMetric
}

// Collect implements required collect function for all prometheus collectors
func (collector *bambulabsCollector) Collect(ch chan<- prometheus.Metric) {

	var broker = "192.168.1.32"
	var port = 8883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("bblp")
	opts.SetPassword("06859d8f")
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
	defer client.Disconnect(250)
	token.Wait()
	time.Sleep(time.Second)
	//fmt.Printf("\nHumidity: %s", data.Print.Ams.Ams[0].Humidity)

	m1 := prometheus.MustNewConstMetric(collector.humidityMetric, prometheus.GaugeValue, 1)
	ch <- m1

}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//fmt.Printf("Payload %s\n", msg.Payload())
	s := msg.Payload()
	data := BambuLabsX1C{}
	json.Unmarshal([]byte(s), &data)
	//fmt.Printf("\nHumidity: %s", data.Print.Ams.Ams[0].Humidity)
	//fmt.Printf()
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %+v", err)
}

func main() {

	bambulabs := newBambulabsCollector()
	prometheus.MustRegister(bambulabs)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))

}

func newTLSConfig() *tls.Config {
	return &tls.Config{InsecureSkipVerify: true}
}

func sub(client mqtt.Client) {
	// Subscribe to the LWT connection status
	topic := "device/00M00A2B0809765/report"
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
