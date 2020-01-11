package main

import (
	"errors"
	"flag"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"log"
	"net/url"
)

var (
	flagTrackMovementsHeightFrom = flag.Uint64("h", 1, "height to start scan account movements")
	flagTrackMovementsAsset      = flag.String("a", "WAVES", "asset account movements")
	flagNodeEndpoint             = flag.String("n", "https://nodes.wavesnodes.com", "node endpoint to connect with api")
)

func validateFlags() error {
	// check asset:
	_, err := proto.NewOptionalAssetFromString(*flagTrackMovementsAsset)
	if err != nil {
		return err
	}

	// check height:
	if *flagTrackMovementsHeightFrom < 1 {
		return errors.New("height must be starts from 1")
	}

	// check node:
	_, err = url.Parse(*flagNodeEndpoint)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	err := validateFlags()
	if err != nil {
		log.Fatal(err)
	}

	ms, err := NewMovements()
	if err != nil {
		log.Fatal(err)
	}

	err = ms.Parse()
	if err != nil {
		log.Fatal(err)
	}
}
