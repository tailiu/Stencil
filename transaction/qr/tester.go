package qr

import (
	"fmt"
)

func (self QR) TestQuery() {
	q := `INSERT INTO customer (c_id, c_d_id, c_w_id, c_first, c_middle, c_last, c_street_1, c_street_2, c_city, c_state, c_zip, c_phone, c_since, c_credit, c_credit_lim, c_discount, c_balance, c_ytd_payment, c_payment_cnt, c_delivery_cnt, c_data, mark_delete, col33, col32, col27, col38) 
		  VALUES
		  (1, 1, 1, 'rb.#y2u\\t(_p', 'OE', 'BARBARBAR', ':XP13%zkV.Ni4LW', 'jZ)VoU1\\3:6%pGb&', 'bx,UCn[(jh]0Y<8', '(v', '562711111', '1276005133733447', '2018-10-17 16:09:46+00:00', 'GC', 50000.0, 0.26560458540916, 
		  '-99.0', 10.0, 1, 0, 'I1&Wq\\$>n+ubgy$Y(y?ovCL!@kq92@;R./l@O*s^:b"imuK+[MuLECW$pOP9r+=fZn#PyA:W2=+/^Efqbq#D9i8|^;Dx(\\wfx>a#{\\jw{Nzy.hs,:q)H] Uwr^BR^wLSYo(;jUoFGxa4PTZwO?/"sizF,,H#vFI&Rr)K;SQnI[<d2nU[MrB_dx=4nsU[4jMta&Sup#lU*CMzF=RP#N*@$[$(-.H1iL9vZa+G0DAk(L59pV(y&>6wu>s{\\-uTLlS,yKaI;7U8c9a!P996Rz8&7Q=jOUpf*1SYgGBqXvNC@q7xjf?-ZW)G@HTz"]DB)y@+ZdNcano>V@%:1tXg7^%IU^4m?9:txNe:h"2cP
		  w!<y3"M-#i*7lWDp', true, 'dataforcol33', 'dataforcol32', 'dataforcol27', 'dataforcol42')`
	q = `
		INSERT INTO item (i_id, i_im_id, i_name, i_price, i_data, mark_delete) VALUES
		(1, 4762, '#PyA:W2=+/^Efqbq#D9i', 1.4249650239944, '8|^;Dx(\\wfx>a#{\\jw{Nzy.hs,:q)H] Uwr^BR^w', false)
	`

	_ = `UPDATE item SET i_id = 129188`
	_ = `UPDATE item SET i_id = 129188 WHERE mark_delete = false AND i_name = 'zain'`
	_ = `UPDATE item SET i_id = 129188, i_data = 'blahblah', col31 = '1233' WHERE ark_delete = false AND i_name = 'zain'`
	_ = `UPDATE customer SET c_id = 129188, c_data = 'blahblah', col38 = '1233' WHERE mark_delete = false AND c_name = 'zain'`
	_ = `DELETE FROM customer WHERE mark_delete = false AND c_name = 'zain'`
	_ = `DELETE FROM customer`

	_ = `SELECT c_id, c_data, c_d_id, col38, col27 FROM customer WHERE c_id = '123'`

	_ = `SELECT * FROM customer  WHERE c_id = '1234' AND c_data = 'aw23de'`

	_ = `SELECT * FROM customer WHERE c_id = '5' `

	q = `SELECT * FROM customer`

	_ = `SELECT Customer.* FROM Customer WHERE c_id = '5' AND Customer.mark_delete != 'true'`

	_ = `SELECT History.* FROM History  JOIN Customer ON History.H_C_ID = Customer.C_ID AND History.H_C_D_ID = Customer.C_D_ID AND History.H_C_W_ID = Customer.C_W_ID WHERE Customer.c_id = '5' AND History.mark_delete != 'true'`

	_ = `SELECT Orderr.* FROM Orderr  JOIN Customer ON Orderr.O_C_ID = Customer.C_ID AND Orderr.O_D_ID = Customer.C_D_ID AND Orderr.O_W_ID = Customer.C_W_ID WHERE Customer.c_id = '5' AND Orderr.mark_delete != 'true'`

	_ = `SELECT New_Order.* FROM New_Order  JOIN Orderr ON New_Order.NO_O_ID = Orderr.O_ID AND New_Order.NO_D_ID = Orderr.O_D_ID AND New_Order.NO_W_ID = Orderr.O_W_ID JOIN Customer ON Orderr.O_C_ID = Customer.C_ID AND Orderr.O_D_ID = Customer.C_D_ID AND Orderr.O_W_ID = Customer.C_W_ID WHERE Customer.c_id = '5' AND New_Order.mark_delete != 'true'`

	_ = `SELECT Order_Line.* FROM Order_Line  JOIN Orderr ON Order_Line.OL_O_ID = Orderr.O_ID AND Order_Line.OL_D_ID = Orderr.O_D_ID AND Order_Line.OL_W_ID = Orderr.O_W_ID JOIN Customer ON Orderr.O_C_ID = Customer.C_ID AND Orderr.O_D_ID = Customer.C_D_ID AND Orderr.O_W_ID = Customer.C_W_ID WHERE Customer.c_id = '5' AND Order_Line.mark_delete != 'true'`

	fmt.Println("------------------------------------------------------------------------------")
	fmt.Println("*QUERY:", q)
	for i, q := range self.Resolve(q) {
		fmt.Println("******************************************************************************")
		fmt.Println(i+1, ":", q)
	}
	fmt.Println("------------------------------------------------------------------------------")
}
