addAuthorForm = document.getElementById("addAuthorForm")

addAuthorForm.addEventListener("submit", event => {
    event.preventDefault()
    const formData = new FormData()
    formData.append("firstName", document.getElementById("authorFirstNameInput").value)
    formData.append("lastName", document.getElementById("authorLastNameInput").value)
    formData.append("bio", document.getElementById("authorBioInput").value)
    formData.append("headshot", document.getElementById("authorHeadshotUpload").files[0])

    console.log("FORM SUBMITTED")

    sendDataToBackend(formData, "/authors/auth/create", "POST")
        .then(response => console.log(response))
        .catch(err => console.error(err))
    
})