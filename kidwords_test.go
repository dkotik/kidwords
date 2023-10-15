package kidwords

import "fmt"

// func TestIntTransformations(t *testing.T) {
// 	cases := []int64{0, 9, 16, 32, 999999, 38729387428974, 2374761653249823, 88999999}
// 	for _, i := range cases {
// 		t.Run(fmt.Sprintf("transforming integer: %d", i), func(t *testing.T) {
// 			b := &bytes.Buffer{}
// 			if err := WriteInt(b, i); err != nil {
// 				t.Fatal(err)
// 			}
// 			legitWords(t, b.String())
// 			t.Log("words:", b.String())
//
// 			j, err := ReadInt(bytes.NewReader(b.Bytes()))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			if j != i {
// 				t.Fatalf("%d does not match %d", j, i)
// 			}
// 		})
// 	}
// }

func ExampleFromBytes() {
	fmt.Println(
		FromBytes([]byte("marvel")),
	)
	// Output: hole gold hush item half hint <nil>
}

func ExampleToBytes() {
	b, err := ToBytes("  hole gold hush item half hint ")
	fmt.Println(string(b), err)
	// Output: marvel <nil>
}
