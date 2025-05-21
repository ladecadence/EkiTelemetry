# EkiTelemetry

GUI desktop software to receive telemetry and SSDV images from the EKI high altitude balloon missions.


## Building the project

You'll need Go (>1.24) and a C compiler (for CGo).

Clone or download the repository and build the project:

```
$ git clone https://github.com/ladecadence/EkiTelemetry
$ go mod tidy
$ go build cmd/ekitelemetry.go
```

You'll have a ekitelemetry binary in the project folder.

## Running

You'll need a compatible LoRa receiver running the EkiTelemetryReceiver firmware connected to the computer. (https://github.com/ladecadence/EkiTelemetryReceiver)
Then go to the config tab, select the serial port where the receiver is connected, and a folder to store the data log and the received images. Save the configuration and close and open the program again.


