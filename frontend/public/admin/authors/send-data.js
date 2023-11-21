async function sendDataToBackend(formData, endpoint) {
    url = "https://localhost:8080" + endpoint
    const options = {
        method: "POST",
        credentials: "include",
        body: formData,
    }

    response = await fetch(url, options)
    if (!response.ok) {
        throw new Error("Not OK")
    }
    responseJSON = await response.json()
    return responseJSON
}


addAuthorForm = document.getElementById("addAuthorForm")

addAuthorForm.addEventListener("submit", event => {
    event.preventDefault()
    const formData = new FormData()
    formData.append("firstName", document.getElementById("authorFirstNameInput").value)
    formData.append("lastName", document.getElementById("authorLastNameInput").value)
    formData.append("bio", document.getElementById("authorBioInput").value)
    formData.append("headshot", document.getElementById("authorHeadshotUpload").files[0])


    sendDataToBackend(formData, "/authors/auth/create")
        .then(response => console.log(response))
        .catch(err => console.error(err))
    
})