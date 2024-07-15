package dao

import (
	"os"

	"github.com/nebstudio/sealos/service/exceptionmonitor/api"

	"github.com/nebstudio/sealos/controllers/pkg/database/cockroach"
)

var (
	CK *cockroach.Cockroach
)

func InitCockroachDB() error {
	var err error
	os.Setenv("LOCAL_REGION", api.LOCALREGION)

	CK, err = cockroach.NewCockRoach(os.Getenv("GlobalCockroachURI"), os.Getenv("LocalCockroachURI"))
	if err != nil {
		return err
	}
	return nil
}
