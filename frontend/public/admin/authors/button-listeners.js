async function deleteAuthor(id) {
    const res = await fetch(`https://localhost:8082/admin/authors/delete_author.php?id=${id}`, {method: "DELETE"})
    console.log(res)
}

const buttons = document.getElementsByTagName("button");
for (const button of buttons) {
    console.log(button)
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