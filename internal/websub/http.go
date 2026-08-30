package websub

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mmcdole/gofeed"

	"github.com/yseto/gion-go/internal/config"
	"github.com/yseto/gion-go/internal/crawler"
	abstractDB "github.com/yseto/gion-go/internal/db"
	"github.com/yseto/gion-go/internal/db/db"
)

func CreateHandler(cfg *config.Config, dbConn *sqlx.DB) http.Handler {
	mux := http.NewServeMux()

	h := handler{db: abstractDB.New(dbConn)}

	mux.HandleFunc("GET /api/webhook/pubsub", h.getCallback)
	mux.HandleFunc("POST /api/webhook/pubsub", h.postCallback)

	return mux
}

type handler struct {
	db abstractDB.DB
}

func (h *handler) getCallback(w http.ResponseWriter, r *http.Request) {
	topicURL := r.URL.Query().Get("hub.topic")
	if topicURL == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	challenge := r.URL.Query().Get("hub.challenge")
	if challenge == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mode := r.URL.Query().Get("hub.mode")
	if !slices.Contains([]string{"subscribe", "unsubscribe"}, mode) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !topicURLs.contains(topicURL) {
		slog.Error("!topicURLs.contains")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if mode == "unsubscribe" {
		if !topicURLs.delete(topicURL) {
			slog.Error("!topicURLs.delete")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		slog.Info("getCallback: unsubscribe", slog.String("topicURL", topicURL))

		w.Header().Add("content-type", "text/plain")
		w.Write([]byte(challenge))
		return
	}

	// ここからは登録、更新

	leaseSeconds, err := strconv.Atoi(r.URL.Query().Get("hub.lease_seconds"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var (
		ctx  = r.Context()
		tx   = h.db.MustBegin()
		item *db.Websub
	)

	item, err = tx.WebsubByTopicURL(ctx, topicURL)
	if err != nil && err != sql.ErrNoRows { // エントリが存在しない以外のエラー
		slog.Error("tx.WebsubByTopicURL", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		tx.Rollback()
		return
	}

	if item == nil {
		// 新規作成
		if err := tx.InsertWebsub(ctx, topicURL, time.Now().Add(time.Second*time.Duration(leaseSeconds))); err != nil {
			slog.Error("tx.InsertWebsub", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			tx.Rollback()
			return
		}
	} else {
		// 期限切れ更新の場合
		item.ExpiredAt = time.Now().Add(time.Second * time.Duration(leaseSeconds))
		if err := tx.UpdateWebsub(ctx, *item); err != nil {
			slog.Error("tx.UpdateWebsub", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			tx.Rollback()
			return
		}
	}

	if !topicURLs.delete(topicURL) {
		slog.Error("!topicURLs.delete")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		slog.Error("tx.Commit", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	slog.Info("getCallback: subscribe", slog.String("topicURL", topicURL))
	w.Header().Add("content-type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte(challenge))
}

func (h *handler) postCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.WriteHeader(200)
	w.Write([]byte("OK"))

	doc, err := gofeed.NewParser().Parse(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed parse", slog.String("error", err.Error()))
		return
	}

	if len(doc.Items) == 0 {
		return
	}

	tx := h.db.MustBegin()

	websub, err := tx.WebsubByTopicURL(ctx, doc.FeedLink)
	if err != nil {
		tx.Rollback()
		slog.ErrorContext(ctx, "postCallback: WebsubByTopicURL", slog.String("error", err.Error()))
		return
	}

	if err = processDoc(ctx, tx, doc, websub.ID); err != nil {
		tx.Rollback()
		slog.ErrorContext(ctx, "postCallback: processDoc", slog.String("error", err.Error()))
		return
	}

	tx.Commit()
}

func processDoc(ctx context.Context, tx *db.ClientTxn, doc *gofeed.Feed, websubID uint64) error {
	feedRow, err := tx.FeedByWebsubID(ctx, websubID)
	if err != nil {
		slog.ErrorContext(ctx, "processDoc: FeedByID", slog.String("error", err.Error()))
		return err
	}

	var latestEntryAt *time.Time
	if ts := feedRow.GetLatestEntryAt(); ts != nil {
		latestEntryAt = ts
	}

	for _, item := range doc.Items {
		if item.PublishedParsed == nil {
			continue
		}

		// UNIX epoch で表現できない時刻のエントリは取り込まない
		if item.PublishedParsed.Unix() < 0 {
			continue
		}

		// フィードの記事データからフィードの最終更新時間を更新する
		// latestEntryAt が nil の場合も更新する
		if latestEntryAt == nil || item.PublishedParsed.After(*latestEntryAt) {
			latestEntryAt = item.PublishedParsed
		}

		// 古すぎる更新を排除したいので、updatedとpublishedの差が12時間以上あれば、スキップする
		if item.UpdatedParsed != nil && item.UpdatedParsed.Sub(*item.PublishedParsed) > 12*time.Hour {
			continue
		}

		if err := crawler.SaveEntries(ctx, tx, feedRow, item, item.Link); err != nil {
			return err
		}
	}

	// 通常のクローラーと同程度に整える
	if latestEntryAt != nil {
		feedRow.SetLatestEntryAt(*latestEntryAt)
	}
	feedRow.HTTPStatus = "200"
	feedRow.SetRevisitNull()
	feedRow.Cache = db.CacheJson{}

	return tx.UpdateFeedWithLatestEntryAt(ctx, *feedRow)
}
