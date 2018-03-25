package id3

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"image"
	"bytes"
	"image/png"
	"image/jpeg"
)

var V23Frames = map[string]DeclaredFrame{
	"AENC": {"AENC", "Audio encryption", FrameText},
	"APIC": {"APIC", "Attached picture", FrameImage},
	"COMM": {"COMM", "Comments", FrameText},
	"COMR": {"COMR", "Commercial frame", FrameText},
	"ENCR": {"ENCR", "Encryption method registration", FrameText},
	"EQUA": {"EQUA", "Equalization", FrameText},
	"ETCO": {"ETCO", "Event timing codes", FrameText},
	"GEOB": {"GEOB", "General encapsulated object", FrameText},
	"GRID": {"GRID", "Group identification registration", FrameText},
	"IPLS": {"IPLS", "Involved people list", FrameText},
	"LINK": {"LINK", "Linked information", FrameText},
	"MCDI": {"MCDI", "Music CD identifier", FrameText},
	"MLLT": {"MLLT", "MPEG location lookup table", FrameText},
	"OWNE": {"OWNE", "Ownership frame", FrameText},
	"PRIV": {"PRIV", "Private frame", FrameText},
	"PCNT": {"PCNT", "Play counter", FrameText},
	"POPM": {"POPM", "Popularimeter", FrameText},
	"POSS": {"POSS", "Position synchronisation frame", FrameText},
	"RBUF": {"RBUF", "Recommended buffer size", FrameText},
	"RVAD": {"RVAD", "Relative volume adjustment", FrameText},
	"RVRB": {"RVRB", "Reverb", FrameText},
	"SYLT": {"SYLT", "Synchronized lyric/text", FrameText},
	"SYTC": {"SYTC", "Synchronized tempo codes", FrameText},
	"TALB": {"TALB", "Album/Movie/Show title", FrameText},
	"TBPM": {"TBPM", "BPM (beats per minute)", FrameText},
	"TCOM": {"TCOM", "Composer", FrameText},
	"TCON": {"TCON", "Content type", FrameGenre},
	"TCOP": {"TCOP", "Copyright message", FrameText},
	"TDAT": {"TDAT", "Date", FrameText},
	"TDLY": {"TDLY", "Playlist delay", FrameText},
	"TENC": {"TENC", "Encoded by", FrameText},
	"TEXT": {"TEXT", "Lyricist/Text writer", FrameText},
	"TFLT": {"TFLT", "File type", FrameText},
	"TIME": {"TIME", "Time", FrameText},
	"TIT1": {"TIT1", "Content group description", FrameText},
	"TIT2": {"TIT2", "Title/songname/content description", FrameText},
	"TIT3": {"TIT3", "Subtitle/Description refinement", FrameText},
	"TKEY": {"TKEY", "Initial key", FrameText},
	"TLAN": {"TLAN", "Language(s)", FrameText},
	"TLEN": {"TLEN", "Length", FrameText},
	"TMED": {"TMED", "Media type", FrameText},
	"TOAL": {"TOAL", "Original album/movie/show title", FrameText},
	"TOFN": {"TOFN", "Original filename", FrameText},
	"TOLY": {"TOLY", "Original lyricist(s)/text writer(s)", FrameText},
	"TOPE": {"TOPE", "Original artist(s)/performer(s)", FrameText},
	"TORY": {"TORY", "Original release year", FrameText},
	"TOWN": {"TOWN", "File owner/licensee", FrameText},
	"TPE1": {"TPE1", "Lead performer(s)/Soloist(s)", FrameText},
	"TPE2": {"TPE2", "Band/orchestra/accompaniment", FrameText},
	"TPE3": {"TPE3", "Conductor/performer refinement", FrameText},
	"TPE4": {"TPE4", "Interpreted, remixed, or otherwise modified by", FrameText},
	"TPOS": {"TPOS", "Part of a set", FrameText},
	"TPUB": {"TPUB", "Publisher", FrameText},
	"TRCK": {"TRCK", "Track number/Position in set", FrameTrackPosition},
	"TRDA": {"TRDA", "Recording dates", FrameText},
	"TRSN": {"TRSN", "Internet radio station name", FrameText},
	"TRSO": {"TRSO", "Internet radio station owner", FrameText},
	"TSIZ": {"TSIZ", "Size", FrameText},
	"TSRC": {"TSRC", "ISRC (international standard recording code)", FrameText},
	"TSSE": {"TSSE", "Software/Hardware and settings used for encoding", FrameText},
	"TYER": {"TYER", "Year", FrameText},
	"TXXX": {"TXXX", "User defined text information frame", FrameText},
	"UFID": {"UFID", "Unique file identifier", FrameText},
	"USER": {"USER", "Terms of use", FrameText},
	"USLT": {"USLT", "Unsychronized lyric/text transcription", FrameText},
	"WCOM": {"WCOM", "Commercial information", FrameText},
	"WCOP": {"WCOP", "Copyright/Legal information", FrameText},
	"WOAF": {"WOAF", "Official audio file webpage", FrameText},
	"WOAR": {"WOAR", "Official artist/performer webpage", FrameText},
	"WOAS": {"WOAS", "Official audio source webpage", FrameText},
	"WORS": {"WORS", "Official internet radio station homepage", FrameText},
	"WPAY": {"WPAY", "Payment", FrameText},
	"WPUB": {"WPUB", "Publishers official webpage", FrameText},
	"WXXX": {"WXXX", "User defined URL link frame", FrameText},
	// iTunes
	"TCMP": {"TCMP", "Part of a compilation", FrameBool},
	// extra
	"WFED": {"WFED", "Podcast URL", FrameText},
}

// V23 is ID3v2.3 tag reader
type V23 struct {
	frames map[string]interface{}
}

// NewID3V23 will read file and return id3v2.3 tag reader
func NewID3V23(f io.ReadSeeker) (*V23, error) {
	// header
	header := make([]byte, V2HeaderSize)
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

	if header[3] != 3 {
		return nil, errors.New("file id3v2 version is not 2.3.0")
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

	tag := new(V23)
	tag.frames = frames

	return tag, nil
}

func (tag V23) ListFrames() []string {
	keys := make([]string, 0, len(tag.frames))
	for k := range tag.frames {
		keys = append(keys, k)
	}

	return keys
}

func (tag V23) TextFrame(id string) string {
	if frm, ok := tag.frames[id]; ok {
		return trim(string(frm.([]uint8)))
	}

	return ""
}

func (tag V23) NumberFrame(id string) int {
	if frm, ok := tag.frames[id]; ok {
		n, err := strconv.Atoi(trim(string(frm.([]uint8))))
		if err == nil {
			return n
		}
	}

	return 0
}

func (tag V23) TrackPositionFrame(id string) (int, int) {
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

func (tag V23) GenreFrame(id string) string {
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

func (tag V23) LangFrame(id string) (string, string) {
	if frm, ok := tag.frames[id]; ok {
		t := trim(string(frm.([]uint8)))
		return t[:3], t[3:]
	}

	return "", ""
}

func (tag V23) ImageFrame(id string) (image.Image, error) {
	if frm, ok := tag.frames[id]; ok {
		// enc := frm.([]byte)[0]
		mimeType := ""
		// desc := ""
		var data []byte
		// picType := byte(0)
		for i := 1; i < len(frm.([]byte)); i++ {
			if frm.([]byte)[i] != 0 {
				continue
			}
			mimeType = string(frm.([]byte)[1:i])
			// picType = frm.([]byte)[i]
			for j := i; j < len(frm.([]byte)); j++ {
				if frm.([]byte)[j] != 0 {
					continue
				}

				// desc = toUTF8(frm.([]byte)[i:j],Encodings[enc])
				data = frm.([]byte)[j:]
				break
			}
			break
		}

		switch mimeType {
		case "image/jpeg":
			return jpeg.Decode(bytes.NewReader(data))
		case "image/png":
			return png.Decode(bytes.NewReader(data))
		default:
			return nil, errors.New("invalid image type")
		}
	}

	return nil, errors.New("frame not found")
}
