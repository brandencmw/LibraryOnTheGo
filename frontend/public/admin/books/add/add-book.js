
const addBookForm = document.getElementById("addBookForm")

addBookForm.addEventListener("submit", event => {
    event.preventDefault()

    const data = new FormData()
    data.append("title", document.getElementById("addBookTitle").value)
    data.append("synopsis", document.getElementById("addBookDescription").value)
    data.append("publishDate", document.getElementById("addBookPublishDate").value)
    data.append("pageCount", document.getElementById("addBookPageCount").value)

    const categories = document.getElementById("addBookCategories").value.split(",")
    data.append("categories", categories)

    const authors = document.getElementById("addBookAuthors").value.split(",")
    data.append("authors", authors)

    data.append("cover", document.getElementById("addBookCover").files[0])

    sendDataToBackend(data, "/books/auth/create", "POST")
        .catch(err => {
            console.log(err)
        })
        .then(res => {
            console.log(res)
        })
})

// type addBookRequest struct {
// 	Title       string               `form:"title" binding:"required"`
// 	Synopsis    string               `form:"synopsis"`
// 	PublishDate string               `form:"publishDate" binding:"required"`
// 	PageCount   int                  `form:"pageCount" binding:"required"`
// 	Categories  []string             `form:"categories"`
// 	Authors     []string             `form:"authors"`
// 	Cover       multipart.FileHeader `form:"cover" binding:"required"`
// }