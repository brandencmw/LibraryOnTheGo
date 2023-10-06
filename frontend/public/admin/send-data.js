async function sendDataToBackend(formData, endpoint) {
    url = "https://localhost:8080" + endpoint

    response = await fetch(url, {method: "POST", credentials: "include", body: formData})
    if (!response.ok) {
        throw new Error("Not OK")
    }
    responseJSON = await response.json()
    return responseJSON
}


addAuthorForm = document.getElementById("addAuthorForm")
addBookForm = document.getElementById("addBookForm")

addAuthorForm.addEventListener("submit", event => {
    event.preventDefault()
    const formData = new FormData()
    formData.append("firstName", document.getElementById("authorFirstNameInput").value)
    formData.append("lastName", document.getElementById("authorLastNameInput").value)
    formData.append("bio", document.getElementById("authorBioInput").value)
    formData.append("headshot", document.getElementById("authorHeadshotUpload").files[0])


    sendDataToBackend(formData, "/authors/create")
        .then(response => console.log(response))
        .catch(err => console.error(err))
    
})