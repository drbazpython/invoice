package database

import (
	//"drbaz.com/invoice/cmd/database"
	"github.com/charmbracelet/log"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

//Post ...
type Post struct {
	gorm.Model
	Title string
	Desc string
	Likes uint
}
//DB ...
//var DB gorm.DB	
//Err ...
//var Err error



//DBInit ... 
func DBInit() *gorm.DB {
	//logging.Logger.Debug("Initializing database")
	var DB, Err = gorm.Open(sqlite.Open("invoice.db"), &gorm.Config{})
	if Err != nil {
		log.Fatal("Failed to connect to database")
	}
	err1 := DB.AutoMigrate(&Post{})
	if err1 != nil {
		log.Fatal("Error initialising database")
	}
	//logging.Logger.Debug("DB", "name", DB.Name())
	return DB
}

//DBCreateItem ...
func DBCreateItem(db *gorm.DB, newPost Post) Post {
	//logging.Logger.Debug("DBCreateItem")
	if res := db.Create(&newPost); res.Error != nil{
		panic(res.Error)
	}

	return newPost
}

//DBGetAllItems ...
func DBGetAllItems(db *gorm.DB) *gorm.DB {
	//logging.Logger.Debug("DBGetAllItems")
	var posts []Post
	records := db.Find(&posts)
	if records.Error != nil {
		log.Error(records.Error)
		return nil
	}
	//logging.Logger.Info("DB","rows affed ",records.RowsAffected)
	
	// for _, post := range posts {
	// 	logging.Logger.Info(post.Title)

	// }

	return records
}

//DBGetFirstPost ...
func DBGetFirstPost(db *gorm.DB) Post {
	//logging.Logger.Debug("DBGetItem")
	var post Post
    result := db.First(&post)
    if result.Error != nil {
        panic("failed to retrieve post: " + result.Error.Error())
    }
	//logging.Logger.Debug(post.Title)
	return post
}
//DBGetItemByID ...
func DBGetItemByID(db *gorm.DB, id int) Post{
	var post Post
	result := db.Where("ID = ?", id).Find(&post)
	if result.Error != nil {
    	panic("failed to retrieve post: " + result.Error.Error())
	}
	//logging.Logger.Debug(post)
	return post
}

//DBGetItem ...
func DBGetItem(db * gorm.DB, title string) []Post {
	var posts []Post
	result := db.Where("Title = ?", title).Find(&posts)
	if result.Error != nil {
    	panic("failed to retrieve post: " + result.Error.Error())
	}

	// for _, post := range posts {
    // 	logging.Logger.Debug("DBGetItem","Post ID", post.ID)
	// 	logging.Logger.Debug("DBGetItem","Title", post.Title)
	// 	logging.Logger.Debug("DBGetItem","Description", post.Desc) 
	// 	logging.Logger.Debug("DBGetItem", "Likes", post.Likes)
	// }
	return posts
	
}

//DBUpdateItem ...
func DBUpdateItem(db *gorm.DB, post1 Post, post2 Post) bool {
	// Retrieve the record you want to update
    result := db.First(&post1)
    if result.Error != nil {
        panic("failed to retrieve post: " + result.Error.Error())
    }
	post1.Title = post2.Title
	post1.Desc = post2.Desc
	post1.Likes = post2.Likes
    // Save the changes back to the database
    result = db.Save(&post1)
    if result.Error != nil {
        panic("failed to update post: " + result.Error.Error())
    }

    log.Debug("Post updated successfully")
	return true
}

//DBDeleteItem ...
func DBDeleteItem(db *gorm.DB, id int) bool {
	var post Post
	post =  DBGetItemByID(db,id) // post to delete
	// delete permanently
	result := db.Unscoped().Delete(&post)
	if result.Error != nil {
    	panic("failed to delete post: " + result.Error.Error())
	} else if result.RowsAffected == 0 {
    	panic("no post record was deleted")
	} else {
    	log.Debug("Post deleted successfully")
		return true
	}
	}


