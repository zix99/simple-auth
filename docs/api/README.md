# REST API

::: warning
The REST API can give you access to all your users and accounts, which means you can pretty much do anything
to anyone.  If you use shared-secret access, guard the secret carefully.
:::

[[toc]]

The *simple-auth* REST API allows you to do everything (and more) that the UI allows a user to do.  Internally,
when a user logs-in to the UI, it gives them access to REST API (but of course, only to access their specific
details).

A second way to access the REST API is via **shared-secret** credentials.  This gives the holder
of the secret access to do anything to any account credentials.

## Accessing the API

Accessing the API is as simple as making an HTTP request to the resource, with two headers:

* `Authorization` Add your `SharedKey` to the authorization header
* `X-Account-UUID` will specify which account resource you're operating on (for endpoints that require it)

For example, your request might look something like this:

```http
GET /api/v1/account HTTP/1.1
Authorization: SharedKey my-shared-secret
X-Account-UUID: c270e7e0-47a2-11eb-b378-0242ac130002
```

### Example API Call

Here's an example of getting details of a user's account via `/api/v1/account` endpoint.

<<< @/examples/rest-api/getAccount.js

## Full API Docs

You can read more about the exposed API calls in the <a :href="`${$themeConfig.docsUrl}/apidocs`">API Documentation</a>
