async function deleteBook(id) {
    const xhr = new XMLHttpRequest();
    xhr.open("DELETE", `delete_book.php?id=${id}`)
    xhr.onload = function() {
        if (xhr.status != 200) {
            throw new Error(xhr.statusText)
        } else {
            window.location.reload()
        }
    }
    xhr.send()
}

const buttons = document.getElementsByTagName("button")
for (let i = 0; i < buttons.length; i++) {
    buttonType = buttons[i].id.split("-")[0]
    if (buttonType == "update") {
        buttons[i].addEventListener("click", event => {
            const buttonID = event.target.id
            const bookID = buttonID.split("-")[1]
            window.location.href = `https://localhost:8082/admin/books/update?id=${bookID}` 
        })
    } else if (buttonType == "delete") {
        buttons[i].addEventListener("click", event => {
            const buttonID = event.target.id
            const bookID = buttonID.split("-")[1]
            deleteBook(bookID).catch(err => {console.error(err)})
        })
    }
}