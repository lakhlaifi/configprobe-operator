package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

func (r *ConfigProbeReconciler) StartHTTPServer() {
	log := ctrl.Log.WithName("controllers").WithName("ConfigProbe")

	port := ":8082"
	probe := "/probe/"
	log.Info("Starting server", "path", probe, "kind", "probe", "port", port)

	http.HandleFunc(probe, r.handleHTTPRequest)
	http.ListenAndServe(port, nil)
}

func (r *ConfigProbeReconciler) handleHTTPRequest(w http.ResponseWriter, req *http.Request) {

	log := ctrl.Log.WithName("controllers").WithName("ConfigProbe")

	// Extract namespace from the path
	pathParts := strings.Split(req.URL.Path, "/")
	if len(pathParts) != 3 || pathParts[1] != "probe" {
		http.Error(w, "Invalid URL format. Expected /probe/<namespace>", http.StatusBadRequest)
		return
	}
	// namespace := pathParts[2]

	// Parse the query parameters
	target := req.URL.Query().Get("target")
	module := req.URL.Query().Get("module")

	// Check if target and module are provided
	if target == "" || module == "" {
		http.Error(w, "Both target and module query parameters are required", http.StatusBadRequest)
		return
	}

	// Send a GET request to another server
	response, err := http.Get(fmt.Sprintf("http://blackbox-westeurope--4m2tpov.wonderfulforest-14e15c11.westeurope.azurecontainerapps.io/probe?target=%s&module=%s", target, module))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send request to other server: %v", err), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Check the status code of the response
	if response.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Received non-OK response from other server: %s", response.Status), http.StatusInternalServerError)
		return
	}

	// Read the response from the other server
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response from other server: %v", err), http.StatusInternalServerError)
		return
	}

	// Send the response back to the original requester
	w.WriteHeader(http.StatusOK)
	log.Info("Request processed successfully", "status", http.StatusOK, "target", target, "module", module)
	w.Write(responseBody)
}
