package routes

import (
	"github.com/gin-gonic/gin"
	"qq-music-api/controllers"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", controllers.HomePage)

	r.GET("/user/cookie", controllers.Cookie)
	r.POST("/user/setCookie", controllers.SetCookie)
	r.GET("/user/getCookie", controllers.GetCookie)
	r.POST("/user/refresh", controllers.Refresh)

	r.GET("/search", controllers.Search)
	r.GET("/search/hot", controllers.HotSearch)
	r.GET("/search/quick", controllers.QuickSearch)

	r.GET("/album", controllers.AlbumInfo)
	r.GET("/album/songs", controllers.AlbumSongs)

	r.GET("/lyric", controllers.GetLyric)

	r.GET("/mv", controllers.GetMVInfo)
	r.GET("/mv/url", controllers.GetMVUrl)
	r.GET("/mv/category", controllers.GetMVCategory)
	r.GET("/mv/list", controllers.GetMVList)
	r.POST("/mv/like", controllers.LikeMV)

	r.GET("/new/songs", controllers.GetNewSongs)
	r.GET("/new/album", controllers.GetNewAlbum)
	r.GET("/new/mv", controllers.GetNewMV)

	r.GET("radio", controllers.GetRadio)
	r.GET("radio/category", controllers.GetRadioCategory)

	r.GET("/recommend/playlist", controllers.GetRecommendPlaylist)
	r.GET("/recommend/playlist/u", controllers.GetRecommendPlaylistByUser)
	r.GET("/recommend/daily", controllers.GetRecommendDaily)
	r.GET("/recommend/banner", controllers.GetRecommendBanner)

	r.GET("/singer/desc", controllers.GetSingerDesc)
	r.GET("/singer/album", controllers.GetSingerAlbum)
	r.GET("/singer/songs", controllers.GetSingerSongs)

	r.GET("/song", controllers.GetSongDetail)
	r.GET("/song/url", controllers.GetSongDownloadURL)
	r.GET("/song/urls", controllers.GetSongPlayURL)
}
