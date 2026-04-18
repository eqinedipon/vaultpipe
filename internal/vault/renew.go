package vault

import (
	"context"
	"log"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Renewer manages periodic renewal of a Vault token.
type Renewer struct {
	client *vaultapi.Client
	logger *log.Logger
}

// NewRenewer creates a Renewer for the given Vault API client.
func NewRenewer(client *vaultapi.Client, logger *log.Logger) *Renewer {
	return &Renewer{client: client, logger: logger}
}

// Start begins renewing the token in the background until ctx is cancelled.
// It looks up the current token's TTL and renews at half-life intervals.
func (r *Renewer) Start(ctx context.Context) {
	go func() {
		for {
			ttl, err := r.tokenTTL()
			if err != nil {
				r.logger.Printf("vault/renew: failed to look up token TTL: %v", err)
				return
			}
			if ttl <= 0 {
				// Non-expiring token; nothing to renew.
				return
			}
			sleep := ttl / 2
			select {
			case <-ctx.Done():
				return
			case <-time.After(sleep):
				if err := r.renew(); err != nil {
					r.logger.Printf("vault/renew: token renewal failed: %v", err)
					return
				}
				r.logger.Printf("vault/renew: token renewed successfully")
			}
		}
	}()
}

func (r *Renewer) tokenTTL() (time.Duration, error) {
	secret, err := r.client.Auth().Token().LookupSelf()
	if err != nil {
		return 0, err
	}
	ttlRaw, ok := secret.Data["ttl"]
	if !ok {
		return 0, nil
	}
	ttlJSON, ok := ttlRaw.(float64)
	if !ok {
		return 0, nil
	}
	return time.Duration(ttlJSON) * time.Second, nil
}

func (r *Renewer) renew() error {
	_, err := r.client.Auth().Token().RenewSelf(0)
	return err
}
