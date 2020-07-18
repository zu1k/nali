package common

const (
	// RedirectMode1 [IP][0x01][国家和地区信息的绝对偏移地址]
	RedirectMode1 = 0x01
	// RedirectMode2 [IP][0x02][信息的绝对偏移][...] or [IP][国家][...]
	RedirectMode2 = 0x02
)
