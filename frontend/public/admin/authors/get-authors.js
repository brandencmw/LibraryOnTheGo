async function getAuthors() {
    const url = "https://localhost:8080/authors?includeimages=false"

    let response = await fetch(url, {method: "GET"})
    if (response.status != 200) {
        throw new Error("There was a problem fetching from the server")
    }
    let json = await response.json()
    console.log(json)
    return json.authors
}

async function deleteAuthor(id) {
    const url = "https://localhost:8080/authors/auth/delete?id="+id
    options = {
        method: "DELETE",
        credentials: "include",
        body: JSON.stringify({"id": parseInt(id)})
    }
    let response = await fetch(url, options)
    if (response.status != 200) {
        throw new Error(`Failed to delete author with ID ${id}`)
    }
}

function insertRowContent(row, author) {
    let idCell = row.insertCell(0)
    idCell.textContent = author.id

    let fNameCell = row.insertCell(1)
    fNameCell.textContent = author.firstName

    let lNameCell = row.insertCell(2)
    lNameCell.textContent = author.lastName

    let actionCell = row.insertCell(3)

    let deleteButton = document.createElement("button")
    deleteButton.textContent = "Delete"
    deleteButton.id = `delete-${author.id}`
    deleteButton.addEventListener("click", event => {
        let buttonID = event.target.id
        let authorID = buttonID.split("-")[1]
        deleteAuthor(authorID)
    })

    let updateButton = document.createElement("button")
    updateButton.textContent = "Update"
    updateButton.id = `update-${author.id}`
    updateButton.addEventListener("click", event => {
        let buttonID = event.target.id
        let authorID = buttonID.split("-")[1]
        window.location.href = `/admin/authors/update?id=${authorID}`
    })

    actionCell.appendChild(updateButton)
    actionCell.appendChild(deleteButton)
}

getAuthors().catch(err => {
    console.log(err)
}).then(authors => {
    const tableBody = document.getElementById("authorTable").getElementsByTagName("tbody")[0]
    authors.forEach(author => {
        let row = tableBody.insertRow(-1)
        insertRowContent(row, author)
    });
})