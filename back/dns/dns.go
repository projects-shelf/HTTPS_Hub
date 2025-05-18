package dns

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"golang.org/x/net/publicsuffix"
)

func UpsertARecord() error {
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if apiToken == "" {
		return fmt.Errorf("CLOUDFLARE_API_TOKEN environment variable required")
	}

	recordName := os.Getenv("DOMAIN") // e.g., local.example.com
	if recordName == "" {
		return fmt.Errorf("DOMAIN environment variable required")
	}

	ip := os.Getenv("LOCAL_IP") // e.g., 192.168.0.2
	if ip == "" {
		return fmt.Errorf("LOCAL_IP environment variable required")
	}

	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return fmt.Errorf("failed to create Cloudflare API client: %w", err)
	}

	zoneName, err := publicsuffix.EffectiveTLDPlusOne(recordName)
	if err != nil {
		return fmt.Errorf("failed to extract zone name: %w", err)
	}
	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return fmt.Errorf("failed to get zone ID: %w", err)
	}

	proxied := false // No proxy
	ttl := 1         // Auto

	// local.example.com & *.local.example.com
	for _, name := range []string{recordName, "*." + recordName} {
		if err := upsertA(api, zoneID, name, ip, ttl, proxied); err != nil {
			return err
		}
	}

	return nil
}

func upsertA(api *cloudflare.API, zoneID, name, ip string, ttl int, proxied bool) error {
	records, _, err := api.ListDNSRecords(context.Background(),
		cloudflare.ZoneIdentifier(zoneID),
		cloudflare.ListDNSRecordsParams{
			Type: "A",
			Name: name,
		})
	if err != nil {
		return fmt.Errorf("failed to fetch DNS records for %s: %w", name, err)
	}

	if len(records) > 0 {
		rec := records[0]
		if rec.Content == ip {
			fmt.Printf("No update needed for A record: %s\n", name)
			return nil
		}
		rec.Content = ip
		rec.Proxied = &proxied
		rec.TTL = ttl

		_, err := api.UpdateDNSRecord(context.Background(),
			cloudflare.ZoneIdentifier(zoneID),
			cloudflare.UpdateDNSRecordParams{
				ID:      rec.ID,
				Type:    "A",
				Name:    rec.Name,
				Content: rec.Content,
				TTL:     rec.TTL,
				Proxied: rec.Proxied,
			})
		if err != nil {
			return fmt.Errorf("failed to update A record %s: %w", name, err)
		}
		fmt.Printf("Updated A record: %s -> %s\n", name, ip)
		return nil
	}

	_, err = api.CreateDNSRecord(context.Background(),
		cloudflare.ZoneIdentifier(zoneID),
		cloudflare.CreateDNSRecordParams{
			Type:    "A",
			Name:    name,
			Content: ip,
			TTL:     ttl,
			Proxied: &proxied,
		})
	if err != nil {
		return fmt.Errorf("failed to create A record %s: %w", name, err)
	}
	fmt.Printf("Created A record: %s -> %s\n", name, ip)
	return nil
}
