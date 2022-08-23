package plugin

import (
	"github.com/fragforce/fragcenter/lib/logs"
	"github.com/google/uuid"
)

type Cell interface {
	logs.LAble       // Add .L()
	GUID() uuid.UUID // Returns the GUID

}
