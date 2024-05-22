package types

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/unidoc/unipdf/v3/creator"
)

type Client struct {
	Creator      *creator.Creator
	UniDocAPIKey string
}

func NewClient(creator *creator.Creator, uniDocAPIKey string) *Client {
	return &Client{creator, uniDocAPIKey}
}

func (c *Client) Save() (string, error) {
	fname := uuid.New().String()

	err := c.Creator.WriteToFile(fmt.Sprintf("%s.pdf", fname))
	if err != nil {
		return "", err
	}

	return fname, nil
}

type CellStyle struct {
	HAlignment creator.CellHorizontalAlignment
	VAlignment creator.CellVerticalAlignment

	Indent float64
}
