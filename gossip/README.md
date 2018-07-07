#### Run
```shell
$ go build

$ touch neighbors_3000.dat

$ ./gossip help
```

#### Command:
*	getaddr: dùng khi 1 peer muốn yêu cầu danh sách các hàng xóm của peer khác.
*	addr: dùng để phản hồi lại getaddr, trả về danh sách hàng xóm của peer.
*	version: dùng để gửi metadata (height) của data set hiện tại của peer cho peer khác.
*	getdata: dùng khi 1 peer muốn yêu cầu lấy data từ peer khác.
*	data: dùng để phản hồi lại getdata, trả về data set của peer.

#### Flow:
1. Seed node: 3000 (luôn chạy)
2. Node mới start:
  1. Gửi cho seed node getaddr
    1. Seed node phản hồi lại addr
  2.	Node lưu lại danh sách các hàng xóm
  3.	Gửi cho tất cả hàng xóm version
    1.	Hàng xóm phản hồi lại version nếu height của hàng xóm dài hơn height của node.
    2.	Hàng xóm phản hồi lại getdata nếu height ngắn hơn height của node gửi.
  4.	Sau khi gửi getdata, bên nhận sẽ phản hồi lại data.
  5.	Lưu data set nhận được vào db local.

#### Cần cải thiện:
* Khi 2 node cùng cập nhật data mới, height data set mới bằng nhau, chưa xử lý.
* Do lệnh version được đặt trong hàm handleAddr, nên khi một node broadcast data mới lên mạng, các node khác nếu không trong danh sách hàng xóm của node broadcast, sẽ cần restart mới nhận được data mới. (do sau khi startnode, node gửi getaddr tới seednode, handleAddr mà seednode phản hồi, sau đó mới gửi version đến seednode và các peer khác để cập nhật, mà đã nhận từ seednode thì chắc chắn là mới nhất).
*	2.2: Cần lưu random, hiện tại đang chạy vòng lặp từ 0 – n.
*	2.3: Lệnh version đang được gửi trong hàm handleAddr. Mình cảm thấy chưa được ổn.
*	2.5: Khi áp dụng vào blockchain, cần validate trước khi lưu.
