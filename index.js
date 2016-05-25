const aws = require('aws-sdk');
const iot = require('aws-iot-device-sdk');

const region = 'us-east-1';
const pool = 'us-east-1:594b1376-9d0e-4adb-9f91-1a3fae53d247';
const counterId = 'counter';

let shadow;

aws.config.region = region;
aws.config.credentials = new aws.CognitoIdentityCredentials({ IdentityPoolId: pool });
aws.config.credentials.get((err, data) => {
  if (err) {
    return;
  }

  shadow = iot.thingShadow({
    region: region,
    clientId: 'count-client-' + (Math.floor((Math.random() * 100000) + 1)),
    protocol: 'wss',
    accessKeyId: aws.config.credentials.accessKeyId,
    secretKey: aws.config.credentials.secretAccessKey,
    sessionToken: aws.config.credentials.sessionToken
  });

  shadow.on('connect', onShadowConnected);
  shadow.on('status', onShadowStatus);
  shadow.on('delta', onShadowDelta);
});

const onShadowConnected = () => {
  shadow.register(counterId, { persistentSubscribe: true });
  setTimeout(() => shadow.get(counterId), 1000);
};

const onShadowStatus = (name, status, token, data) => {
  console.log('Received IoT status update!');
  document.body.classList.remove('loading');
  document.getElementById('count').textContent = data.state.desired.count;
};

const onShadowDelta = (name, data) => {
  console.log('Received IoT state delta!');
  document.getElementById('count').textContent = data.state.count;
};
