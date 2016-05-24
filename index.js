const aws = require('aws-sdk');
const iot = require('aws-iot-device-sdk');

const region = 'us-east-1';
const pool = 'us-east-1:594b1376-9d0e-4adb-9f91-1a3fae53d247';

aws.config.region = region;
aws.config.credentials = new aws.CognitoIdentityCredentials({ IdentityPoolId: pool });
aws.config.credentials.get((err, data) => {
  if (err) {
    return;
  }

  const credentials = aws.config.credentials;
  const shadow = iot.thingShadow({
   region: region,
   clientId: 'count-client-' + (Math.floor((Math.random() * 100000) + 1)),
   protocol: 'wss',
   maximumReconnectTimeMs: 8000,
   debug: true,
   accessKeyId: credentials.accessKeyId,
   secretKey: credentials.secretAccessKey,
   sessionToken: credentials.sessionToken
  });
});
