# short url generator

the core are [hashids](https://hashids.org/) and [go-hashids](https://github.com/speps/go-hashids)

# Deployment

### https://hub.docker.com/r/uvaa/short-url

Docker deployment is recommended, the simple command:

    docker run -d -p 8080:80 -v /data:/app/data -e origin=https://abcd.com uvaa/short-url

**_NOTE: origin must be defined_**

# Environment Variables

When you start the image, you can adjust the configuration of the instance by passing one or more environment variables on the docker run command line

### `origin`

This variable is mandatory and specifies the `public orgin` that will be set for the return `short url`, via `https://abc.com`

### `salt`

Optional variable. If not set, the program will automatically generate a random string when it is first started and stored in the database

### `minlength`

Optional variable. Use to define the minimum length of an index string to generate short url. `Default is 5`

**_NOTE: `salt` & `minlength` If once generated, they cannot be modified_**

# Use it

Be sure to use the `POST` request to get `Short URL`

    curl -X POST 'https://abc.com?url=https://www.google.com'

    {"msg":"success","succ":true,"url":"https://abc.com/xpnOk"}
