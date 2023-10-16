package store

import "testing"

func TestParseArgonHash(t *testing.T) {
	_, err := ParseArgonHash(`$argon2id$v=19$m=16,t=2,p=1$TktaRkpkSnlqazVvbjRLUQ$fO7XzlspHhv2xxSEkY05Eg`)
	if err != nil {
		t.Fatal(err)
	}
}
