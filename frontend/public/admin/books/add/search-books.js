books = []

function searchBook(title, author) {
    const reqTitle = encodeURI(title)
    const reqAuthor = encodeURI(author)

    const xhr = new XMLHttpRequest();
    xhr.open("GET", `search_books.php?title=${reqTitle}&author=${reqAuthor}`)
    xhr.setRequestHeader("Content-type", "application/json")
    xhr.onload = function() {
        console.log(xhr.responseText)
        const res = JSON.parse(xhr.responseText)
        const resultTableBody = document.getElementById("googleSearchResults").getElementsByTagName("tbody")[0]
        books = res.books
        for (let i = 0; i < books.length; i++) {
            let row = resultTableBody.insertRow(-1)
            insertRowContent(row, books[i])
            row.id = "searchresult-" + i
            row.addEventListener("click", event => {
                const targetRow = event.target.closest("tr");
                if (targetRow) {
                    id = targetRow.id.split("-")[1]
                    populateAddForm(id)
                }
            })
        }
    }
    xhr.send()
}

function insertRowContent(row, book) {
    const titleCell = row.insertCell(0)
    titleCell.textContent = book.title

    const authorsCell = row.insertCell(1)
    authorsCell.textContent = book.authors

    const publishDateCell = row.insertCell(2)
    publishDateCell.textContent = book.publishDate

    const descriptionCell = row.insertCell(3)
    descriptionCell.textContent = book.description.substring(0, Math.min(200, book.description.length))
}

function populateAddForm(id) {
    document.getElementById("addBookTitle").value = books[id].title
    document.getElementById("addBookAuthors").value = books[id].authors
    document.getElementById("addBookPublishDate").value =  books[id].publishDate
    document.getElementById("addBookPageCount").value = books[id].pageCount
    document.getElementById("addBookCategories").value = books[id].categories
    document.getElementById("addBookDescription").value = books[id].description
}

const searchForm = document.getElementById("googleBookSearchForm")

searchForm.addEventListener("submit", event => {
    event.preventDefault()

    const searchTitle = document.getElementById("googleBookTitle").value
    const searchAuthor = document.getElementById("googleBookAuthor").value

    searchBook(searchTitle, searchAuthor)
})