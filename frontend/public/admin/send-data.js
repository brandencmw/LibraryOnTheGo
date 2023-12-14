async function sendDataToBackend(formData, endpoint, method) {
    const url = "https://localhost:8080" + endpoint
    const options = {
        method: method,
        credentials: "include",
        body: formData,
    }

    
    const response = await fetch(url, options)
    const responseJSON = await response.json()
    if (!response.ok) {
        throw new Error(responseJSON)
    }

    return responseJSON
}