package id3

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

// V1HeaderSize is size of ID3v2.2, ID3v2.3 and ID3v2.4 tag header
const V2HeaderSize = 10

// V2FrameHeaderSize is size of ID3v2.2 tag frame header
const V2FrameHeaderSize = 6

type FrameType int

const (
	TypeUnknown FrameType = iota
	TypeUniqueFileIdentifier
	TypeTextInformation
	TypeURLLink
	TypeInvolvedPeopleList
	TypeMusicCDIdentifier
	TypeEventTimingCodes
	TypeMPEGLocationLookupTable
	TypeSyncedTempoCodes
	TypeUnsychronisedLyricsOrTextTranscription
	TypeSynchronisedLyricsOrText
	TypeComments
	TypeRelativeVolumeAdjustment
	TypeEqualisation
	TypeReverb
	TypeAttachedPicture
	TypeGeneralEncapsulatedObject
	TypePlayCounter
	TypePopularimeter
	TypeRecommendedBufferSize
	TypeEncryptedMetaFrame
	TypeAudioEncryption
	TypeLinkedInformation
)

type Frame interface {
	ID() string
	Size() int
}

type frameBase struct {
	id   string
	size int
}

func (f frameBase) ID() string {
	return f.id
}

func (f frameBase) Size() int {
	return f.size
}

type UnknownFrame struct {
	frameBase
	data []byte
}

func (f UnknownFrame) Data() []byte {
	return f.data
}

type UniqueFileIdentifierFrame struct {
	frameBase
	ownerIdentifier string
	identifier      []byte
}

func (f UniqueFileIdentifierFrame) OwnerIdentifier() string {
	return f.ownerIdentifier
}

func (f UniqueFileIdentifierFrame) Identifier() []byte {
	return f.identifier
}

type TextInformationFrame struct {
	frameBase
	encoding Encoding
	text     string
}

func (f TextInformationFrame) Text() string {
	return f.text
}

type InvolvedPeopleListFrame struct {
	frameBase
	encoding   Encoding
	peopleList []string
}

func (f InvolvedPeopleListFrame) PeopleList() []string {
	return f.peopleList
}

type URLLinkFrame struct {
	frameBase
	url string
}

func (f URLLinkFrame) URL() string {
	return f.url
}

type MusicCDIdentifierFrame struct {
	frameBase
	cdTOC []byte
}

func (f MusicCDIdentifierFrame) CDTOC() []byte {
	return f.cdTOC
}

type TimeStampFormat byte

type EventTimingCodesFrame struct {
	frameBase
	timeStampFormat TimeStampFormat
}

func (f EventTimingCodesFrame) TimeStampFormat() TimeStampFormat {
	return f.timeStampFormat
}

// 4.7.   MPEG location lookup table

// 4.8.   Synced tempo codes

// 4.9.   Unsychronised lyrics/text transcription

// 4.10.   Synchronised lyrics/text

// 4.11.   Comments

// 4.12.   Relative volume adjustment

// 4.13.   Equalisation

// 4.14.   Reverb

type PictureType int

const (
	PictureTypeOther PictureType = iota
	PictureType32x32
	PictureTypeOtherFileIcon
	PictureTypeCoverFront
	PictureTypeCoverBack
	PictureTypeLeafletPage
	PictureTypeMedia
	PictureTypeLeadArtist
	PictureTypeArtist
	PictureTypeConductor
	PictureTypeBandOrOrchestra
	PictureTypeComposer
	PictureTypeLyricist
	PictureTypeRecordingLocation
	PictureTypeDuringRecording
	PictureTypeDuringPerformance
	PictureTypeMovieOrVideoScreenCapture
	PictureTypeABrightColouredFish
	PictureTypeIllustration
	PictureTypeBandOrArtistLogotype
	PictureTypePublisherOrStudioLogotype
)

type AttachedPictureFrame struct {
	frameBase
	textEncoding Encoding
	imageFormat  string
	pictureType  PictureType
	description  string
	pictureData  []byte
}

func (f AttachedPictureFrame) Image() (image.Image, error) {
	switch f.imageFormat {
	case "JPG":
		return jpeg.Decode(bytes.NewReader(f.pictureData))
	case "PNG":
		return png.Decode(bytes.NewReader(f.pictureData))
	default:
		return nil, errors.New("invalid image format")
	}
}

func (f AttachedPictureFrame) Description() string {
	return f.description
}

// 4.16.   General encapsulated object

// 4.17.   Play counter

// 4.18.   Popularimeter

// 4.19.   Recommended buffer size

// 4.20.   Encrypted meta frame

// 4.21.   Audio encryption

// 4.22.   Linked information

type DeclaredFrame struct {
	ID          string
	Description string
	Type        FrameType
}

var V22DeclaredFrames = map[string]DeclaredFrame{
	"BUF": {"BUF", "Recommended buffer size", TypeUnknown},
	"CNT": {"CNT", "Play counter", TypeUnknown},
	"COM": {"COM", "Comments", TypeUnknown},
	"CRA": {"CRA", "Audio encryption", TypeUnknown},
	"CRM": {"CRM", "Encrypted meta frame", TypeUnknown},
	"ETC": {"ETC", "Event timing codes", TypeUnknown},
	"EQU": {"EQU", "Equalization", TypeUnknown},
	"GEO": {"GEO", "General encapsulated object", TypeUnknown},
	"IPL": {"IPL", "Involved people list", TypeInvolvedPeopleList},
	"LNK": {"LNK", "Linked information", TypeUnknown},
	"MCI": {"MCI", "Music CD Identifier", TypeUnknown},
	"MLL": {"MLL", "MPEG location lookup table", TypeUnknown},
	"PIC": {"PIC", "Attached picture", TypeAttachedPicture},
	"POP": {"POP", "Popularimeter", TypeUnknown},
	"REV": {"REV", "Reverb", TypeUnknown},
	"RVA": {"RVA", "Relative volume adjustment", TypeUnknown},
	"SLT": {"SLT", "Synchronized lyric/text", TypeUnknown},
	"STC": {"STC", "Synced tempo codes", TypeUnknown},

	"TAL": {"TAL", "Album/Movie/Show title", TypeTextInformation},
	"TBP": {"TBP", "BPM (Beats Per Minute)", TypeTextInformation},
	"TCM": {"TCM", "Composer", TypeTextInformation},
	"TCO": {"TCO", "Content type", TypeTextInformation},
	"TCR": {"TCR", "Copyright message", TypeTextInformation},
	"TDA": {"TDA", "Date", TypeTextInformation},
	"TDY": {"TDY", "Playlist delay", TypeTextInformation},
	"TEN": {"TEN", "Encoded by", TypeTextInformation},
	"TFT": {"TFT", "File type", TypeTextInformation},
	"TIM": {"TIM", "Time", TypeTextInformation},
	"TKE": {"TKE", "Initial key", TypeTextInformation},
	"TLA": {"TLA", "Language(s)", TypeTextInformation},
	"TLE": {"TLE", "Length", TypeTextInformation},
	"TMT": {"TMT", "Media type", TypeTextInformation},
	"TOA": {"TOA", "Original artist(s)/performer(s)", TypeTextInformation},
	"TOF": {"TOF", "Original filename", TypeTextInformation},
	"TOL": {"TOL", "Original Lyricist(s)/text writer(s)", TypeTextInformation},
	"TOR": {"TOR", "Original release year", TypeTextInformation},
	"TOT": {"TOT", "Original album/Movie/Show title", TypeTextInformation},
	"TP1": {"TP1", "Lead artist(s)/Lead performer(s)/Soloist(s)/Performing group", TypeTextInformation},
	"TP2": {"TP2", "Band/Orchestra/Accompaniment", TypeTextInformation},
	"TP3": {"TP3", "Conductor/Performer refinement", TypeTextInformation},
	"TP4": {"TP4", "Interpreted, remixed, or otherwise modified by", TypeTextInformation},
	"TPA": {"TPA", "Part of a set", TypeTextInformation},
	"TPB": {"TPB", "Publisher", TypeTextInformation},
	"TRC": {"TRC", "ISRC (International Standard Recording Code)", TypeTextInformation},
	"TRD": {"TRD", "Recording dates", TypeTextInformation},
	"TRK": {"TRK", "Track number/Position in set", TypeTextInformation},
	"TSI": {"TSI", "Size", TypeTextInformation},
	"TSS": {"TSS", "Software/hardware and settings used for encoding", TypeTextInformation},
	"TT1": {"TT1", "Content group description", TypeTextInformation},
	"TT2": {"TT2", "Title/Songname/Content description", TypeTextInformation},
	"TT3": {"TT3", "Subtitle/Description refinement", TypeTextInformation},
	"TXT": {"TXT", "Lyricist/text writer", TypeTextInformation},
	"TXX": {"TXX", "User defined text information frame", TypeUnknown},
	"TYE": {"TYE", "Year", TypeTextInformation},
	"TCP": {"TCP", "Part of a compilation", TypeUnknown}, // iTunes

	"UFI": {"UFI", "Unique file identifier", TypeUniqueFileIdentifier},
	"ULT": {"ULT", "Unsychronized lyric/text transcription", TypeUnknown},

	"WAF": {"WAF", "Official audio file webpage", TypeURLLink},
	"WAR": {"WAR", "Official artist/performer webpage", TypeURLLink},
	"WAS": {"WAS", "Official audio source webpage", TypeURLLink},
	"WCM": {"WCM", "Commercial information", TypeURLLink},
	"WCP": {"WCP", "Copyright/Legal information", TypeURLLink},
	"WPB": {"WPB", "Publishers official webpage", TypeURLLink},
	"WXX": {"WXX", "User defined URL link frame", TypeUnknown},
}

// V22 is ID3v2.2 tag reader
type V22 struct {
	frames []Frame
}

// NewID3V22 will read file and return id3v2.2 tag reader
func NewID3V22(f io.ReadSeeker) (*V22, error) {
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

	if header[3] != 2 {
		return nil, errors.New("file id3v2 version is not 2.2.0")
	}

	frames := make([]Frame, 0)
	framesSize := int(uint32(header[9]) + uint32(header[8])<<8 + uint32(header[7])<<16 + uint32(header[6])<<32)

	for t := 0; t < framesSize; {
		frameHeader := make([]byte, V2FrameHeaderSize)
		n, err = f.Read(frameHeader)
		if err != nil {
			return nil, err
		}
		// FIXME
		if frameHeader[0]+frameHeader[1]+frameHeader[2] == 0 {
			break
		}
		t += n

		frameID := string(frameHeader[:3])
		frameSize := int(uint32(frameHeader[5]) + uint32(frameHeader[4])<<8 + uint32(frameHeader[3])<<16)
		frameBody := make([]byte, frameSize)
		n, err = f.Read(frameBody)
		if err != nil {
			return nil, err
		}
		t += n

		frameBase := frameBase{
			id:   frameID,
			size: frameSize,
		}

		df, ok := V22DeclaredFrames[string(frameID)]
		if !ok {
			frame := UnknownFrame{
				frameBase: frameBase,
				data:      frameBody,
			}
			frames = append(frames, frame)
			continue
		}

		switch df.Type {
		case TypeTextInformation:
			frame := TextInformationFrame{
				frameBase: frameBase,
				encoding:  Encoding(frameBody[0]),
				text:      toUTF8(frameBody[1:], Encoding(frameBody[0])),
			}
			frames = append(frames, frame)
		case TypeURLLink:
			frame := URLLinkFrame{
				frameBase: frameBase,
				url:       string(frameBody),
			}
			frames = append(frames, frame)
		case TypeAttachedPicture:
			frame := AttachedPictureFrame{
				frameBase:    frameBase,
				textEncoding: Encoding(frameBody[0]),
				imageFormat:  string(frameBody[1:4]),
				pictureType:  PictureType(frameBody[4]),
			}
			for i := 5; i < frameSize; i++ {
				if frameBody[i] == 0 {
					frame.description = toUTF8(frameBody[5:i], frame.textEncoding)
					frame.pictureData = frameBody[i+1:]
					break
				}
			}
			frames = append(frames, frame)
		default:
			frame := UnknownFrame{
				frameBase: frameBase,
				data:      frameBody,
			}
			frames = append(frames, frame)
		}
	}

	tag := new(V22)
	tag.frames = frames
	return tag, nil
}

func (tag V22) Frames() []Frame {
	return tag.frames
}
