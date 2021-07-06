If you have a device that sends and receives OSC through SLIPSerial, you can use this program to bridge it with other programs on the host computer.

```
Usage of ./ttyosc:
  -baud int
    	serial connection's baudrate (default 115200)
  -l string
    	local address: send your osc here (default ":6020")
  -r string
    	remote address: osc from serial is sent here (default ":6010")
  -t string
    	serial device (default "/dev/ttyUSB0")
  -v	verbose
```
