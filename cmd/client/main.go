package main

import (
	"fmt"
	"os"

	"github.com/bryant-rh/srew/cmd/client/cmd"
	"k8s.io/klog/v2"
)

func main() {
	defer klog.Flush()
	cmd := cmd.NewCmd()
	if err := cmd.Execute(); err != nil {
		if klog.V(1).Enabled() {
			klog.Fatalf("%+v", err) // with stack trace
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
