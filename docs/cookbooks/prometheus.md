# Prometheus

*Simple-auth* supports [Prometheus](https://prometheus.io/) monitoring built in, and exposes a `/metrics` endpoint (if enabled).

## Enabling

To enable, set `web.prometheus` to `true`.

## Metrics

In addition to the default metrics exposed by [golang client](https://prometheus.io/docs/guides/go-application/), *simple-auth* exposes the following metrics:

### sa_requests_total{code,host,method,url}

Request total, labeled by response, method, and url

```
# HELP sa_requests_total How many HTTP requests processed, partitioned by status code and HTTP method.
# TYPE sa_requests_total counter
sa_requests_total{code="200",host="localhost:9002",method="GET",url="/"} 2
```

### sa_email_sends{template,success}

Email sends, with `template` and `success` (true/false) labels

```
# HELP sa_email_sends Email sending metrics
# TYPE sa_email_sends counter
sa_email_sends{success="true",template="welcome"} 1
```

### sa_auth{type,success,errCode}

Keeps track of all authenticators (eg. vouch, simple API, or OAuth2), with `type` (source of auth), `success` (true/false), `errCode`

```
# HELP sa_auth Authentication counter
# TYPE sa_auth counter
sa_auth{errCode="nil",success="true",type="oauth2:code"} 1
sa_auth{errCode="nil",success="true",type="oauth2:token"} 1
```

### sa_local_login{success}

Local login counter with `success` label

```
# HELP sa_local_login Counter for local login
# TYPE sa_local_login counter
sa_local_login{success="true"} 1
```

### sa_session_create{source}

Count how many UI sessions have been issued

```
# HELP sa_session_create Session creation counter
# TYPE sa_session_create counter
sa_session_create{source="login"} 1
sa_session_create{source="oidc"} 1
```

### sa_account_create{}

Count account creations

```
# HELP sa_account_create Account creation counter
# TYPE sa_account_create counter
sa_account_create 1
```
