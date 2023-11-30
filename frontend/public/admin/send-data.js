async function sendDataToBackend(formData, endpoint, method) {
    url = "https://localhost:8080" + endpoint
    const options = {
        method: method,
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