package ipdb

import (
	"embed"
	"fmt"
	"log/slog"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"

	"github.com/chaitin/ModelKit/backend/domain"
)

//go:embed ip2region.xdb
var ipdbFiles embed.FS

type IPDB struct {
	searcher *xdb.Searcher
	logger   *slog.Logger
}

func NewIPDB(logger *slog.Logger) (*IPDB, error) {
	cBuff, err := xdb.LoadContentFromFS(ipdbFiles, "ip2region.xdb")
	if err != nil {
		return nil, fmt.Errorf("load xdb index failed: %w", err)
	}
	searcher, err := xdb.NewWithBuffer(cBuff)
	if err != nil {
		return nil, fmt.Errorf("new xdb reader failed: %w", err)
	}
	return &IPDB{searcher: searcher, logger: logger.With("module", "ipdb")}, nil
}

func (a *IPDB) Lookup(ip string) (*domain.IPAddress, error) {
	region, err := a.searcher.SearchByStr(ip)
	if err != nil {
		return nil, fmt.Errorf("search ip failed: %w", err)
	}
	ipInfo := strings.Split(region, "|")
	if len(ipInfo) != 5 {
		return nil, fmt.Errorf("invalid ip info: %s", region)
	}
	country := ipInfo[0]
	province := ipInfo[2]
	city := ipInfo[3]
	if country == "0" {
		country = "未知"
	}
	if province == "0" {
		province = "未知"
	}
	if city == "0" {
		city = "未知"
	}
	return &domain.IPAddress{
		IP:       ip,
		Country:  country,
		Province: province,
		City:     city,
	}, nil
}
