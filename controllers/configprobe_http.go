package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *ConfigProbeReconciler) StartHTTPServer() {
	log := ctrl.Log.WithName("controllers").WithName("ConfigProbe")

	port := ":8082"
	probe := "/probe/"
	log.Info("Starting server", "path", probe, "kind", "probe", "port", port)

	http.HandleFunc(probe, r.handleHTTPRequest)
	http.ListenAndServe(port, nil)
}

func getBboxForRegion(region string) (string, error) {
	data, err := ioutil.ReadFile("/config/blackbox/bbox_urls")
	if err != nil {
		return "BBOX DataSource: ", err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		if len(parts) == 2 && parts[0] == region {
			return parts[1], nil
		}
	}
	return "", fmt.Errorf("Region not found")
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
	region := req.URL.Query().Get("region")

	// Check if target and module are provided
	if target == "" || module == "" {
		http.Error(w, "Both target and module query parameters are required", http.StatusBadRequest)
		return
	}

	// /config/blackbox/bbox_urls
	// Get the bbox for the region
	bbox, err := getBboxForRegion(region)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get host for region: %v", err), http.StatusInternalServerError)
		return
	}
	response, err := http.Get(fmt.Sprintf("%s/probe?target=%s&module=%s", bbox, target, module))

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
