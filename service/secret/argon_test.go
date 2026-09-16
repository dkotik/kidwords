package secret

import "testing"

func TestParseArgonHash(t *testing.T) {
	_, err := ParseArgonHash(`$argon2id$v=19$m=16,t=2,p=1$TktaRkpkSnlqazVvbjRLUQ$fO7XzlspHhv2xxSEkY05Eg`)
	if err != nil {
		t.Fatal(err)
	}

	secret, err := NewArgonHash([]byte(`testSecret`))
	if err != nil {
		t.Fatal(err)
	}

	ok, err := secret.Match([]byte(`testSecret`))
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("secret does not match")
	}
}
