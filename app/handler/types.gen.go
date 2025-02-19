// Package handler provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
package handler

import (
	pin "github.com/yseto/gion-go/internal/pin"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// AsRead 既読情報
type AsRead struct {
	FeedID uint64 `json:"feed_id"`
	Serial uint64 `json:"serial"`
}

// Authorization ログイン情報
type Authorization struct {
	Autoseen bool   `json:"autoseen"`
	Token    string `json:"token"`
}

// Category カテゴリ一覧
type Category struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// CategoryAndUnreadEntryCount カテゴリごとに未読記事数
type CategoryAndUnreadEntryCount struct {
	Count uint64 `json:"count"`

	// Id category ID
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// CategorySubscription defines model for CategorySubscription.
type CategorySubscription struct {
	CategoryId uint64 `json:"category_id"`

	// HttpStatus 最終アクセス時のレスポンスコード
	HttpStatus string `json:"http_status"`

	// Id フィードID
	FeedID uint64 `json:"id"`

	// Siteurl フィード配信元サイトURL
	Siteurl string `json:"siteurl"`
	Title   string `json:"title"`
}

// ExamineFeed フィード探索におけるフィード詳細
type ExamineFeed struct {
	Date  string `json:"date"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

// ExamineSubscription フィード探索
type ExamineSubscription struct {
	PreviewFeed []ExamineFeed `json:"preview_feed"`
	Success     bool          `json:"success"`
	Title       string        `json:"title"`
	URL         string        `json:"url"`
}

// PinnedItem ピン止めしたアイテム
type PinnedItem struct {
	FeedId   uint64         `json:"feed_id"`
	Serial   uint64         `json:"serial"`
	Title    string         `json:"title"`
	UpdateAt pin.UpdateTime `json:"update_at"`
	Url      string         `json:"url"`
}

// Profile 個人設定
type Profile struct {
	Autoseen           bool   `json:"autoseen"`
	EntryCount         uint64 `json:"entryCount"`
	OnLoginSkipPinList bool   `json:"onLoginSkipPinList"`
	SubstringLength    uint64 `json:"substringLength"`
}

// SimpleResult defines model for SimpleResult.
type SimpleResult struct {
	Result string `json:"result"`
}

// Subscription カテゴリおよび購読一覧
type Subscription struct {
	// Id カテゴリID
	CategoryID uint64 `json:"id"`

	// Name カテゴリ名
	Name string `json:"name"`

	// Subscription カテゴリに属するフィード一覧
	Subscription []CategorySubscription `json:"subscription"`
}

// UnreadEntry カテゴリに属した未読記事一覧
type UnreadEntry struct {
	DateEpoch      uint64       `json:"date_epoch"`
	Description    string       `json:"description"`
	FeedId         uint64       `json:"feed_id"`
	Readflag       pin.ReadFlag `json:"readflag"`
	Serial         uint64       `json:"serial"`
	SiteTitle      string       `json:"site_title"`
	SubscriptionId uint64       `json:"subscription_id"`
	Title          string       `json:"title"`
	Url            string       `json:"url"`
}

// RegisterCategoryJSONBody defines parameters for RegisterCategory.
type RegisterCategoryJSONBody struct {
	Name string `json:"name"`
}

// ChangeSubscriptionJSONBody defines parameters for ChangeSubscription.
type ChangeSubscriptionJSONBody struct {
	Category uint64 `json:"category"`

	// Id Feed ID
	Id uint64 `json:"id"`
}

// ExamineSubscriptionJSONBody defines parameters for ExamineSubscription.
type ExamineSubscriptionJSONBody struct {
	// Url Site URL
	Url string `json:"url"`
}

// LoginJSONBody defines parameters for Login.
type LoginJSONBody struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

// OpmlImportJSONBody defines parameters for OpmlImport.
type OpmlImportJSONBody struct {
	// Xml Opml xml document
	Xml string `json:"xml"`
}

// SetPinJSONBody defines parameters for SetPin.
type SetPinJSONBody struct {
	FeedId   uint64       `json:"feed_id"`
	Readflag pin.ReadFlag `json:"readflag"`
	Serial   uint64       `json:"serial"`
}

// SetAsReadJSONBody defines parameters for SetAsRead.
type SetAsReadJSONBody = []AsRead

// RegisterSubscriptionJSONBody defines parameters for RegisterSubscription.
type RegisterSubscriptionJSONBody struct {
	Category uint64 `json:"category"`

	// Rss RSS Feed URL
	Rss string `json:"rss"`

	// Title Site Title
	Title string `json:"title"`

	// Url Site URL
	Url string `json:"url"`
}

// UpdatePasswordJSONBody defines parameters for UpdatePassword.
type UpdatePasswordJSONBody struct {
	Password    string `json:"password"`
	PasswordOld string `json:"password_old"`
	Passwordc   string `json:"passwordc"`
}

// RegisterCategoryJSONRequestBody defines body for RegisterCategory for application/json ContentType.
type RegisterCategoryJSONRequestBody RegisterCategoryJSONBody

// ChangeSubscriptionJSONRequestBody defines body for ChangeSubscription for application/json ContentType.
type ChangeSubscriptionJSONRequestBody ChangeSubscriptionJSONBody

// ExamineSubscriptionJSONRequestBody defines body for ExamineSubscription for application/json ContentType.
type ExamineSubscriptionJSONRequestBody ExamineSubscriptionJSONBody

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody LoginJSONBody

// OpmlImportJSONRequestBody defines body for OpmlImport for application/json ContentType.
type OpmlImportJSONRequestBody OpmlImportJSONBody

// SetPinJSONRequestBody defines body for SetPin for application/json ContentType.
type SetPinJSONRequestBody SetPinJSONBody

// UpdateProfileJSONRequestBody defines body for UpdateProfile for application/json ContentType.
type UpdateProfileJSONRequestBody = Profile

// SetAsReadJSONRequestBody defines body for SetAsRead for application/json ContentType.
type SetAsReadJSONRequestBody = SetAsReadJSONBody

// RegisterSubscriptionJSONRequestBody defines body for RegisterSubscription for application/json ContentType.
type RegisterSubscriptionJSONRequestBody RegisterSubscriptionJSONBody

// UpdatePasswordJSONRequestBody defines body for UpdatePassword for application/json ContentType.
type UpdatePasswordJSONRequestBody UpdatePasswordJSONBody
