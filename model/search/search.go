package search

type Response struct {
	Code    int    `json:"code"`
	Data    Data   `json:"data"`
	Message string `json:"message"`
	Notice  string `json:"notice"`
	SubCode int    `json:"subcode"`
	Time    int64  `json:"time"`
	Tips    string `json:"tips"`
}

type Data struct {
	Keyword   string   `json:"keyword"`
	Priority  int      `json:"priority"`
	QC        []string `json:"qc"`
	Semantic  Semantic `json:"semantic"`
	Song      SongData `json:"song"`
	Tab       int      `json:"tab"`
	TagList   []string `json:"taglist"`
	TotalTime int      `json:"totaltime"`
	ZhiDa     ZhiDa    `json:"zhida"`
}

type Semantic struct {
	CurNum   int      `json:"curnum"`
	CurPage  int      `json:"curpage"`
	List     []string `json:"list"`
	TotalNum int      `json:"totalnum"`
}

type SongData struct {
	CurNum   int    `json:"curnum"`
	CurPage  int    `json:"curpage"`
	List     []Song `json:"list"`
	TotalNum int    `json:"totalnum"`
}

type Song struct {
	AlbumID          int      `json:"albumid"`
	AlbumMID         string   `json:"albummid"`
	AlbumName        string   `json:"albumname"`
	AlbumNameHiLight string   `json:"albumname_hilight"`
	AlertID          int      `json:"alertid"`
	ChineseSinger    int      `json:"chinesesinger"`
	DocID            string   `json:"docid"`
	Format           string   `json:"format"`
	Grp              []string `json:"grp"`
	Interval         int      `json:"interval"`
	IsOnly           int      `json:"isonly"`
	Lyric            string   `json:"lyric"`
	LyricHiLight     string   `json:"lyric_hilight"`
	MsgID            int      `json:"msgid"`
	NewStatus        int      `json:"newStatus"`
	NT               int64    `json:"nt"`
	Pay              Pay      `json:"pay"`
	Preview          Preview  `json:"preview"`
	PubTime          int64    `json:"pubtime"`
	Pure             int      `json:"pure"`
	Singer           []Singer `json:"singer"`
	Size128          int      `json:"size128"`
	Size320          int      `json:"size320"`
	SizeAPE          int      `json:"sizeape"`
	SizeFLAC         int      `json:"sizeflac"`
	SizeOGG          int      `json:"sizeogg"`
	SongID           int      `json:"songid"`
	SongMID          string   `json:"songmid"`
	SongName         string   `json:"songname"`
	SongNameHiLight  string   `json:"songname_hilight"`
	SongURL          string   `json:"songurl"`
	Stream           int      `json:"stream"`
	Switch           int64    `json:"switch"`
	T                int      `json:"t"`
	Tag              int      `json:"tag"`
	Type             int      `json:"type"`
	Version          int      `json:"ver"`
	VID              string   `json:"vid"`
}

type Pay struct {
	PayAlbum      int `json:"payalbum"`
	PayAlbumPrice int `json:"payalbumprice"`
	PayDownload   int `json:"paydownload"`
	PayInfo       int `json:"payinfo"`
	PayPlay       int `json:"payplay"`
	PayTrackMouth int `json:"paytrackmouth"`
	PayTrackPrice int `json:"paytrackprice"`
}

type Preview struct {
	TryBegin int `json:"trybegin"`
	TryEnd   int `json:"tryend"`
	TrySize  int `json:"trysize"`
}

type Singer struct {
	ID          int    `json:"id"`
	MID         string `json:"mid"`
	Name        string `json:"name"`
	NameHiLight string `json:"name_hilight"`
}

type ZhiDa struct {
	ChineseSinger int `json:"chinesesinger"`
	Type          int `json:"type"`
}
