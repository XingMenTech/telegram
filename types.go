package telegram

import (
	"io"
	"log"
)

const (
	MessageTypeText     = "text"     //文本消息
	MessageTypeAudio    = "audio"    //语音
	MessageTypeDocument = "document" //
	MessageTypeSticker  = "sticker"  //表情
	MessageTypeVideo    = "video"    //视频
	MessageTypePhoto    = "photo"    // 图片
)

type Store interface {
	RPush(value string) error
	BLPop() (string, error)
	Size() int64
	Close() error
}

type Log struct {
	*log.Logger
}

func NewLog(out io.Writer) *Log {
	d := new(Log)
	d.Logger = log.New(out, "[TG]", log.LstdFlags)
	return d
}

// User represents a Telegram user or bot
type User struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

// Chat represents a Telegram chat
type Chat struct {
	ID                    int64            `json:"id"`
	Type                  string           `json:"type"` // “private”, “group”, “supergroup”, or “channel”
	Title                 string           `json:"title,omitempty"`
	Username              string           `json:"username,omitempty"`
	FirstName             string           `json:"first_name,omitempty"`
	LastName              string           `json:"last_name,omitempty"`
	Photo                 *ChatPhoto       `json:"photo,omitempty"`
	Description           string           `json:"description,omitempty"`
	InviteLink            string           `json:"invite_link,omitempty"`
	PinnedMessage         *Message         `json:"pinned_message,omitempty"`
	Permissions           *ChatPermissions `json:"permissions,omitempty"`
	SlowModeDelay         int              `json:"slow_mode_delay,omitempty"`
	MessageAutoDeleteTime int              `json:"message_auto_delete_time,omitempty"`
	HasProtectedContent   bool             `json:"has_protected_content,omitempty"`
	StickerSetName        string           `json:"sticker_set_name,omitempty"`
	CanSetStickerSet      bool             `json:"can_set_sticker_set,omitempty"`
	LinkedChatID          int64            `json:"linked_chat_id,omitempty"`
	Location              *ChatLocation    `json:"location,omitempty"`
}

// Message represents a message
type Message struct {
	MessageID             int                `json:"message_id"`
	From                  *User              `json:"from,omitempty"`
	Date                  int                `json:"date"`
	Chat                  Chat               `json:"chat"`
	ForwardFrom           *User              `json:"forward_from,omitempty"`
	ForwardFromChat       *Chat              `json:"forward_from_chat,omitempty"`
	ForwardFromMessageID  int                `json:"forward_from_message_id,omitempty"`
	ForwardSignature      string             `json:"forward_signature,omitempty"`
	ForwardDate           int                `json:"forward_date,omitempty"`
	ReplyToMessage        *Message           `json:"reply_to_message,omitempty"`
	ViaBot                *User              `json:"via_bot,omitempty"`
	EditDate              int                `json:"edit_date,omitempty"`
	MediaGroupID          string             `json:"media_group_id,omitempty"`
	AuthorSignature       string             `json:"author_signature,omitempty"`
	Text                  string             `json:"text,omitempty"`
	Entities              []MessageEntity    `json:"entities,omitempty"`
	CaptionEntities       []MessageEntity    `json:"caption_entities,omitempty"`
	Audio                 *Audio             `json:"audio,omitempty"`
	Document              *Document          `json:"document,omitempty"`
	Photo                 []PhotoSize        `json:"photo,omitempty"`
	Sticker               *Sticker           `json:"sticker,omitempty"`
	Video                 *Video             `json:"video,omitempty"`
	Voice                 *Voice             `json:"voice,omitempty"`
	VideoNote             *VideoNote         `json:"video_note,omitempty"`
	Caption               string             `json:"caption,omitempty"`
	Contact               *Contact           `json:"contact,omitempty"`
	Location              *Location          `json:"location,omitempty"`
	Venue                 *Venue             `json:"venue,omitempty"`
	Poll                  *Poll              `json:"poll,omitempty"`
	Dice                  *Dice              `json:"dice,omitempty"`
	NewChatMembers        []User             `json:"new_chat_members,omitempty"`
	LeftChatMember        *User              `json:"left_chat_member,omitempty"`
	NewChatTitle          string             `json:"new_chat_title,omitempty"`
	NewChatPhoto          []PhotoSize        `json:"new_chat_photo,omitempty"`
	DeleteChatPhoto       bool               `json:"delete_chat_photo,omitempty"`
	GroupChatCreated      bool               `json:"group_chat_created,omitempty"`
	SupergroupChatCreated bool               `json:"supergroup_chat_created,omitempty"`
	ChannelChatCreated    bool               `json:"channel_chat_created,omitempty"`
	MigrateToChatID       int64              `json:"migrate_to_chat_id,omitempty"`
	MigrateFromChatID     int64              `json:"migrate_from_chat_id,omitempty"`
	PinnedMessage         *Message           `json:"pinned_message,omitempty"`
	Invoice               *Invoice           `json:"invoice,omitempty"`
	SuccessfulPayment     *SuccessfulPayment `json:"successful_payment,omitempty"`
	ConnectedWebsite      string             `json:"connected_website,omitempty"`
}

// MessageEntity represents one special entity in a text message
type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	URL    string `json:"url,omitempty"`
	User   *User  `json:"user,omitempty"`
}

// PhotoSize represents one size of a photo or a file/sticker thumbnail
type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size,omitempty"`
}

// Audio represents an audio file to be treated as music by the Telegram clients
type Audio struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Duration     int        `json:"duration"`
	Performer    string     `json:"performer,omitempty"`
	Title        string     `json:"title,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	FileSize     int        `json:"file_size,omitempty"`
	Thumb        *PhotoSize `json:"thumb,omitempty"`
}

// Document represents a general file
type Document struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Thumb        *PhotoSize `json:"thumb,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	FileSize     int        `json:"file_size,omitempty"`
}

// Video represents a video file
type Video struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Duration     int        `json:"duration"`
	Thumb        *PhotoSize `json:"thumb,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	FileSize     int        `json:"file_size,omitempty"`
}

// Voice represents a voice note
type Voice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
}

// VideoNote represents a video message
type VideoNote struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Length       int        `json:"length"`
	Duration     int        `json:"duration"`
	Thumb        *PhotoSize `json:"thumb,omitempty"`
	FileSize     int        `json:"file_size,omitempty"`
}

// Contact represents a phone contact
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	UserID      int64  `json:"user_id,omitempty"`
}

// Location represents a point on the map
type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// Venue represents a venue
type Venue struct {
	Location     Location `json:"location"`
	Title        string   `json:"title"`
	Address      string   `json:"address"`
	FoursquareID string   `json:"foursquare_id,omitempty"`
}

// Poll represents a native poll
type Poll struct {
	ID                    string       `json:"id"`
	Question              string       `json:"question"`
	Options               []PollOption `json:"options"`
	TotalVoterCount       int          `json:"total_voter_count"`
	IsClosed              bool         `json:"is_closed"`
	IsAnonymous           bool         `json:"is_anonymous"`
	Type                  string       `json:"type"`
	AllowsMultipleAnswers bool         `json:"allows_multiple_answers"`
}

// PollOption represents an option in a poll
type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}

// Dice represents an animated emoji that displays a random value
type Dice struct {
	Emoji string `json:"emoji"`
	Value int    `json:"value"`
}

// ChatPhoto represents a chat photo
type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`
	BigFileID         string `json:"big_file_id"`
	BigFileUniqueID   string `json:"big_file_unique_id"`
}

// ChatPermissions describes actions that a non-administrator user is allowed to take in a chat
type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages,omitempty"`
	CanSendMediaMessages  bool `json:"can_send_media_messages,omitempty"`
	CanSendPolls          bool `json:"can_send_polls,omitempty"`
	CanSendOtherMessages  bool `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"`
	CanChangeInfo         bool `json:"can_change_info,omitempty"`
	CanInviteUsers        bool `json:"can_invite_users,omitempty"`
	CanPinMessages        bool `json:"can_pin_messages,omitempty"`
}

// ChatLocation represents a location to which a chat is connected
type ChatLocation struct {
	Location Location `json:"location"`
	Address  string   `json:"address"`
}

// Update represents an incoming update
type Update struct {
	UpdateID          int64              `json:"update_id"`
	Message           *Message           `json:"message,omitempty"`
	EditedMessage     *Message           `json:"edited_message,omitempty"`
	ChannelPost       *Message           `json:"channel_post,omitempty"`
	EditedChannelPost *Message           `json:"edited_channel_post,omitempty"`
	CallbackQuery     *CallbackQuery     `json:"callback_query,omitempty"`
	ShippingQuery     *ShippingQuery     `json:"shipping_query,omitempty"`
	PreCheckoutQuery  *PreCheckoutQuery  `json:"pre_checkout_query,omitempty"`
	Poll              *Poll              `json:"poll,omitempty"`
	PollAnswer        *PollAnswer        `json:"poll_answer,omitempty"`
	MyChatMember      *ChatMemberUpdated `json:"my_chat_member,omitempty"`
	ChatMember        *ChatMemberUpdated `json:"chat_member,omitempty"`
	ChatJoinRequest   *ChatJoinRequest   `json:"chat_join_request,omitempty"`
}

// CallbackQuery represents an incoming callback query from a callback button in an inline keyboard
type CallbackQuery struct {
	ID              string   `json:"id"`
	From            User     `json:"from"`
	Message         *Message `json:"message,omitempty"`
	InlineMessageID string   `json:"inline_message_id,omitempty"`
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data,omitempty"`
	GameShortName   string   `json:"game_short_name,omitempty"`
}

// ShippingQuery contains information about an incoming shipping query
type ShippingQuery struct {
	ID              string          `json:"id"`
	From            User            `json:"from"`
	InvoicePayload  string          `json:"invoice_payload"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}

// PreCheckoutQuery contains information about an incoming pre-checkout query
type PreCheckoutQuery struct {
	ID               string     `json:"id"`
	From             User       `json:"from"`
	Currency         string     `json:"currency"`
	TotalAmount      int        `json:"total_amount"`
	InvoicePayload   string     `json:"invoice_payload"`
	ShippingOptionID string     `json:"shipping_option_id,omitempty"`
	OrderInfo        *OrderInfo `json:"order_info,omitempty"`
}

// PollAnswer represents an answer of a user in a non-anonymous poll
type PollAnswer struct {
	PollID    string `json:"poll_id"`
	User      User   `json:"user"`
	OptionIDs []int  `json:"option_ids"`
}

// ChatMemberUpdated represents changes in the status of a chat member
type ChatMemberUpdated struct {
	Chat          Chat            `json:"chat"`
	From          User            `json:"from"`
	Date          int             `json:"date"`
	OldChatMember ChatMember      `json:"old_chat_member"`
	NewChatMember ChatMember      `json:"new_chat_member"`
	InviteLink    *ChatInviteLink `json:"invite_link,omitempty"`
}

// ChatJoinRequest represents a join request sent to a chat
type ChatJoinRequest struct {
	Chat       Chat            `json:"chat"`
	From       User            `json:"from"`
	Date       int             `json:"date"`
	Bio        string          `json:"bio,omitempty"`
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
}

// ChatMember represents a chat member
type ChatMember struct {
	User   User   `json:"user"`
	Status string `json:"status"` // "creator", "administrator", "member", "restricted", "left", "kicked"
}

// ChatInviteLink represents an invite link for a chat
type ChatInviteLink struct {
	InviteLink              string `json:"invite_link"`
	Creator                 User   `json:"creator"`
	CreatesJoinRequest      bool   `json:"creates_join_request"`
	IsPrimary               bool   `json:"is_primary"`
	IsRevoked               bool   `json:"is_revoked"`
	Name                    string `json:"name,omitempty"`
	ExpireDate              int    `json:"expire_date,omitempty"`
	MemberLimit             int    `json:"member_limit,omitempty"`
	PendingJoinRequestCount int    `json:"pending_join_request_count,omitempty"`
}

// OrderInfo represents information about an order
type OrderInfo struct {
	Name            string           `json:"name,omitempty"`
	PhoneNumber     string           `json:"phone_number,omitempty"`
	Email           string           `json:"email,omitempty"`
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
}

// ShippingAddress represents a shipping address
type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

// Invoice contains basic information about an invoice
type Invoice struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	StartParameter string `json:"start_parameter"`
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
}

// SuccessfulPayment contains basic information about a successful payment
type SuccessfulPayment struct {
	Currency                string     `json:"currency"`
	TotalAmount             int        `json:"total_amount"`
	InvoicePayload          string     `json:"invoice_payload"`
	ShippingOptionID        string     `json:"shipping_option_id,omitempty"`
	OrderInfo               *OrderInfo `json:"order_info,omitempty"`
	TelegramPaymentChargeID string     `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID string     `json:"provider_payment_charge_id"`
}

// File represents a file ready to be downloaded
type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
}

// InputMedia represents the content of a media message to be sent
type InputMedia struct {
	Type      string `json:"type"`
	Media     string `json:"media"`
	Caption   string `json:"caption,omitempty"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// Sticker represents a sticker
type Sticker struct {
	FileID       string        `json:"file_id"`
	FileUniqueID string        `json:"file_unique_id"`
	Width        int           `json:"width"`
	Height       int           `json:"height"`
	IsAnimated   bool          `json:"is_animated"`
	IsVideo      bool          `json:"is_video"`
	Thumb        *PhotoSize    `json:"thumb,omitempty"`
	Emoji        string        `json:"emoji,omitempty"`
	SetName      string        `json:"set_name,omitempty"`
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
	FileSize     int           `json:"file_size,omitempty"`
}

// MaskPosition describes the position on faces where a mask should be placed by default
type MaskPosition struct {
	Point  string  `json:"point"`
	XShift float64 `json:"x_shift"`
	YShift float64 `json:"y_shift"`
	Scale  float64 `json:"scale"`
}
