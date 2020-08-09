package seed

import (
	"log"

	"github.com/jinzhu/gorm"

	"github.com/daguito81/fullstack/api/models"
)

var users = []models.User{
	{
		Nickname: "Dago Romer",
		Email:    "dago@gmail.com",
		Password: "password",
	},
	{
		Nickname: "Cristy Pereira",
		Email:    "cristy@gmail.com",
		Password: "alsopassword",
	},
}

var posts = []models.Post{
	{
		Title:   "Title 1",
		Content: "Hello World 1",
	},
	{
		Title:   "Title 2",
		Content: "Hello World 2",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table %v", err)
	}
	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key errror: %v", err)
	}
	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}

}
