#  AWS SES Gateway


[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Docker](https://github.com/askrella/whatsapp-chatgpt/actions/workflows/docker.yml/badge.svg)
![Docker AMD64](https://img.shields.io/badge/docker-amd64-blue)
![Docker ARM64](https://img.shields.io/badge/docker-arm64-green)
![Build](https://img.shields.io/github/actions/workflow/status/netologist/aws-ses-gateway/docker.yml?branch=main)


We created this project as a new version of aws-ses-local, which doesn't seem to be maintained for a few years.
Our goal is to provide more features, small containers and be more accurate than the alternatives.

# :gear: Getting Started

## Running the Docker Container

```bash
docker run -p 8081:8081 ghcr.io/netologist/aws-ses-gateway:1.0.0
```

## Usage with NodeJS


<details>
<summary>JavaScript/TypeScript for the V2 API with the V3 SDK (recommended)</summary>

```typescript
import { SESv2Client, SendEmailCommand } from "@aws-sdk/client-sesv2"

const ses = new SESv2Client({
    endpoint: 'http://localhost:8005',
    region: 'aws-ses-v2-local',
    credentials: { accessKeyId: 'ANY_STRING', secretAccessKey: 'ANY_STRING' },
});
await ses.send(new SendEmailCommand({
    FromEmailAddress: 'sender@example.com',
    Destination: { ToAddresses: ['receiver@example.com'] },
    Content: {
        Simple: {
            Subject: { Data: 'This is the subject' },
            Body: { Text: { Data: 'This is the email contents' } },
        }
    },
}))
```

</details>

<details>
<summary>JavaScript/TypeScript for the V1 API with the V3 SDK</summary>

```typescript
import { SES, SendEmailCommand } from '@aws-sdk/client-ses'

const ses = new SES({
    endpoint: 'http://localhost:8005',
    region: 'aws-ses-v2-local',
    credentials: { accessKeyId: 'ANY_STRING', secretAccessKey: 'ANY_STRING' },
})
await ses.send(new SendEmailCommand({
    Source: 'sender@example.com',
    Destination: { ToAddresses: ['receiver@example.com'] },
    Message: {
        Subject: { Data: 'This is the subject' },
        Body: { Text: { Data: 'This is the email contents' } },
    },
}))
```

</details>

<details>
<summary>JavaScript/TypeScript for the V2 API with the V2 SDK</summary>

```typescript
import AWS from 'aws-sdk'

const ses = new AWS.SESV2({
    endpoint: 'http://localhost:8005',
    region: 'aws-ses-v2-local',
    credentials: { accessKeyId: 'ANY_STRING', secretAccessKey: 'ANY_STRING' },
})
ses.sendEmail({
    FromEmailAddress: 'sender@example.com',
    Destination: { ToAddresses: ['receiver@example.com'] },
    Content: {
        Simple: {
            Subject: { Data: 'This is the subject' },
            Body: { Text: { Data: 'This is the email contents' } },
        }
    },
})
```

</details>

<details>
<summary>JavaScript/TypeScript with nodemailer for the V1 raw API with the V3 SDK</summary>

```typescript
import * as aws from '@aws-sdk/client-ses'

const ses = new aws.SES({
    endpoint: 'http://localhost:8005',
    region: 'aws-ses-v2-local',
    credentials: { accessKeyId: 'ANY_STRING', secretAccessKey: 'ANY_STRING' },
})
const transporter = nodemailer.createTransport({ SES: { ses, aws } })

await transporter.sendMail({
    from: 'sender@example.com',
    to: ['receiver@example.com'],
    subject: 'This is the subject',
    text: 'This is the email contents',
    attachments: [{
        filename: `some-file.pdf`,
        contentType: 'application/pdf',
        content: Buffer.from(pdfBytes),
    }],
})
```

</details>

Using another language or version? Submit a PR to update this list :)

## Manual testing

```
docker-compose up --build
```

# :warning: License
Distributed under the MIT License. See LICENSE.txt for more information.

# :handshake: Contact Us

In case you need professional support, feel free to <a href="mailto:contact@netologist.org">contact us</a>
