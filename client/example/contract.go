package example

import "github.com/PlatONE_Network/PlatONE-SDK-Go/client"

type Contract struct {
	con client.IContract
}

func (c Contract) test() {
	c.con.IsFuncNameInContract("test")
}
