/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	//"flag"
	"log"
	"os"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/k8snetworkplumbingwg/multi-networkpolicy-iptables/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
)

const logFlushFreqFlagName = "log-flush-frequency"

var logFlushFreq = pflag.Duration(logFlushFreqFlagName, 5*time.Second, "Maximum number of seconds between log flushes")

// KlogWriter serves as a bridge between the standard log package and the glog package.
type KlogWriter struct{}

// Write implements the io.Writer interface.
func (writer KlogWriter) Write(data []byte) (n int, err error) {
	klog.InfoDepth(1, string(data))
	return len(data), nil
}

func initLogs() {
	log.SetOutput(KlogWriter{})
	log.SetFlags(0)
	go wait.Forever(klog.Flush, *logFlushFreq)
}

func main() {
	go func() {log.Println(http.ListenAndServe("0.0.0.0:6060", nil)) }()
	initLogs()
	defer klog.Flush()
	opts := server.NewOptions()

	cmd := &cobra.Command{
		Use:  "multi-networkpolicy-node",
		Long: `TBD`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := opts.Run(); err != nil {
				klog.Exit(err)
			}
		},
	}
	opts.AddFlags(cmd.Flags())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
