package db

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/yseto/gion-go/internal/db/db"
	"github.com/yseto/gion-go/internal/pin"
)

type UserScopedDB interface {
	CategoryByID(ctx context.Context, ID uint64) (*db.Category, error)
	Category(ctx context.Context) ([]*db.Category, error)
	CategoryAndUnreadEntryCount(ctx context.Context) ([]*db.UnreadEntryCount, error)
	DeleteCategory(ctx context.Context, id uint64) error

	UnreadEntryByCategory(ctx context.Context, categoryID uint64) ([]*db.EntryDetail, error)
	UpdateEntrySeen(ctx context.Context, feedID, serial uint64) error

	PinnedItems(ctx context.Context) ([]*db.PinnedItem, error)
	RemovePinnedItem(ctx context.Context) error

	Subscriptions(ctx context.Context) ([]*db.SubscriptionForUser, error)
	FeedsByCategoryID(ctx context.Context, categoryID uint64) ([]*db.Feed, error)
	DeleteSubscription(ctx context.Context, feedID uint64) error
	UpdateSubscription(ctx context.Context, feedID, categoryID uint64) error

	Profile(ctx context.Context) (*db.UserProfile, error)
	UpdateProfile(ctx context.Context, item db.UserProfile) error
	User(ctx context.Context) (*db.User, error)
	UpdateUserDigest(ctx context.Context, digest string) error

	MustBegin() *db.UserClientTxn
}

type UserClientTxn interface {
	SubscriptionByFeedID(ctx context.Context, feedID uint64) (*db.Subscription, error)
	InsertSubscription(ctx context.Context, feedID, categoryID uint64) error
	CategoryByName(ctx context.Context, name string) (*db.Category, error)
	InsertCategory(ctx context.Context, name string) error

	UpdateEntry(ctx context.Context, feedID, serial uint64, readflag pin.ReadFlag) error
}

func NewUserScopedDB(conn *sqlx.DB, userID uint64) UserScopedDB {
	return &db.UserClient{DB: conn, UserID: userID}
}

type DB interface {
	FeedsByID(ctx context.Context, ids []uint64) ([]*db.Feed, error)
	Feeds(ctx context.Context) ([]*db.Feed, error)

	PurgeReadEntry(ctx context.Context) error

	UserByName(ctx context.Context, name string) (*db.User, error)

	UpdateUserLastLogin(ctx context.Context, userID uint64) error

	InsertUser(ctx context.Context, username, password string) error

	MustBegin() *db.ClientTxn
}

type ClientTxn interface {
	FeedByID(ctx context.Context, id uint64) (*db.Feed, error)
	UpdateFeed(ctx context.Context, item db.Feed) error
	UpdateFeedRSSUrl(ctx context.Context, item db.Feed) error
	UpdateFeedSiteUrl(ctx context.Context, item db.Feed) error
	UpdateFeedWithLatestEntryAt(ctx context.Context, item db.Feed) error
	FeedByWebsubID(ctx context.Context, websubID uint64) (*db.Feed, error)

	UpdateNextSerial(ctx context.Context, feedID uint64) error
	GetNextSerial(ctx context.Context, feedID uint64) (*uint64, error)
	SubscriptionByFeedID(ctx context.Context, feedID uint64) ([]*db.Subscription, error)

	InsertEntry(ctx context.Context, userID, feedID, serial, subscriptionID uint64, pubdate time.Time) error
	InsertStory(ctx context.Context, feedID, serial uint64, title, description, url string) error
	DeleteFeed(ctx context.Context, feedID uint64) error
	ExistEntry(ctx context.Context, feedID, serial uint64) (uint64, error)
	Stories(ctx context.Context) ([]*db.Story, error)
	DeleteStory(ctx context.Context, feedID, serial uint64) error

	FeedByUrl(ctx context.Context, feedUrl, siteUrl string) (*db.Feed, error)
	InsertFeed(ctx context.Context, feedUrl, siteUrl, title string, websubID *uint64) error

	Websubs(ctx context.Context) ([]*db.Websub, error)
	WebsubByID(ctx context.Context, id uint64) (*db.Websub, error)
	WebsubByTopicURL(ctx context.Context, topicURL string) (*db.Websub, error)
	InsertWebsub(ctx context.Context, topicURL string, expiresAt time.Time) error
	UpdateWebsub(ctx context.Context, item db.Websub) error
	DeleteWebsub(ctx context.Context, item db.Websub) error
}

func New(conn *sqlx.DB) DB {
	return &db.Client{DB: conn}
}
