package funcmap

import "git.zxq.co/ripple/rippleapi/common"

// Has returns whether priv1 has all 1 bits of priv2, aka priv1 & priv2 == priv2
func Has(priv1 interface{}, priv2 float64) bool {
	var p1 uint64
	switch priv1 := priv1.(type) {
	case common.UserPrivileges:
		p1 = uint64(priv1)
	case float64:
		p1 = uint64(priv1)
	case int:
		p1 = uint64(priv1)
	}
	return p1&uint64(priv2) == uint64(priv2)
}
