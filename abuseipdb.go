package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

// Only return reports within the last x amount of days. Default is 30.
var maxAgeInDays = "90"

// AbuseIPDB holds information about an IP address from abuseipdb.com database.
type AbuseIPDB struct {
	Data struct {
		AbuseConfidenceScore int  `json:"abuseConfidenceScore"`
		TotalReports         int  `json:"totalReports"`
		IsWhitelisted        bool `json:"isWhitelisted"`
	} `json:"data"`
}

// Check fills in AbuseIPDB data for a given IP address. It gets the data from
// api.abuseipdb.com/api/v2/check (docs.abuseipdb.com/#check-endpoint).
func (a *AbuseIPDB) Check(ipaddr net.IP) error {
	apiKey, err := getConfigValue("ABUSEIPDB_API_KEY")
	if err != nil {
		return fmt.Errorf("can't call API: %w", err)
	}

	headers := map[string]string{
		"Key":          apiKey,
		"Accept":       "application/json",
		"Content-Type": "application/x-www-form-urlencoded",
	}

	queryParams := map[string]string{
		"ipAddress":    ipaddr.String(),
		"maxAgeInDays": maxAgeInDays,
	}

	resp, err := makeAPIcall("https://api.abuseipdb.com/api/v2/check", headers, queryParams)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("calling API: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(a); err != nil {
		return err
	}

	return nil
}

// IsOK returns true if the IP address is not considered suspicious.
func (a *AbuseIPDB) IsOK() bool {
	return a.Data.TotalReports == 0 || a.Data.IsWhitelisted || a.Data.AbuseConfidenceScore <= 25
}
