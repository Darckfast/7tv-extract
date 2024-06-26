package types

type Emotes struct {
	ID            string      `json:"id"`
	Platform      string      `json:"platform"`
	Username      string      `json:"username"`
	DisplayName   string      `json:"display_name"`
	LinkedAt      int64       `json:"linked_at"`
	EmoteCapacity int         `json:"emote_capacity"`
	EmoteSetID    interface{} `json:"emote_set_id"`
	EmoteSet      struct {
		ID         string        `json:"id"`
		Name       string        `json:"name"`
		Flags      int           `json:"flags"`
		Tags       []interface{} `json:"tags"`
		Immutable  bool          `json:"immutable"`
		Privileged bool          `json:"privileged"`
		Emotes     []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Flags     int    `json:"flags"`
			Timestamp int64  `json:"timestamp"`
			ActorID   string `json:"actor_id"`
			Data      struct {
				ID        string   `json:"id"`
				Name      string   `json:"name"`
				Flags     int      `json:"flags"`
				Lifecycle int      `json:"lifecycle"`
				State     []string `json:"state"`
				Listed    bool     `json:"listed"`
				Animated  bool     `json:"animated"`
				Owner     struct {
					ID          string   `json:"id"`
					Username    string   `json:"username"`
					DisplayName string   `json:"display_name"`
					AvatarURL   string   `json:"avatar_url"`
					Style       struct{} `json:"style"`
					Roles       []string `json:"roles"`
				} `json:"owner"`
				Host struct {
					URL   string `json:"url"`
					Files []struct {
						Name       string `json:"name"`
						StaticName string `json:"static_name"`
						Width      int    `json:"width"`
						Height     int    `json:"height"`
						FrameCount int    `json:"frame_count"`
						Size       int    `json:"size"`
						Format     string `json:"format"`
					} `json:"files"`
				} `json:"host"`
			} `json:"data"`
		} `json:"emotes"`
		EmoteCount int `json:"emote_count"`
		Capacity   int `json:"capacity"`
		Owner      struct {
			ID          string   `json:"id"`
			Username    string   `json:"username"`
			DisplayName string   `json:"display_name"`
			AvatarURL   string   `json:"avatar_url"`
			Style       struct{} `json:"style"`
			Roles       []string `json:"roles"`
		} `json:"owner"`
	} `json:"emote_set"`
	User struct {
		ID          string   `json:"id"`
		Username    string   `json:"username"`
		DisplayName string   `json:"display_name"`
		CreatedAt   int64    `json:"created_at"`
		AvatarURL   string   `json:"avatar_url"`
		Style       struct{} `json:"style"`
		Editors     []struct {
			ID          string `json:"id"`
			Permissions int    `json:"permissions"`
			Visible     bool   `json:"visible"`
			AddedAt     int64  `json:"added_at"`
		} `json:"editors"`
		Roles       []string `json:"roles"`
		Connections []struct {
			ID            string      `json:"id"`
			Platform      string      `json:"platform"`
			Username      string      `json:"username"`
			DisplayName   string      `json:"display_name"`
			LinkedAt      int64       `json:"linked_at"`
			EmoteCapacity int         `json:"emote_capacity"`
			EmoteSetID    interface{} `json:"emote_set_id"`
			EmoteSet      struct {
				ID         string        `json:"id"`
				Name       string        `json:"name"`
				Flags      int           `json:"flags"`
				Tags       []interface{} `json:"tags"`
				Immutable  bool          `json:"immutable"`
				Privileged bool          `json:"privileged"`
				Capacity   int           `json:"capacity"`
				Owner      interface{}   `json:"owner"`
			} `json:"emote_set"`
		} `json:"connections"`
	} `json:"user"`
}

type ShortEmoteList struct {
	EmoteName  string `json:"emote_name"`
	FullUrl    string `json:"full_url"`
	Extension  string `json:"extension"`
	IsAnimated bool   `json:"is_animated"`
	FullPath   string
	DirPath    string
	OutputPath string
	Size       int
}

type DiscordEmotes []struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	User struct {
		ID                   string      `json:"id"`
		Username             string      `json:"username"`
		Avatar               string      `json:"avatar"`
		Discriminator        string      `json:"discriminator"`
		PublicFlags          int         `json:"public_flags"`
		Flags                int         `json:"flags"`
		Banner               interface{} `json:"banner"`
		AccentColor          interface{} `json:"accent_color"`
		GlobalName           string      `json:"global_name"`
		AvatarDecorationData interface{} `json:"avatar_decoration_data"`
		BannerColor          interface{} `json:"banner_color"`
		Clan                 interface{} `json:"clan"`
	} `json:"user"`
	Roles         []interface{} `json:"roles"`
	RequireColons bool          `json:"require_colons"`
	Managed       bool          `json:"managed"`
	Animated      bool          `json:"animated"`
	Available     bool          `json:"available"`
}
