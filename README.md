
# BAMBULABS EXPORTER
This is an exporter for all the data peeps that want to know all the things about their awesome BambuLabs 3D Printer. This was only tested on the X1C.

`Supported BambuLabs X1 Firmware:
1.04.01.00`


## GO, DOCKER, & PROMETHEUS ⚡ Powered
This is an MQTT Exporter powered by Go & Docker. 
https://hub.docker.com/r/aetrius/bambulabs-exporter


### Prometheus Metrics Available
- `*annotates recent changes or additions`

| Metric   | Description | Examples |
| ------------- | ------------- |  ------------- |
| ams_humidity  | Humdity of the Enclosure, includes the AMS Number 0-many  |ams_humidity_metric{ams_number="0"} 4 |
| ams_temp  | *Temperature of the AMS, includes the AMS Number 0-many | ams_temp_metric{ams_number="0"} 30.7 |
| ams_tray_color | *Filament color in the tray of the AMS, includes the AMS Number 0-many & Tray Numbers 0-4 | ams_tray_color_metric{ams_number="0",tray_color="AF7933FF",tray_number="0",tray_type="PLA"} 1 |
| ams_bed_temp | *Temperature of the AMS bed, includes the AMS Number 0-many & Tray Numbers 0-4 | ams_bed_temp_metric{ams_number="0",tray_number="0"} 0 |
| layer_number | GCode Layer number  | |
| print_error | Print Error Code Detected  | |
| wifi_signal | Wifi Signal in dBm  | |
| big_fan1_speed | Big1 Fan Speed  | big_fan1_speed_metric 0 |
| big_fan2_speed | Big2 Fan Speed  | big_fan2_speed_metric 0 |
| chamber_temper | Temperature of the Bambu Enclosure  | chamber_temper_metric 30 |
| cooling_fan_speed | Print Head Cooling Fan Speed  | cooling_fan_speed_metric 0 |
| fail_reason | Failure Print Reason Code  | fail_reason_metric 0 |
| fan_gear | Fan Gear   | fan_gear_metric 0 |
| mc_percent | Print Progress in Percentage  | mc_percent_metric 36 |
| mc_print_error_code | Print Progress Error Code | mc_print_error_code_metric 0 |
| mc_print_stage | Print Progress Stage | mc_print_stage_metric 2 |
| mc_print_sub_stage | Print Progress Sub Stage | mc_print_sub_stage_metric 4 |
| mc_remaining_time | Print Progress Remaining Time in minutes  | mc_remaining_time_metric 1973 |
| nozzle_target_temper |Nozzle Target Temperature Metric | nozzle_target_temper_metric 0 |
| nozzle_temper | Nozzle Temperature Metric | nozzle_temper_metric 221 |
| print_error | Print Error reported by the Control board | print_error_metric 0 |
---

## Steps to run the exporter
Step 0: [Prereqs](#step-0-prereqs)

Step 1: [Create the env file](#step-1-env-file)

Step 2: [Clone the Repo](#step-2-clone-the-repo)

Step 3: [Run Docker Compose](#step-3-run-docker-compose)

---

## Step 0: Prereqs
This project assumes you have a Grafana/Prometheus Setup. You would point your Prometheus instance to the (host:9101) endpoint. This is not a tutorial on Prometheus / Grafana. Click [here](README-FULLSTACK.md) for a full stack that includes Prometheus &  Grafana for this.

This program/container would run on a virtual host, raspberry pi, or a computer that has access to the Bambu printer. IT is possible to port forward your printer and host this in AWS or offpremise.
- Install Git (only for windows)
- Install Docker
- Install Docker-Compose

---

## Step 1 env File
Create an .env file.
Add the Printer IP you configured when you setup your printer.
Add the Printer Password from the Printer Network Settings Menu.
Add the MQTT_TOPIC for your printer. This can be achived by loading up an MQTT Application such as MQTT Explorer. 
- You must Enable (TLS), use the protocol mqtt://, add the port 8883, username bblp, and the password on your printer. 
- *Please note you can regenerate a password on the device manually.
- Collect the MQTT_TOPIC by expanding the "Device", "Serial Number", "Report". The result should look similar to "device/00M00A2B08124765/report"

```
# EXAMPLE .ENV FILE
BAMBU_PRINTER_IP=""
USERNAME="bblp"
PASSWORD=""
MQTT_TOPIC="device/00M00A2B08124765/report"
```


## Step 2 Clone the repo

```
git clone https://github.com/Aetrius/Bambulabs-Exporter.git
```

## Step 3 Run Docker Compose
```
cd Bambulabs-Exporter
docker-compose up -d
```

---

## (Important Notes)
You will need to likely run an MQTT program to test your connection. You can pull the password from the printer interface manually, or reset it on the printer itself.


### Prometheus Ingestion
Setup prometheus to scrape the node and setup the ports to pull from port 9101.



### Bugs
- 3/4/2023 - Exporter loses connection during firmware upgrade or powercycle. (Temp solution to restart the docker container).

---

### Feature Changes
- 3/4/2023 - Added new Metrics ams_humidity, ams_temp, ams_tray_color, ams_bed_temp. These include ams number and tray numbers to be dynamic depending on how many AMS's are included. Will push new container to dockerhub later today 3/4/23
- 2/28/2023 - Initial Metrics released. Further re-work needed to account for all the useful metrics available.
---

### Future Development
- Add Kubernetes Configs
- Add Grafana Dashboard changes for AMS

---

### Credit
```
Give me a shout if you utilize this code base (Anywhere!)
```

---

### Support Questions 

```
tylerwbennet@gmail.com
```
