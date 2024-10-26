## Controller Response Examples

```go
c.JSON(http.StatusOK, gin.H{
	"message": "User added successfully",
	"note":    note,
	"type":    "POST",
	"user":    user,
})


c.HTML(http.StatusOK, "index.tmpl", gin.H{
	"username": "SUCCESS",
	"password": "password",
	"error":    nil,
})
```