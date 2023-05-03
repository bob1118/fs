package run_time

import (
	"fmt"
	"strings"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/esl/eventsocket"
)

//mod_fifo default queue is cool_fifo@$${domain}
//define queue fifomember@fifos and fifoconsumer@fifos;
//fifo member manage function, fifo_member add/fifo_member del

// FifoMemberManage
// func FifoMemberManage(c *eventsocket.Connection, originate string, is bool) (e error) {
// 	var myerr error
// 	var apicmd string
// 	var op string

// 	condition := fmt.Sprintf("member_string='%s'", originate)
// 	if fifomembers, err := db.SelectFifomembersWithCondition(condition); err != nil {
// 		fmt.Println(err)
// 		return err
// 	} else {
// 		for _, fifomember := range fifomembers {
// 			if is {
// 				op = "fifo_member add"
// 			} else {
// 				op = "fifo_member del"
// 			}
// 			apicmd = fmt.Sprintf("api %s %s %s %s %s %s", op, fifomember.Fname, fifomember.Mstring, fifomember.Msimo, fifomember.Mtimeout, fifomember.Mlag)
// 			if _, err := c.Send(apicmd); err != nil {
// 				fmt.Println(err)
// 				myerr = err
// 			}
// 		}
// 	}
// 	return myerr
// }

// FifoMemberAdd
func FifoMemberAdd(c *eventsocket.Connection, originates ...string) error {
	var myerr error
	var apicmd string
	var condition string

	if len(originates) > 0 {
		for _, originate := range originates {
			condition += fmt.Sprintf("member_string='%s' or", originate)
		}
		condition = strings.TrimSuffix(condition, `or`)
		if fifomembers, err := db.SelectFifomembersWithCondition(condition); err != nil {
			fmt.Println(err)
			return err
		} else {
			for _, fifomember := range fifomembers {
				apicmd = fmt.Sprintf("api fifo_member add %s %s %s %s %s", fifomember.Fname, fifomember.Mstring, fifomember.Msimo, fifomember.Mtimeout, fifomember.Mlag)
				if _, err := c.Send(apicmd); err != nil {
					fmt.Println(err)
					myerr = err
				}
			}
		}
	} else {
		myerr = fmt.Errorf("FifoMemberAdd(%v,%v), len(originates)= 0", c, originates)
	}
	return myerr
}

// FifoMemberDel
func FifoMemberDel(c *eventsocket.Connection, originates ...string) error {
	var myerr error
	var apicmd string
	var condition string

	if len(originates) > 0 {
		for _, originate := range originates {
			condition += fmt.Sprintf("member_string='%s' or", originate)
		}
		condition = strings.TrimSuffix(condition, `or`)
		if fifomembers, err := db.SelectFifomembersWithCondition(condition); err != nil {
			fmt.Println(err)
			return err
		} else {
			for _, fifomember := range fifomembers {
				apicmd = fmt.Sprintf("api fifo_member del %s %s %s %s %s", fifomember.Fname, fifomember.Mstring, fifomember.Msimo, fifomember.Mtimeout, fifomember.Mlag)
				if _, err := c.Send(apicmd); err != nil {
					fmt.Println(err)
					myerr = err
				}
			}
		}
	} else {
		myerr = fmt.Errorf("FifoMemberDel(%v,%v), len(originates)= 0", c, originates)
	}
	return myerr
}
