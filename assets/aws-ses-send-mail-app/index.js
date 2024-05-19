const { SES, SendEmailCommand }  = require('@aws-sdk/client-ses');

const ses = new SES({
    endpoint: process.env.AWS_SES_ENDPOINT,  
    region: 'aws-ses-v2-local',
    credentials: { accessKeyId: 'ANY_STRING', secretAccessKey: 'ANY_STRING' },
})

ses.send(new SendEmailCommand({
    Source: 'sender@example.com',
    Destination: { ToAddresses: ['receiver@example.com'] },
    Message: {
        Subject: { Data: 'This is the subject' },
        Body: { Text: { Data: 'This is the email contents' } },
    },
})).then(x => console.log(x));
