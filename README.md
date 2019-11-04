# dev-router

A router for debugging and developing web apps.

Start a HTTP/HTTPS server and it will modify the requests your send, and the responses you get, according to the rule you specified.

Use cases:
- You only want to replace some static JS and CSS files with local files for an existing website.
- You want to direct one set of requests to site A, and another set of requests to site B.
- You want to manually delay some request by a certain amount of seconds.
- And many more ...

Why?
- Yes, similar (or all) functions can be achieved using other tools (e.g. Charles, Postman, etc.).
- But the flexibility is best with custom code in a proper programming language like this dev-router.
- Free. Plus, what could beat a great learning opportunity by hand-crafting your own router?

## Build, config and run

For example:

Get a released bin from [the releases](https://github.com/bradleyzhou/dev-router/releases), or build your own by `go install`.

Then:
```bash
# for a http server
dev-router -conf=example/config.json

# or if you need a https/h2 server
dev-router -conf=example/config-https.json
```

Then visit http://localhost:8081 or https://localhost:8082 to see the effect.

If you are getting certificate errors, you might want to trust the certificate you prepared(self-signed), or just get a valid certificate. Just search the web for the steps to get a valid certificate or trust the self-signed certificate.

## Custom rules

You can build on top of the example config.json.
