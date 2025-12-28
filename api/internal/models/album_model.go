package models

import (
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID          uint       `gorm:"type:int;primaryKey"`
	Name        string     `gorm:"type:varchar(255);not null"`
	Description string     `gorm:"type:text;null"`
	UserID      *uuid.UUID `gorm:"type:char(36);null"`                             // Foreign Key, nullable
	User        User       `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"` // Relation
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Many-to-Many Relation with Media
	Media      []Media      `gorm:"many2many:album_media;constraint:OnDelete:CASCADE;"`
	AlbumMedia []AlbumMedia `gorm:"foreignKey:AlbumID"`

	// Computed field for media count
	MediaCount int `gorm:"-" json:"mediaCount"`
}

type AlbumMedia struct {
	AlbumID uint  `gorm:"primaryKey"`
	MediaID uint  `gorm:"primaryKey"`
	IsCover bool  `gorm:"default:false"`
	Album   Album `gorm:"foreignKey:AlbumID;constraint:OnDelete:CASCADE;"`
	Media   Media `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE;"`
}

/*
Note: gorm have problems with cascading deletes in many2many relations.
To ensure that deleting an album also removes its associations in the album_media table,
we need to manually define the AlbumMedia model and set up the foreign key constraints with ON DELETE CASCADE.
The following SQL commands should be run in the database to update the foreign key constraint:

```sql
ALTER TABLE album_media DROP FOREIGN KEY fk_albums_album_media;
ALTER TABLE album_media
ADD CONSTRAINT fk_albums_album_media
FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE;
```
*/
