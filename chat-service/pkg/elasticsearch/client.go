package elasticsearch

import (
    "bytes"
    "encoding/json"
    "fmt"
    "time"
    "github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
    URL        string
    MaxRetries int
    RetryDelay time.Duration
}

type Client struct {
    es *elasticsearch.Client
}

func NewClient(cfg Config) (*Client, error) {
    if cfg.MaxRetries == 0 {
        cfg.MaxRetries = 5 // default retries
    }
    if cfg.RetryDelay == 0 {
        cfg.RetryDelay = 5 * time.Second // default delay
    }

    config := elasticsearch.Config{
        Addresses: []string{cfg.URL},
        RetryOnStatus: []int{502, 503, 504}, 
        MaxRetries:    cfg.MaxRetries,
        RetryBackoff: func(i int) time.Duration {
            return cfg.RetryDelay * time.Duration(i+1)
        },
    }

    var client *elasticsearch.Client
    var err error

    // Retry loop for connecting to Elasticsearch because it takes time to start in dockerized instances
    for i := 0; i < cfg.MaxRetries; i++ {
        client, err = elasticsearch.NewClient(config)
        if err != nil {
            fmt.Printf("Attempt %d: failed to create Elasticsearch client: %v\n", i+1, err)
            time.Sleep(cfg.RetryDelay)
            continue
        }

        res, err := client.Info()
        if err != nil {
            fmt.Printf("Attempt %d: failed to ping Elasticsearch: %v\n", i+1, err)
            time.Sleep(cfg.RetryDelay)
            continue
        }
        defer res.Body.Close()

        if res.IsError() {
            fmt.Printf("Attempt %d: Elasticsearch connection error, status: %s\n", i+1, res.Status())
            time.Sleep(cfg.RetryDelay)
            continue
        }

        fmt.Println("Successfully connected to Elasticsearch!")
        return &Client{es: client}, nil
    }

    return nil, fmt.Errorf("could not connect to Elasticsearch after %d attempts: %w", cfg.MaxRetries, err)
}

func (c *Client) Index(index string, id string, document interface{}) error {
    var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(document); err != nil {
        return fmt.Errorf("failed to encode document: %w", err)
    }

    res, err := c.es.Index(
        index,
        &buf,
        c.es.Index.WithDocumentID(id),
        c.es.Index.WithRefresh("true"), 
    )
    if err != nil {
        return fmt.Errorf("failed to index document: %w", err)
    }
    defer res.Body.Close()

    if res.IsError() {
        var errorMap map[string]interface{}
        if err := json.NewDecoder(res.Body).Decode(&errorMap); err != nil {
            return fmt.Errorf("failed to decode error response: %w", err)
        }
        return fmt.Errorf("failed to index document: %v", errorMap)
    }

    return nil
}

func (c *Client) Search(index string, query map[string]interface{}) ([]byte, error) {
    var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(query); err != nil {
        return nil, fmt.Errorf("failed to encode query: %w", err)
    }

    res, err := c.es.Search(
        c.es.Search.WithIndex(index),
        c.es.Search.WithBody(&buf),
        c.es.Search.WithPretty(),
        c.es.Search.WithSize(10),
        c.es.Search.WithTimeout(30 * time.Second),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to execute search: %w", err)
    }
    defer res.Body.Close()

    if res.IsError() {
        var errorMap map[string]interface{}
        if err := json.NewDecoder(res.Body).Decode(&errorMap); err != nil {
            return nil, fmt.Errorf("failed to decode error response: %w", err)
        }
        return nil, fmt.Errorf("search failed: %v", errorMap)
    }

    var buf2 bytes.Buffer
    if _, err := buf2.ReadFrom(res.Body); err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }

    return buf2.Bytes(), nil
}
