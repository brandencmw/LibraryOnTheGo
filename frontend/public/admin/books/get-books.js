async function getBooks() {
    const url = "https://localhost:8080/books?includeimages=false"

    let response = await fetch(url, {method: "GET"})
    let json = await response.json()
    if (response.status != 200) {
        throw new Error(json.error)
    }
    console.log(json)
    return json.books
}

// async function deleteBook(id) {
//     console.log(id)
//     const url = "https://localhost:8080/authors/auth/delete?id="+id
//     options = {
//         method: "DELETE",
//         credentials: "include",
//         body: JSON.stringify({})
//     }
//     let response = await fetch(url, options)
//     console.log(response)
//     if (response.status != 200) {
//         throw new Error(`Failed to delete author with ID ${id}`)
//     }
//     return response.json()
// }

function insertRowContent(row, book) {
    let idCell = row.insertCell(0)
    idCell.textContent = book.id

    let titleCell = row.insertCell(1)
    titleCell.textContent = book.title

    let authorCell = row.insertCell(2)
    authorCell.textContent = book.authors.join(", ")

    let actionCell = row.insertCell(3)

    let deleteButton = document.createElement("button")
    deleteButton.textContent = "Delete"
    deleteButton.id = `delete-${book.id}`
    deleteButton.addEventListener("click", event => {
        let buttonID = event.target.id
        let bookID = buttonID.split("-")[1]
        deleteAuthor(bookID).catch(err => {
            console.log(err)
        }).then(res => {
            console.log(res)
        })
    })

    let updateButton = document.createElement("button")
    updateButton.textContent = "Update"
    updateButton.id = `update-${book.id}`
    updateButton.addEventListener("click", event => {
        let buttonID = event.target.id
        let bookID = buttonID.split("-")[1]
        window.location.href = `/admin/authors/update?id=${bookID}`
    })

    actionCell.appendChild(updateButton)
    actionCell.appendChild(deleteButton)
}

getBooks().catch(err => {
    console.log(err)
}).then(books => {
    const tableBody = document.getElementById("bookTable").getElementsByTagName("tbody")[0]
    books.forEach(book => {
        let row = tableBody.insertRow(-1)
        insertRowContent(row, book)
    });
})