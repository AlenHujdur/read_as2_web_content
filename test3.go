import (
	"encoding/base64"
	"fmt"
	"strings"
)

func main() {
	data := `2006	;,13088	;,0E11E45010DE4283AA81968E0006FC45	;,33	;,3	;,DIV_1.pdf	;,1	;,5592	;,JVBERi0xLjIgDQoxIDAgb2JqDQo8PA0KL1R5cGUgL0NhdGFsb2cNCi9QYWdlcyAzIDAgUg0KPj4N
CmVuZG9iag0KMiAwIG9iag0KPDwNCi9Qcm9kdWNlciAoUG93ZXJQZGYgdmVyc2lvbiAwLjkpDQov
Q3JlYXRpb25EYXRlIChEOjIwMTcwODE2MDkzNTA4KQ0KL01vZERhdGUgKEQ6MjAxNzA4MTYwOTM1
MDgpDQo+Pg0KZW5kb2JqDQozIDAgb2JqDQo8PA0KL1R5cGUgL1BhZ2VzDQovS2lkcyBbNCAwIFIg
XQ0KL0NvdW50IDENCj4+DQplbmRvYmoNCjQgMCBvYmoNCjw8DQovVHlwZSAvUGFnZQ0KL1BhcmVu
dCAzIDAgUg0KL01lZGlhQm94IFswIDAgNTk2IDg0MiBdDQovUmVzb3VyY2VzIDw8DQovRm9udCA8
PA0KL0YwIDYgMCBSDQovRjEgNyAwIFINCj4+DQovWE9iamVjdCA8PA0KPj4NCi9Qcm9jU2V0IFsv
UERGIC9UZXh0IC9JbWFnZUMgXQ0KPj4NCi9Db250ZW50cyA1IDAgUg0KPj4NCmVuZG9iag0KNSAw
IG9iag0KPDwNCi9MZW5ndGggMzAxNw0KL0ZpbHRlciBbXQ0KPj4NCnN0cmVhbQ0KL0YwIDEwIFRm
MDAwMDAzNjg3IDAwMDAwIG4NCnRyYWlsZXINCjw8DQovU2l6ZSA4DQovUm9vdCAxIDAgUg0KL0lu
Zm8gMiAwIFINCj4+DQpzdGFydHhyZWYNCjM4MzgNCiUlRU9GDQo=
	;,2#EndOfLine#`

	// Split the data into columns based on the separator " ;,"
	columns := strings.Split(data, ";,")

	// Retrieve the last column value
	lastColumn := columns[len(columns)-1]

	// Decode the Base64 encoded content
	decodedContent, err := base64.StdEncoding.DecodeString(lastColumn)
	if err != nil {
		fmt.Println("Failed to decode Base64 content:", err)
		return
	}

	fmt.Println("Decoded Content:", string(decodedContent))
}
