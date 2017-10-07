package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	ga "google.golang.org/api/analyticsreporting/v4"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Load command line flags using kingpin
var (
	// debug is disabled by default
	debug = kingpin.Flag("debug", "Enable debug mode.").Short('d').Bool()
	// Check README for detailed instructions for obtaining a JSON key file
	keyfile = kingpin.Flag("keyfile", "Path to JSON key file.").Short('k').Required().String()
	// This is the Analytics view ID from which to retrieve data.
	viewID = kingpin.Flag("view-id", "Google Analytics View ID.").Short('v').Required().String()
)

func init() {
	kingpin.Parse()

	log.SetOutput(os.Stdout)
	// Default is log.InfoLevel
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	// Log command line flags _after_ setting the log level
	log.WithFields(
		log.Fields{"debug": *debug, "keyfile": *keyfile, "viewID": *viewID},
	).Debug("Parsed flags using kingpin")
}

func main() {
	defer TimeTrack(time.Now(), "Main func")
	// kingpin.Version("0.0.1") // Not required
	log.Debug("Setting up Google Analytics reporting service")

	svc, err := makeReportSvc()

	if err != nil {
		log.WithError(err).Panic("Error while creating Google Analytics Reporting Service")
	}

	res, err := getReport(svc)

	if err != nil {
		log.WithError(err).Panic("GET request to analyticsreporting/v4 returned error")
	}

	if res.HTTPStatusCode != 200 {
		log.WithField(
			"HTTPStatusCode", res.HTTPStatusCode,
		).Panic("Did not get expected HTTP response code")
	}

	log.Info("Got response from analytics reporting")

	printResponse(res)
}

// makeReportSvc initializes and returns an authorized
// Analytics Reporting API V4 service object.
func makeReportSvc() (*ga.Service, error) {
	defer TimeTrack(time.Now(), "Make reporting service")
	// Your credentials should be obtained from the Google
	// Developer Console (https://console.developers.google.com).
	// Navigate to your project, then see the "Credentials" page
	// under "APIs & Auth".
	// To create a service account client, click "Create new Client ID",
	// select "Service Account", and click "Create Client ID". A JSON
	// key file will then be downloaded to your computer.
	data, err := ioutil.ReadFile(*keyfile)

	if err != nil {
		log.Error("Failed to load credentials for Google Analytics")
		return nil, err
	}

	log.WithField("keyfile", keyfile).Debug("Read key file")

	conf, err := google.JWTConfigFromJSON(data, ga.AnalyticsReadonlyScope)

	if err != nil {
		log.Error("Failed to create JWT config from JSON creds")
		return nil, err
	}

	log.Debug("Created jwt config")

	// Initiate an http.Client. The following GET request will be
	// authorized and authenticated on the behalf of
	// your service account.
	var netClient *http.Client
	if *debug {
		ctx := context.WithValue(
			context.Background(),
			oauth2.HTTPClient,
			&http.Client{Transport: &logTransport{http.DefaultTransport}},
		)
		netClient = conf.Client(ctx)
	} else {
		netClient = conf.Client(oauth2.NoContext)
	}

	log.Debug("Created authentication capable HTTP client")

	// Construct the Analytics Reporting service object.
	svc, err := ga.New(netClient)

	if err != nil {
		log.Error("Failed to create Google Analytics Reporting Service")
		return nil, err
	}

	log.Info("Created Google Analytics Reporting Service object")

	return svc, nil
}

// getReport queries the Analytics Reporting API V4 using
// the Analytics Reporting API V4 service object.
// It returns the Analytics Reporting API V4 response
func getReport(svc *ga.Service) (*ga.GetReportsResponse, error) {
	defer TimeTrack(time.Now(), "GET Analytics Report")
	// A GetReportsRequest instance is a batch request
	// which can have a maximum of 5 requests
	req := &ga.GetReportsRequest{
		// Our request contains only one request
		// So initialise the slice with one ga.ReportRequest object
		ReportRequests: []*ga.ReportRequest{
			// Create the ReportRequest object.
			{
				ViewId: *viewID,
				DateRanges: []*ga.DateRange{
					// Create the DateRange object.
					{StartDate: "7daysAgo", EndDate: "today"},
				},
				Metrics: []*ga.Metric{
					// Create the Metrics object.
					{Expression: "ga:sessions"},
				},
				Dimensions: []*ga.Dimension{
					{Name: "ga:country"},
				},
			},
		},
	}

	log.Info("Doing GET request from analytics reporting")
	// Call the BatchGet method and return the response.
	return svc.Reports.BatchGet(req).Do()
}

// printResponse parses and prints the Analytics Reporting API V4 response.
func printResponse(res *ga.GetReportsResponse) {
	defer TimeTrack(time.Now(), "Print Analytics Report")
	log.Info("Printing Response from analytics reporting")
	for _, report := range res.Reports {
		header := report.ColumnHeader
		dimHdrs := header.Dimensions
		metricHdrs := header.MetricHeader.MetricHeaderEntries
		rows := report.Data.Rows

		if rows == nil {
			log.WithField("viewID", viewID).Info("No data found for given view.")
		}

		for _, row := range rows {
			dims := row.Dimensions
			metrics := row.Metrics

			for i := 0; i < len(dimHdrs) && i < len(dims); i++ {
				log.Infof("%s: %s", dimHdrs[i], dims[i])
			}

			for _, metric := range metrics {
				// We have only 1 date range in the example
				// So it'll always print "Date Range (0)"
				// log.Infof("Date Range (%d)", idx)
				for j := 0; j < len(metricHdrs) && j < len(metric.Values); j++ {
					log.Infof("%s: %s", metricHdrs[j].Name, metric.Values[j])
				}
			}
		}
	}
	log.Info("Completed printing response")
}
