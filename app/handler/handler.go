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

func (*ApiServer) Profile(ctx context.Context, request ProfileRequestObject) (ProfileResponseObject, error) {
	pin, err := DBUserFromContext(ctx).Profile()
	if err != nil {
		return Profile400Response{}, nil
	}
	return Profile200JSONResponse(Profile{
		Autoseen:  pin.AutoSeen,
		Nopinlist: pin.NoPinList,
		Numentry:  pin.EntryCount,
		Numsubstr: pin.SubstringLength,
	}), nil
}

func (*ApiServer) Categories(ctx context.Context, request CategoriesRequestObject) (CategoriesResponseObject, error) {
	cats, err := DBUserFromContext(ctx).Category()
	if err != nil {
		return Categories400Response{}, nil
	}

	items := []Category{}
	for _, i := range cats {
		items = append(items, Category{
			ID:   i.ID,
			Name: i.Name,
		})

	}
	return Categories200JSONResponse(items), nil
}

func (*ApiServer) CategoryAndUnreadEntryCount(ctx context.Context, request CategoryAndUnreadEntryCountRequestObject) (CategoryAndUnreadEntryCountResponseObject, error) {
	cat, err := DBUserFromContext(ctx).CategoryAndUnreadEntryCount()
	if err != nil {
		return CategoryAndUnreadEntryCount400Response{}, nil
	}

	items := []CategoryAndUnreadEntryCount{}
	for _, i := range cat {
		items = append(items, CategoryAndUnreadEntryCount{
			ID:    i.ID,
			Name:  i.Name,
			Count: i.Count,
		})
	}

	return CategoryAndUnreadEntryCount200JSONResponse(items), nil
}
