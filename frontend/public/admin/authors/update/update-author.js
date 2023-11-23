async function getAuthorToUpdate(id) {
    const url = "https://localhost:8080/authors?id=" + id
    const response = await fetch(url, {method: "GET"})
    if (!response.ok) {
        throw new Error(`Failed to get author with ID ${id}`)
    }
    const json = await response.json()
    return json.author
}

function getImageExtension(filename) {
    return filename.split(".")[1]
}

const params = new URLSearchParams(window.location.search)
const strID = params.get("id")
if (strID != null) {
    const authorID = parseInt(strID)
    getAuthorToUpdate(authorID)
        .then(author => {
            console.log(author)
            let fNameField = document.getElementById("authorFirstName")
            let lNameField = document.getElementById("authorLastName")
            let bioField = document.getElementById("authorBio")
            let headshotPreview = document.getElementById("headshotPreview")

            fNameField.value = author.firstName
            lNameField.value = author.lastName
            bioField.value = author.bio

            extension = getImageExtension(author.headshot.name)
            headshotPreview.src = `data:image/${extension};base64,${author.headshot.content}`


        })
        .catch(err => console.log(err))
} else {
    console.log("Must have ID")
}