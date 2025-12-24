
# Identity

 - Phone
   - Currently supports login using phone number.
   - The phone number is verified using the twilio service.

 - OAuth2
   - The [KubeRack Platform](https://github.com/kuberack/platforms/blob/main/oauth2-proxy/oauth2.md) supports an OAuth2 proxy.
   - This proxy supports extraction of information like preferred usernames and groups. Those details can then be forwarded as HTTP headers to  upstream applications.
