
# BAMBULABS EXPORTER (In progress...)
This is an exporter for all the data peeps that want to know all the things about their fancy 3D BambuLabs Printer. This is only tested on the X1 Carbon as I only have one. Feel free to ship me a P1P and I'll test that out too!


## DOCKER ⚡ Powered eventually...
This is an Exporter for SNMP utilized for pulling data out and providing an endpoint to scrape data via Docker container endpoints.

## (Important Notes)
You will need to likely run an MQTT program to test your connection. You can pull the password from the printer interface manually, or reset it on the printer itself.

## Dev
Building the exporter, I'm using simple go commands to run this for now. Eventually, I'll dockerize the entire project.


#### Prometheus Ingestion
Setup prometheus to scrape the node and setup the ports to pull from port 9101.

###
`Give me a shout if you utilize this code base (Anywhere!)
`
###

### Tyler Bennet Git https://github.com/Aetrius tylerwbennet@gmail.com
