// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/context"

	"istio.io/fortio/fhttp"
	"istio.io/fortio/periodic"
)

const (
	fetchInterval   = 5
)

var (
	webAddress      = flag.String("listen_port", ":9103", "Port on which to expose metrics and web interface.")

	metricsSuite *p8sMetricsSuite

	gcsClient  *storage.Client
	httpClient = &http.Client{}

)



type p8sMetricsSuite struct {
	succeededBuilds   *prometheus.SummaryVec
	failedBuilds      *prometheus.SummaryVec
	succeededBuildMax *prometheus.GaugeVec
	succeededBuildMin *prometheus.GaugeVec
	codeCoverage      *prometheus.GaugeVec
}





func newP8sMetricsSuite() *p8sMetricsSuite {
	return &p8sMetricsSuite{
		succeededBuilds: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "succeeded_build_durations_seconds",
				Help: "Succeeded build durations seconds.",
			},
			[]string{"build_job", "repo"},
		),


	}
}

func init() {
	metricsSuite = newP8sMetricsSuite()
	metricsSuite.registerMetricVec()
}

func main() {
	flag.Parse()

	var err error
	gcsClient, err = storage.NewClient(context.Background())
	if err != nil {
		log.Printf("Failed to create a gcs client, %v", err)
	}

	go func() {
		for {

			time.Sleep(time.Duration(fetchInterval) * time.Minute)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Print(http.ListenAndServe(*webAddress, nil))

}

func (m *p8sMetricsSuite) registerMetricVec() {
	prometheus.MustRegister(m.succeededBuilds)
	prometheus.MustRegister(m.succeededBuildMax)
	prometheus.MustRegister(m.succeededBuildMin)
	prometheus.MustRegister(m.failedBuilds)
	prometheus.MustRegister(m.codeCoverage)
}



func (j *job) publishCIMetrics() error {

metricsSuite.failedBuilds.WithLabelValues(j.jobName, j.repoName).Observe(t)

metricsSuite.succeededBuildMin.WithLabelValues(j.jobName, j.repoName).Set(min)

}

func load() error {
	url := "fortio.release02.istio.webinf.info"
	opts := fhttp.HTTPRunnerOptions{
		RunnerOptions: periodic.RunnerOptions{
			QPS:        10,
			Duration:   1 * time.Minute,
			NumThreads: 8,
		},
		URL: url,
	}

	res, err := fhttp.RunHTTPTest(&opts)
	if err != nil {
		fatalf(t, "Generating traffic via fortio failed: %v", err)
	}

}

func record





