package controllers

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/zackartz/time-tracker/models"
	"os"
	"time"
)

type Client struct {
	DB *gorm.DB
	ts models.Timestamp
}

func (c *Client) Initialize() {
	err := c.DB.Debug().AutoMigrate(&models.Timestamp{}).Error
	if err != nil {
		panic(err)
	}
}

func (c *Client) Start() {
	c.ts.StartTime = time.Now()
}

func (c *Client) End(comment, category string) error {
	c.ts.EndTime = time.Now()
	c.ts.Category = category
	c.ts.Comment = comment
	return c.DB.Debug().Create(&c.ts).Error
}

func (c *Client) Export(category string) {
	var ts []models.Timestamp
	err := c.DB.Debug().Model(models.Timestamp{}).Where("category = ? and start_time between ? and ?", category, "2020-07-26", "2020-08-08").Find(&ts).Error
	if err != nil {
		fmt.Println("Problem with SQL", err)
	}
	f, err := os.Create(fmt.Sprintf("/home/zack/Documents/time-out-%s.txt", time.Now().String()))
	if err != nil {
		panic(fmt.Sprintf("could not create file %v", err))
	}

	var total time.Time
	for i := range ts {
		total = total.Add(ts[i].EndTime.Sub(ts[i].StartTime))
		l, err := f.WriteString(fmt.Sprintf("%s - %v, %s \n", ts[i].StartTime.Format("January 2, 2006"), ts[i].EndTime.Sub(ts[i].StartTime), ts[i].Comment))
		if err != nil {
			fmt.Println(err)
			err := f.Close()
			if err != nil {
				panic("fuck")
			}
			break
		}
		fmt.Println(l, "bytes written successfully")
	}
	l, err := f.WriteString(fmt.Sprintf("Total : %vhr(s) %dm", total.Hour(), total.Minute()))
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		panic("fuck")
	}
}
