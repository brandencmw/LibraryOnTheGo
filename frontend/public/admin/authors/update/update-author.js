async function getAuthorToUpdate(id) {
    const url = "https://localhost:8080/authors?id=" + id
    const response = await fetch(url, {method: "GET"})
    if (!response.ok) {
        throw new Error(`Failed to get author with ID ${id}`)
    }
    const json = await response.json()
    return json.author
}

const params = new URLSearchParams(window.location.search)
const strID = params.get("id")
let valueMap = new Map()
if (strID != null) {
    const authorID = parseInt(strID)
    getAuthorToUpdate(authorID)
        .then(author => {
            let fNameField = document.getElementById("authorFirstName")
            fNameField.value = author.firstName
            valueMap.set(fNameField.id, author.firstName)

            let lNameField = document.getElementById("authorLastName")
            lNameField.value = author.lastName
            valueMap.set(lNameField.id, author.lastName)

            let bioField = document.getElementById("authorBio")
            bioField.value = author.bio
            valueMap.set(bioField.id, author.bio)

            let headshotPreview = document.getElementById("headshotPreview")
            headshotPreview.src = `https://library-pictures.s3.amazonaws.com/${author.headshotKey}`
        })
        .catch(err => console.log(err))
} else {
    console.log("Must have ID")
}

function goBack() {
    window.history.back()
}

const updateAuthorForm = document.getElementById("updateAuthorForm")

updateAuthorForm.addEventListener("submit", event => {
    event.preventDefault()
    const formData = new FormData()
    const firstName = document.getElementById("authorFirstName").value
    const lastName = document.getElementById("authorLastName").value
    const bio = document.getElementById("authorBio").value
    const headshot = document.getElementById("authorHeadshot")
    let fields = 0

    let urlParams = new URLSearchParams(window.location.href)
    formData.append("id", urlParams.get("id"))

    if (headshot.files.length > 0) {
        fields++
        formData.append("headshot", headshot.files[0])
    }
    if (firstName != valueMap.get("authorFirstName")) {
        fields++
        formData.append("firstName", firstName)
    }
    if (lastName != valueMap.get("authorLastName")) {
        fields++
        formData.append("lastName", lastName)
    }
    if (bio != valueMap.get("authorBio")) {
        fields++
        formData.append("bio", bio)
    }

    if (fields > 0) {
        sendDataToBackend(formData, "/authors/auth/update", "PUT")
            .then(response => console.log(response))
            .catch(err => console.error(err))
    } else {
        console.log("No entries have changed")
    }  
})

const backButton = document.getElementById("backButton")
backButton.addEventListener("click", _ => goBack())