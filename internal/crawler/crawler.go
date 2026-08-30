package crawler

import (
	"context"
	"database/sql"
	"log/slog"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/yseto/gion-go/internal/client"
	"github.com/yseto/gion-go/internal/config"
	"github.com/yseto/gion-go/internal/crawler/revisit"
	"github.com/yseto/gion-go/internal/db"
	dbType "github.com/yseto/gion-go/internal/db/db"
)

const (
	crawlerFactor      = 30
	minRevisitDuration = 30 * time.Minute
	maxRevisitDuration = 7 * 24 * time.Hour

	titleLength = 80
)

func toDbTypeCacheJson(u *client.Response) dbType.CacheJson {
	if u == nil {
		return dbType.CacheJson{}
	}
	return dbType.CacheJson{Etag: u.Cache.Etag, Modified: u.Cache.Modified}
}

func Task(ctx context.Context, feedID uint64) error {
	logger := slog.Default().With("task", "crawler").With("feedID", feedID)

	logger.Info("Start Task")

	cfg, _ := config.ReadConfig()

	dbConn, err := db.Open(cfg)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	dbc := db.New(dbConn)
	tx := dbc.MustBegin()

	feedRow, err := tx.FeedByID(ctx, feedID)
	if err != nil && err == sql.ErrNoRows {
		tx.Rollback()
		return nil
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	revisitManager := revisit.New(time.Now(), crawlerFactor, minRevisitDuration, maxRevisitDuration)

	res, err := client.GetWithCache(
		ctx,
		feedRow.URL,
		client.Cache{
			Etag:     feedRow.Cache.Etag,
			Modified: feedRow.Cache.Modified,
		},
	)

	// 結果が得られない場合、次の対象を処理する
	if err != nil || res.StatusCode >= 400 && res.StatusCode < 600 {
		feedRow.HTTPStatus = "404"
		feedRow.SetRevisitAt(revisitManager.RevisitAtMaxRevisitDuration())
		feedRow.Cache = toDbTypeCacheJson(res)
		logger.Error("failed request", "type", "request", "url", feedRow.URL, "error", err)
		if err = tx.UpdateFeed(ctx, *feedRow); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}

	// リダイレクト先を保存する
	if res.Location != "" {
		logger.Info("redirect", "type", "feed", "url", feedRow.URL, "newURL", res.Location)
		feedRow.URL = res.Location
		if err = tx.UpdateFeedRSSUrl(ctx, *feedRow); err != nil {
			tx.Rollback()
			return err
		}

		if resSite, err := client.GetRedirectURL(feedRow.SiteURL); err == nil {
			logger.Info("redirect", "type", "site", "url", feedRow.SiteURL, "newURL", resSite)
			feedRow.SiteURL = resSite
			if err = tx.UpdateFeedSiteUrl(ctx, *feedRow); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// 304 であることを記録
	if res.StatusCode == 304 {
		feedRow.HTTPStatus = "304"
		logger.Info("Not Modified")

		if ts := feedRow.GetLatestEntryAt(); ts != nil {
			feedRow.SetRevisitAt(revisitManager.RevisitAt(*ts))
		} else {
			feedRow.SetRevisitAt(revisitManager.RevisitAtMaxRevisitDuration())
		}

		// 304 応答で、ヘッダによるレスポンス返却がある場合にのみ上書きをする
		if res.Cache.Etag != "" || res.Cache.Modified != "" {
			feedRow.Cache = toDbTypeCacheJson(res)
		}

		if err = tx.UpdateFeed(ctx, *feedRow); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}

	feed, err := gofeed.NewParser().ParseString(res.Content)
	if err != nil {
		feedRow.Cache = toDbTypeCacheJson(res)
		feedRow.HTTPStatus = strconv.Itoa(res.StatusCode)
		feedRow.SetRevisitAt(revisitManager.RevisitAtMaxRevisitDuration())
		logger.Error("failed parse", "error", err)
		if err = tx.UpdateFeed(ctx, *feedRow); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}

	// logger.Info(strconv.Itoa(res.StatusCode))

	var items []*gofeed.Item
	for i := range feed.Items {
		if feed.Items[i].PublishedParsed == nil {
			continue
		}
		items = append(items, feed.Items[i])
	}

	// 日付に関するフィールドがすべてのアイテムにない場合、エラーとする
	if len(items) == 0 && len(items) != len(feed.Items) {
		feedRow.Cache = toDbTypeCacheJson(res)
		feedRow.HTTPStatus = strconv.Itoa(res.StatusCode)
		feedRow.SetRevisitAt(revisitManager.RevisitAtMaxRevisitDuration())
		logger.Error("missing <pubdate>", "type", "feed")
		if err = tx.UpdateFeed(ctx, *feedRow); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}

	// フィードのエントリを日付順 新 -> 古 にする
	sort.SliceStable(items, func(i, j int) bool { return items[i].PublishedParsed.After(*items[j].PublishedParsed) })

	toleranceTime := time.Now().Add(7 * 24 * time.Hour)

	base, err := url.Parse(feedRow.URL)
	if err != nil {
		return err
	}

	// レコードのフィールドに有効な時刻がなければ現在時刻を用いる
	var existedLatestEntryAt = time.Now()

	// 取り込み済みエントリの日時、新規の場合はnil
	var latestEntryAt *time.Time

	if ts := feedRow.GetLatestEntryAt(); ts != nil {
		existedLatestEntryAt = *ts
		latestEntryAt = ts
	}

	for i := range items {
		// 遠い未来のエントリは取り込まない
		if toleranceTime.Before(*items[i].PublishedParsed) {
			continue
		}

		// UNIX epoch で表現できない時刻のエントリは取り込まない
		if items[i].PublishedParsed.Unix() < 0 {
			continue
		}

		// フィードの記事データからフィードの最終更新時間を更新する
		// latestEntryAt が nil の場合も更新する
		if latestEntryAt == nil || items[i].PublishedParsed.After(*latestEntryAt) {
			latestEntryAt = items[i].PublishedParsed
		}

		// 新しいもののみを取り込む XXX デバッグ時は以下を抑止
		if items[i].PublishedParsed.Before(existedLatestEntryAt) || items[i].PublishedParsed.Equal(existedLatestEntryAt) {
			continue
		}

		u, err := url.Parse(items[i].Link)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err = SaveEntries(ctx, tx, feedRow, items[i], base.ResolveReference(u).String()); err != nil {
			tx.Rollback()
			return err
		}
	}

	feedRow.Cache = toDbTypeCacheJson(res)
	feedRow.HTTPStatus = strconv.Itoa(res.StatusCode)
	if latestEntryAt == nil {
		feedRow.SetRevisitAt(revisitManager.RevisitAtMaxRevisitDuration())
	} else {
		feedRow.SetLatestEntryAt(*latestEntryAt)
		feedRow.SetRevisitAt(revisitManager.RevisitAt(*latestEntryAt))
	}

	if err = tx.UpdateFeedWithLatestEntryAt(ctx, *feedRow); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	logger.Info("Updated Feed")
	return nil
}

func substringByBytesLength(s string, blen int) string {
	var output string
	for _, i := range []rune(s) {
		if len(output)+len(string(i)) > blen {
			break
		}
		output += string(i)
	}
	return output
}

func SaveEntries(ctx context.Context, tx *dbType.ClientTxn, feedRow *dbType.Feed, item *gofeed.Item, link string) error {
	if err := tx.UpdateNextSerial(ctx, feedRow.ID); err != nil {
		return err
	}
	serial, err := tx.GetNextSerial(ctx, feedRow.ID)
	if err != nil {
		return err
	}

	subs, err := tx.SubscriptionByFeedID(ctx, feedRow.ID)
	if err != nil {
		return err
	}
	for j := range subs {
		if err = tx.InsertEntry(ctx, subs[j].UserID, feedRow.ID, *serial, subs[j].ID, *item.PublishedParsed); err != nil {
			return err
		}
	}

	return tx.InsertStory(ctx, feedRow.ID, *serial, substringByBytesLength(item.Title, titleLength), substringByBytesLength(item.Description, 255), link)
}
