package influxdb

import (
	"fmt"
	"log"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

var (
	I *Client
)

type Client struct {
	dbName, username, password, addr, precision string
	Session                                     client.Client
	reader                                      client.Query
	writer                                      client.BatchPoints
}

func NewClient(addr, username, password, dbName, precision string) (*Client, error) {
	c := &Client{
		precision: precision,
		dbName:    dbName,
	}
	if c.precision == "" {
		c.precision = "s" // 默认设置为秒
	}
	var err error
	c.Session, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
		// Precision: precision,
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

// 存数据   这个数据库自动扩展字段
func (c *Client) WriteDB(table string, tags map[string]string, fields map[string]interface{}) error {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  c.dbName,
		Precision: c.precision,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	// Create a point and add to batch
	// tags := map[string]string{"mem": "mem-total","mem2":"mem2-total"}
	// fields := map[string]interface{}{
	// 	"all": 4096,
	// 	"used": 3308,
	// }
	pt, err := client.NewPoint(table, tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
		return err
	}
	bp.AddPoint(pt)
	// Write the batch
	if err := c.Session.Write(bp); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// QueryDB convenience function to query the database
func (c *Client) QueryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: c.dbName,
	}
	if response, err := c.Session.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

// 创建数据库
func (c *Client) CreateDB(db string) error {
	_, err := c.QueryDB(fmt.Sprintf("CREATE DATABASE %s", db))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
