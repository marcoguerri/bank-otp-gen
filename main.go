package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"flag"
	"os"
	"strings"
	"time"
)

var f6939b []byte
var f6940b string
var c_global []byte

var a = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var e []byte
var iv = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

var sc_sac = ""
var sc_k2 = ""
var pin = ""

var naes int64 = 0

func time_attempts(minTimeDelta int, maxTimeDelta int) []uint32 {
	now := uint32(time.Now().UnixMilli() / 3600000)
	attempts := []uint32{}
	for i := minTimeDelta; i < maxTimeDelta; i++ {
		attempts = append(attempts, now+uint32(i))
		attempts = append(attempts, now-uint32(i))
	}
	return attempts
}

func copy_part_of_slice_and_xor_with_constant(in []byte, i2 int, i3 int) []byte {
	// Implements smali/n2/a/a/a.smali|c()
	b2 := make([]byte, i2)
	copy(b2[0:i2], in[i3:])
	s := "AliasLabAliasLab"
	b3 := make([]byte, i2)
	if len([]byte(s)) >= i2 {
		for i := 0; i < i2; i++ {
			b3[i] = b2[i] ^ []byte(s)[i]
		}
	}
	return b3
}

func copy_part_of_slice_and_xor_with_constant_second_variant(in []byte, i2 int, i3 int) []byte {
	// This is exactly identical to `copy_part_of_slice_and_xor_with_constant` except that the
	// constant is different
	b2 := make([]byte, i2)
	copy(b2[0:i2], in[i3:])
	s := "L=T=1W:JCFLSKH3B"
	b3 := make([]byte, i2)
	if len([]byte(s)) >= i2 {
		for i := 0; i < i2; i++ {
			b3[i] = b2[i] ^ []byte(s)[i]
		}
	}
	return b3
}

func combine_key_with_salt_initialize_globals(s string, b []byte) {
	// Implements smali/n2/a/a/f.smali|<init>
	out := make([]byte, 16)
	i2 := uint32(0)
	for _, c := range strings.ToUpper(s) {
		i2 = (i2 * 31) + uint32(c)
	}
	binary.BigEndian.PutUint32(out[0:], i2)
	copy(out[4:], b[0:12])

	key := out
	c2, c3 := encrypt_nulltext_and_shift_right(out)
	f6939b = c2
	c_global = c3
	e = key
}

func xor_trailing_bytes_c2_c3_and_encrypt_with_e_wrapper(fragment string) string {
	// Implements smali/n2/a/a/f.smali|a()
	res := xor_trailing_bytes_c2_c3_and_encrypt_with_e([]byte(fragment))
	return hex.EncodeToString(res)
}

func shift_right_byte_slice(in []byte) []byte {
	// Implements smali/n2/a/a/e.smali|c()
	l := len(in)
	out := make([]byte, l)
	copy(out, in)

	z := false

	if in[0]&128 != 0 {
		z = true
	}

	b2 := out[l-1]
	out[l-1] = out[l-1] & 255 << 1
	i3 := l - 2

	for {
		if i3 < 0 {
			break
		}
		b4 := out[i3]
		out[i3] = (b2 & 255 >> 7) | (out[i3] & 255 << 1)
		i3 -= 1
		b2 = b4
	}
	if z == true {
		len2 := l - 1
		out[len2] = out[len2] ^ byte(0x87)
	}
	return out

}

func xor_common_lenght_arrays_copy_rest(b1 []byte, b2 []byte) []byte {
	// Implements smali/n2/a/a/c.smali|d()
	i2 := 0
	if len(b1) >= len(b2) {
		b3 := make([]byte, len(b1))
		for {
			if i2 >= len(b2) {
				break
			}
			b3[i2] = b1[i2] ^ b2[i2]
			i2 += 1
		}

		copy(b3[len(b2):len(b1)], b1[len(b2):len(b1)])
		return b3
	}

	b4 := make([]byte, len(b2))
	for {
		if i2 >= len(b1) {
			break
		}
		b4[i2] = b1[i2] ^ b2[i2]
		i2 += 1
	}

	copy(b4[len(b1):len(b2)], b2[len(b1):len(b2)])

	return b4

}

func split_slice(in []byte, l int) [][]byte {

	// Implements smali/n2/a/a/c.smali|e()

	out := make([][]byte, 0)
	start := 0
	for {
		if start >= len(in) {
			break
		}
		if start+l > len(in) {
			out = append(out, in[start:])
			break
		}
		out = append(out, in[start:start+l])
		start = start + l
	}

	return out

}

func hex_to_byte(s string) []byte {
	// Implements smali/n2/a/a/c.smali|b()
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func invert_slice(in []byte) []byte {
	// Implements smali/n2/a/a/c.smali|f()
	out := make([]byte, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = in[len(in)-1-i]
	}
	return out
}

func encrypt_nulltext_and_shift_right(key []byte) ([]byte, []byte) {
	// Implements smali/n2/a/a/e.smali | b()
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	text := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	out := make([]byte, len(text))
	c.Encrypt(out, text)
	naes += 1

	c2 := shift_right_byte_slice(out) // this is saved as f6939b (global)

	c3 := shift_right_byte_slice(c2) // this is saved as c (global)

	return c2, c3

}

func concatenate_slices(l [][]byte) []byte {
	// Implements smali/n2/a/a/c.smali|c()
	i2 := 0
	for i := 0; i < len(l); i++ {
		i2 += len(l[i])
	}

	out := make([]byte, i2)
	i3 := 0
	for i4 := 0; i4 < len(l); i4++ {
		b := l[i4]
		copy(out[i3:], b)
		i3 += len(l[i4])
	}
	return out
}

func xor_trailing_bytes_c2_c3_and_encrypt_with_e(b []byte) []byte {
	// Implements smali/n2/a/a/e.smali|a()
	arrayList := split_slice(b, 16)
	size := len(arrayList) - 1
	b3 := arrayList[size]

	if len(b3) == 16 {
		mod := xor_common_lenght_arrays_copy_rest(b3, f6939b)
		arrayList[size] = mod
	} else {
		d2 := xor_common_lenght_arrays_copy_rest(b3, a)
		d2[len(b3)] = 0x80
		d3 := xor_common_lenght_arrays_copy_rest(d2, c_global)
		arrayList[size] = d3
	}

	plaintext := concatenate_slices(arrayList)

	ciphertext := make([]byte, len(plaintext))
	block, _ := aes.NewCipher(e)
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)
	naes += 1

	b4 := make([]byte, 16)
	copy(b4, split_slice(ciphertext, 16)[size][0:len(b4)])

	return b4

}

func xor_with_encrypted_nulltext_i_times(keyString string, b []byte, i int) []byte {
	// Implements smali/n2/a/a/a.smali|b()
	key := hex_to_byte(keyString)
	c2, c3 := encrypt_nulltext_and_shift_right(key)
	f6939b = c2
	c_global = c3
	e = key

	b1 := xor_trailing_bytes_c2_c3_and_encrypt_with_e(b)
	b2 := make([]byte, 16)

	copy(b2, b1[0:16])

	b3 := b1
	for i3 := 2; i3 <= i; i3++ {
		b3 = xor_trailing_bytes_c2_c3_and_encrypt_with_e(b3)
		for i4 := 0; i4 < 16; i4++ {
			b2[i4] = b2[i4] ^ b3[i4]
		}
	}
	copy(b1, b3[0:16])
	b4 := make([]byte, 32)

	copy(b4, b2[0:16])
	a2 := xor_trailing_bytes_c2_c3_and_encrypt_with_e(b1)
	for i5 := 0; i5 < 16; i5++ {
		b2[i5] = b2[i5] ^ a2[i5]
	}
	copy(b4[16:], b2[0:16])
	return b4
}

func n2_a_smali_a(keyCombined string, encrypted []byte, variant int) []byte {
	// Implements java_src/n2/a/a/a.java|a()
	if len(keyCombined) != 64 {
		panic("string not 64 bytes long")
	}
	var c2 []byte
	var c3 []byte
	// Variant is chosen based on whether we need to decrypt QR code or Transaction data
	if variant == 0 {
		c2 = copy_part_of_slice_and_xor_with_constant(encrypted, 16, 0)
	} else {
		c2 = copy_part_of_slice_and_xor_with_constant_second_variant(encrypted, 16, 0)
	}
	newKey := xor_with_encrypted_nulltext_i_times(keyCombined, c2, 5000)

	if variant == 0 {
		c3 = copy_part_of_slice_and_xor_with_constant(encrypted, 8, 16)
	} else {
		c3 = copy_part_of_slice_and_xor_with_constant_second_variant(encrypted, 8, 16)
	}


	bArr2 := make([]byte, len(encrypted)-8-16)
	copy(bArr2, encrypted[24:24+len(encrypted)-8-16])
	arrayList := split_slice(bArr2, 16)
	size := len(arrayList)
	arrayList2 := make([][]byte, 0)

	if size > 256 {
		panic("cannot serialize integer with more than 8 bits")
	}

	for i2 := 0; i2 < size; i2++ {
		bArr3 := make([]byte, 16)
		f := invert_slice(xor_common_lenght_arrays_copy_rest(invert_slice([]byte{byte(i2)}), bArr3))
		copy(f, c3[0:8])
		c, err := aes.NewCipher(newKey)
		if err != nil {
			panic(err)
		}
		out := make([]byte, len(f))
		c.Encrypt(out, f)
		naes += 1
		arrayList2 = append(arrayList2, out)
	}

	arrayList3 := make([][]byte, 0)
	for i3 := 0; i3 < size; i3++ {
		bArr4 := arrayList[i3]
		bArr5 := arrayList2[i3]
		l := len(bArr4)
		bArr6 := make([]byte, l)
		for i4 := 0; i4 < l; i4++ {
			bArr6[i4] = bArr4[i4] ^ bArr5[i4]
		}
		arrayList3 = append(arrayList3, bArr6)
	}
	return concatenate_slices(arrayList3)

}

func v2_c_smali_b(qrcode string, minTimeDelta int, maxTimeDelta int) string {
	// Implements smali_classes4/v2/a/a/c.smali|b()
	sc_sac_b, err := hex.DecodeString(sc_sac)
	if err != nil {
		panic(err)
	}

	combine_key_with_salt_initialize_globals(strings.ToUpper(sc_k2[0:5]), sc_sac_b)
	key := xor_trailing_bytes_c2_c3_and_encrypt_with_e_wrapper(strings.ToUpper(sc_k2[0:5]))

	qrcodeb, err := hex.DecodeString(qrcode)
	if err != nil {
		panic(err)
	}

	sc_k2_b, err := hex.DecodeString(sc_k2)
	if err != nil {
		panic(err)
	}

	for _, delta := range time_attempts(minTimeDelta, maxTimeDelta) {
		deltaB := make([]byte, 4)
		binary.BigEndian.PutUint32(deltaB, delta)
		xored := xor_common_lenght_arrays_copy_rest(sc_k2_b, deltaB)
		xoredS := hex.EncodeToString(xored)
		keyCombined := "6ab392fd02" + key + xoredS

		dec := n2_a_smali_a(keyCombined, qrcodeb, 0)
		if strings.HasPrefix(string(dec), "<txt>") {
			return string(dec)
		}
	}

	panic("qr code could not be decrypted")

}

func main() {
	minTimeDelta := flag.Int("minTimeDelta", -1, "Minimum value for time delta")
	maxTimeDelta := flag.Int("maxTimeDelta", 1, "Maximumvalue for time delta")
	qrcode := flag.String("qrcode", "", "QR code")

	flag.Parse()
	if *minTimeDelta > *maxTimeDelta {
		panic("minTimeDelta > maxTimeDelta")
	}

	fmt.Printf("Reading secret file from stdin\n")
	var secrets map[string]string
	secretFile := ""
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		secretFile += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(secretFile), &secrets); err != nil {
		panic(err)
	}
	sc_sac = secrets["sc_sac"]
	sc_k2 = secrets["sc_k2"]
	pin = secrets["pin"]

	qrcodePlaintext := v2_c_smali_b(*qrcode, *minTimeDelta, *maxTimeDelta)

    fmt.Printf("qrcode: %s\n", qrcodePlaintext)

	start := strings.Index(qrcodePlaintext, "}")
	end := strings.Index(qrcodePlaintext[1:], "<")

	transactionData := qrcodePlaintext[start+1 : end+1]

	raw, err := base64.StdEncoding.DecodeString(transactionData)
	if err != nil {
		panic(err)
	}

	sc_sac_b, err := hex.DecodeString(sc_sac)
	if err != nil {
		panic(err)
	}
	combine_key_with_salt_initialize_globals(pin, sc_sac_b)
	key := xor_trailing_bytes_c2_c3_and_encrypt_with_e_wrapper(pin)

	sc_k2_b, err := hex.DecodeString(sc_k2)
	if err != nil {
		panic(err)
	}

	for _, delta := range time_attempts(*minTimeDelta, *maxTimeDelta) {
		deltaB := make([]byte, 4)
		binary.BigEndian.PutUint32(deltaB, delta)

		xored := xor_common_lenght_arrays_copy_rest(sc_k2_b, deltaB)
		xoredS := hex.EncodeToString(xored)
		keyCombined := key + xoredS + "8000000000"

		plaintext := n2_a_smali_a(keyCombined, raw, 1)
		if strings.Index(string(plaintext), "#") != -1 {
			if len(plaintext)-strings.Index(string(plaintext), "#") != 7 {
				continue
			}
			fmt.Printf("plaintext: %s\n", string(plaintext))
			fmt.Printf("delta is %d\n", delta)
			fmt.Printf("naes is %d\n", naes)
			break
		}
	}
}
