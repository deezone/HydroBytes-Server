package main

import (
	// Standard packages
	"context"       // https://golang.org/pkg/context/
	"encoding/json" // https://golang.org/pkg/encoding/json/
	"log"           // https://golang.org/pkg/log/
	"net/http"      // https://golang.org/pkg/net/http/
	"os"
	"os/signal"     // https://golang.org/src/os/signal/doc.go
	"syscall"
	"time"
)

/**
 * StationTypes is a type of station in the automated garden system.
 *
 * Note: use of "stuct tags" (ex: `json:"id"`) to manage the names of properties to be lowercase and snake_case. Due to
 * the use of case for visibility in Go "id" rather and "Id" would result in the value being excluded in the JSON
 * response as the encoding/json package is external to this package.
 */
type StationTypes struct {
	Id			int    `json:"id"`
	Name		string `json:"name"`
	Description	string `json:"description"`
}

// Main entry point for program.
func main() {

	// Define constants that includes their type to ensure a "kind" decarlation is not used resulting in much
	// larger memory space use. See https://education.ardanlabs.com/courses/take/ultimate-syntax/lessons/13570526-constants-pt-2
	const (
		readTimeout time.Duration   = 5 * time.Second
		writeTimeout time.Duration  = 5 * time.Second
		cancelTimeout time.Duration = 5 * time.Second
	)

	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	// =========================================================================
	// Start API Service

	/**
	 * Convert the ListStationTypes function to a type that implements http.Handler
	 * See https://education.ardanlabs.com/courses/take/ultimate-syntax/lessons/13570357-type-conversions for details
	 * on "types" and conversions.
	 *
	 * Details on using HandlerFunc which is an adapter
	 * "to allow the use of ordinary functions as HTTP handlers"
	 * https://golang.org/pkg/net/http/#HandlerFunc
	 *
	 * Note this is an implementation of the same functionality that http.ListenAndServe() provides but is made
	 * available in a channel.
	 * https://golang.org/src/net/http/server.go?s=97511:97566#L3108
	 */
	api := http.Server{
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(ListStationTypes),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a buffered channel so the goroutine can exit
	// if we don't collect this error.
	//
	serverErrors := make(chan error, 1)

	// Start a server listening ("bind") on port 8000 and responding using handler called ListStationTypes()
	// https://golang.org/pkg/net/http/#Server.ListenAndServe
	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	/**
	 * Make a channel to listen for an interrupt or terminate signal from the OS. Use a buffered channel because the
	 * signal package requires it. Note the value of "1", this limits the capacity to one to prevent processes from
	 * staying alive if more than one thread is added to the channel which would prevent termination.
	 *
	 * Listening for os.Interrupt and syscall.SIGTERM events.
	 */
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	// Note two active channels: serverErrors and shutdown as defined above
	select {
	case err := <-serverErrors:
		log.Fatalf("error: listening and serving: %s", err)

	case <-shutdown:
		log.Println("main : Start shutdown")

		/**
		 * Give outstanding requests a deadline for completion. Context defines the Context type, which carries
		 * deadlines, cancellation signals, and other request-scoped values across API boundaries and between
		 * processes.
		 * https://golang.org/pkg/context/
		 */
		ctx, cancel := context.WithTimeout(context.Background(), cancelTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cancelTimeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main : could not stop server gracefully : %v", err)
		}
	}
}

/**
 * ListStationTypes is a basic HTTP Handler that lists all of the station types in the HydroByte Automated Garden system.
 * A Handler is also know as a "controller" in the MVC pattern.
 *
 * The parameters must follow the signature defined in the http.HandlerFunc() adapter
 * (https://golang.org/src/net/http/server.go?s=97511:97566#L2034) to convert this method into a HTTP handler type.
 *
 * Note: If you open localhost:8000 in your browser, you may notice double requests being made. This happens because
 * the browser sends a request in the background for a website favicon. More the reason to use Postman to test!
 */
func ListStationTypes(w http.ResponseWriter, r *http.Request) {
	list := []StationTypes{
		{Id: 1, Name: "Base", Description: "Coordinator for all station types - monitor, command and control. Access point to public Intenet."},
		{Id: 2, Name: "Water", Description: "Management of water resources. Controls water levels in resavour and impliments irrigation."},
		{Id: 3, Name: "Plant", Description: "Monitors and reports plant health."},
	}

	// https://golang.org/pkg/encoding/json/#Marshal
	data, err := json.Marshal(list)
	if err != nil {
		log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// https://golang.org/pkg/net/http/#Request.Write
	if _, err := w.Write(data); err != nil {
		log.Println("error writing result", err)
	}
}
