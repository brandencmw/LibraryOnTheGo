function goBack() {
    window.history.back()
}

const updateBookForm = document.getElementById("updateBookForm")

updateBookForm.addEventListener("submit", event => {
    event.preventDefault()
    const formData = new FormData()
    const synopsis = document.getElementById("bookSynopsis").value
    const cover = document.getElementById("bookCover")
    let fields = 0

    let urlParams = new URLSearchParams(window.location.search)
    formData.append("id", urlParams.get("id"))

    if (cover.files.length > 0) {
        fields++
        formData.append("cover", cover.files[0])
    }
    if (synopsis != valueMap.get("bookSynopsis")) {
        fields++
        formData.append("synopsis", synopsis)
    }

    if (fields > 0) {
        sendDataToBackend(formData, "/books/auth/update", "PUT")
            .then(response => console.log(response))
            .catch(err => console.error(err))
    } else {
        console.log("No entries have changed")
    }  
})

let valueMap = new Map()
const synopsisField = document.getElementById("bookSynopsis")
valueMap.set(synopsisField.id, synopsisField.value)

const backButton = document.getElementById("backButton")
backButton.addEventListener("click", _ => goBack())