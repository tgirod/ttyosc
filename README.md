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

Here is an example for the arduino, taken from https://github.com/CNMAT/OSC/

```
#include <OSCBoards.h>
#include <OSCBundle.h>
#include <SLIPEncodedSerial.h>

SLIPEncodedSerial SLIPSerial(Serial);

void setup() {
  SLIPSerial.begin(115200);
  pinMode(13, OUTPUT);
  digitalWrite(13, LOW);
}

void loop() {
  OSCBundle bndl;
  int size;

  // receive a bundle
  while (!SLIPSerial.endofPacket()) {
    if ((size = SLIPSerial.available()) > 0) {
      while (size--)
        bndl.fill(SLIPSerial.read());
    }
  }

  if (!bndl.hasError()) {
    static int32_t sequencenumber = 0;
    // we can sneak an addition onto the end of the bundle
    bndl.add("/micros").add(
        (int32_t)micros()); // (int32_t) is the type of OSC Integers
    bndl.add("/sequencenumber").add(sequencenumber++);
    bndl.add("/digital/5").add(digitalRead(5) == HIGH);
    bndl.add("/lsb").add((sequencenumber & 1) == 1);
    SLIPSerial.beginPacket(); // mark the beginning of the OSC Packet
    bndl.send(SLIPSerial);
    SLIPSerial.endPacket();
  }
}
```

