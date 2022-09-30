package db

import (
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/yseto/gion-go/db/db"
)

type UserScopedDB interface {
	CategoryByID(ID uint64) (*db.Category, error)
	Category() ([]*db.Category, error)
	CategoryAndUnreadEntryCount() ([]*db.UnreadEntryCount, error)
	DeleteCategory(id uint64) error

	UnreadEntryByCategory(categoryID uint64) ([]*db.EntryDetail, error)
	UpdateEntrySeen(feedID, serial uint64) error

	PinnedItems() ([]*db.PinnedItem, error)
	RemovePinnedItem() error

	Subscriptions() ([]*db.SubscriptionForUser, error)
	FeedsByCategoryID(categoryID uint64) ([]*db.Feed, error)
	DeleteSubscription(feedID uint64) error
	UpdateSubscription(feedID, categoryID uint64) error

	Profile() (*db.UserProfile, error)
	UpdateProfile(item db.UserProfile) error
	User() (*db.User, error)
	UpdateUserDigest(digest string) error

	MustBegin() *db.UserClientTxn
}

type UserClientTxn interface {
	FeedByUrl(feedUrl, siteUrl string) (*db.Feed, error)
	InsertFeed(feedUrl, siteUrl, title string) error
	SubscriptionByFeedID(feedID uint64) (*db.Subscription, error)
	InsertSubscription(feedID, categoryID uint64) error
	CategoryByName(name string) (*db.Category, error)
	InsertCategory(name string) error

	UpdateEntry(feedID, serial uint64, readflag db.ReadFlag) error
}

func NewUserScopedDB(conn *sqlx.DB, userID uint64) UserScopedDB {
	return &db.UserClient{conn, userID}
}

type DB interface {
	FeedsByID(ids []uint64) ([]*db.Feed, error)
	FeedsByTerm(term uint64) ([]*db.Feed, error)
	Feeds() ([]*db.Feed, error)

	PurgeReadEntry() error

	UserByName(name string) (*db.User, error)

	UpdateUserLastLogin(userID uint64) error

	InsertUser(username, password string) error

	MustBegin() *db.ClientTxn
}

type ClientTxn interface {
	FeedByID(id uint64) (*db.Feed, error)
	UpdateFeed(item db.Feed) error
	UpdateFeedRSSUrl(item db.Feed) error
	UpdateFeedWithPubdate(item db.Feed) error

	UpdateNextSerial(feedID uint64) error
	GetNextSerial(feedID uint64) (*uint64, error)
	SubscriptionByFeedID(feedID uint64) ([]*db.Subscription, error)

	InsertEntry(userID, feedID, serial, subscriptionID uint64, pubdate time.Time) error
	InsertStory(feedID, serial uint64, title, description, url string) error
	DeleteFeed(feedID uint64) error
	ExistEntry(feedID, serial uint64) (uint64, error)
	Stories() ([]*db.Story, error)
	DeleteStory(feedID, serial uint64) error
}

func New(conn *sqlx.DB) DB {
	return &db.Client{conn}
}
