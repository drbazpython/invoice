package database_test

//
import (
	"testing"
	"drbaz.com/invoice/cmd/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var Posts = []database.Post{
		{
			Title: "title_1",
			Desc: "desc_1",
			Likes: 1,
		},
		{
			Title: "title_2",
			Desc: "desc_2",
			Likes: 2,
		},
	}

func TestDBInit(t *testing.T){
	db := database.DBInit()

	// delete all rows
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&database.Post{})
	
	var actualName string = db.Name()
	assert.Equal(t,"sqlite",actualName,"they should be equal")
}

func TestDBCreateItem(t *testing.T){
	db := database.DBInit()
	count1 := database.DBGetAllItems(db).RowsAffected
	newPost := Posts[0]
	_ = database.DBCreateItem(db, newPost)
	count2 := database.DBGetAllItems(db).RowsAffected
	assert.Greater(t, count2, count1)
}

func TestDBGetFirstPost(t *testing.T){
	db := database.DBInit()
	post := database.DBGetFirstPost(db)
	assert.Equal(t, Posts[0].Title,post.Title)
}
 
func TestDBGetItemBYID(t *testing.T){
	db := database.DBInit()
	//get id of first record
	firstPost := database.DBGetFirstPost(db)
	var id = int(firstPost.ID)
	post := database.DBGetItemByID(db,id)
	assert.Equal(t, Posts[0].Title, post.Title)
}

func TestDBGetItem(t *testing.T){
	db := database.DBInit()
	posts := database.DBGetItem(db, Posts[0].Title)
	assert.Equal(t, Posts[0].Title, posts[0].Title)
}

func TestDBUpdateItem(t *testing.T){
	db := database.DBInit()
	var post1 = Posts[0]
	var post2 = Posts[1]

	assert.True(t,database.DBUpdateItem(db,post1, post2 ))
}

func TestDBDeleteItem(t *testing.T){
	db := database.DBInit()
	firstPost := database.DBGetFirstPost(db)
	assert.True(t,database.DBDeleteItem(db,int(firstPost.ID)))
}

