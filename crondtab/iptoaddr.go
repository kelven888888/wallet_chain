package crondtab

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"net"
	"time"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/global"
)

var iptoaddr = func() {
	var acclog []model.MAccesslog
	datetimenow := time.Now().AddDate(0, 0, -11).Format("2006-01-02")
	err := global.SHOP_DB.Model(model.MAccesslog{}).Where("create_at<?", datetimenow).Delete(model.MAccesslog{}).Error
	if err != nil {
		global.SHOP_LOG.Log(2, err.Error())
	}
	global.SHOP_DB.Where("address='' or address is Null").Find(&acclog)
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		global.SHOP_LOG.Log(2, err.Error())
	}
	defer db.Close()
	for _, v := range acclog {
		// If you are using strings that may be invalid, check that ip is not nil

		ip := net.ParseIP(v.Ip)
		//ip = net.ParseIP("112.97.203.222")
		record, err := db.City(ip)
		if err != nil {
			global.SHOP_LOG.Log(2, err.Error())
		}
		//fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
		//if len(record.Subdivisions) > 0 {
		//	fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
		//}
		//fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
		//fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
		//fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
		//fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
		// Output:
		// Portuguese (BR) city name: Londres
		// English subdivision name: England
		// Russian country name: Великобритания
		// ISO country code: GB
		// Time zone: Europe/London
		// Coordinates: 51.5142, -0.0931

		//fmt.Println("中文结果")
		//fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["zh-CN"])
		//if len(record.Subdivisions) > 0 {
		//	fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["zh-CN"])
		//}
		address := fmt.Sprintf("%s_%s", record.Country.Names["zh-CN"], record.City.Names["zh-CN"])
		if len(record.Subdivisions) > 0 {
			address = fmt.Sprintf("%s_%s", address, record.Subdivisions[0].Names["zh-CN"])
		}
		if len(record.Subdivisions) > 0 {
			global.SHOP_DB.Model(model.MAccesslog{}).Where("id=?", v.Id).Updates(model.MAccesslog{
				City:        record.City.Names["zh-CN"],
				Country:     record.Country.Names["zh-CN"],
				Subdivision: record.Subdivisions[0].Names["zh-CN"],
				Address:     address,
			})
		} else {
			global.SHOP_DB.Model(model.MAccesslog{}).Where("id=?", v.Id).Updates(model.MAccesslog{
				City:    record.City.Names["zh-CN"],
				Country: record.Country.Names["zh-CN"],
				Address: address,
			})
		}
		//fmt.Printf("Russian country name: %v\n", record.Country.Names["zh-CN"])

	}

}
