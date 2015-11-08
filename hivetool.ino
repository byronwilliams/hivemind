/** Time **/
#include <Time.h>

#define TIME_HEADER  "T"   // Header tag for serial time sync message
#define TIME_REQUEST  7    // ASCII bell character requests a time sync message

/** Temperature **/
#include <OneWire.h>
#include <DallasTemperature.h>

// Setup a oneWire instance to communicate with any OneWire devices (not just Maxim/Dallas temperature ICs)
OneWire wireOne(10);

DallasTemperature tempBusOn10(&wireOne);

const int ledPin = 13;
float temps[8] = {-127, -127, -127, -127, -127, -127, -127, -127};
DeviceAddress  MyDS18B20Addresses[8];

void setup(void)
{
  Serial.begin(9600); //Begin serial communication
  Serial.println("start: true"); //Print a message
  pinMode(ledPin, OUTPUT);
  digitalWrite(ledPin, LOW);
  setSyncProvider( requestSync);  //set function to call when sync required
  Serial.println("Waiting for sync message");

  while(timeStatus()==timeNotSet) {
    requestSync();
    processSyncMessage();
    delay(250);
  }

  pinMode(ledPin, HIGH);
  tempBusOn10.begin();
}

void loop(void)
{
  int startTime = millis();

  int sensorCount = tempBusOn10.getDeviceCount();
  tempBusOn10.requestTemperatures();
  Serial.println("log:dev_count=" + String(sensorCount));

  int curSensor;

  for(curSensor = 0; curSensor < sensorCount; curSensor++) {
      tempBusOn10.getAddress(MyDS18B20Addresses[curSensor], curSensor);

      float s10 = tempBusOn10.getTempCByIndex(curSensor);

      if(s10 != temps[curSensor]) {
          temps[curSensor] = s10;
          writeTemp(addressToString(MyDS18B20Addresses[curSensor]), s10);
      }
  }

  int endTime = millis();
  int elapsed = endTime - startTime;
  int remaining = 2500 - elapsed;

  Serial.println("elapsed:" + String(elapsed));

  if(remaining > 0) {
    delay(remaining);
  }
}

void writeTemp(String sensorId, float temp) {
  writeToSerialWithTime("temp," + sensorId + "," + String(temp));
}

void writeToSerialWithTime(String s) {
  Serial.println(String(now()) + ":" + s);
}

void processSyncMessage() {
  unsigned long pctime;
  const unsigned long DEFAULT_TIME = 1357041600; // Jan 1 2013

  if(Serial.find(TIME_HEADER)) {
     pctime = Serial.parseInt();
     if( pctime >= DEFAULT_TIME) { // check the integer is a valid time (greater than Jan 1 2013)
       setTime(pctime); // Sync Arduino clock to the time received on the serial port
     }
  }
}

String addressToString(DeviceAddress deviceAddress) {
    static char return_me[18];
    static char *hex = "0123456789ABCDEF";
    uint8_t i, j;

    for (i=0, j=0; i<8; i++)
    {
         return_me[j++] = hex[deviceAddress[i] / 16];
         return_me[j++] = hex[deviceAddress[i] & 15];
    }
    return_me[j] = '\0';

    return String(return_me);
}

time_t requestSync()
{
  Serial.println("cfg:time_unset");
  return 0; // the time will be sent later in response to serial mesg
}
