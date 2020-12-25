# Marvin - Home automation middleware

Marvin was created to connect Amazon Alexa`s Smarthome capabilities with various smarthome hardware vendors and systems like Homematic and Philips Hue.
This is basically only a code repository without installable binaries. So you will need some technical knowledge to use this project for setting up your own voice controlled smarthome. Anyway, i've added some brief documentation to give you something to start with.
I'm looking forward to get in contact for questions, feedback or collaboration!

## Prerequisites

- Some device to use Amazon Alexa (Echo, Echo Dot, etc.) and requirements from [go-alexa](https://github.com/betom84/go-alexa#requirements)
- For Homematic devices you need a CCU2, CCU3 or RaspberryMatic including [xmlapi](https://github.com/homematic-community/XML-API)
- For Philips Hue lights you need a Philips Hue Bridge (optionally a [Hue developer account](https://developers.meethue.com/develop/get-started-2/) to access API documentation)
- E.g. a Raspberry Pi to run the application
- Either a static IP or some dynamic DNS to connect with Marvin from a AWS Lambda Function
- Basic [go](https://golang.org/) knowledge to compile the project
- [Optional marvin-ui](https://github.com/betom84/marvin-ui)

## Setup

### Prepare Amazon Alexa

That's the tricky part, you need to follow Amazon's [Steps to Build a Smart Home Skill](https://developer.amazon.com/de-DE/docs/alexa/smarthome/understand-the-smart-home-skill-api.html). We will use the Lambda function to implement a proxy which forwards all alexa smarthome events to your local running Marvin instance.

#### Example Node.js 10.x index.handler implementation

```javascript
var https = require("https");

exports.handler = function (event, context) {
  console.log("Input", event);

  var options = {
    hostname: process.env.HOSTNAME,
    path: process.env.PATH,
    port: 443,
    method: "POST",
    rejectUnauthorized: false,
    headers: {
      Authorization: "Basic " + process.env.AUTH_TOKEN,
    },
  };

  var request = https.request(options, function (response) {
    var body = "";
    response.on("data", function (d) {
      body += d;
    });
    response.on("end", function () {
      console.log("Response: " + body);

      context.succeed(JSON.parse(body));
    });

    response.on("error", function (e) {
      console.log("Got error: " + e.message);
    });
  });

  request.write(JSON.stringify(event));
  request.end();
};
```

This example uses environment variables you also need to declare for your lambda function.

| Variable   | Description                                                                                                                                                      | Example                      |
| ---------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------- |
| PATH       | URI for directive requests, not important since Marvin doesn't distinguish between request URIs                                                                  | `/alexa/smarthome/directive` |
| HOSTNAME   | Either a static IP or some dynamic DNS to connect Marvin (e.g. some myfritz.net address if you are using AVM Fritz!Box)                                          | -                            |
| AUTH_TOKEN | Base64 encoded [Basic-Authentication](https://de.wikipedia.org/wiki/HTTP-Authentifizierung#Basic_Authentication) token (Make sure Marvin is running with HTTPS!) | `d2lraTpwZWRpYQ==`           |

### Install on Raspberry Pi

- checkout this repository using Mac or Linux (Makefile will not work on Windows) with [go](https://golang.org/dl/) installed
- add your individual [config.json](#configuration-configjson) and [endpoints.json](#devices-endpointsjson) to project folder
- create folder `/home/pi/marvin` on raspberry pi
- run `make deploy-rpi RPI_HOST=<ip or name>` to compile and copy resources to raspberry pi
- on raspberry pi, copy `/home/pi/marvin/resources/marvin.sh` to `/etc/init.d/marvin` to install start scripts
- run `sudo update-rc.d marvin defaults` to install rc-scripts for system startup
- restart raspberry pi

> Note: Resources are transfered using raspberry user `pi`, edit Makefile in case that doesn't suites you.

### Web-Interface

A very rudimentary UI to check log and current state is included in this project.
Check [optional marvin-ui](https://github.com/betom84/marvin-ui) for more UI features.

### Configuration (config.json)

```json
{
  "alexaServerPort": 6443,
  "uiServerPort": 8080,
  "uiRoot": "webapp",
  "log": "marvin.log",
  "endpoints": "config/endpoints.json",
  "validationEnabled": false,
  "sslCertificate": "./resources/ssl/certificate.pem",
  "sslKey": "./resources/ssl/private-key.pem",
  "amazonClientID": "amzn1.application-oa2-client....",
  "amazonClientSecret": "...",
  "restrictedUser": "john@mail.com",
  "basicAuthUser": "user",
  "basicAuthPassword": "password",
  "philipsHueHost": "philips-hue",
  "philipsHueUser": "$ANY_ENV_VARIABLE",
  "homematicHost": "homematic-ccu3"
}
```

| Key                  | Type        | Description                                                                                      | Default                         |
| -------------------- | ----------- | ------------------------------------------------------------------------------------------------ | ------------------------------- |
| `alexaServerPort`    | Port        | Port to listen for alexa directives (optional)                                                   | 6443                            |
| `uiServerPort`       | Port        | Port to serve web-ui (optional)                                                                  | 8081                            |
| `uiRoot`             | Folder      | UI web-root folder (optional)                                                                    | ./webapp                        |
| `log`                | File        | Logfile (optional)                                                                               | stdout                          |
| `endpoints`          | File        | Device configuration (required)                                                                  | ./config/endpoints.json         |
| `validationEnabled`  | Boolean     | Enable schema validation for alexa responses (debugging)                                         | false                           |
| `validationSchema`   | File        | [Schema](https://github.com/alexa/alexa-smarthome/wiki/Validation-Schemas) to use for validation | ./resources/schema.json         |
| `sslCertificate`     | File        | SSL certificate to use HTTPS                                                                     | ./resources/ssl/certificate.pem |
| `sslKey`             | File        | SSL certificate private key                                                                      | ./resources/ssl/private-key.pem |
| `amazonClientID`     | String      | Amazon Client ID                                                                                 | \$MARVIN_AMAZON_CLIENT_ID       |
| `amazonClientSecret` | String      | Amazon Client Secret                                                                             | \$MARVIN_AMAZON_CLIENT_SECRET   |
| `restrictedUser`     | String      | Restrict access to certain amazon user profiles                                                  | \$MARVIN_RESTRICTED_USER        |
| `basicAuthUser`      | String      | Username for basic authentication to use with AWS Lambda proxy function                          | \$MARVIN_BASIC_AUTH_USER        |
| `basicAuthPassword`  | String      | Password for basic authentication to use with AWS Lambda proxy function                          | \$MARVIN_BASIC_AUTH_PASSWORD    |
| `philipsHueHost`     | Hostname/IP | Hostname/IP of Philips Hue Bridge                                                                | philips-hue                     |
| `philipsHueUser`     | String      | User to access Philips Hue Bridge                                                                | \$MARVIN_PHILIPSHUE_USER        |
| `homematicHost`      | Hostname/IP | Hostname/IP of Homematic CCU/RaspberryMatic                                                      | homematic-ccu3                  |

> Note: Since we are using Basic-Auth to authorize requests it's highly recommed to configure SSL to use HTTPS!

### Devices (endpoints.json)

Devices (Endpoints) are defined in a JSON configuration file which is based on [Amazon Alexa API](https://developer.amazon.com/de-DE/docs/alexa/device-apis/alexa-discovery.html#discover-directive).

```json
[
  {
    "endpointId": "homematic-1834",
    "friendlyName": "My gorgeous switch",
    "description": "The best switch they ever made.",
    "manufacturerName": "Homematic",
    "displayCategories": ["OTHER"],
    "cookie": {
      "type": "homematic",
      "id": "1834",
      "name": "Switcherino"
    },
    "capabilities": [
      {
        "type": "AlexaInterface",
        "interface": "Alexa.PowerController",
        "version": "3",
        "properties": {
          "supported": [
            {
              "name": "powerState"
            }
          ],
          "proactivelyReported": false,
          "retrievable": true
        }
      }
    ]
  },
  {
    "endpointId": "hue-3",
    "friendlyName": "My shiny light",
    "description": "The best light they ever made.",
    "manufacturerName": "Philips Hue",
    "displayCategories": ["LIGHT"],
    "cookie": {
      "type": "hue",
      "id": "3",
      "name": "Lighterino"
    },
    "capabilities": [
      {
        "type": "AlexaInterface",
        "interface": "Alexa.PowerController",
        "version": "3",
        "properties": {
          "supported": [
            {
              "name": "powerState"
            }
          ],
          "proactivelyReported": false,
          "retrievable": true
        }
      }
    ]
  }
]
```

`cookie` section is used by Marvin to distinguish device vendors. Currently `homematic` and `hue` are supported for `type`.
`id` field is used for device id required to idenify the device in vendor API (use [xmlapi](https://github.com/homematic-community/XML-API) or [Hue API](https://developers.meethue.com/develop/get-started-2/) to check for device ids).
