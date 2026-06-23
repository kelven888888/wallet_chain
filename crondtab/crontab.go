package crondtab

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func Initcrond() {

	s := gocron.NewScheduler(time.UTC)
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("Error loading location: %s", err)
	}
	s.ChangeLocation(location)

	s.Every(30).Second().Do(iptoaddr)

	//s.Every(1).Day().At("16:00").Do(Initquanaccount)

	s.StartAsync()

}

func WalletAddCreate() {
	var wallet WalletServer
	wallet.CrondCreateAdd()

}
func WalletWithdraw() {
	var wallet WalletServer
	wallet.CrondWithPassToUdun()

}
