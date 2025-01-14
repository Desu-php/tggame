package image

import (
	"fmt"
	"strings"
	"example.com/v2/config"
)

type Image struct {
	aws *config.Aws
	url string
}

func NewImage(
	aws *config.Aws,
) (*Image, error) {

	url, err := generateCustomURL(aws.Bucket, aws.Endpoint)

	if err != nil {
		return nil, fmt.Errorf("NewImage err %w", err)
	}

	return &Image{aws: aws, url: url}, nil
}

func (i *Image) Url(path string) string {
	return i.url + "/" + path
}

func generateCustomURL(bucket, endpoint string) (string, error) {
	if !strings.HasPrefix(endpoint, "https://") {
		return "", fmt.Errorf("invalid endpoint: must start with 'https://'")
	}

	endpointWithoutProtocol := strings.TrimPrefix(endpoint, "https://")

	customURL := fmt.Sprintf("https://%s.%s", bucket, endpointWithoutProtocol)

	return customURL, nil
}
