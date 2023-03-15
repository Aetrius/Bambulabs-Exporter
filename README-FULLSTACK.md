
## Tutorial (for Monitoring stack (prometheus + grafana) + Bambulabs exporter)
- [https://www.youtube.com/watch?v=VQSAEKGYJBQ](https://www.youtube.com/watch?v=VQSAEKGYJBQ)
- (Use the steps below if you just want to run the exporter and assuming you have a Prometheus/Grafana stack already).
- If your monitor-compose.yml file will not start without a docker network. Run the following command `docker network create bambulabs-exporter_default`.


## Steps to run the exporter
Step 0: [Prereqs](#step-0-prereqs)

Step 1: [Create the env file](#step-1-env-file)

Step 2: [Clone the Repo](#step-2-clone-the-repo)

Step 3: [Run Docker Compose](#step-3-run-docker-compose)

---

## Step 0: Prereqs
This program/container would run on a virtual host, raspberry pi, or a computer that has access to the BambuLabs printer. IT is possible to port forward your printer and host this in AWS or off-premise.
- Install Git (only for windows)
- Install Docker
- Install Docker-Compose

Setup Bambulabs Exporter (Docker Network)
`docker network create bambulabs-exporter_default`

---

## Step 1: env File
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


## Step 2: Clone the repo

```
git clone https://github.com/Aetrius/Bambulabs-Exporter.git
```

## Step 3: Run Docker Compose
```
cd Bambulabs-Exporter
docker-compose -f ./monitoring/monitoring-compose.yml up -d
```
