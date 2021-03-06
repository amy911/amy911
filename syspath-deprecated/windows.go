// +build windows

package syspath

import (
	"path/filepath"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/amy911/amy911/onfail"
)

const dwFlags uint32 =
	0x00008000 | // KF_FLAG_CREATE
	0x00000800 | // KF_FLAG_INIT
	0x00000000   // KF_FLAG_DEFAULT

var (
	OLE32 = syscall.MustLoadDLL("OLE32.DLL")
	CoTaskMemFree = OLE32.MustFindProc("CoTaskMemFree")
	SHELL32 = syscall.MustLoadDLL("SHELL32.DLL")
	SHGetKnownFolderPathW = SHELL32.MustFindProc("SHGetKnownFolderPathW")
)

func initSysPath(sp *SysPath, vendor, application string) {
	sp.sCache = filepath.Join(GetKnownFolderPath(&FOLDERID_InternetCache), vendor, application)
	sp.sConfig = filepath.Join(GetKnownFolderPath(&FOLDERID_RoamingAppData), vendor, application)
	sp.sData = filepath.Join(GetKnownFolderPath(&FOLDERID_LocalAppData), vendor, application)
	sp.sDesktop = GetKnownFolderPath(&FOLDERID_Desktop)
	sp.sDocuments = GetKnownFolderPath(&FOLDERID_Documents)
	sp.sDownloads = GetKnownFolderPath(&FOLDERID_Downloads)
	sp.sHome = GetKnownFolderPath(&FOLDERID_Profile)
	sp.sPictures = GetKnownFolderPath(&FOLDERID_Pictures)
	sp.sScreenshots = GetKnownFolderPath(&FOLDERID_Screenshots)
}

func GetKnownFolderPath(rfid unsafe.Pointer, onFail ...onfail.OnFail) (path string) {
	var out *uint16
	if 0 != SHGetKnownFolderPathW.Call(rfid, dwFlags, 0, &out) {
		onfail.Fail("SHGetKnownFolderPathW reported an error", rfid, onfail.Panic, onFail...)
		return
	}
	len16 := 0
	for p := out; *p != 0; p++ {
		len16++
	}
	buf := make([]uint16, len16, len16)
	for idx := 0; idx < len16; idx++ {
		buf[idx] = out[idx]
	}
	CoTaskMemFree.Call(out)
	path = string(utf16.Decode(buf))
	return
}
