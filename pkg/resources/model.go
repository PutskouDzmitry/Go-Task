package resources

import (
	"MommyCO/internal/config"
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"
)

type ResourceClient struct {
	basePath string
}

func NewResourceClient(cfg *config.Config) (*ResourceClient, error) {
	return &ResourceClient{
		basePath: cfg.Resources,
	}, nil
}

func (c *ResourceClient) Init(dirs ...string) {
	for _, dir := range dirs {
		log.Println(dir)
		if !strings.HasSuffix(dir, "/") {
			panic("dir path must have \"/\" suffix")
		}
		path := c.basePath + dir
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (c *ResourceClient) Store(dir, filename string, data []byte) error {
	if !strings.HasSuffix(dir, "/") {
		panic("dir path must have \"/\" suffix")
	}
	path := c.basePath + dir + filename
	err := os.WriteFile(path, data, os.ModePerm)
	return err
}

func (c *ResourceClient) Get(dir, filename string) ([]byte, error) {
	if !strings.HasSuffix(dir, "/") {
		panic("dir path must have \"/\" suffix")
	}

	path := c.basePath + dir + filename

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *ResourceClient) Delete(dir, filename string) error {
	if !strings.HasSuffix(dir, "/") {
		panic("dir path must have \"/\" suffix")
	}
	path := c.basePath + dir + filename
	err := os.Remove(path)
	return err
}

func (c *ResourceClient) BasePath() string {
	return c.basePath
}
