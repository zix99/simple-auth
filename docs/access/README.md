# Access Layer

The access layer allows users to access content via *simple-auth*.

This primarily takes two forms:
* [Same-domain cookie](cookie) is when *simple-auth* acts as a login-provider, and a downstream service validates the user via signed cookie or other API mechanism
* [Gateway](gateway) is when *simple-auth* sits between the user and what they're trying to access as a portal

Alternatively, you can use [simple-auth's authenticators](../authenticators) to enable things like nginx authentication
