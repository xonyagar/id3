package id3

import (
	"errors"
	"fmt"
	"io"
	"image"
	"image/jpeg"
	"bytes"
	"image/png"
	"strings"
	"strconv"
)

// V1HeaderSize is size of ID3v2.2, ID3v2.3 and ID3v2.4 tag header
const V2HeaderSize = 10

// V2FrameHeaderSize is size of ID3v2.2 tag frame header
const V2FrameHeaderSize = 6

type FrameType int

const (
	FrameText          FrameType = iota
	FrameNumber
	FrameImage
	FrameGenre
	FrameTrackPosition
	FrameLang
	FrameBool
)

type DeclaredFrame struct {
	ID          string
	Description string
	Type        FrameType
}

var V22Frames = map[string]DeclaredFrame{
	"BUF": {"BUF", "Recommended buffer size", FrameText},
	"CNT": {"CNT", "Play counter", FrameText},
	"COM": {"COM", "Comments", FrameLang},
	"CRA": {"CRA", "Audio encryption", FrameText},
	"CRM": {"CRM", "Encrypted meta frame", FrameText},
	"ETC": {"ETC", "Event timing codes", FrameText},
	"EQU": {"EQU", "Equalization", FrameText},
	"GEO": {"GEO", "General encapsulated object", FrameText},
	"IPL": {"IPL", "Involved people list", FrameText},
	"LNK": {"LNK", "Linked information", FrameText},
	"MCI": {"MCI", "Music CD Identifier", FrameText},
	"MLL": {"MLL", "MPEG location lookup table", FrameText},
	"PIC": {"PIC", "Attached picture", FrameImage},
	"POP": {"POP", "Popularimeter", FrameText},
	"REV": {"REV", "Reverb", FrameText},
	"RVA": {"RVA", "Relative volume adjustment", FrameText},
	"SLT": {"SLT", "Synchronized lyric/text", FrameText},
	"STC": {"STC", "Synced tempo codes", FrameText},
	"TAL": {"TAL", "Album/Movie/Show title", FrameText},
	"TBP": {"TBP", "BPM (Beats Per Minute)", FrameText},
	"TCM": {"TCM", "Composer", FrameText},
	"TCO": {"TCO", "Content type", FrameGenre},
	"TCR": {"TCR", "Copyright message", FrameText},
	"TDA": {"TDA", "Date", FrameText},
	"TDY": {"TDY", "Playlist delay", FrameText},
	"TEN": {"TEN", "Encoded by", FrameText},
	"TFT": {"TFT", "File type", FrameText},
	"TIM": {"TIM", "Time", FrameText},
	"TKE": {"TKE", "Initial key", FrameText},
	"TLA": {"TLA", "Language(s)", FrameText},
	"TLE": {"TLE", "Length", FrameNumber},
	"TMT": {"TMT", "Media type", FrameText},
	"TOA": {"TOA", "Original artist(s)/performer(s)", FrameText},
	"TOF": {"TOF", "Original filename", FrameText},
	"TOL": {"TOL", "Original Lyricist(s)/text writer(s)", FrameText},
	"TOR": {"TOR", "Original release year", FrameText},
	"TOT": {"TOT", "Original album/Movie/Show title", FrameText},
	"TP1": {"TP1", "Lead artist(s)/Lead performer(s)/Soloist(s)/Performing group", FrameText},
	"TP2": {"TP2", "Band/Orchestra/Accompaniment", FrameText},
	"TP3": {"TP3", "Conductor/Performer refinement", FrameText},
	"TP4": {"TP4", "Interpreted, remixed, or otherwise modified by", FrameText},
	"TPA": {"TPA", "Part of a set", FrameText},
	"TPB": {"TPB", "Publisher", FrameText},
	"TRC": {"TRC", "ISRC (International Standard Recording Code)", FrameText},
	"TRD": {"TRD", "Recording dates", FrameText},
	"TRK": {"TRK", "Track number/Position in set", FrameTrackPosition},
	"TSI": {"TSI", "Size", FrameText},
	"TSS": {"TSS", "Software/hardware and settings used for encoding", FrameText},
	"TT1": {"TT1", "Content group description", FrameText},
	"TT2": {"TT2", "Title/Songname/Content description", FrameText},
	"TT3": {"TT3", "Subtitle/Description refinement", FrameText},
	"TXT": {"TXT", "Lyricist/text writer", FrameText},
	"TXX": {"TXX", "User defined text information frame", FrameText},
	"TYE": {"TYE", "Year", FrameText},
	"UFI": {"UFI", "Unique file identifier", FrameText},
	"ULT": {"ULT", "Unsychronized lyric/text transcription", FrameLang},
	"WAF": {"WAF", "Official audio file webpage", FrameText},
	"WAR": {"WAR", "Official artist/performer webpage", FrameText},
	"WAS": {"WAS", "Official audio source webpage", FrameText},
	"WCM": {"WCM", "Commercial information", FrameText},
	"WCP": {"WCP", "Copyright/Legal information", FrameText},
	"WPB": {"WPB", "Publishers official webpage", FrameText},
	"WXX": {"WXX", "User defined URL link frame", FrameText},
	// iTunes
	"TCP": {"TCP", "Part of a compilation", FrameBool},
}

// V22 is ID3v2.2 tag reader
type V22 struct {
	frames map[string]interface{}
}

// NewID3V22 will read file and return id3v2.2 tag reader
func NewID3V22(f io.ReadSeeker) (*V22, error) {
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

	if header[3] != 2 {
		return nil, errors.New("file id3v2 version is not 2.2.0")
	}

	frames := map[string]interface{}{}
	frmsSize := int(header[9]) + int(header[8])*128 + int(header[7])*128*128 + int(header[6])*128*128*128

	for t := 0; t < frmsSize; {
		frmHeader := make([]byte, V2FrameHeaderSize)
		n, err = f.Read(frmHeader)
		if err != nil {
			return nil, err
		}
		if frmHeader[0]+frmHeader[1]+frmHeader[2] == 0 {
			break
		}
		t += n

		frmSize := uint32(frmHeader[5]) + uint32(frmHeader[4])<<8 + uint32(frmHeader[3])<<16
		frmBody := make([]byte, frmSize)
		n, err = f.Read(frmBody)
		if err != nil {
			return nil, err
		}
		t += n
		frames[string(frmHeader[0:3])] = frmBody
	}

	tag := new(V22)
	tag.frames = frames
	return tag, nil
}

func (tag V22) ListFrames() []string {
	keys := make([]string, 0, len(tag.frames))
	for k := range tag.frames {
		keys = append(keys, k)
	}

	return keys
}

func (tag V22) TextFrame(id string) string {
	if frm, ok := tag.frames[id]; ok {
		return trim(string(frm.([]uint8)))
	}

	return ""
}

func (tag V22) BoolFrame(id string) bool {
	if frm, ok := tag.frames[id]; ok {
		return trim(string(frm.([]uint8))) == "1"
	}

	return false
}

func (tag V22) LangFrame(id string) (string, string) {
	if frm, ok := tag.frames[id]; ok {
		t := trim(string(frm.([]uint8)))
		return t[:3], t[3:]
	}

	return "", ""
}

func (tag V22) NumberFrame(id string) int {
	if frm, ok := tag.frames[id]; ok {
		n, err := strconv.Atoi(trim(string(frm.([]uint8))))
		if err == nil {
			return n
		}
	}

	return 0
}

func (tag V22) TrackPositionFrame(id string) (int, int) {
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

func (tag V22) GenreFrame(id string) string {
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

func (tag V22) ImageFrame(id string) (image.Image, error) {
	if frm, ok := tag.frames[id]; ok {
		switch string(frm.([]byte)[1:4]) {
		case "JPG":
			return jpeg.Decode(bytes.NewReader(frm.([]byte)[6:]))
		case "PNG":
			return png.Decode(bytes.NewReader(frm.([]byte)[6:]))
		default:
			return nil, errors.New("invalid image type")
		}
	}

	return nil, errors.New("frame not found")
}
