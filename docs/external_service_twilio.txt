

# Twilio Service

Ref: 
  - [iam api](https://www.twilio.com/docs/iam/api)
  - [verification api](https://github.com/twilio/twilio-go?tab=readme-ov-file#test-your-installation)


  - Base url
    - https://api.twilio.com/2010-04-01

  - Authentication
    - Twilio supports HTTP Basic authentication
    - Create an API key either in the Twilio Console
    - The API key is the username and the API key secret is the password.
    - Twilio recommends using API keys for authentication in production apps. 
    - For local testing, you can use Account SID as the username and your
      Auth token as the password. 
    - You can find your Account SID and Auth Token in the Twilio Console
      - SID       : []
      - Auth Token: []
    - Test out the key, and the key secret
      ```
            curl -G https://api.twilio.com/2010-04-01/Accounts -u $TWILIO_API_KEY:$TWILIO_API_KEY_SECRET
      ```
  - Verification
    - Essential, the user supplies a phone number on which to receive an OTP
    - The twilio go client is then invoked to create a verification object
    ([link](https://www.twilio.com/docs/verify/api/verification))
    - After some delay the user will receive an OTP
    - This OTP is then passed to the verfication check method of the twilio 
      client.([link](https://www.twilio.com/docs/verify/api/verification-check))

