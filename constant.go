package main

import "fmt"

func main() {
	i3 := 0
	s := ""
	s2 := ""
	s3 := ""
	s4 := ""
	s5 := ""
	for {
		i4 := 138
		if i3 >= 3 {
			break
		}
		if i3 == 0 {
			i4 = 130
		} else if i3 != 1 {
			if i3 != 2 {
				i4 = 0
			} else {
				i4 = 166
			}
		}

		s = s + string(i4/2)

		i3++
	}

	i5 := 0
	for {
		i6 := 206
		if i5 < 17 {
			switch i5 {
			case 0:
				i6 = 130
			case 1:
				fallthrough
			case 4:
				i6 = 138
			case 2:
				i6 = 166
			case 3:
				fallthrough
			case 7:
				i6 = 94
			case 5:
				i6 = 134
			case 6:
				i6 = 132
			case 8:
				i6 = 156
			case 9:
				i6 = 222
			case 10:
				i6 = 160
			case 11:
				i6 = 194
			case 12:
				fallthrough
			case 13:
				i6 = 200
			case 14:
				i6 = 210
			case 15:
				i6 = 220
			case 16:
			default:
				i6 = 0
			}
			s2 = s2 + string(i6/2)
			i5++
		} else {
			i2 := 0
			for i7 := 0; i7 < 17; i7++ {
				switch i7 {
				case 0:
					i2 = 130
				case 1:
					i2 = 138
				case 2:
					i2 = 166
				case 3:
					fallthrough
				case 7:
					i2 = 94
				case 4:
					fallthrough
				case 6:
					i2 = 134
				case 5:
					i2 = 132
				case 8:
					i2 = 156
				case 9:
					i2 = 222
				case 10:
					i2 = 160
				case 11:
					i2 = 194
				case 12:
					fallthrough
				case 13:
					i2 = 200
				case 14:
					i2 = 210
				case 15:
					i2 = 220
				case 16:
					i2 = 206
				default:
					i2 = 0
				}
				s3 = s3 + string(i2/2)
			}
			i8 := 0
			for {
				i9 := 152
				if i8 < 16 {
					switch i8 {
					case 0:
						fallthrough
					case 8:
						i9 = 130
					case 1:
						fallthrough
					case 9:
						i9 = 216
					case 2:
						fallthrough
					case 10:
						i9 = 210
					case 3:
						fallthrough
					case 6:
						fallthrough
					case 11:
						fallthrough
					case 14:
						i9 = 194
					case 4:
						fallthrough
					case 12:
						i9 = 230
					case 5:
						fallthrough
					case 13:
						// break
					case 7:
						fallthrough
					case 15:
						i9 = 196
					default:
						i9 = 0
					}
					s4 = s4 + string(i9/2)
					i8 += 1
				} else {
					for i10 := 0; i10 < 16; i10++ {
						i11 := 122
						switch i10 {
						case 0:
							fallthrough
						case 10:
							i11 = 152
						case 1:
							fallthrough
						case 3:
							// break;
						case 2:
							i11 = 168
						case 4:
							i11 = 98
						case 5:
							i11 = 174
						case 6:
							i11 = 116
						case 7:
							i11 = 148
						case 8:
							i11 = 134
						case 9:
							i11 = 140
						case 11:
							i11 = 166
						case 12:
							i11 = 150
						case 13:
							i11 = 144
						case 14:
							i11 = 102
						case 15:
							i11 = 132
						default:
							i11 = 0
						}
						s5 = s5 + string(i11/2)
					}
                    fmt.Printf("s %s\n", s)
                    fmt.Printf("s2 %s\n", s2)
                    fmt.Printf("s3 %s\n", s3)
                    fmt.Printf("s4 %s\n", s4)
                    fmt.Printf("s5 %s\n", s5)
                    return 
				}
			}
		}
	}
}
