![official ioki Hackday project](https://img.shields.io/badge/official-ioki%20Hackday%20project-%23000)

# TicTxtToe

An Tic Tac Toe client which can be played over SMS.

https://user-images.githubusercontent.com/64888250/192966997-e1fa0aa6-2e28-42bd-ab32-39c38c993997.mp4

The project was created at the [Hackday 2022](https://stefma.medium.com/a1d14341e3f2).

## Setup

We use [twilio](https://www.twilio.com/) to send the SMS.
Therefore you need to create a twilio account and set up the ["Mesage/SMS"](https://www.twilio.com/messaging/sms) feature.

> **Note**: For testing purpose you can use the trial version of twilio.

Afterwards you have to paste the credentials in a file named `config.json`:
```json
{
    "username": "TWILIO_USERNAME",
    "password": "TWILIO_PASSOWRD",
    "phonenumber": "TWILIO_PHONENUMBER"
}
```

Now you can run the server:
```
go run main.go
```

You may want to use [ngrok](https://ngrok.com/) to expose the them to the internet.

The URL has to be set as a [SMS webhook for your phonenumber](https://www.twilio.com/docs/serverless/functions-assets/quickstart/receive-sms#set-a-function-as-a-webhook) in the twilio console to receive the SMS someone sends to the twilio phonenumber.

> **Note**: In case you use the trial account you have to verify all phonenumbers who sends a SMS to your twilio phonenumber first. See also the [twilio docs](https://support.twilio.com/hc/en-us/articles/223180048-Adding-a-Verified-Phone-Number-or-Caller-ID-with-Twilio) for this.
