# example-analyticsreporting
Example for Google Analytics Reporting API v4 in golang using the [official client](https://github.com/google/google-api-go-client).

This example was written using the official examples for [java](https://developers.google.com/analytics/devguides/reporting/core/v4/quickstart/service-java), [python](https://developers.google.com/analytics/devguides/reporting/core/v4/quickstart/service-py)

## How to run
  * Enable the API - Using steps described in https://developers.google.com/analytics/devguides/reporting/core/v4/quickstart/service-java,
    create a private key `client_secrets.json` and download it to your file system.
    Also add the new service account to the Google Analytics Account with [Read & Analyze](https://support.google.com/analytics/answer/2884495) permission.
  * Clone the repo and install dependencies. [Glide](https://github.com/Masterminds/glide) is preferred.
    ```
    $ mkdir -p $GOPATH/src/github.com/sudo-suhas && cd $_
    $ git clone https://github.com/sudo-suhas/example-analyticsreporting.git
    $ cd example-analyticsreporting

    # preferred method
    $ glide install

    # or use `go get`
    $ go get .

    ```
  * Build the executable binary.
    ```
    # This should generate an executable binary in the current folder
    # Example example-analyticsreporting.exe on windows
    $ go build .

    ```
  * You need to pass the location of the `client_secrets.json` file
    and the Google analytics view ID from the command line.
    You can use the [Account Explorer](https://ga-dev-tools.appspot.com/account-explorer/) to find a View ID.
    Additionally, you can pass the flag `--debug` for verbose logging.

    Usage:
    ```
    $ ./example-analyticsreporting.exe --help
    usage: hello_analytics.exe --keyfile=KEYFILE --view-id=VIEW-ID [<flags>]

    Flags:
          --help             Show context-sensitive help (also try --help-long and
                            --help-man).
      -d, --debug            Enable debug mode.
      -k, --keyfile=KEYFILE  Path to JSON key file.
      -v, --view-id=VIEW-ID  Google Analytics View ID.

    $ ./example-analyticsreporting.exe --keyfile=E:\creds\client_secrets.json -view-id=299792458 --debug

    ```

## File Structure
  * `hello_analytics.go` - This is the main file. It does the following:
    - Create a Google Analytics Reporting API v4 service client
    - Execute a GET analytics report request
    - Parse and print the response using [`logrus`](https://github.com/Sirupsen/logrus/blob/master/logrus.go).
  * `debug.go` - This is copied from https://github.com/google/google-api-go-client/blob/master/examples/debug.go.
    It is used to log the HTTP request and response to `os.Stdout` in debug mode.
  * `util.go` - This has a simple utility function for tracking function execution time in debug mode.
