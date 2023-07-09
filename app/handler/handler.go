package handler

type ApiServer struct{}

func NewApiServer() *ApiServer {
	return &ApiServer{}
}

var _ StrictServerInterface = (*ApiServer)(nil)

func (*ApiServer) PinnedItems(ctx context.Context, request PinnedItemsRequestObject) (PinnedItemsResponseObject, error) {
	pins, err := DBUserFromContext(ctx).PinnedItems()
	if err != nil {
		return PinnedItems400Response{}, nil
	}

	items := []PinnedItem{}
	for _, i := range pins {
		items = append(items, PinnedItem{
			FeedId:   i.EntryFeedID,
			Serial:   i.EntrySerial,
			Title:    i.Title,
			UpdateAt: i.EntryUpdateAt,
			Url:      i.URL,
		})
	}

	return PinnedItems200JSONResponse(items), nil
}
