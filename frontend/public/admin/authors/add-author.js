
const addAuthorForm = document.getElementById("addAuthorForm")

addAuthorForm.addEventListener("submit", event => {
    event.preventDefault()
    const formData = new FormData()
    formData.append("firstName", document.getElementById("authorFirstName").value)
    formData.append("lastName", document.getElementById("authorLastName").value)
    formData.append("bio", document.getElementById("authorBio").value)
    formData.append("headshot", document.getElementById("authorHeadshot").files[0])

    sendDataToBackend(formData, "/authors/auth/create", "POST")
        .then(response => console.log(response))
        .catch(err => console.error(err))
    
})