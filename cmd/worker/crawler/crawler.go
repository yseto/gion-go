package crawler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mmcdole/gofeed"

	"github.com/yseto/gion-go/config"
	"github.com/yseto/gion-go/db"
	dbType "github.com/yseto/gion-go/db/db"
	"github.com/yseto/gion-go/internal/client"
	"github.com/yseto/gion-go/internal/crawler/revisit"
)

const (
	TypeCrawler = "crawler"

	crawlerFactor      = 30
	minRevisitDuration = 30 * time.Minute
	maxRevisitDuration = 7 * 24 * time.Hour
)

type CrawlerPayload struct {
	FeedID uint64
}

func NewCrawlerTask(feedID uint64) (*asynq.Task, error) {
	payload, err := json.Marshal(CrawlerPayload{feedID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeCrawler, payload), nil
}

func toDbTypeCacheJson(u *client.Response) dbType.CacheJson {
	if u == nil {
		return dbType.CacheJson{}
	}
	return dbType.CacheJson{Etag: u.Cache.Etag, Modified: u.Cache.Modified}
}

func HandleCrawlerTask(ctx context.Context, t *asynq.Task) error {
	cfg := ctx.Value("cfg").(*config.Config)

	var p CrawlerPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	dbConn, err := db.Open(cfg)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	dbc := db.New(dbConn)
	tx := dbc.MustBegin()

	feedRow, err := tx.FeedByID(p.FeedID)
	if err != nil && err == sql.ErrNoRows {
		tx.Rollback()
		return nil
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	revisitManager := revisit.New(time.Now(), crawlerFactor, minRevisitDuration, maxRevisitDuration)

	fmt.Printf(">>> %s\n", feedRow.URL)
	res, err := client.GetWithCache(
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
		fmt.Printf("<<< ERR %s : %v\n", feedRow.URL, err)
		if err = tx.UpdateFeed(*feedRow); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}

	// リダイレクト先を保存する
	if res.Location != "" {
		fmt.Printf("<<< 301 [Feed] %s -> %s\n", feedRow.URL, res.Location)
		feedRow.URL = res.Location
		if err = tx.UpdateFeedRSSUrl(*feedRow); err != nil {
			tx.Rollback()
			return err
		}

		if resSite, err := client.GetRedirectURL(feedRow.SiteURL); err == nil {
			fmt.Printf("<<< 301 [Site] %s -> %s\n", feedRow.SiteURL, resSite)
			feedRow.SiteURL = resSite
			if err = tx.UpdateFeedSiteUrl(*feedRow); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	fmt.Printf("<<< %3d %s\n", res.StatusCode, feedRow.URL)

	// 304 であることを記録
	if res.StatusCode == 304 {
		feedRow.HTTPStatus = "304"

		if ts := feedRow.GetLatestEntryAt(); ts != nil {
			feedRow.SetRevisitAt(revisitManager.RevisitAt(*ts))
		} else {
			feedRow.SetRevisitAt(revisitManager.RevisitAtMaxRevisitDuration())
		}

		// 304 応答で、ヘッダによるレスポンス返却がある場合にのみ上書きをする
		if res.Cache.Etag != "" || res.Cache.Modified != "" {
			feedRow.Cache = toDbTypeCacheJson(res)
		}
		fmt.Println("<<< 304")
		if err = tx.UpdateFeed(*feedRow); err != nil {
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
		fmt.Printf("<<< ERR Parser: %v\n", err)
		if err = tx.UpdateFeed(*feedRow); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}

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
		fmt.Printf("<<< ERR %s : missing <pubdate>\n", feedRow.URL)
		if err = tx.UpdateFeed(*feedRow); err != nil {
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

		if err := tx.UpdateNextSerial(feedRow.ID); err != nil {
			tx.Rollback()
			return err
		}
		serial, err := tx.GetNextSerial(feedRow.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		subs, err := tx.SubscriptionByFeedID(feedRow.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for j := range subs {
			if err = tx.InsertEntry(subs[j].UserID, feedRow.ID, *serial, subs[j].ID, *items[i].PublishedParsed); err != nil {
				tx.Rollback()
				return err
			}
		}

		u, err := url.Parse(items[i].Link)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err = tx.InsertStory(feedRow.ID, *serial, SubstringByBytesLength(items[i].Title, TitleLength), SubstringByBytesLength(items[i].Description, 255), base.ResolveReference(u).String()); err != nil {
			tx.Rollback()
			return err
		}
	}

	fmt.Printf("UPDATE FEED INFO: feed_id:%d\n", feedRow.ID)

	feedRow.Cache = toDbTypeCacheJson(res)
	feedRow.HTTPStatus = strconv.Itoa(res.StatusCode)
	if latestEntryAt == nil {
		feedRow.SetRevisitAt(revisitManager.RevisitAtMaxRevisitDuration())
	} else {
		feedRow.SetLatestEntryAt(*latestEntryAt)
		feedRow.SetRevisitAt(revisitManager.RevisitAt(*latestEntryAt))
	}

	if err = tx.UpdateFeedWithLatestEntryAt(*feedRow); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
