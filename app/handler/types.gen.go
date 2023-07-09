// Package handler provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
package handler

import (
	"github.com/yseto/gion-go/db/db"
	pin "github.com/yseto/gion-go/internal/pin"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for DeleteSubscriptionFormdataBodySubscription.
const (
	DeleteSubscriptionFormdataBodySubscriptionCategory DeleteSubscriptionFormdataBodySubscription = "category"
	DeleteSubscriptionFormdataBodySubscriptionEntry    DeleteSubscriptionFormdataBodySubscription = "entry"
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
	FeedId   uint64    `json:"feed_id"`
	Serial   uint64    `json:"serial"`
	Title    string    `json:"title"`
	UpdateAt db.MyTime `json:"update_at"`
	Url      string    `json:"url"`
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
	Name interface{} `json:"name"`

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

// ChangeSubscriptionFormdataBody defines parameters for ChangeSubscription.
type ChangeSubscriptionFormdataBody struct {
	Category uint64 `form:"category" json:"category"`

	// Id Feed ID
	Id uint64 `form:"id" json:"id"`
}

// DeleteSubscriptionFormdataBody defines parameters for DeleteSubscription.
type DeleteSubscriptionFormdataBody struct {
	Id uint64 `form:"id" json:"id"`

	// Subscription choose type
	Subscription DeleteSubscriptionFormdataBodySubscription `form:"subscription" json:"subscription"`
}

// DeleteSubscriptionFormdataBodySubscription defines parameters for DeleteSubscription.
type DeleteSubscriptionFormdataBodySubscription string

// ExamineSubscriptionFormdataBody defines parameters for ExamineSubscription.
type ExamineSubscriptionFormdataBody struct {
	// Url Site URL
	Url string `form:"url" json:"url"`
}

// LoginFormdataBody defines parameters for Login.
type LoginFormdataBody struct {
	Id       string `form:"id" json:"id"`
	Password string `form:"password" json:"password"`
}

// OpmlImportFormdataBody defines parameters for OpmlImport.
type OpmlImportFormdataBody struct {
	// Xml Opml xml document
	Xml string `form:"xml" json:"xml"`
}

// RegisterCategoryFormdataBody defines parameters for RegisterCategory.
type RegisterCategoryFormdataBody struct {
	Name string `form:"name" json:"name"`
}

// RegisterSubscriptionFormdataBody defines parameters for RegisterSubscription.
type RegisterSubscriptionFormdataBody struct {
	Category uint64 `form:"category" json:"category"`

	// Rss RSS Feed URL
	Rss string `form:"rss" json:"rss"`

	// Title Site Title
	Title string `form:"title" json:"title"`

	// Url Site URL
	Url string `form:"url" json:"url"`
}

// SetAsReadJSONBody defines parameters for SetAsRead.
type SetAsReadJSONBody = []AsRead

// SetPinFormdataBody defines parameters for SetPin.
type SetPinFormdataBody struct {
	FeedId   uint64       `form:"feed_id" json:"feed_id"`
	Readflag pin.ReadFlag `form:"readflag" json:"readflag"`
	Serial   uint64       `form:"serial" json:"serial"`
}

// UnreadEntryFormdataBody defines parameters for UnreadEntry.
type UnreadEntryFormdataBody struct {
	Category uint64 `form:"category" json:"category"`
}

// UpdatePasswordFormdataBody defines parameters for UpdatePassword.
type UpdatePasswordFormdataBody struct {
	Password    string `form:"password" json:"password"`
	PasswordOld string `form:"password_old" json:"password_old"`
	Passwordc   string `form:"passwordc" json:"passwordc"`
}

// ChangeSubscriptionFormdataRequestBody defines body for ChangeSubscription for application/x-www-form-urlencoded ContentType.
type ChangeSubscriptionFormdataRequestBody ChangeSubscriptionFormdataBody

// DeleteSubscriptionFormdataRequestBody defines body for DeleteSubscription for application/x-www-form-urlencoded ContentType.
type DeleteSubscriptionFormdataRequestBody DeleteSubscriptionFormdataBody

// ExamineSubscriptionFormdataRequestBody defines body for ExamineSubscription for application/x-www-form-urlencoded ContentType.
type ExamineSubscriptionFormdataRequestBody ExamineSubscriptionFormdataBody

// LoginFormdataRequestBody defines body for Login for application/x-www-form-urlencoded ContentType.
type LoginFormdataRequestBody LoginFormdataBody

// OpmlImportFormdataRequestBody defines body for OpmlImport for application/x-www-form-urlencoded ContentType.
type OpmlImportFormdataRequestBody OpmlImportFormdataBody

// RegisterCategoryFormdataRequestBody defines body for RegisterCategory for application/x-www-form-urlencoded ContentType.
type RegisterCategoryFormdataRequestBody RegisterCategoryFormdataBody

// RegisterSubscriptionFormdataRequestBody defines body for RegisterSubscription for application/x-www-form-urlencoded ContentType.
type RegisterSubscriptionFormdataRequestBody RegisterSubscriptionFormdataBody

// SetAsReadJSONRequestBody defines body for SetAsRead for application/json ContentType.
type SetAsReadJSONRequestBody = SetAsReadJSONBody

// SetPinFormdataRequestBody defines body for SetPin for application/x-www-form-urlencoded ContentType.
type SetPinFormdataRequestBody SetPinFormdataBody

// UpdateProfileJSONRequestBody defines body for UpdateProfile for application/json ContentType.
type UpdateProfileJSONRequestBody = Profile

// UnreadEntryFormdataRequestBody defines body for UnreadEntry for application/x-www-form-urlencoded ContentType.
type UnreadEntryFormdataRequestBody UnreadEntryFormdataBody

// UpdatePasswordFormdataRequestBody defines body for UpdatePassword for application/x-www-form-urlencoded ContentType.
type UpdatePasswordFormdataRequestBody UpdatePasswordFormdataBody
