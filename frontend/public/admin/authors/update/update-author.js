let valueMap = new Map()
let fNameField = document.getElementById("authorFirstName")
valueMap.set(fNameField.id, fNameField.value)

let lNameField = document.getElementById("authorLastName")
valueMap.set(lNameField.id, lNameField.value)

let bioField = document.getElementById("authorBio")
valueMap.set(bioField.id, bioField.value)

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

    let urlParams = new URLSearchParams(window.location.search)
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
            .then(response => {
                console.log(response)
                location.reload()
            })
            .catch(err => console.error(err))
    } else {
        console.log("No entries have changed")
    }  
})

const backButton = document.getElementById("backButton")
backButton.addEventListener("click", _ => goBack())