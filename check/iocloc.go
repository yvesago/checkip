package check

import (
	"context"
	"encoding/json"
	"errors"
	"os/user"
	"path/filepath"

	"fmt"
	//"log"
	"net"
	"time"

	"github.com/logrusorgru/aurora"
	//"github.com/mitchellh/go-homedir"
	"github.com/oschwald/geoip2-golang"
)

type loc struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	Name    string `json:"name"`
	ASN     string `json:"asn"`
}

func (l loc) Json() ([]byte, error) {
	return json.Marshal(l)
}

func (l loc) Summary() string {
	au := aurora.NewAurora(true)
	flag, _ := emocode(l.Country)
	//return fmt.Sprintf("%s%s (%s)%s %s", l.IP, l.Name, au.Green(l.Country), flag, l.ASN)
	//fmt.Printf("%s%s (%s)%s %s", l.IP, l.Name, au.Green(l.Country), flag, l.ASN)
	return fmt.Sprintf("%s (%s)%s %s", l.IP, au.Green(l.Country), flag, l.ASN)
}

func emocode(x string) (string, error) {
	if len(x) != 2 {
		return "", errors.New("country code must be two letters")
	}
	if x[0] < 'A' || x[0] > 'Z' || x[1] < 'A' || x[1] > 'Z' {
		return "", errors.New("invalid country code")
	}
	return string(0x1F1E6+rune(x[0])-'A') + string(0x1F1E6+rune(x[1])-'A'), nil
}

/*func getConfigConfDir(path string) string {
	d, e := homedir.Expand(path)
	if e != nil {
		log.Fatalf("failed to get home dir: error=%v", e)
	}
	return d
}*/

// IOCLoc print my loc format
func IOCLoc(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "IOCLoc",
		Type: TypeInfo,
	}

	timeout := time.Duration(500) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()
	var r net.Resolver

	l := loc{IP: ipaddr.String()}
	usr, err := user.Current()
    if err != nil {
        return Check{}, err
    }
    ConfDir := filepath.Join(usr.HomeDir, ".scan/")
	geoip := false
	DbCity, err := geoip2.Open(ConfDir + "/GeoLite2-City.mmdb")
	var DbASN *geoip2.Reader
	if err == nil {
		geoip = true
	} else {
		fmt.Println("Warning:", fmt.Sprintf("missing geoip dbs in «%s»\n\n", ConfDir))
	}
	if geoip {
		DbASN, _ = geoip2.Open(ConfDir + "/GeoLite2-ASN.mmdb")
	}

	rec, err := DbCity.City(ipaddr)
	l.Country = rec.Country.IsoCode

	name := ""
	names, err := r.LookupAddr(ctx, ipaddr.String())
	if err == nil && len(names) > 0 {
		name = " (" + names[0] + ")"
	}
	l.Name = name

	asn, _ := DbASN.ASN(ipaddr)
	l.ASN = fmt.Sprintf(" AS%d", asn.AutonomousSystemNumber)
	l.ASN += " - " + asn.AutonomousSystemOrganization

	result.IpAddrInfo = l

	return result, nil
}
