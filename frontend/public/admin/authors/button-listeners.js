async function deleteAuthor(id) {
    const xhr = new XMLHttpRequest();
    xhr.open("DELETE", `delete_author.php?id=${id}`)
    xhr.onload = function() {
        if (xhr.status != 200) {
            throw new Error(xhr.statusText)
        } else {
            window.location.reload()
        }
    }
    xhr.send()
}

const buttons = document.getElementsByTagName("button");
for (const button of buttons) {
    const buttonPrefix = button.id.split("-")[0]
    if (buttonPrefix == "update") {
        button.addEventListener("click", _ => {
            const id = button.id.split("-")[1]
            window.location.href = `/admin/authors/update?id=${id}`
        })
    } else if (buttonPrefix == "delete") {
        button.addEventListener("click", _ => {
            console.log("delete")
            const id = button.id.split("-")[1]
            deleteAuthor(id)
        })
    }
}