package v24

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xonyagar/id3/lib"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"regexp"
)

// HeaderSize is size of ID3v2.4 tag header
const HeaderSize = 10

// FrameHeaderSize is size of ID3v2.4 tag frame header
const FrameHeaderSize = 10

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

	TypeTermOfUse
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
	encoding lib.Encoding
	text     string
}

func (f TextInformationFrame) Text() string {
	return f.text
}

type TermOfUseFrame struct {
	frameBase
	textEncoding  lib.Encoding
	language      string
	theActualText string
}

func (f TermOfUseFrame) Language() string {
	return f.language
}

func (f TermOfUseFrame) TheActualText() string {
	return f.theActualText
}

type InvolvedPeopleListFrame struct {
	frameBase
	encoding   lib.Encoding
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

type UnsynchronisedLyricsOrTextTranscriptionFrame struct {
	frameBase
	textEncoding      lib.Encoding
	language          string
	contentDescriptor string
	lyricsOrText      string
}

func (f UnsynchronisedLyricsOrTextTranscriptionFrame) Language() string {
	return f.language
}

func (f UnsynchronisedLyricsOrTextTranscriptionFrame) ContentDescriptor() string {
	return f.contentDescriptor
}

func (f UnsynchronisedLyricsOrTextTranscriptionFrame) LyricsOrText() string {
	return f.lyricsOrText
}

// 4.10.   Synchronised lyrics/text

type CommentsFrame struct {
	frameBase
	textEncoding            lib.Encoding
	language                string
	shortContentDescription string
	theActualText           string
}

func (f CommentsFrame) Language() string {
	return f.language
}

func (f CommentsFrame) ShortContentDescription() string {
	return f.shortContentDescription
}

func (f CommentsFrame) TheActualText() string {
	return f.theActualText
}

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
	textEncoding lib.Encoding
	mimeType     string
	pictureType  PictureType
	description  string
	pictureData  []byte
}

func (f AttachedPictureFrame) Image() (image.Image, error) {
	switch f.mimeType {
	case "image/jpeg":
		return jpeg.Decode(bytes.NewReader(f.pictureData))
	case "image/png":
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
type PopularimeterFrame struct {
	frameBase
	emailToUser string
	rating      uint8
	counter     int
}

func (f PopularimeterFrame) EmailToUser() string {
	return f.emailToUser
}

func (f PopularimeterFrame) Rating() uint8 {
	return f.rating
}

func (f PopularimeterFrame) Counter() int {
	return f.counter
}

// 4.19.   Recommended buffer size

// 4.20.   Encrypted meta frame

// 4.21.   Audio encryption

// 4.22.   Linked information

type DeclaredFrame struct {
	ID          string
	Description string
	Type        FrameType
}

var DeclaredFrames = map[string]DeclaredFrame{

	"AENC": {"AENC", "Audio encryption", TypeUnknown},
	"APIC": {"APIC", "Attached picture", TypeAttachedPicture},
	"ASPI": {"ASPI", "Audio seek point index", TypeUnknown},
	"COMM": {"COMM", "Comments", TypeUnknown},
	"COMR": {"COMR", "Commercial frame", TypeUnknown},
	"ENCR": {"ENCR", "Encryption method registration", TypeUnknown},
	"EQU2": {"EQU2", "Equalisation (2)", TypeUnknown},
	"ETCO": {"ETCO", "Event timing codes", TypeUnknown},
	"GEOB": {"GEOB", "General encapsulated object", TypeUnknown},
	"GRID": {"GRID", "Group identification registration", TypeUnknown},
	"LINK": {"LINK", "Linked information", TypeUnknown},
	"MCDI": {"MCDI", "Music CD identifier", TypeUnknown},
	"MLLT": {"MLLT", "MPEG location lookup table", TypeUnknown},
	"OWNE": {"OWNE", "Ownership frame", TypeUnknown},
	"PRIV": {"PRIV", "Private frame", TypeUnknown},
	"PCNT": {"PCNT", "Play counter", TypeUnknown},
	"POPM": {"POPM", "Popularimeter", TypePopularimeter},
	"POSS": {"POSS", "Position synchronisation frame", TypeUnknown},
	"RBUF": {"RBUF", "Recommended buffer size", TypeUnknown},
	"RVA2": {"RVA2", "Relative volume adjustment (2)", TypeUnknown},
	"RVRB": {"RVRB", "Reverb", TypeUnknown},
	"SEEK": {"SEEK", "Seek frame", TypeUnknown},
	"SIGN": {"SIGN", "Signature frame", TypeUnknown},
	"SYLT": {"SYLT", "Synchronised lyric/text", TypeUnknown},
	"SYTC": {"SYTC", "Synchronised tempo codes", TypeUnknown},

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

	"TXXX": {"TXXX", "User defined text information frame", TypeUnknown},

	"UFID": {"UFID", "Unique file identifier", TypeUnknown},
	"USER": {"USER", "Terms of use", TypeUnknown},
	"USLT": {"USLT", "Unsynchronised lyric/text transcription", TypeUnknown},
	"WCOM": {"WCOM", "Commercial information", TypeUnknown},
	"WCOP": {"WCOP", "Copyright/Legal information", TypeUnknown},
	"WOAF": {"WOAF", "Official audio file webpage", TypeUnknown},
	"WOAR": {"WOAR", "Official artist/performer webpage", TypeUnknown},
	"WOAS": {"WOAS", "Official audio source webpage", TypeUnknown},
	"WORS": {"WORS", "Official Internet radio station homepage", TypeUnknown},
	"WPAY": {"WPAY", "Payment", TypeUnknown},
	"WPUB": {"WPUB", "Publishers official webpage", TypeUnknown},
	"WXXX": {"WXXX", "User defined URL link frame", TypeUnknown},
	// iTunes
	"TCMP": {"TCMP", "Part of a compilation", TypeUnknown},
}

// V24 is ID3v2.4 tag reader
type V24 struct {
	frames                    []Frame
	Size                      int
	UnsynchronisationFlag     bool
	ExtendedHeaderFlag        bool
	ExperimentalIndicatorFlag bool
	FooterPresentFlag         bool
}

// New will read file and return id3v2.4 tag reader
func New(f io.ReadSeeker) (*V24, error) {
	header := make([]byte, HeaderSize)
	n, err := f.Read(header)
	if err != nil {
		return nil, err
	}

	if n != HeaderSize {
		return nil, fmt.Errorf("must read '%d' bytes, but read '%d'", HeaderSize, n)
	}

	if string(header[:3]) != "ID3" {
		return nil, errors.New("no id3v2 tag at the end of file")
	}

	if header[3] != 4 {
		return nil, errors.New("file id3v2 version is not 2.4.0")
	}

	frames := make([]Frame, 0)
	framesSize := lib.ByteToInt(header[6:10])
	flag := header[5]

	for t := 0; t < framesSize; {
		frameHeader := make([]byte, FrameHeaderSize)
		n, err = f.Read(frameHeader)
		if err != nil {
			return nil, err
		}
		t += n

		frameID := string(frameHeader[:4])
		if !regexp.MustCompile(`^[0-9A-Z]+$`).MatchString(frameID) {
			break
		}

		frameSize := lib.ByteToInt(frameHeader[4:8])
		// TODO: Frame Flags
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

		df, ok := DeclaredFrames[string(frameID)]
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
				encoding:  lib.Encodings[frameBody[0]],
				text:      lib.ToUTF8(frameBody[1:], lib.Encodings[frameBody[0]]),
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
				textEncoding: lib.Encodings[frameBody[0]],
			}
			for i := 1; i < frameSize; i++ {
				if frameBody[i] == 0 {
					frame.mimeType = string(frameBody[1:i])
					frame.pictureType = PictureType(frameBody[i+1])

					for j := i + 2; j < frameSize; j+=frame.textEncoding.Size {
						if frameBody[j] == 0 {
							frame.description = lib.ToUTF8(frameBody[i+2:j], frame.textEncoding)
							frame.pictureData = frameBody[j+frame.textEncoding.Size:]

							break
						}
					}

					break
				}
			}
			frames = append(frames, frame)
		case TypeUnsychronisedLyricsOrTextTranscription:
			frame := UnsynchronisedLyricsOrTextTranscriptionFrame{
				frameBase:    frameBase,
				textEncoding: lib.Encodings[frameBody[0]],
				language:     string(frameBody[1:4]),
			}

			for i := 4; i < frameSize; i+=frame.textEncoding.Size {
				if frameBody[i] == 0 {
					frame.contentDescriptor = lib.ToUTF8(frameBody[4:i], frame.textEncoding)
					frame.lyricsOrText = lib.ToUTF8(frameBody[i+frame.textEncoding.Size:], frame.textEncoding)
					break
				}
			}
			frames = append(frames, frame)
		case TypeComments:
			frame := CommentsFrame{
				frameBase:    frameBase,
				textEncoding: lib.Encodings[frameBody[0]],
				language:     string(frameBody[1:4]),
			}

			for i := 4; i < frameSize; i+=frame.textEncoding.Size {
				if frameBody[i] == 0 {
					frame.shortContentDescription = lib.ToUTF8(frameBody[4:i], frame.textEncoding)
					frame.theActualText = lib.ToUTF8(frameBody[i+frame.textEncoding.Size:], frame.textEncoding)
					break
				}
			}
			frames = append(frames, frame)
		case TypeTermOfUse:
			frame := TermOfUseFrame{
				frameBase:     frameBase,
				textEncoding:  lib.Encodings[frameBody[0]],
				language:      string(frameBody[1:4]),
				theActualText: lib.ToUTF8(frameBody[4:], lib.Encodings[frameBody[0]]),
			}
			frames = append(frames, frame)
		case TypePopularimeter:
			frame := PopularimeterFrame{
				frameBase:     frameBase,
			}
			for i:= 0; i < framesSize; i++ {
				if frameBody[i] == 0 {
					frame.emailToUser = string(frameBody[:i])
					frame.rating = frameBody[i+1]
					frame.counter = lib.ByteToInt(frameBody[i+2:])
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

	tag := new(V24)
	tag.frames = frames
	tag.Size = framesSize
	// Flags
	tag.UnsynchronisationFlag = flag&128 == 1
	tag.ExtendedHeaderFlag = flag&64 == 1
	tag.ExperimentalIndicatorFlag = flag&32 == 1
	tag.FooterPresentFlag = flag&16 == 1
	return tag, nil
}

func (tag V24) Frames() []Frame {
	return tag.frames
}