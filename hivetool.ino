/** Time **/
#include <Time.h>

#define TIME_HEADER  "T"   // Header tag for serial time sync message
#define TIME_REQUEST  7    // ASCII bell character requests a time sync message

/** Temperature **/
#include <OneWire.h>
#include <DallasTemperature.h>

// Setup a oneWire instance to communicate with any OneWire devices (not just Maxim/Dallas temperature ICs)
OneWire wireOne(10);
OneWire wireTwo(11);

DallasTemperature sensor10(&wireOne);
DallasTemperature sensor11(&wireTwo);

const int ledPin = 13;
float lastTempWireOne = -127;
float lastTempWireTwo = -127;

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
  sensor10.begin();
  sensor11.begin();
}

void loop(void)
{
  int startTime = millis();


  // Send the command to get temperatures
  if(true) {
      sensor10.requestTemperatures();
      float s10 = sensor10.getTempCByIndex(0);
      if(s10 != lastTempWireOne) {
          lastTempWireOne = s10;
          writeTemp("s10", s10);
      }

      float s11 = sensor10.getTempCByIndex(1);
      if(s11 != lastTempWireTwo) {
          lastTempWireTwo = s11;
          writeTemp("s11", s11);
      }
  } else {
      sensor10.requestTemperatures();
      float s10 = sensor10.getTempCByIndex(0);
      if(s10 != lastTempWireTwo) {
          lastTempWireOne = s10;
          writeTemp("10", s10);
      }

      sensor11.requestTemperatures();
      float s11 = sensor11.getTempCByIndex(0);
      if(s11 != lastTempWireTwo) {
          lastTempWireTwo = s11;
          writeTemp("11", s11);
      }
  }


  int endTime = millis();
  int elapsed = endTime - startTime;
  int remaining = 15000 - elapsed;

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

time_t requestSync()
{
  Serial.println("cfg:time_unset");
  return 0; // the time will be sent later in response to serial mesg
}
