package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"net/url"

	"github.com/linkvault/internal/store"
)

type LinkService struct {
	store store.LinkStore
}

func NewLinkService(store *store.PostgresLinkStore) *LinkService {
	return &LinkService{store}
}

func (svc *LinkService) CreateShortLink(ctx context.Context, longUrl string) (*store.Link, error) {
	userID := ctx.Value("userID").(string)

	shortCode, err := shortLink(longUrl)
	if err != nil {
		return nil, err
	}

	var link store.Link
	link.Longurl = longUrl
	link.Shortcode = shortCode
	link.ClickCount = 0
	link.CreatedBy = userID

	newLink, err := svc.store.SaveShortLink(ctx, link)
	if err != nil {
		return nil, err
	}

	return &newLink, nil
}

func (svc *LinkService) AllLinks(ctx context.Context) ([]store.Link, error) {
	userID := ctx.Value("userID").(string)
	links, err := svc.store.GetAllLinks(ctx, userID)
	if err != nil {
		return nil, err
	}
	return links, err
}

func shortLink(long_url string) (string, error) {
	_, err := url.ParseRequestURI(long_url)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write([]byte(long_url))
	hash := hasher.Sum(nil)
	encoded := base64.URLEncoding.EncodeToString(hash)

	return encoded[:7], nil
}
