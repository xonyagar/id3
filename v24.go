package id3

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strconv"
	"strings"
)

var V24Frames = map[string]DeclaredFrame{
	"AENC": {"AENC", "Audio encryption", TypeTextInformation},
	"APIC": {"APIC", "Attached picture", TypeAttachedPicture},
	"ASPI": {"ASPI", "Audio seek point index", TypeTextInformation},
	"COMM": {"COMM", "Comments", TypeTextInformation},
	"COMR": {"COMR", "Commercial frame", TypeTextInformation},
	"ENCR": {"ENCR", "Encryption method registration", TypeTextInformation},
	"EQU2": {"EQU2", "Equalisation (2)", TypeTextInformation},
	"ETCO": {"ETCO", "Event timing codes", TypeTextInformation},
	"GEOB": {"GEOB", "General encapsulated object", TypeTextInformation},
	"GRID": {"GRID", "Group identification registration", TypeTextInformation},
	"LINK": {"LINK", "Linked information", TypeTextInformation},
	"MCDI": {"MCDI", "Music CD identifier", TypeTextInformation},
	"MLLT": {"MLLT", "MPEG location lookup table", TypeTextInformation},
	"OWNE": {"OWNE", "Ownership frame", TypeTextInformation},
	"PRIV": {"PRIV", "Private frame", TypeTextInformation},
	"PCNT": {"PCNT", "Play counter", TypeTextInformation},
	"POPM": {"POPM", "Popularimeter", TypeTextInformation},
	"POSS": {"POSS", "Position synchronisation frame", TypeTextInformation},
	"RBUF": {"RBUF", "Recommended buffer size", TypeTextInformation},
	"RVA2": {"RVA2", "Relative volume adjustment (2)", TypeTextInformation},
	"RVRB": {"RVRB", "Reverb", TypeTextInformation},
	"SEEK": {"SEEK", "Seek frame", TypeTextInformation},
	"SIGN": {"SIGN", "Signature frame", TypeTextInformation},
	"SYLT": {"SYLT", "Synchronised lyric/text", TypeTextInformation},
	"SYTC": {"SYTC", "Synchronised tempo codes", TypeTextInformation},
	"TALB": {"TALB", "Album/Movie/Show title", TypeTextInformation},
	"TBPM": {"TBPM", "BPM (beats per minute)", TypeTextInformation},
	"TCOM": {"TCOM", "Composer", TypeTextInformation},
	"TCON": {"TCON", "Content type", TypeTextInformation},
	"TCOP": {"TCOP", "Copyright message", TypeTextInformation},
	"TDEN": {"TDEN", "Encoding time", TypeTextInformation},
	"TDLY": {"TDLY", "Playlist delay", TypeTextInformation},
	"TDOR": {"TDOR", "Original release time", TypeTextInformation},
	"TDRC": {"TDRC", "Recording time", TypeTextInformation},
	"TDRL": {"TDRL", "Release time", TypeTextInformation},
	"TDTG": {"TDTG", "Tagging time", TypeTextInformation},
	"TENC": {"TENC", "Encoded by", TypeTextInformation},
	"TEXT": {"TEXT", "Lyricist/Text writer", TypeTextInformation},
	"TFLT": {"TFLT", "File type", TypeTextInformation},
	"TIPL": {"TIPL", "Involved people list", TypeTextInformation},
	"TIT1": {"TIT1", "Content group description", TypeTextInformation},
	"TIT2": {"TIT2", "Title/songname/content description", TypeTextInformation},
	"TIT3": {"TIT3", "Subtitle/Description refinement", TypeTextInformation},
	"TKEY": {"TKEY", "Initial key", TypeTextInformation},
	"TLAN": {"TLAN", "Language(s)", TypeTextInformation},
	"TLEN": {"TLEN", "Length", TypeTextInformation},
	"TMCL": {"TMCL", "Musician credits list", TypeTextInformation},
	"TMED": {"TMED", "Media type", TypeTextInformation},
	"TMOO": {"TMOO", "Mood", TypeTextInformation},
	"TOAL": {"TOAL", "Original album/movie/show title", TypeTextInformation},
	"TOFN": {"TOFN", "Original filename", TypeTextInformation},
	"TOLY": {"TOLY", "Original lyricist(s)/text writer(s)", TypeTextInformation},
	"TOPE": {"TOPE", "Original artist(s)/performer(s)", TypeTextInformation},
	"TOWN": {"TOWN", "File owner/licensee", TypeTextInformation},
	"TPE1": {"TPE1", "Lead performer(s)/Soloist(s)", TypeTextInformation},
	"TPE2": {"TPE2", "Band/orchestra/accompaniment", TypeTextInformation},
	"TPE3": {"TPE3", "Conductor/performer refinement", TypeTextInformation},
	"TPE4": {"TPE4", "Interpreted, remixed, or otherwise modified by", TypeTextInformation},
	"TPOS": {"TPOS", "Part of a set", TypeTextInformation},
	"TPRO": {"TPRO", "Produced notice", TypeTextInformation},
	"TPUB": {"TPUB", "Publisher", TypeTextInformation},
	"TRCK": {"TRCK", "Track number/Position in set", TypeTextInformation},
	"TRSN": {"TRSN", "Internet radio station name", TypeTextInformation},
	"TRSO": {"TRSO", "Internet radio station owner", TypeTextInformation},
	"TSOA": {"TSOA", "Album sort order", TypeTextInformation},
	"TSOP": {"TSOP", "Performer sort order", TypeTextInformation},
	"TSOT": {"TSOT", "Title sort order", TypeTextInformation},
	"TSRC": {"TSRC", "ISRC (international standard recording code)", TypeTextInformation},
	"TSSE": {"TSSE", "Software/Hardware and settings used for encoding", TypeTextInformation},
	"TSST": {"TSST", "Set subtitle", TypeTextInformation},
	"TXXX": {"TXXX", "User defined text information frame", TypeTextInformation},
	"UFID": {"UFID", "Unique file identifier", TypeTextInformation},
	"USER": {"USER", "Terms of use", TypeTextInformation},
	"USLT": {"USLT", "Unsynchronised lyric/text transcription", TypeUnknown},
	"WCOM": {"WCOM", "Commercial information", TypeTextInformation},
	"WCOP": {"WCOP", "Copyright/Legal information", TypeTextInformation},
	"WOAF": {"WOAF", "Official audio file webpage", TypeTextInformation},
	"WOAR": {"WOAR", "Official artist/performer webpage", TypeTextInformation},
	"WOAS": {"WOAS", "Official audio source webpage", TypeTextInformation},
	"WORS": {"WORS", "Official Internet radio station homepage", TypeTextInformation},
	"WPAY": {"WPAY", "Payment", TypeTextInformation},
	"WPUB": {"WPUB", "Publishers official webpage", TypeTextInformation},
	"WXXX": {"WXXX", "User defined URL link frame", TypeTextInformation},
	// iTunes
	"TCMP": {"TCMP", "Part of a compilation", TypeUnknown},
}

// V24 is ID3v2.4 tag reader
type V24 struct {
	frames map[string]interface{}
}

// NewID3V24 will read file and return id3v2.3 tag reader
func NewID3V24(f io.ReadSeeker) (*V24, error) {
	// header
	header := make([]byte, 10)
	n, err := f.Read(header)
	if err != nil {
		return nil, err
	}

	if n != V2HeaderSize {
		return nil, fmt.Errorf("must read '%d' bytes, but read '%d'", V2HeaderSize, n)
	}

	if string(header[:3]) != "ID3" {
		return nil, errors.New("no id3v2 tag at the end of file")
	}

	if header[3] != 4 {
		return nil, errors.New("file id3v2 version is not 2.4.0")
	}

	frmsSize := uint32(header[9]) + uint32(header[8])<<8 + uint32(header[7])<<16 + uint32(header[6])<<32

	// frames
	frames := map[string]interface{}{}
	for t := 0; t < int(frmsSize); {
		frmHeader := make([]byte, 10)
		n, err = f.Read(frmHeader)
		if err != nil {
			return nil, err
		}
		if frmHeader[0]+frmHeader[1]+frmHeader[2]+frmHeader[3] == 0 {
			break
		}
		t += n

		frmSize := uint32(frmHeader[7]) + uint32(frmHeader[6])<<8 + uint32(frmHeader[5])<<16 + uint32(frmHeader[4])<<32

		frmBody := make([]byte, frmSize)
		n, err = f.Read(frmBody)
		if err != nil {
			return nil, err
		}
		t += n
		frames[string(frmHeader[0:4])] = frmBody
	}

	tag := new(V24)
	tag.frames = frames

	return tag, nil
}

func (tag V24) ListFrames() []string {
	keys := make([]string, 0, len(tag.frames))
	for k := range tag.frames {
		keys = append(keys, k)
	}

	return keys
}

func (tag V24) TextFrame(id string) string {
	if frm, ok := tag.frames[id]; ok {
		return trim(string(frm.([]uint8)))
	}

	return ""
}

func (tag V24) NumberFrame(id string) int {
	if frm, ok := tag.frames[id]; ok {
		n, err := strconv.Atoi(trim(string(frm.([]uint8))))
		if err == nil {
			return n
		}
	}

	return 0
}

func (tag V24) TrackPositionFrame(id string) (int, int) {
	if frm, ok := tag.frames[id]; ok {
		s := strings.Split(trim(string(frm.([]uint8))), "/")
		trk, err := strconv.Atoi(s[0])
		if err != nil {
			return 0, 0
		}

		if len(s) == 2 {
			psn, err := strconv.Atoi(s[1])
			if err != nil {
				return trk, 0
			}

			return trk, psn
		}

		return trk, 0
	}

	return 0, 0
}

func (tag V24) GenreFrame(id string) string {
	if frm, ok := tag.frames[id]; ok {
		gids := strings.Trim(trim(string(frm.([]uint8))), "()")
		gid, err := strconv.Atoi(gids)
		if err == nil {
			if gid < len(V1Genres) {
				return V1Genres[gid]
			}
		}
	}

	return ""
}

func (tag V24) LangFrame(id string) (string, string) {
	if frm, ok := tag.frames[id]; ok {
		enc := frm.([]byte)[0]
		lang := string(frm.([]byte)[1:4])
		data := frm.([]byte)[4:]
		return lang, toUTF8(data, Encodings[enc])
	}

	return "", ""
}

func (tag V24) ImageFrame(id string) (image.Image, error) {
	if frm, ok := tag.frames[id]; ok {
		mimeType := ""
		for i := 1; i < len(frm.([]byte)); i++ {
			if frm.([]byte)[i] != 0 {
				continue
			}
			mimeType = string(frm.([]byte)[1:i])
			break
		}
		switch mimeType {
		case "image/jpeg":
			return jpeg.Decode(bytes.NewReader(frm.([]byte)[len(mimeType)+4:]))
		case "image/png":
			return png.Decode(bytes.NewReader(frm.([]byte)[len(mimeType)+4:]))
		default:
			return nil, errors.New("invalid image type")
		}
	}

	return nil, errors.New("frame not found")
}
