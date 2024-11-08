package database

impott (
    "context"
    "fmt"

    "github.com/danmuck/the_cookie_jar/pkg/api/models"
    "github.com/gin-gonic/gin"
)

// This is just a simple mapping 
var ValidFormats = map[string]bool{
    "jpeg": true,
    "png":  true,
    "webp": true,
    "gif":  true,
 }
 
 // AddMedia adds new media to database. Media types have a UUID, an asssociated username, and a format, which is a map of valid formats to boolean values.
 func AddMedia(uuid string, username string, format string) error {
    if !ValidFormats[format] {
        return fmt.Errorf("unsupported format: %s", format)
    }
 
    media := &models.Media{
        ID:       uuid,
        Username: username, 
        Format:   format,
    }
 
    _, err := GetCollection("media").InsertOne(context.TODO(), media)
    if err != nil {
        return fmt.Errorf("failed to upload media: %w", err)
    }
 
    return nil
 }
 //This searches for a media upload by id. Returns the media if it is found.
 func GetMedia(id string) (*models.Media, error) {
    var media *models.Media
    err := GetCollection("media").FindOne(context.TODO(), gin.H{"_id": id}).Decode(&media)
    return media, err
 }
 //This will search the DB for the media by id, returning an error if it is not found.
 func DeleteMedia(id string) error {
    result, err := GetCollection("media").DeleteOne(context.TODO(), gin.H{"_id": id})
    if err != nil {
        return err
    }
    if result.DeletedCount == 0 {
        return fmt.Errorf("media not found: %s", id)
    }
    return nil
 }